// Package save provides MCP server save strategies
package save

import (
	"context"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/service/mcp"
)

// DatabaseSaveStrategy Database类型保存策略
type DatabaseSaveStrategy struct {
	*AbstractMcpServerSaveStrategy
	dsnConverter *mcp.McpServerDBConfigDsnConverter
	validator    *mcp.McpServerDBConfigValidator
}

// NewDatabaseSaveStrategy 创建Database保存策略
func NewDatabaseSaveStrategy(
	configMapHelper *mcp.McpServerConfigMapHelper,
	routeService RouteServiceInterface,
	consumerService ConsumerServiceInterface,
) *DatabaseSaveStrategy {
	return &DatabaseSaveStrategy{
		AbstractMcpServerSaveStrategy: NewAbstractMcpServerSaveStrategy(configMapHelper, routeService, consumerService),
		dsnConverter:                  mcp.NewMcpServerDBConfigDsnConverter(),
		validator:                     mcp.NewMcpServerDBConfigValidator(),
	}
}

// Support 是否支持该MCP服务器
func (s *DatabaseSaveStrategy) Support(mcpServer *model.McpServer) bool {
	return mcpServer != nil && mcpServer.Type == model.McpServerTypeDatabase
}

// Save 保存MCP服务器
func (s *DatabaseSaveStrategy) Save(ctx context.Context, mcpServer *model.McpServer) (*model.McpServer, error) {
	// 验证数据库配置
	if err := s.validator.ValidateMcpServer(mcpServer); err != nil {
		return nil, err
	}

	return s.AbstractMcpServerSaveStrategy.Save(ctx, mcpServer)
}

// SaveWithAuthorization 带授权保存MCP服务器
func (s *DatabaseSaveStrategy) SaveWithAuthorization(ctx context.Context, mcpServer *model.McpServer) (*model.McpServer, error) {
	// 验证数据库配置
	if err := s.validator.ValidateMcpServer(mcpServer); err != nil {
		return nil, err
	}

	return s.AbstractMcpServerSaveStrategy.SaveWithAuthorization(ctx, mcpServer)
}

// DoSaveMcpServerConfig 保存MCP服务器配置
func (s *DatabaseSaveStrategy) DoSaveMcpServerConfig(ctx context.Context, mcpServer *model.McpServer) error {
	// 初始化配置
	if err := s.GetConfigMapHelper().InitMcpServerConfig(ctx); err != nil {
		return err
	}

	// 添加匹配规则
	matchList := s.GetConfigMapHelper().GenerateMatchList(mcpServer)
	if err := s.GetConfigMapHelper().AddMatchList(ctx, matchList); err != nil {
		return err
	}

	// 转换DSN
	dsn, err := s.dsnConverter.ConvertToDsn(mcpServer.DBConfig, mcpServer.DBType)
	if err != nil {
		return err
	}

	// 添加服务器配置
	server := &model.McpServerConfigMapServer{
		Name: mcpServer.Name,
		Config: map[string]interface{}{
			"dsn":    dsn,
			"dbType": string(mcpServer.DBType),
		},
	}

	return s.GetConfigMapHelper().AddServer(ctx, server)
}
