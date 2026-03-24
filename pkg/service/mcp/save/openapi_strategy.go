// Package save provides MCP server save strategies
package save

import (
	"context"

	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/service/mcp"
)

// OpenApiSaveStrategy OpenAPI类型保存策略
type OpenApiSaveStrategy struct {
	*AbstractMcpServerSaveStrategy
}

// NewOpenApiSaveStrategy 创建OpenAPI保存策略
func NewOpenApiSaveStrategy(
	configMapHelper *mcp.McpServerConfigMapHelper,
	routeService RouteServiceInterface,
	consumerService ConsumerServiceInterface,
) *OpenApiSaveStrategy {
	return &OpenApiSaveStrategy{
		AbstractMcpServerSaveStrategy: NewAbstractMcpServerSaveStrategy(configMapHelper, routeService, consumerService),
	}
}

// Support 是否支持该MCP服务器
func (s *OpenApiSaveStrategy) Support(mcpServer *model.McpServer) bool {
	return mcpServer != nil && mcpServer.Type == model.McpServerTypeOpenApi
}

// Save 保存MCP服务器
func (s *OpenApiSaveStrategy) Save(ctx context.Context, mcpServer *model.McpServer) (*model.McpServer, error) {
	return s.AbstractMcpServerSaveStrategy.Save(ctx, mcpServer)
}

// SaveWithAuthorization 带授权保存MCP服务器
func (s *OpenApiSaveStrategy) SaveWithAuthorization(ctx context.Context, mcpServer *model.McpServer) (*model.McpServer, error) {
	return s.AbstractMcpServerSaveStrategy.SaveWithAuthorization(ctx, mcpServer)
}

// DoSaveMcpServerConfig 保存MCP服务器配置
func (s *OpenApiSaveStrategy) DoSaveMcpServerConfig(ctx context.Context, mcpServer *model.McpServer) error {
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
	if mcpServer.RawConfigurations != "" {
		// 解析YAML配置
		server.Config = map[string]interface{}{
			"rawConfigurations": mcpServer.RawConfigurations,
		}
	}

	return s.GetConfigMapHelper().AddServer(ctx, server)
}
