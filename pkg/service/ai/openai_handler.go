// Package ai provides AI-related services for Higress Admin SDK.
package ai

import (
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
)

const (
	// OpenAI默认配置
	openaiDefaultServiceProtocol = "https"
	openaiDefaultServiceDomain   = "api.openai.com"
	openaiDefaultServicePort     = 443
	openaiDefaultContextPath     = "/v1"
)

// OpenaiLlmProviderHandler OpenAI提供商处理器
type OpenaiLlmProviderHandler struct {
	DefaultLlmProviderHandler
}

// NewOpenaiLlmProviderHandler 创建OpenAI处理器
func NewOpenaiLlmProviderHandler() *OpenaiLlmProviderHandler {
	return &OpenaiLlmProviderHandler{
		DefaultLlmProviderHandler: *NewDefaultLlmProviderHandlerWithContextPath(
			model.LlmProviderTypeOpenai,
			openaiDefaultServiceDomain,
			openaiDefaultServicePort,
			openaiDefaultServiceProtocol,
			openaiDefaultContextPath,
		),
	}
}

func init() {
	RegisterHandler(NewOpenaiLlmProviderHandler())
}
