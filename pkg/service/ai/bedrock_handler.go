// Package ai provides AI-related services for Higress Admin SDK.
package ai

import (
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
)

// BedrockLlmProviderHandler AWS Bedrock提供商处理器
type BedrockLlmProviderHandler struct {
	DefaultLlmProviderHandler
}

// NewBedrockLlmProviderHandler 创建AWS Bedrock处理器
func NewBedrockLlmProviderHandler() *BedrockLlmProviderHandler {
	return &BedrockLlmProviderHandler{
		DefaultLlmProviderHandler: *NewDefaultLlmProviderHandler(
			model.LlmProviderTypeBedrock,
			"", // Bedrock使用区域特定端点
			443,
			"https",
		),
	}
}

// getProviderEndpoints 获取提供商端点 - Bedrock需要从配置中获取区域
func (h *BedrockLlmProviderHandler) getProviderEndpoints(providerConfig map[string]interface{}) []model.LlmProviderEndpoint {
	// 从配置中获取AWS区域
	region := getString(providerConfig, "awsRegion")
	if region == "" {
		region = "us-east-1" // 默认区域
	}

	// 构建Bedrock端点
	return []model.LlmProviderEndpoint{
		{
			Protocol:    "https",
			Address:     "bedrock-runtime." + region + ".amazonaws.com",
			Port:        443,
			ContextPath: "/",
		},
	}
}

// NeedSyncRouteAfterUpdate 更新后是否需要同步路由
func (h *BedrockLlmProviderHandler) NeedSyncRouteAfterUpdate() bool {
	return true
}
