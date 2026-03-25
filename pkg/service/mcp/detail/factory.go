// Package detail provides MCP server detail strategies
package detail

import (
	"context"
	"fmt"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/service/mcp"
)

// AbstractMcpServerDetailStrategy 抽象详情策略
type AbstractMcpServerDetailStrategy struct {
	helper          *mcp.McpServerHelper
	configMapHelper *mcp.McpServerConfigMapHelper
}

// NewAbstractMcpServerDetailStrategy 创建抽象详情策略
func NewAbstractMcpServerDetailStrategy(configMapHelper *mcp.McpServerConfigMapHelper) *AbstractMcpServerDetailStrategy {
	return &AbstractMcpServerDetailStrategy{
		helper:          mcp.NewMcpServerHelper(),
		configMapHelper: configMapHelper,
	}
}

// Query 查询MCP服务器详情
func (s *AbstractMcpServerDetailStrategy) Query(ctx context.Context, name string) (*model.McpServer, error) {
	server, err := s.configMapHelper.GetServer(ctx, name)
	if err != nil {
		return nil, err
	}

	result := &model.McpServer{
		Name: server.Name,
	}

	if server.Config != nil {
		if rawConfig, ok := server.Config["rawConfigurations"].(string); ok {
			result.RawConfigurations = rawConfig
		}
		if dsn, ok := server.Config["dsn"].(string); ok {
			result.DBConfig = &model.McpServerDBConfig{}
			// 简化处理，实际应该解析DSN
			_ = dsn
		}
	}

	return result, nil
}

// OpenApiDetailStrategy OpenAPI详情策略
type OpenApiDetailStrategy struct {
	*AbstractMcpServerDetailStrategy
}

// NewOpenApiDetailStrategy 创建OpenAPI详情策略
func NewOpenApiDetailStrategy(configMapHelper *mcp.McpServerConfigMapHelper) *OpenApiDetailStrategy {
	return &OpenApiDetailStrategy{
		AbstractMcpServerDetailStrategy: NewAbstractMcpServerDetailStrategy(configMapHelper),
	}
}

// Support 是否支持
func (s *OpenApiDetailStrategy) Support(mcpServer *model.McpServer) bool {
	return mcpServer != nil && mcpServer.Type == model.McpServerTypeOpenApi
}

// DatabaseDetailStrategy Database详情策略
type DatabaseDetailStrategy struct {
	*AbstractMcpServerDetailStrategy
}

// NewDatabaseDetailStrategy 创建Database详情策略
func NewDatabaseDetailStrategy(configMapHelper *mcp.McpServerConfigMapHelper) *DatabaseDetailStrategy {
	return &DatabaseDetailStrategy{
		AbstractMcpServerDetailStrategy: NewAbstractMcpServerDetailStrategy(configMapHelper),
	}
}

// Support 是否支持
func (s *DatabaseDetailStrategy) Support(mcpServer *model.McpServer) bool {
	return mcpServer != nil && mcpServer.Type == model.McpServerTypeDatabase
}

// DirectRoutingDetailStrategy DirectRoute详情策略
type DirectRoutingDetailStrategy struct {
	*AbstractMcpServerDetailStrategy
}

// NewDirectRoutingDetailStrategy 创建DirectRoute详情策略
func NewDirectRoutingDetailStrategy(configMapHelper *mcp.McpServerConfigMapHelper) *DirectRoutingDetailStrategy {
	return &DirectRoutingDetailStrategy{
		AbstractMcpServerDetailStrategy: NewAbstractMcpServerDetailStrategy(configMapHelper),
	}
}

// Support 是否支持
func (s *DirectRoutingDetailStrategy) Support(mcpServer *model.McpServer) bool {
	return mcpServer != nil && mcpServer.Type == model.McpServerTypeDirectRoute
}

// McpServerDetailStrategyFactory 详情策略工厂
type McpServerDetailStrategyFactory struct {
	strategies []McpServerDetailStrategy
}

// NewMcpServerDetailStrategyFactory 创建详情策略工厂
func NewMcpServerDetailStrategyFactory(configMapHelper *mcp.McpServerConfigMapHelper) *McpServerDetailStrategyFactory {
	return &McpServerDetailStrategyFactory{
		strategies: []McpServerDetailStrategy{
			NewOpenApiDetailStrategy(configMapHelper),
			NewDatabaseDetailStrategy(configMapHelper),
			NewDirectRoutingDetailStrategy(configMapHelper),
		},
	}
}

// GetService 获取策略
func (f *McpServerDetailStrategyFactory) GetService(mcpServer *model.McpServer) McpServerDetailStrategy {
	for _, strategy := range f.strategies {
		if strategy.Support(mcpServer) {
			return strategy
		}
	}
	return nil
}

// GetServiceByType 根据类型获取策略
func (f *McpServerDetailStrategyFactory) GetServiceByType(serverType model.McpServerTypeEnum) (McpServerDetailStrategy, error) {
	for _, strategy := range f.strategies {
		tempServer := &model.McpServer{Type: serverType}
		if strategy.Support(tempServer) {
			return strategy, nil
		}
	}
	return nil, fmt.Errorf("no detail strategy found for type: %s", serverType)
}
