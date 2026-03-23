// Package service provides business services for the SDK
package service

import (
	"context"
	"strings"

	"github.com/Jayj1997/higress-admin-sdk-golang/internal/kubernetes"
	"github.com/Jayj1997/higress-admin-sdk-golang/internal/kubernetes/crd/mcp"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/errors"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
)

// ProxyServerService 代理服务器管理服务接口
type ProxyServerService interface {
	// List 列出代理服务器
	List(ctx context.Context, query *model.CommonPageQuery) (*model.PaginatedResult[model.ProxyServer], error)

	// Get 根据名称获取代理服务器
	Get(ctx context.Context, name string) (*model.ProxyServer, error)

	// Add 添加代理服务器
	Add(ctx context.Context, server *model.ProxyServer) (*model.ProxyServer, error)

	// Update 更新代理服务器
	Update(ctx context.Context, server *model.ProxyServer) (*model.ProxyServer, error)

	// Delete 删除代理服务器
	Delete(ctx context.Context, name string) error
}

// ProxyServerServiceImpl 代理服务器服务实现
type ProxyServerServiceImpl struct {
	kubernetesClient *kubernetes.KubernetesClientService
	modelConverter   *kubernetes.KubernetesModelConverter
}

// NewProxyServerService 创建代理服务器服务实例
func NewProxyServerService(
	kubernetesClient *kubernetes.KubernetesClientService,
	modelConverter *kubernetes.KubernetesModelConverter,
) ProxyServerService {
	return &ProxyServerServiceImpl{
		kubernetesClient: kubernetesClient,
		modelConverter:   modelConverter,
	}
}

// List 列出代理服务器
func (s *ProxyServerServiceImpl) List(ctx context.Context, query *model.CommonPageQuery) (*model.PaginatedResult[model.ProxyServer], error) {
	mcpBridge, err := s.kubernetesClient.GetMcpBridge(ctx, DefaultMcpBridgeName)
	if err != nil {
		return nil, errors.NewBusinessError("Error occurs when getting McpBridge: " + err.Error())
	}

	servers := make([]model.ProxyServer, 0)
	if mcpBridge != nil && mcpBridge.Spec != nil && len(mcpBridge.Spec.Proxies) > 0 {
		resourceVersion := ""
		if mcpBridge.Metadata != nil {
			resourceVersion = mcpBridge.Metadata.ResourceVersion
		}

		for _, proxy := range mcpBridge.Spec.Proxies {
			server := s.convertProxyToProxyServer(proxy)
			server.Version = resourceVersion
			servers = append(servers, *server)
		}
	}

	// 应用分页
	total := len(servers)
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

	pagedData := servers[start:end]
	return model.NewPaginatedResult(pagedData, total, pageNum, pageSize), nil
}

// Get 根据名称获取代理服务器
func (s *ProxyServerServiceImpl) Get(ctx context.Context, name string) (*model.ProxyServer, error) {
	mcpBridge, err := s.kubernetesClient.GetMcpBridge(ctx, DefaultMcpBridgeName)
	if err != nil {
		return nil, errors.NewBusinessError("Error occurs when getting McpBridge: " + err.Error())
	}

	if mcpBridge == nil || mcpBridge.Spec == nil {
		return nil, nil
	}

	for _, proxy := range mcpBridge.Spec.Proxies {
		// 使用Host作为名称标识
		if proxy.Host == name {
			server := s.convertProxyToProxyServer(proxy)
			if mcpBridge.Metadata != nil {
				server.Version = mcpBridge.Metadata.ResourceVersion
			}
			return server, nil
		}
	}

	return nil, nil
}

// Add 添加代理服务器
func (s *ProxyServerServiceImpl) Add(ctx context.Context, server *model.ProxyServer) (*model.ProxyServer, error) {
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
	for _, proxy := range mcpBridge.Spec.Proxies {
		if proxy.Host == server.Host {
			return nil, errors.NewResourceConflictError("ProxyServer", server.Host)
		}
	}

	// 添加新的proxy
	proxy := s.convertProxyServerToProxy(server)
	mcpBridge.Spec.Proxies = append(mcpBridge.Spec.Proxies, proxy)

	// 保存
	if mcpBridge.Metadata.ResourceVersion == "" {
		_, err = s.kubernetesClient.CreateMcpBridge(ctx, mcpBridge)
	} else {
		_, err = s.kubernetesClient.UpdateMcpBridge(ctx, mcpBridge)
	}

	if err != nil {
		if strings.Contains(err.Error(), "409") || strings.Contains(err.Error(), "Conflict") {
			return nil, errors.NewResourceConflictError("ProxyServer", server.Host)
		}
		return nil, errors.NewBusinessError("Error occurs when adding ProxyServer: " + err.Error())
	}

	return server, nil
}

// Update 更新代理服务器
func (s *ProxyServerServiceImpl) Update(ctx context.Context, server *model.ProxyServer) (*model.ProxyServer, error) {
	mcpBridge, err := s.kubernetesClient.GetMcpBridge(ctx, DefaultMcpBridgeName)
	if err != nil {
		return nil, errors.NewBusinessError("Error occurs when getting McpBridge: " + err.Error())
	}

	if mcpBridge == nil || mcpBridge.Spec == nil {
		return nil, errors.NewNotFoundError("ProxyServer", server.Name)
	}

	// 查找并更新
	found := false
	for i, proxy := range mcpBridge.Spec.Proxies {
		if proxy.Host == server.Name {
			mcpBridge.Spec.Proxies[i] = s.convertProxyServerToProxy(server)
			found = true
			break
		}
	}

	if !found {
		return nil, errors.NewNotFoundError("ProxyServer", server.Name)
	}

	// 保存
	_, err = s.kubernetesClient.UpdateMcpBridge(ctx, mcpBridge)
	if err != nil {
		if strings.Contains(err.Error(), "409") || strings.Contains(err.Error(), "Conflict") {
			return nil, errors.NewResourceConflictError("ProxyServer", server.Name)
		}
		return nil, errors.NewBusinessError("Error occurs when updating ProxyServer: " + err.Error())
	}

	return server, nil
}

// Delete 删除代理服务器
func (s *ProxyServerServiceImpl) Delete(ctx context.Context, name string) error {
	mcpBridge, err := s.kubernetesClient.GetMcpBridge(ctx, DefaultMcpBridgeName)
	if err != nil {
		return errors.NewBusinessError("Error occurs when getting McpBridge: " + err.Error())
	}

	if mcpBridge == nil || mcpBridge.Spec == nil {
		return errors.NewNotFoundError("ProxyServer", name)
	}

	// 查找并删除
	found := false
	newProxies := make([]*mcp.V1ProxyConfig, 0)
	for _, proxy := range mcpBridge.Spec.Proxies {
		if proxy.Host == name {
			found = true
			continue
		}
		newProxies = append(newProxies, proxy)
	}

	if !found {
		return errors.NewNotFoundError("ProxyServer", name)
	}

	mcpBridge.Spec.Proxies = newProxies

	// 保存
	_, err = s.kubernetesClient.UpdateMcpBridge(ctx, mcpBridge)
	if err != nil {
		return errors.NewBusinessError("Error occurs when deleting ProxyServer: " + err.Error())
	}

	return nil
}

// convertProxyToProxyServer 将V1ProxyConfig转换为ProxyServer
func (s *ProxyServerServiceImpl) convertProxyToProxyServer(proxy *mcp.V1ProxyConfig) *model.ProxyServer {
	server := &model.ProxyServer{
		Host: proxy.Host,
	}

	if proxy.Port > 0 {
		server.Port = int(proxy.Port)
	}

	if proxy.Type != "" {
		server.Protocol = proxy.Type
	}

	// 使用Host作为名称
	server.Name = proxy.Host

	return server
}

// convertProxyServerToProxy 将ProxyServer转换为V1ProxyConfig
func (s *ProxyServerServiceImpl) convertProxyServerToProxy(server *model.ProxyServer) *mcp.V1ProxyConfig {
	proxy := &mcp.V1ProxyConfig{
		Host: server.Host,
	}

	if server.Port > 0 {
		proxy.Port = uint32(server.Port)
	}

	if server.Protocol != "" {
		proxy.Type = server.Protocol
	}

	proxy.Enabled = true

	return proxy
}
