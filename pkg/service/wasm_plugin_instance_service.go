// Package service provides business services for the SDK
package service

import (
	"context"

	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
)

// WasmPluginInstanceScope 定义插件实例的作用域类型
type WasmPluginInstanceScope string

const (
	// WasmPluginInstanceScopeGlobal 全局作用域
	WasmPluginInstanceScopeGlobal WasmPluginInstanceScope = "GLOBAL"
	// WasmPluginInstanceScopeDomain 域名作用域
	WasmPluginInstanceScopeDomain WasmPluginInstanceScope = "DOMAIN"
	// WasmPluginInstanceScopeRoute 路由作用域
	WasmPluginInstanceScopeRoute WasmPluginInstanceScope = "ROUTE"
	// WasmPluginInstanceScopeService 服务作用域
	WasmPluginInstanceScopeService WasmPluginInstanceScope = "SERVICE"
)

// WasmPluginInstanceService WASM插件实例服务接口
// 注意：此接口将在里程碑7中完整实现
type WasmPluginInstanceService interface {
	// List 列出插件实例
	List(ctx context.Context, scope WasmPluginInstanceScope, target string) ([]model.WasmPluginInstance, error)

	// Get 获取插件实例
	Get(ctx context.Context, id string) (*model.WasmPluginInstance, error)

	// Add 添加插件实例
	Add(ctx context.Context, instance *model.WasmPluginInstance) (*model.WasmPluginInstance, error)

	// Update 更新插件实例
	Update(ctx context.Context, instance *model.WasmPluginInstance) (*model.WasmPluginInstance, error)

	// Delete 删除插件实例
	Delete(ctx context.Context, id string) error

	// DeleteAll 删除指定作用域和目标的所有插件实例
	DeleteAll(ctx context.Context, scope WasmPluginInstanceScope, target string) error
}

// MockWasmPluginInstanceService 用于测试的Mock实现
type MockWasmPluginInstanceService struct{}

// NewMockWasmPluginInstanceService 创建Mock服务实例
func NewMockWasmPluginInstanceService() WasmPluginInstanceService {
	return &MockWasmPluginInstanceService{}
}

func (s *MockWasmPluginInstanceService) List(ctx context.Context, scope WasmPluginInstanceScope, target string) ([]model.WasmPluginInstance, error) {
	return []model.WasmPluginInstance{}, nil
}

func (s *MockWasmPluginInstanceService) Get(ctx context.Context, id string) (*model.WasmPluginInstance, error) {
	return nil, nil
}

func (s *MockWasmPluginInstanceService) Add(ctx context.Context, instance *model.WasmPluginInstance) (*model.WasmPluginInstance, error) {
	return instance, nil
}

func (s *MockWasmPluginInstanceService) Update(ctx context.Context, instance *model.WasmPluginInstance) (*model.WasmPluginInstance, error) {
	return instance, nil
}

func (s *MockWasmPluginInstanceService) Delete(ctx context.Context, id string) error {
	return nil
}

func (s *MockWasmPluginInstanceService) DeleteAll(ctx context.Context, scope WasmPluginInstanceScope, target string) error {
	return nil
}
