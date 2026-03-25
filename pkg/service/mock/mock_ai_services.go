// Package mock provides mock implementations for testing purposes
package mock

import (
	"context"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/service"
)

// MockAiRouteService AiRouteService的Mock实现
type MockAiRouteService struct{}

// NewMockAiRouteService 创建Mock AiRouteService实例
func NewMockAiRouteService() service.AiRouteService {
	return &MockAiRouteService{}
}

// List 列出所有AI路由
func (s *MockAiRouteService) List(ctx context.Context) ([]model.AiRoute, error) {
	return []model.AiRoute{}, nil
}

// Get 获取AI路由详情
func (s *MockAiRouteService) Get(ctx context.Context, name string) (*model.AiRoute, error) {
	return nil, nil
}

// Add 添加AI路由
func (s *MockAiRouteService) Add(ctx context.Context, route *model.AiRoute) (*model.AiRoute, error) {
	return route, nil
}

// Update 更新AI路由
func (s *MockAiRouteService) Update(ctx context.Context, route *model.AiRoute) (*model.AiRoute, error) {
	return route, nil
}

// Delete 删除AI路由
func (s *MockAiRouteService) Delete(ctx context.Context, name string) error {
	return nil
}

// MockLlmProviderService LlmProviderService的Mock实现
type MockLlmProviderService struct{}

// NewMockLlmProviderService 创建Mock LlmProviderService实例
func NewMockLlmProviderService() service.LlmProviderService {
	return &MockLlmProviderService{}
}

// List 列出所有LLM提供商
func (s *MockLlmProviderService) List(ctx context.Context) ([]model.LlmProvider, error) {
	return []model.LlmProvider{}, nil
}

// Get 获取LLM提供商详情
func (s *MockLlmProviderService) Get(ctx context.Context, name string) (*model.LlmProvider, error) {
	return nil, nil
}

// Add 添加LLM提供商
func (s *MockLlmProviderService) Add(ctx context.Context, provider *model.LlmProvider) (*model.LlmProvider, error) {
	return provider, nil
}

// Update 更新LLM提供商
func (s *MockLlmProviderService) Update(ctx context.Context, provider *model.LlmProvider) (*model.LlmProvider, error) {
	return provider, nil
}

// Delete 删除LLM提供商
func (s *MockLlmProviderService) Delete(ctx context.Context, name string) error {
	return nil
}
