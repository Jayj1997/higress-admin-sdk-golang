// Package ai provides AI-related services for Higress Admin SDK.
package ai

import (
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model"
)

// AzureLlmProviderHandler Azure OpenAI提供商处理器
type AzureLlmProviderHandler struct {
	DefaultLlmProviderHandler
}

// NewAzureLlmProviderHandler 创建Azure处理器
func NewAzureLlmProviderHandler() *AzureLlmProviderHandler {
	// Azure使用自定义域名，默认配置为空
	return &AzureLlmProviderHandler{
		DefaultLlmProviderHandler: *NewDefaultLlmProviderHandler(
			model.LlmProviderTypeAzure,
			"", // Azure使用自定义域名
			443,
			"https",
		),
	}
}

// GetProviderEndpoints 获取提供商端点 - Azure需要从配置中获取
func (h *AzureLlmProviderHandler) GetProviderEndpoints(providerConfig map[string]interface{}) []model.LlmProviderEndpoint {
	// Azure从配置中获取域名
	domain := getString(providerConfig, "resourceName")
	if domain == "" {
		return nil
	}

	// 构建Azure OpenAI端点
	return []model.LlmProviderEndpoint{
		{
			Protocol:    "https",
			Address:     domain + ".openai.azure.com",
			Port:        443,
			ContextPath: "/",
		},
	}
}
