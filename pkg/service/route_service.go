// Package service provides business services for the SDK
package service

import (
	"context"
	"strings"

	"github.com/Jayj1997/higress-admin-sdk-golang/internal/kubernetes"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/errors"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
	networkingv1 "k8s.io/api/networking/v1"
)

// RouteService 路由管理服务接口
type RouteService interface {
	// List 列出路由
	List(ctx context.Context, query *model.RoutePageQuery) (*model.PaginatedResult[model.Route], error)

	// Get 根据路由名称获取路由详情
	Get(ctx context.Context, routeName string) (*model.Route, error)

	// Add 添加新路由
	Add(ctx context.Context, route *model.Route) (*model.Route, error)

	// Update 更新路由
	Update(ctx context.Context, route *model.Route) (*model.Route, error)

	// Delete 删除路由
	Delete(ctx context.Context, routeName string) error
}

// RouteServiceImpl 路由服务实现
type RouteServiceImpl struct {
	kubernetesClient      *kubernetes.KubernetesClientService
	modelConverter        *kubernetes.KubernetesModelConverter
	wasmPluginInstanceSvc WasmPluginInstanceService
}

// NewRouteService 创建路由服务实例
func NewRouteService(
	kubernetesClient *kubernetes.KubernetesClientService,
	modelConverter *kubernetes.KubernetesModelConverter,
	wasmPluginInstanceSvc WasmPluginInstanceService,
) RouteService {
	return &RouteServiceImpl{
		kubernetesClient:      kubernetesClient,
		modelConverter:        modelConverter,
		wasmPluginInstanceSvc: wasmPluginInstanceSvc,
	}
}

// List 列出路由
func (s *RouteServiceImpl) List(ctx context.Context, query *model.RoutePageQuery) (*model.PaginatedResult[model.Route], error) {
	ingresses, err := s.listIngresses(ctx, query)
	if err != nil {
		return nil, err
	}

	if len(ingresses) == 0 {
		return model.NewPaginatedResult([]model.Route{}, 0, 1, 10), nil
	}

	// 转换为Route列表
	routes := make([]model.Route, 0, len(ingresses))
	for i := range ingresses {
		route, err := s.modelConverter.IngressToRoute(&ingresses[i])
		if err != nil {
			continue
		}
		routes = append(routes, *route)
	}

	// 应用分页
	total := len(routes)
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

	pagedData := routes[start:end]
	return model.NewPaginatedResult(pagedData, total, pageNum, pageSize), nil
}

// Get 根据路由名称获取路由详情
func (s *RouteServiceImpl) Get(ctx context.Context, routeName string) (*model.Route, error) {
	ingress, err := s.kubernetesClient.GetIngress(ctx, routeName)
	if err != nil {
		return nil, errors.NewBusinessError("Error occurs when reading the Ingress with name: " + routeName)
	}

	if ingress == nil {
		return nil, nil
	}

	return s.modelConverter.IngressToRoute(ingress)
}

// Add 添加新路由
func (s *RouteServiceImpl) Add(ctx context.Context, route *model.Route) (*model.Route, error) {
	ingress, err := s.modelConverter.RouteToIngress(route)
	if err != nil {
		return nil, err
	}

	newIngress, err := s.kubernetesClient.CreateIngress(ctx, ingress)
	if err != nil {
		if strings.Contains(err.Error(), "409") || strings.Contains(err.Error(), "AlreadyExists") {
			return nil, errors.NewResourceConflictError("Route", route.Name)
		}
		return nil, errors.NewBusinessError("Error occurs when adding a new route: " + err.Error())
	}

	return s.modelConverter.IngressToRoute(newIngress)
}

// Update 更新路由
func (s *RouteServiceImpl) Update(ctx context.Context, route *model.Route) (*model.Route, error) {
	ingress, err := s.modelConverter.RouteToIngress(route)
	if err != nil {
		return nil, err
	}

	newIngress, err := s.kubernetesClient.UpdateIngress(ctx, ingress)
	if err != nil {
		if strings.Contains(err.Error(), "409") || strings.Contains(err.Error(), "Conflict") {
			return nil, errors.NewResourceConflictError("Route", route.Name)
		}
		return nil, errors.NewBusinessError("Error occurs when updating the route: " + err.Error())
	}

	return s.modelConverter.IngressToRoute(newIngress)
}

// Delete 删除路由
func (s *RouteServiceImpl) Delete(ctx context.Context, routeName string) error {
	err := s.kubernetesClient.DeleteIngress(ctx, routeName)
	if err != nil {
		return errors.NewBusinessError("Error occurs when deleting the Ingress with name: " + routeName)
	}

	// 删除关联的插件实例
	if s.wasmPluginInstanceSvc != nil {
		s.wasmPluginInstanceSvc.DeleteAll(ctx, model.WasmPluginInstanceScopeRoute, routeName)
	}

	return nil
}

// listIngresses 列出Ingress资源
func (s *RouteServiceImpl) listIngresses(ctx context.Context, query *model.RoutePageQuery) ([]networkingv1.Ingress, error) {
	// 根据查询条件获取Ingress
	if query != nil && query.DomainName != "" {
		// 按域名过滤
		ingresses, err := s.kubernetesClient.ListIngresses(ctx)
		if err != nil {
			return nil, errors.NewBusinessError("Error occurs when listing Ingresses: " + err.Error())
		}

		// 过滤包含指定域名的Ingress
		filtered := make([]networkingv1.Ingress, 0)
		for i := range ingresses {
			ingress := &ingresses[i]
			if s.ingressHasDomain(ingress, query.DomainName) {
				filtered = append(filtered, *ingress)
			}
		}
		return filtered, nil
	}

	// 获取所有Ingress
	ingresses, err := s.kubernetesClient.ListIngresses(ctx)
	if err != nil {
		return nil, errors.NewBusinessError("Error occurs when listing Ingresses: " + err.Error())
	}

	return ingresses, nil
}

// ingressHasDomain 检查Ingress是否包含指定域名
func (s *RouteServiceImpl) ingressHasDomain(ingress *networkingv1.Ingress, domainName string) bool {
	if ingress == nil || len(ingress.Spec.Rules) == 0 {
		return false
	}

	for _, rule := range ingress.Spec.Rules {
		if rule.Host == domainName {
			return true
		}
	}

	return false
}
