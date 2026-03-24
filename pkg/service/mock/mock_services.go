// Package mock provides mock implementations for testing purposes
package mock

import (
	"context"

	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/service"
)

// MockWasmPluginService WasmPluginService的Mock实现
type MockWasmPluginService struct{}

// NewMockWasmPluginService 创建Mock WasmPluginService实例
func NewMockWasmPluginService() service.WasmPluginService {
	return &MockWasmPluginService{}
}

// List 列出WASM插件
func (s *MockWasmPluginService) List(ctx context.Context, name, version string, builtIn *bool) ([]model.WasmPlugin, error) {
	return []model.WasmPlugin{}, nil
}

// Get 获取WASM插件详情
func (s *MockWasmPluginService) Get(ctx context.Context, name string) (*model.WasmPlugin, error) {
	return nil, nil
}

// GetSchema 获取插件配置Schema
func (s *MockWasmPluginService) GetSchema(ctx context.Context, name string) (map[string]interface{}, error) {
	return nil, nil
}

// MockConsumerService ConsumerService的Mock实现
type MockConsumerService struct{}

// NewMockConsumerService 创建Mock ConsumerService实例
func NewMockConsumerService() service.ConsumerService {
	return &MockConsumerService{}
}

// List 列出所有消费者
func (s *MockConsumerService) List(ctx context.Context) ([]model.Consumer, error) {
	return []model.Consumer{}, nil
}

// Get 获取消费者详情
func (s *MockConsumerService) Get(ctx context.Context, name string) (*model.Consumer, error) {
	return nil, nil
}

// Add 添加消费者
func (s *MockConsumerService) Add(ctx context.Context, consumer *model.Consumer) (*model.Consumer, error) {
	return consumer, nil
}

// Update 更新消费者
func (s *MockConsumerService) Update(ctx context.Context, consumer *model.Consumer) (*model.Consumer, error) {
	return consumer, nil
}

// Delete 删除消费者
func (s *MockConsumerService) Delete(ctx context.Context, name string) error {
	return nil
}
