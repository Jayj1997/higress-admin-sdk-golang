// Package mock provides mock implementations for testing purposes
package mock

import (
	"context"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/service"
)

// MockWasmPluginService WasmPluginService的Mock实现
type MockWasmPluginService struct{}

// NewMockWasmPluginService 创建Mock WasmPluginService实例
func NewMockWasmPluginService() service.WasmPluginService {
	return &MockWasmPluginService{}
}

// List 列出WASM插件
func (s *MockWasmPluginService) List(ctx context.Context, query *model.WasmPluginPageQuery) (*model.PaginatedResult[model.WasmPlugin], error) {
	return model.NewPaginatedResult([]model.WasmPlugin{}, 0, 1, 10), nil
}

// Get 获取WASM插件详情
func (s *MockWasmPluginService) Get(ctx context.Context, name, language string) (*model.WasmPlugin, error) {
	return nil, nil
}

// GetConfig 获取插件配置Schema
func (s *MockWasmPluginService) GetConfig(ctx context.Context, name, language string) (*model.WasmPluginConfig, error) {
	return nil, nil
}

// GetReadme 获取插件README文档
func (s *MockWasmPluginService) GetReadme(ctx context.Context, name, language string) (string, error) {
	return "", nil
}

// UpdateBuiltIn 更新内置插件配置
func (s *MockWasmPluginService) UpdateBuiltIn(ctx context.Context, plugin *model.WasmPlugin) (*model.WasmPlugin, error) {
	return plugin, nil
}

// AddCustom 添加自定义插件
func (s *MockWasmPluginService) AddCustom(ctx context.Context, plugin *model.WasmPlugin) (*model.WasmPlugin, error) {
	return plugin, nil
}

// UpdateCustom 更新自定义插件
func (s *MockWasmPluginService) UpdateCustom(ctx context.Context, plugin *model.WasmPlugin) (*model.WasmPlugin, error) {
	return plugin, nil
}

// DeleteCustom 删除自定义插件
func (s *MockWasmPluginService) DeleteCustom(ctx context.Context, name string) error {
	return nil
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

// AddOrUpdate 添加或更新消费者
func (s *MockConsumerService) AddOrUpdate(ctx context.Context, consumer *model.Consumer) (*model.Consumer, error) {
	return consumer, nil
}

// Delete 删除消费者
func (s *MockConsumerService) Delete(ctx context.Context, name string) error {
	return nil
}

// ListAllowLists 列出所有允许列表
func (s *MockConsumerService) ListAllowLists(ctx context.Context) ([]model.AllowList, error) {
	return []model.AllowList{}, nil
}

// GetAllowList 获取指定目标的允许列表
func (s *MockConsumerService) GetAllowList(ctx context.Context, targets map[model.WasmPluginInstanceScope]string) (*model.AllowList, error) {
	return nil, nil
}

// UpdateAllowList 更新允许列表
func (s *MockConsumerService) UpdateAllowList(ctx context.Context, operation model.AllowListOperation, allowList *model.AllowList) error {
	return nil
}
