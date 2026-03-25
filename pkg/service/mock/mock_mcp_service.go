// Package mock provides mock implementations for testing purposes
package mock

import (
	"context"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/service"
)

// MockMcpServerService McpServerService的Mock实现
type MockMcpServerService struct{}

// NewMockMcpServerService 创建Mock McpServerService实例
func NewMockMcpServerService() service.McpServerService {
	return &MockMcpServerService{}
}

// List 列出MCP服务器
func (s *MockMcpServerService) List(ctx context.Context, query *model.McpServerPageQuery) (*model.PaginatedResult[model.McpServer], error) {
	return &model.PaginatedResult[model.McpServer]{}, nil
}

// Get 获取MCP服务器详情
func (s *MockMcpServerService) Get(ctx context.Context, name string) (*model.McpServer, error) {
	return nil, nil
}

// Add 添加MCP服务器
func (s *MockMcpServerService) Add(ctx context.Context, server *model.McpServer) (*model.McpServer, error) {
	return server, nil
}

// Update 更新MCP服务器
func (s *MockMcpServerService) Update(ctx context.Context, server *model.McpServer) (*model.McpServer, error) {
	return server, nil
}

// Delete 删除MCP服务器
func (s *MockMcpServerService) Delete(ctx context.Context, name string) error {
	return nil
}

// ListConsumers 列出MCP服务器的消费者
func (s *MockMcpServerService) ListConsumers(ctx context.Context, query *model.McpServerPageQuery) (*model.PaginatedResult[model.McpServerConsumerDetail], error) {
	return &model.PaginatedResult[model.McpServerConsumerDetail]{}, nil
}

// AddConsumer 添加消费者到MCP服务器
func (s *MockMcpServerService) AddConsumer(ctx context.Context, serverName string, consumer *model.McpServerConsumer) error {
	return nil
}

// RemoveConsumer 从MCP服务器移除消费者
func (s *MockMcpServerService) RemoveConsumer(ctx context.Context, serverName, consumerName string) error {
	return nil
}
