// Package service provides business services for the SDK
package service

import (
	"context"
	"sort"
	"strings"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/internal/kubernetes"
	k8smodel "github.com/Jayj1997/higress-admin-sdk-golang/v2/internal/kubernetes/model"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/errors"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model"
)

// ServiceService 服务管理服务接口
type ServiceService interface {
	// List 列出服务
	List(ctx context.Context, query *model.CommonPageQuery) (*model.PaginatedResult[model.Service], error)
}

// ServiceServiceImpl 服务管理服务实现
type ServiceServiceImpl struct {
	kubernetesClient *kubernetes.KubernetesClientService
}

// NewServiceService 创建服务管理服务实例
func NewServiceService(kubernetesClient *kubernetes.KubernetesClientService) ServiceService {
	return &ServiceServiceImpl{
		kubernetesClient: kubernetesClient,
	}
}

// List 列出服务
func (s *ServiceServiceImpl) List(ctx context.Context, query *model.CommonPageQuery) (*model.PaginatedResult[model.Service], error) {
	// 从Controller获取服务列表
	registryzServices, err := s.kubernetesClient.GatewayServiceList(ctx)
	if err != nil {
		return nil, errors.NewBusinessError("Error occurs when listing services: " + err.Error())
	}

	if len(registryzServices) == 0 {
		return model.NewPaginatedResult([]model.Service{}, 0, 1, 10), nil
	}

	// 获取服务端点
	serviceEndpoints, err := s.kubernetesClient.GatewayServiceEndpoint(ctx)
	if err != nil {
		// 端点获取失败不影响服务列表返回
		serviceEndpoints = make(map[string]map[string]*k8smodel.IstioEndpointShard)
	}

	// 获取MCP Bridge DNS域名映射
	mcpBridgeDomain := s.getMcpBridgeDnsDomain(ctx)

	// 转换服务列表
	services := make([]model.Service, 0)
	for _, registryzService := range registryzServices {
		if registryzService.Attributes == nil {
			continue
		}

		namespace := registryzService.Attributes.Namespace
		// 跳过受保护的命名空间
		if s.kubernetesClient.IsNamespaceProtected(namespace) {
			continue
		}

		name := registryzService.Hostname

		// 获取端点
		endpoints := s.getServiceEndpoints(serviceEndpoints, namespace, name)
		if len(endpoints) == 0 {
			// 尝试从MCP DNS获取端点
			endpoints = s.completeMcpDnsEndpoints(registryzService, mcpBridgeDomain)
		}

		// 处理端口
		if len(registryzService.Ports) == 0 {
			service := model.Service{
				Name:      name,
				Namespace: namespace,
				Endpoints: endpoints,
			}
			services = append(services, service)
		} else {
			// 为每个端口创建服务实例
			for _, port := range registryzService.Ports {
				service := model.Service{
					Name:      name,
					Namespace: namespace,
					Port:      int(port.Port),
					Endpoints: endpoints,
				}
				services = append(services, service)
			}
		}
	}

	// 按名称排序
	sort.Slice(services, func(i, j int) bool {
		if services[i].Namespace != services[j].Namespace {
			return services[i].Namespace < services[j].Namespace
		}
		if services[i].Name != services[j].Name {
			return services[i].Name < services[j].Name
		}
		return services[i].Port < services[j].Port
	})

	// 应用分页
	total := len(services)
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

	pagedData := services[start:end]
	return model.NewPaginatedResult(pagedData, total, pageNum, pageSize), nil
}

// getServiceEndpoints 获取服务端点
func (s *ServiceServiceImpl) getServiceEndpoints(
	serviceEndpoints map[string]map[string]*k8smodel.IstioEndpointShard,
	namespace, name string,
) []string {
	endpoints := make([]string, 0)

	namespaceEndpoints, ok := serviceEndpoints[namespace]
	if !ok {
		return endpoints
	}

	endpointShard, ok := namespaceEndpoints[name]
	if !ok || endpointShard == nil {
		return endpoints
	}

	// 直接从endpointShard获取端点
	for _, endpoint := range endpointShard.Endpoints {
		if endpoint == nil {
			continue
		}
		addr := endpoint.Address
		if addr != "" {
			endpoints = append(endpoints, addr)
		}
	}

	return endpoints
}

// getMcpBridgeDnsDomain 获取MCP Bridge DNS域名映射
func (s *ServiceServiceImpl) getMcpBridgeDnsDomain(ctx context.Context) map[string]string {
	domainMap := make(map[string]string)

	// 获取McpBridge
	mcpBridge, err := s.kubernetesClient.GetMcpBridge(ctx, "default")
	if err != nil || mcpBridge == nil {
		return domainMap
	}

	if mcpBridge.Spec == nil || len(mcpBridge.Spec.Registries) == 0 {
		return domainMap
	}

	// 提取域名映射
	for _, registry := range mcpBridge.Spec.Registries {
		if registry.Name != "" && registry.Domain != "" {
			domainMap[registry.Name] = registry.Domain
		}
	}

	return domainMap
}

// completeMcpDnsEndpoints 补充MCP DNS端点
func (s *ServiceServiceImpl) completeMcpDnsEndpoints(
	registryzService *k8smodel.RegistryzService,
	mcpBridgeDomain map[string]string,
) []string {
	if registryzService == nil || registryzService.Attributes == nil {
		return nil
	}

	namespace := registryzService.Attributes.Namespace
	name := registryzService.Hostname

	// 检查是否是MCP服务
	if !strings.HasPrefix(namespace, "mcp") {
		return nil
	}

	// 从域名映射中获取
	for registryName, domain := range mcpBridgeDomain {
		if strings.Contains(name, registryName) {
			// 构造DNS端点
			endpoint := name + "." + domain
			return []string{endpoint}
		}
	}

	return nil
}
