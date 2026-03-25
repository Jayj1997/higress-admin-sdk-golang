// Package service provides business services for the SDK
package service

import (
	"context"
	"strings"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/internal/kubernetes"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/internal/kubernetes/crd/mcp"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/errors"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model"
)

const (
	// DefaultMcpBridgeName is the default name for McpBridge resource
	DefaultMcpBridgeName = "default"
)

// ServiceSourceService 服务来源管理服务接口
type ServiceSourceService interface {
	// List 列出服务来源
	List(ctx context.Context, query *model.CommonPageQuery) (*model.PaginatedResult[model.ServiceSource], error)

	// Get 根据名称获取服务来源
	Get(ctx context.Context, name string) (*model.ServiceSource, error)

	// Add 添加服务来源
	Add(ctx context.Context, source *model.ServiceSource) (*model.ServiceSource, error)

	// Update 更新服务来源
	Update(ctx context.Context, source *model.ServiceSource) (*model.ServiceSource, error)

	// Delete 删除服务来源
	Delete(ctx context.Context, name string) error
}

// ServiceSourceServiceImpl 服务来源服务实现
type ServiceSourceServiceImpl struct {
	kubernetesClient *kubernetes.KubernetesClientService
	modelConverter   *kubernetes.KubernetesModelConverter
}

// NewServiceSourceService 创建服务来源服务实例
func NewServiceSourceService(
	kubernetesClient *kubernetes.KubernetesClientService,
	modelConverter *kubernetes.KubernetesModelConverter,
) ServiceSourceService {
	return &ServiceSourceServiceImpl{
		kubernetesClient: kubernetesClient,
		modelConverter:   modelConverter,
	}
}

// List 列出服务来源
func (s *ServiceSourceServiceImpl) List(ctx context.Context, query *model.CommonPageQuery) (*model.PaginatedResult[model.ServiceSource], error) {
	mcpBridge, err := s.kubernetesClient.GetMcpBridge(ctx, DefaultMcpBridgeName)
	if err != nil {
		return nil, errors.NewBusinessError("Error occurs when getting McpBridge: " + err.Error())
	}

	sources := make([]model.ServiceSource, 0)
	if mcpBridge != nil && mcpBridge.Spec != nil && len(mcpBridge.Spec.Registries) > 0 {
		resourceVersion := ""
		if mcpBridge.Metadata != nil {
			resourceVersion = mcpBridge.Metadata.ResourceVersion
		}

		for _, registry := range mcpBridge.Spec.Registries {
			source := s.convertRegistryToServiceSource(registry)
			source.Version = resourceVersion
			sources = append(sources, *source)
		}
	}

	// 应用分页
	total := len(sources)
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

	pagedData := sources[start:end]
	return model.NewPaginatedResult(pagedData, total, pageNum, pageSize), nil
}

// Get 根据名称获取服务来源
func (s *ServiceSourceServiceImpl) Get(ctx context.Context, name string) (*model.ServiceSource, error) {
	mcpBridge, err := s.kubernetesClient.GetMcpBridge(ctx, DefaultMcpBridgeName)
	if err != nil {
		return nil, errors.NewBusinessError("Error occurs when getting McpBridge: " + err.Error())
	}

	if mcpBridge == nil || mcpBridge.Spec == nil {
		return nil, nil
	}

	for _, registry := range mcpBridge.Spec.Registries {
		if registry.Name == name {
			source := s.convertRegistryToServiceSource(registry)
			if mcpBridge.Metadata != nil {
				source.Version = mcpBridge.Metadata.ResourceVersion
			}
			return source, nil
		}
	}

	return nil, nil
}

// Add 添加服务来源
func (s *ServiceSourceServiceImpl) Add(ctx context.Context, source *model.ServiceSource) (*model.ServiceSource, error) {
	mcpBridge, err := s.kubernetesClient.GetMcpBridge(ctx, DefaultMcpBridgeName)
	if err != nil {
		return nil, errors.NewBusinessError("Error occurs when getting McpBridge: " + err.Error())
	}

	if mcpBridge == nil {
		// 创建新的McpBridge
		mcpBridge = mcp.NewV1McpBridge()
		mcpBridge.Metadata.Name = DefaultMcpBridgeName
	}

	// 检查是否已存在
	for _, registry := range mcpBridge.Spec.Registries {
		if registry.Name == source.Name {
			return nil, errors.NewResourceConflictError("ServiceSource", source.Name)
		}
	}

	// 添加新的registry
	registry := s.convertServiceSourceToRegistry(source)
	mcpBridge.Spec.Registries = append(mcpBridge.Spec.Registries, registry)

	// 保存
	if mcpBridge.Metadata.ResourceVersion == "" {
		_, err = s.kubernetesClient.CreateMcpBridge(ctx, mcpBridge)
	} else {
		_, err = s.kubernetesClient.UpdateMcpBridge(ctx, mcpBridge)
	}

	if err != nil {
		if strings.Contains(err.Error(), "409") || strings.Contains(err.Error(), "Conflict") {
			return nil, errors.NewResourceConflictError("ServiceSource", source.Name)
		}
		return nil, errors.NewBusinessError("Error occurs when adding ServiceSource: " + err.Error())
	}

	return source, nil
}

// Update 更新服务来源
func (s *ServiceSourceServiceImpl) Update(ctx context.Context, source *model.ServiceSource) (*model.ServiceSource, error) {
	mcpBridge, err := s.kubernetesClient.GetMcpBridge(ctx, DefaultMcpBridgeName)
	if err != nil {
		return nil, errors.NewBusinessError("Error occurs when getting McpBridge: " + err.Error())
	}

	if mcpBridge == nil || mcpBridge.Spec == nil {
		return nil, errors.NewNotFoundError("ServiceSource", source.Name)
	}

	// 查找并更新
	found := false
	for i, registry := range mcpBridge.Spec.Registries {
		if registry.Name == source.Name {
			mcpBridge.Spec.Registries[i] = s.convertServiceSourceToRegistry(source)
			found = true
			break
		}
	}

	if !found {
		return nil, errors.NewNotFoundError("ServiceSource", source.Name)
	}

	// 保存
	_, err = s.kubernetesClient.UpdateMcpBridge(ctx, mcpBridge)
	if err != nil {
		if strings.Contains(err.Error(), "409") || strings.Contains(err.Error(), "Conflict") {
			return nil, errors.NewResourceConflictError("ServiceSource", source.Name)
		}
		return nil, errors.NewBusinessError("Error occurs when updating ServiceSource: " + err.Error())
	}

	return source, nil
}

// Delete 删除服务来源
func (s *ServiceSourceServiceImpl) Delete(ctx context.Context, name string) error {
	mcpBridge, err := s.kubernetesClient.GetMcpBridge(ctx, DefaultMcpBridgeName)
	if err != nil {
		return errors.NewBusinessError("Error occurs when getting McpBridge: " + err.Error())
	}

	if mcpBridge == nil || mcpBridge.Spec == nil {
		return errors.NewNotFoundError("ServiceSource", name)
	}

	// 查找并删除
	found := false
	newRegistries := make([]*mcp.V1RegistryConfig, 0)
	for _, registry := range mcpBridge.Spec.Registries {
		if registry.Name == name {
			found = true
			continue
		}
		newRegistries = append(newRegistries, registry)
	}

	if !found {
		return errors.NewNotFoundError("ServiceSource", name)
	}

	mcpBridge.Spec.Registries = newRegistries

	// 保存
	_, err = s.kubernetesClient.UpdateMcpBridge(ctx, mcpBridge)
	if err != nil {
		return errors.NewBusinessError("Error occurs when deleting ServiceSource: " + err.Error())
	}

	return nil
}

// convertRegistryToServiceSource 将V1RegistryConfig转换为ServiceSource
func (s *ServiceSourceServiceImpl) convertRegistryToServiceSource(registry *mcp.V1RegistryConfig) *model.ServiceSource {
	source := &model.ServiceSource{
		Name: registry.Name,
		Type: registry.Type,
	}

	if registry.Port > 0 {
		port := int(registry.Port)
		source.Port = &port
	}

	// 根据类型设置不同字段
	switch registry.Type {
	case "nacos":
		source.Domain = registry.Domain
		source.Namespace = registry.NacosNamespace
		if len(registry.NacosGroups) > 0 {
			source.Group = registry.NacosGroups[0]
		}
	case "consul":
		source.Domain = registry.Domain
		source.Namespace = registry.ConsulNamespace
	case "eureka":
		source.Domain = registry.EurekaClientServiceUrl
	case "dns":
		source.Domain = registry.DnsDomain
		if registry.DnsPort > 0 {
			port := int(registry.DnsPort)
			source.Port = &port
		}
	case "static":
		if len(registry.StaticAddresses) > 0 {
			source.Domain = registry.StaticAddresses[0]
		}
	default:
		source.Domain = registry.Domain
	}

	return source
}

// convertServiceSourceToRegistry 将ServiceSource转换为V1RegistryConfig
func (s *ServiceSourceServiceImpl) convertServiceSourceToRegistry(source *model.ServiceSource) *mcp.V1RegistryConfig {
	registry := &mcp.V1RegistryConfig{
		Name: source.Name,
		Type: source.Type,
	}

	if source.Port != nil {
		registry.Port = uint32(*source.Port)
	}

	// 根据类型设置不同字段
	switch source.Type {
	case "nacos":
		registry.Domain = source.Domain
		registry.NacosNamespace = source.Namespace
		if source.Group != "" {
			registry.NacosGroups = []string{source.Group}
		}
	case "consul":
		registry.Domain = source.Domain
		registry.ConsulNamespace = source.Namespace
	case "eureka":
		registry.EurekaClientServiceUrl = source.Domain
	case "dns":
		registry.DnsDomain = source.Domain
		if source.Port != nil {
			registry.DnsPort = uint32(*source.Port)
		}
	case "static":
		if source.Domain != "" {
			registry.StaticAddresses = []string{source.Domain}
		}
	default:
		registry.Domain = source.Domain
	}

	return registry
}
