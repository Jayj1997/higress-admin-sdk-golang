// Package service provides business services for the SDK
package service

import (
	"context"
	"strings"

	"github.com/Jayj1997/higress-admin-sdk-golang/internal/kubernetes"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/constant"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/errors"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
	corev1 "k8s.io/api/core/v1"
)

// DomainService 域名管理服务接口
type DomainService interface {
	// List 列出所有域名
	List(ctx context.Context, query *model.CommonPageQuery) (*model.PaginatedResult[model.Domain], error)

	// Get 根据域名名称获取域名详情
	Get(ctx context.Context, domainName string) (*model.Domain, error)

	// Add 添加新域名
	Add(ctx context.Context, domain *model.Domain) (*model.Domain, error)

	// Update 更新域名
	Update(ctx context.Context, domain *model.Domain) (*model.Domain, error)

	// Delete 删除域名
	Delete(ctx context.Context, domainName string) error
}

// DomainServiceImpl 域名服务实现
type DomainServiceImpl struct {
	kubernetesClient      *kubernetes.KubernetesClientService
	modelConverter        *kubernetes.KubernetesModelConverter
	routeService          RouteService
	wasmPluginInstanceSvc WasmPluginInstanceService
}

// NewDomainService 创建域名服务实例
func NewDomainService(
	kubernetesClient *kubernetes.KubernetesClientService,
	modelConverter *kubernetes.KubernetesModelConverter,
	routeService RouteService,
	wasmPluginInstanceSvc WasmPluginInstanceService,
) DomainService {
	return &DomainServiceImpl{
		kubernetesClient:      kubernetesClient,
		modelConverter:        modelConverter,
		routeService:          routeService,
		wasmPluginInstanceSvc: wasmPluginInstanceSvc,
	}
}

// List 列出所有域名
func (s *DomainServiceImpl) List(ctx context.Context, query *model.CommonPageQuery) (*model.PaginatedResult[model.Domain], error) {
	configMaps, err := s.kubernetesClient.ListConfigMaps(ctx, nil)
	if err != nil {
		return nil, errors.NewBusinessError("Error occurs when listing ConfigMap: " + err.Error())
	}

	// 过滤域名ConfigMap
	domainConfigMaps := make([]corev1.ConfigMap, 0)
	for i := range configMaps {
		cm := &configMaps[i]
		if cm.Name != "" && strings.HasPrefix(cm.Name, constant.DomainKeyPrefix) {
			domainConfigMaps = append(domainConfigMaps, *cm)
		}
	}

	// 转换为Domain列表
	domains := make([]model.Domain, 0, len(domainConfigMaps))
	for i := range domainConfigMaps {
		domain, err := s.modelConverter.ConfigMapToDomain(&domainConfigMaps[i])
		if err != nil {
			continue
		}
		domains = append(domains, *domain)
	}

	// 应用分页
	total := len(domains)
	pageNum := 1
	pageSize := 10
	if query != nil {
		pageNum = query.PageNum
		pageSize = query.GetPageSize()
	}

	start := (pageNum - 1) * pageSize
	end := start + pageSize
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	pagedData := domains[start:end]
	return model.NewPaginatedResult(pagedData, total, pageNum, pageSize), nil
}

// Get 根据域名名称获取域名详情
func (s *DomainServiceImpl) Get(ctx context.Context, domainName string) (*model.Domain, error) {
	configMapName := s.modelConverter.DomainNameToConfigMapName(domainName)
	cm, err := s.kubernetesClient.GetConfigMap(ctx, configMapName)
	if err != nil {
		return nil, errors.NewBusinessError("Error occurs when reading the ConfigMap with name: " + configMapName)
	}

	if cm == nil {
		return nil, nil
	}

	return s.modelConverter.ConfigMapToDomain(cm)
}

// Add 添加新域名
func (s *DomainServiceImpl) Add(ctx context.Context, domain *model.Domain) (*model.Domain, error) {
	cm, err := s.modelConverter.DomainToConfigMap(domain)
	if err != nil {
		return nil, err
	}

	newCm, err := s.kubernetesClient.CreateConfigMap(ctx, cm)
	if err != nil {
		if strings.Contains(err.Error(), "409") || strings.Contains(err.Error(), "AlreadyExists") {
			return nil, errors.NewResourceConflictError("Domain", domain.Name)
		}
		return nil, errors.NewBusinessError("Error occurs when adding a new domain: " + err.Error())
	}

	// 同步路由域名配置
	s.syncRouteDomainConfigs(ctx, domain.Name)

	return s.modelConverter.ConfigMapToDomain(newCm)
}

// Update 更新域名
func (s *DomainServiceImpl) Update(ctx context.Context, domain *model.Domain) (*model.Domain, error) {
	cm, err := s.modelConverter.DomainToConfigMap(domain)
	if err != nil {
		return nil, err
	}

	newCm, err := s.kubernetesClient.UpdateConfigMap(ctx, cm)
	if err != nil {
		if strings.Contains(err.Error(), "409") || strings.Contains(err.Error(), "Conflict") {
			return nil, errors.NewResourceConflictError("Domain", domain.Name)
		}
		return nil, errors.NewBusinessError("Error occurs when updating the domain: " + err.Error())
	}

	// 同步路由域名配置
	s.syncRouteDomainConfigs(ctx, domain.Name)

	return s.modelConverter.ConfigMapToDomain(newCm)
}

// Delete 删除域名
func (s *DomainServiceImpl) Delete(ctx context.Context, domainName string) error {
	// 检查是否有路由使用此域名
	query := &model.RoutePageQuery{DomainName: domainName}
	routes, err := s.routeService.List(ctx, query)
	if err != nil {
		return errors.NewBusinessError("Error occurs when checking routes for domain: " + err.Error())
	}

	if len(routes.Data) > 0 {
		return errors.NewValidationError("The domain has routes. Please delete them first.")
	}

	configMapName := s.modelConverter.DomainNameToConfigMapName(domainName)
	err = s.kubernetesClient.DeleteConfigMap(ctx, configMapName)
	if err != nil {
		return errors.NewBusinessError("Error occurs when deleting the ConfigMap with name: " + configMapName)
	}

	// 删除关联的插件实例
	if s.wasmPluginInstanceSvc != nil {
		s.wasmPluginInstanceSvc.DeleteAll(ctx, WasmPluginInstanceScopeDomain, domainName)
	}

	return nil
}

// syncRouteDomainConfigs 同步路由域名配置
func (s *DomainServiceImpl) syncRouteDomainConfigs(ctx context.Context, domainName string) {
	if s.routeService == nil {
		return
	}

	var routes []model.Route

	// 处理默认域名
	if domainName == constant.DefaultDomain {
		allRoutes, err := s.routeService.List(ctx, nil)
		if err != nil {
			return
		}
		// 过滤没有域名的路由
		for _, r := range allRoutes.Data {
			if len(r.Domains) == 0 {
				routes = append(routes, r)
			}
		}
	} else {
		query := &model.RoutePageQuery{DomainName: domainName}
		result, err := s.routeService.List(ctx, query)
		if err != nil {
			return
		}
		routes = result.Data
	}

	// 更新路由
	for _, route := range routes {
		_, err := s.routeService.Update(ctx, &route)
		if err != nil {
			continue
		}
	}
}
