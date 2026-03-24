// Package save provides MCP server save strategies
package save

import (
	"context"

	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/service/mcp"
)

// DirectRoutingSaveStrategy DirectRoute类型保存策略
type DirectRoutingSaveStrategy struct {
	*AbstractMcpServerSaveStrategy
}

// NewDirectRoutingSaveStrategy 创建DirectRoute保存策略
func NewDirectRoutingSaveStrategy(
	configMapHelper *mcp.McpServerConfigMapHelper,
	routeService RouteServiceInterface,
	consumerService ConsumerServiceInterface,
) *DirectRoutingSaveStrategy {
	return &DirectRoutingSaveStrategy{
		AbstractMcpServerSaveStrategy: NewAbstractMcpServerSaveStrategy(configMapHelper, routeService, consumerService),
	}
}

// Support 是否支持该MCP服务器
func (s *DirectRoutingSaveStrategy) Support(mcpServer *model.McpServer) bool {
	return mcpServer != nil && mcpServer.Type == model.McpServerTypeDirectRoute
}

// Save 保存MCP服务器
func (s *DirectRoutingSaveStrategy) Save(ctx context.Context, mcpServer *model.McpServer) (*model.McpServer, error) {
	return s.AbstractMcpServerSaveStrategy.Save(ctx, mcpServer)
}

// SaveWithAuthorization 带授权保存MCP服务器
func (s *DirectRoutingSaveStrategy) SaveWithAuthorization(ctx context.Context, mcpServer *model.McpServer) (*model.McpServer, error) {
	return s.AbstractMcpServerSaveStrategy.SaveWithAuthorization(ctx, mcpServer)
}

// DoSaveMcpServerConfig 保存MCP服务器配置
func (s *DirectRoutingSaveStrategy) DoSaveMcpServerConfig(ctx context.Context, mcpServer *model.McpServer) error {
	// 初始化配置
	if err := s.GetConfigMapHelper().InitMcpServerConfig(ctx); err != nil {
		return err
	}

	// 添加匹配规则
	matchList := s.GetConfigMapHelper().GenerateMatchList(mcpServer)
	if err := s.GetConfigMapHelper().AddMatchList(ctx, matchList); err != nil {
		return err
	}

	// 添加服务器配置
	server := &model.McpServerConfigMapServer{
		Name: mcpServer.Name,
	}

	if mcpServer.DirectRouteConfig != nil {
		server.Config = map[string]interface{}{
			"upstreamProtocol":   mcpServer.DirectRouteConfig.UpstreamProtocol,
			"downstreamProtocol": mcpServer.DirectRouteConfig.DownstreamProtocol,
		}
	}

	return s.GetConfigMapHelper().AddServer(ctx, server)
}
