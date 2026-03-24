// Package service provides business services for the SDK
package service

import (
	"context"

	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
)

// McpServerService MCP服务器管理服务接口
// 注意：此接口将在里程碑10中完整实现
type McpServerService interface {
	// List 列出MCP服务器
	List(ctx context.Context, query *model.McpServerPageQuery) (*model.PaginatedResult[model.McpServer], error)

	// Get 获取MCP服务器详情
	Get(ctx context.Context, name string) (*model.McpServer, error)

	// Add 添加MCP服务器
	Add(ctx context.Context, server *model.McpServer) (*model.McpServer, error)

	// Update 更新MCP服务器
	Update(ctx context.Context, server *model.McpServer) (*model.McpServer, error)

	// Delete 删除MCP服务器
	Delete(ctx context.Context, name string) error

	// ListConsumers 列出MCP服务器的消费者
	ListConsumers(ctx context.Context, query *model.McpServerPageQuery) (*model.PaginatedResult[model.McpServerConsumerDetail], error)

	// AddConsumer 添加消费者到MCP服务器
	AddConsumer(ctx context.Context, serverName string, consumer *model.McpServerConsumer) error

	// RemoveConsumer 从MCP服务器移除消费者
	RemoveConsumer(ctx context.Context, serverName, consumerName string) error
}
