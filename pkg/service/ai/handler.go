// Package ai provides AI-related services for Higress Admin SDK.
package ai

import (
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model/route"
)

// LlmProviderHandler LLM提供商处理器接口
type LlmProviderHandler interface {
	// GetType 获取处理器类型
	GetType() string

	// CreateProvider 创建提供商实例
	CreateProvider() *model.LlmProvider

	// LoadConfig 从配置加载提供商信息
	LoadConfig(provider *model.LlmProvider, configurations map[string]interface{}) bool

	// SaveConfig 保存提供商配置
	SaveConfig(provider *model.LlmProvider, configurations map[string]interface{})

	// NormalizeConfigs 规范化配置
	NormalizeConfigs(configurations map[string]interface{})

	// BuildServiceSource 构建服务来源
	BuildServiceSource(providerName string, providerConfig map[string]interface{}) (*model.ServiceSource, error)

	// BuildUpstreamService 构建上游服务
	BuildUpstreamService(providerName string, providerConfig map[string]interface{}) (*route.UpstreamService, error)

	// GetServiceSourceName 获取服务来源名称
	GetServiceSourceName(providerName string) string

	// GetExtraServiceSources 获取额外的服务来源
	GetExtraServiceSources(providerName string, providerConfig map[string]interface{}, forDelete bool) []model.ServiceSource

	// NeedSyncRouteAfterUpdate 更新后是否需要同步路由
	NeedSyncRouteAfterUpdate() bool

	// GetProviderEndpoints 获取提供商端点
	GetProviderEndpoints(providerConfig map[string]interface{}) []model.LlmProviderEndpoint
}
