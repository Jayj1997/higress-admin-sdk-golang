// Package model provides data models for the SDK
package model

import (
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model/route"
)

// AiRoute AI路由配置
type AiRoute struct {
	// Name 路由名称
	Name string `json:"name,omitempty"`

	// Version 路由版本，更新时需要
	Version string `json:"version,omitempty"`

	// Domains 路由适用的域名列表，为空表示所有域名
	Domains []string `json:"domains,omitempty"`

	// PathPredicate 路径谓词
	PathPredicate *route.RoutePredicate `json:"pathPredicate,omitempty"`

	// HeaderPredicates 头部谓词列表
	HeaderPredicates []route.KeyedRoutePredicate `json:"headerPredicates,omitempty"`

	// UrlParamPredicates URL参数谓词列表
	UrlParamPredicates []route.KeyedRoutePredicate `json:"urlParamPredicates,omitempty"`

	// Upstreams 路由上游列表
	Upstreams []AiUpstream `json:"upstreams,omitempty"`

	// ModelPredicates 模型谓词列表
	ModelPredicates []AiModelPredicate `json:"modelPredicates,omitempty"`

	// AuthConfig 路由认证配置
	AuthConfig *RouteAuthConfig `json:"authConfig,omitempty"`

	// FallbackConfig 路由降级配置
	FallbackConfig *AiRouteFallbackConfig `json:"fallbackConfig,omitempty"`

	// Cors CORS配置
	Cors *route.CorsConfig `json:"cors,omitempty"`

	// HeaderControl 头部控制配置
	HeaderControl *route.HeaderControlConfig `json:"headerControl,omitempty"`

	// ProxyNextUpstream 代理下一上游配置
	ProxyNextUpstream *route.ProxyNextUpstreamConfig `json:"proxyNextUpstream,omitempty"`

	// CustomConfigs 自定义配置（如自定义注解）
	CustomConfigs map[string]string `json:"customConfigs,omitempty"`

	// CustomLabels 自定义标签
	CustomLabels map[string]string `json:"customLabels,omitempty"`
}

// AiUpstream AI上游配置
type AiUpstream struct {
	// ProviderName 提供商名称
	ProviderName string `json:"providerName,omitempty"`

	// Weight 权重
	Weight int `json:"weight,omitempty"`
}

// AiModelPredicate AI模型谓词
type AiModelPredicate struct {
	// Model 模型名称
	Model string `json:"model,omitempty"`

	// MatchType 匹配类型
	MatchType string `json:"matchType,omitempty"`
}

// AiRouteFallbackConfig AI路由降级配置
type AiRouteFallbackConfig struct {
	// Enabled 是否启用
	Enabled bool `json:"enabled,omitempty"`

	// FallbackUpstream 降级上游
	FallbackUpstream *AiUpstream `json:"fallbackUpstream,omitempty"`
}
