// Package mcp provides MCP server related services
package mcp

import (
	"strings"

	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/constant"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model/route"
)

const (
	// McpServerRoutePrefix MCP服务器路由前缀
	McpServerRoutePrefix = "mcp-server-"
	// InternalResourceNameSuffix 内部资源名称后缀
	InternalResourceNameSuffix = "-internal"
)

// McpServerHelper MCP服务器辅助工具
type McpServerHelper struct{}

// NewMcpServerHelper 创建MCP服务器辅助工具
func NewMcpServerHelper() *McpServerHelper {
	return &McpServerHelper{}
}

// McpServerName2RouteName MCP服务器名称转路由名称
func (h *McpServerHelper) McpServerName2RouteName(mcpServerName string) string {
	if strings.HasPrefix(mcpServerName, McpServerRoutePrefix) &&
		strings.HasSuffix(mcpServerName, InternalResourceNameSuffix) {
		return mcpServerName
	}
	return McpServerRoutePrefix + mcpServerName + InternalResourceNameSuffix
}

// RouteName2McpServerName 路由名称转MCP服务器名称
func (h *McpServerHelper) RouteName2McpServerName(routeName string) string {
	if routeName == "" {
		return ""
	}
	if !strings.HasPrefix(routeName, McpServerRoutePrefix) {
		return routeName
	}
	result := strings.TrimPrefix(routeName, McpServerRoutePrefix)
	result = strings.TrimSuffix(result, InternalResourceNameSuffix)
	return result
}

// RouteToMcpServer 路由转MCP服务器
func (h *McpServerHelper) RouteToMcpServer(r *model.Route) *model.McpServer {
	if r == nil {
		return nil
	}

	result := &model.McpServer{
		Name: h.RouteName2McpServerName(r.Name),
	}

	// 转换Services类型 []*route.UpstreamService -> []route.UpstreamService
	if len(r.Services) > 0 {
		services := make([]route.UpstreamService, len(r.Services))
		for i, s := range r.Services {
			if s != nil {
				services[i] = *s
			}
		}
		result.Services = services
	}

	result.Domains = r.Domains

	// 解析自定义配置
	if r.CustomConfigs != nil {
		if desc, ok := r.CustomConfigs[constant.AnnotationResourceDescriptionKey]; ok {
			result.Description = desc
		}
	}

	// 解析自定义标签
	if r.CustomLabels != nil {
		if mcpServerTypeStr, ok := r.CustomLabels[constant.LabelResourceMcpServerTypeKey]; ok {
			result.Type = model.ParseMcpServerTypeEnum(mcpServerTypeStr)
		}
	}

	// 解析认证配置
	if r.AuthConfig != nil {
		consumerAuthInfo := &model.ConsumerAuthInfo{
			AllowedConsumers: r.AuthConfig.AllowedConsumers,
		}
		if r.AuthConfig.Enabled != nil {
			consumerAuthInfo.Enable = *r.AuthConfig.Enabled
		}
		if len(r.AuthConfig.AllowedCredentialTypes) > 0 {
			consumerAuthInfo.Type = r.AuthConfig.AllowedCredentialTypes[0]
		}
		result.ConsumerAuthInfo = consumerAuthInfo
	}

	return result
}

// GenerateMcpServerPath 生成MCP服务器路径
func (h *McpServerHelper) GenerateMcpServerPath(mcpServerName string) string {
	return constant.McpServerPathPre + mcpServerName
}

// IsMcpServerRoute 判断是否为MCP服务器路由
func (h *McpServerHelper) IsMcpServerRoute(customLabels map[string]string) bool {
	if len(customLabels) == 0 {
		return false
	}
	bizType, ok := customLabels[constant.LabelResourceBizTypeKey]
	if !ok {
		return false
	}
	return strings.EqualFold(bizType, constant.LabelMcpServerBizTypeValue)
}
