// Package save provides MCP server save strategies
package save

import (
	"fmt"

	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/service/mcp"
)

// McpServerSaveStrategyFactory 保存策略工厂
type McpServerSaveStrategyFactory struct {
	strategies []McpServerSaveStrategy
}

// NewMcpServerSaveStrategyFactory 创建保存策略工厂
func NewMcpServerSaveStrategyFactory(
	configMapHelper *mcp.McpServerConfigMapHelper,
	routeService RouteServiceInterface,
	consumerService ConsumerServiceInterface,
) *McpServerSaveStrategyFactory {
	return &McpServerSaveStrategyFactory{
		strategies: []McpServerSaveStrategy{
			NewOpenApiSaveStrategy(configMapHelper, routeService, consumerService),
			NewDatabaseSaveStrategy(configMapHelper, routeService, consumerService),
			NewDirectRoutingSaveStrategy(configMapHelper, routeService, consumerService),
		},
	}
}

// GetService 根据MCP服务器类型获取对应策略
func (f *McpServerSaveStrategyFactory) GetService(mcpServer *model.McpServer) McpServerSaveStrategy {
	for _, strategy := range f.strategies {
		if strategy.Support(mcpServer) {
			return strategy
		}
	}
	return nil
}

// GetServiceByType 根据MCP服务器类型获取对应策略
func (f *McpServerSaveStrategyFactory) GetServiceByType(serverType model.McpServerTypeEnum) (McpServerSaveStrategy, error) {
	for _, strategy := range f.strategies {
		// 创建临时对象用于检查类型
		tempServer := &model.McpServer{Type: serverType}
		if strategy.Support(tempServer) {
			return strategy, nil
		}
	}
	return nil, fmt.Errorf("no save strategy found for type: %s", serverType)
}
