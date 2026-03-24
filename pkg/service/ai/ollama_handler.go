// Package ai provides AI-related services for Higress Admin SDK.
package ai

import (
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
)

// OllamaLlmProviderHandler Ollama提供商处理器
type OllamaLlmProviderHandler struct {
	DefaultLlmProviderHandler
}

// NewOllamaLlmProviderHandler 创建Ollama处理器
func NewOllamaLlmProviderHandler() *OllamaLlmProviderHandler {
	return &OllamaLlmProviderHandler{
		DefaultLlmProviderHandler: *NewDefaultLlmProviderHandler(
			model.LlmProviderTypeOllama,
			"localhost", // Ollama默认本地部署
			11434,
			"http",
		),
	}
}

// getProviderEndpoints 获取提供商端点 - Ollama支持自定义地址
func (h *OllamaLlmProviderHandler) getProviderEndpoints(providerConfig map[string]interface{}) []model.LlmProviderEndpoint {
	// 从配置中获取服务地址
	serverHost := getString(providerConfig, "serverHost")
	if serverHost == "" {
		serverHost = "localhost"
	}

	serverPort := getInt(providerConfig, "serverPort")
	if serverPort == 0 {
		serverPort = 11434
	}

	protocol := "http"
	if serverPort == 443 {
		protocol = "https"
	}

	return []model.LlmProviderEndpoint{
		{
			Protocol:    protocol,
			Address:     serverHost,
			Port:        serverPort,
			ContextPath: "/",
		},
	}
}

// NeedSyncRouteAfterUpdate 更新后是否需要同步路由
func (h *OllamaLlmProviderHandler) NeedSyncRouteAfterUpdate() bool {
	return true
}
