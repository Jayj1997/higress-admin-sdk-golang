// Package save provides MCP server save strategies
package save

import (
	"context"

	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
)

// McpServerSaveStrategy MCP服务器保存策略接口
type McpServerSaveStrategy interface {
	// Support 是否支持该MCP服务器
	Support(mcpServer *model.McpServer) bool

	// Save 保存MCP服务器
	Save(ctx context.Context, mcpServer *model.McpServer) (*model.McpServer, error)

	// SaveWithAuthorization 带授权保存MCP服务器
	SaveWithAuthorization(ctx context.Context, mcpServer *model.McpServer) (*model.McpServer, error)
}
