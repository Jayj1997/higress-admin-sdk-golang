// Package ai provides AI-related services for Higress Admin SDK.
package ai

import (
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model"
)

const (
	// 通义千问默认配置
	qwenDefaultServiceProtocol = "https"
	qwenDefaultServiceDomain   = "dashscope.aliyuncs.com"
	qwenDefaultServicePort     = 443
	qwenDefaultContextPath     = "/compatible-mode/v1"
)

// QwenLlmProviderHandler 通义千问提供商处理器
type QwenLlmProviderHandler struct {
	DefaultLlmProviderHandler
}

// NewQwenLlmProviderHandler 创建通义千问处理器
func NewQwenLlmProviderHandler() *QwenLlmProviderHandler {
	return &QwenLlmProviderHandler{
		DefaultLlmProviderHandler: *NewDefaultLlmProviderHandlerWithContextPath(
			model.LlmProviderTypeQwen,
			qwenDefaultServiceDomain,
			qwenDefaultServicePort,
			qwenDefaultServiceProtocol,
			qwenDefaultContextPath,
		),
	}
}
