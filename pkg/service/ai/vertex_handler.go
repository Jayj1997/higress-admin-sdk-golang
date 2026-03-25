// Package ai provides AI-related services for Higress Admin SDK.
package ai

import (
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model"
)

// VertexLlmProviderHandler GCP Vertex AI提供商处理器
type VertexLlmProviderHandler struct {
	DefaultLlmProviderHandler
}

// NewVertexLlmProviderHandler 创建GCP Vertex处理器
func NewVertexLlmProviderHandler() *VertexLlmProviderHandler {
	return &VertexLlmProviderHandler{
		DefaultLlmProviderHandler: *NewDefaultLlmProviderHandler(
			model.LlmProviderTypeVertex,
			"", // Vertex使用项目特定端点
			443,
			"https",
		),
	}
}

// GetProviderEndpoints 获取提供商端点 - Vertex需要从配置中获取项目和区域
func (h *VertexLlmProviderHandler) GetProviderEndpoints(providerConfig map[string]interface{}) []model.LlmProviderEndpoint {
	// 从配置中获取GCP项目和区域
	project := getString(providerConfig, "gcpProject")
	location := getString(providerConfig, "gcpLocation")

	if project == "" || location == "" {
		return nil
	}

	// 构建Vertex AI端点
	return []model.LlmProviderEndpoint{
		{
			Protocol:    "https",
			Address:     location + "-aiplatform.googleapis.com",
			Port:        443,
			ContextPath: "/v1/projects/" + project + "/locations/" + location,
		},
	}
}

// NeedSyncRouteAfterUpdate 更新后是否需要同步路由
func (h *VertexLlmProviderHandler) NeedSyncRouteAfterUpdate() bool {
	return true
}
