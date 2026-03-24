// Package service provides business services for the SDK
package service

import (
	"context"

	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
)

// AiRouteService AI路由管理服务接口
// 注意：此接口将在里程碑8中完整实现
type AiRouteService interface {
	// List 列出所有AI路由
	List(ctx context.Context) ([]model.AiRoute, error)

	// Get 获取AI路由详情
	Get(ctx context.Context, name string) (*model.AiRoute, error)

	// Add 添加AI路由
	Add(ctx context.Context, route *model.AiRoute) (*model.AiRoute, error)

	// Update 更新AI路由
	Update(ctx context.Context, route *model.AiRoute) (*model.AiRoute, error)

	// Delete 删除AI路由
	Delete(ctx context.Context, name string) error
}

// LlmProviderService LLM提供商管理服务接口
// 注意：此接口将在里程碑8中完整实现
type LlmProviderService interface {
	// List 列出所有LLM提供商
	List(ctx context.Context) ([]model.LlmProvider, error)

	// Get 获取LLM提供商详情
	Get(ctx context.Context, name string) (*model.LlmProvider, error)

	// Add 添加LLM提供商
	Add(ctx context.Context, provider *model.LlmProvider) (*model.LlmProvider, error)

	// Update 更新LLM提供商
	Update(ctx context.Context, provider *model.LlmProvider) (*model.LlmProvider, error)

	// Delete 删除LLM提供商
	Delete(ctx context.Context, name string) error
}
