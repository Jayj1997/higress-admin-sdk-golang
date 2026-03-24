// Package detail provides MCP server detail strategies
package detail

import (
	"context"

	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
)

// McpServerDetailStrategy MCP服务器详情策略接口
type McpServerDetailStrategy interface {
	// Support 是否支持该MCP服务器
	Support(mcpServer *model.McpServer) bool

	// Query 查询MCP服务器详情
	Query(ctx context.Context, name string) (*model.McpServer, error)
}
