// Package model provides data models for the SDK
package model

import (
	"strings"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/errors"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model/route"
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

// Validate 验证AI路由配置
func (r *AiRoute) Validate() error {
	if r.Name == "" {
		return errors.NewValidationError("name cannot be blank")
	}
	if len(r.Upstreams) == 0 {
		return errors.NewValidationError("upstreams cannot be empty")
	}
	if r.PathPredicate != nil {
		if err := r.PathPredicate.Validate(); err != nil {
			return err
		}
		// AI路由只支持前缀匹配
		if r.PathPredicate.MatchType != route.MatchTypePrefix {
			return errors.NewValidationError("pathPredicate must be of type prefix")
		}
	}
	for i := range r.HeaderPredicates {
		if err := r.HeaderPredicates[i].Validate(); err != nil {
			return err
		}
		if strings.EqualFold(r.HeaderPredicates[i].Key, ModelRoutingHeader) {
			return errors.NewValidationError("headerPredicates cannot contain the model routing header")
		}
	}
	for i := range r.UrlParamPredicates {
		if err := r.UrlParamPredicates[i].Validate(); err != nil {
			return err
		}
	}
	for i := range r.Upstreams {
		if err := r.Upstreams[i].Validate(); err != nil {
			return err
		}
	}
	// 验证权重总和
	weightSum := 0
	for _, upstream := range r.Upstreams {
		weightSum += upstream.Weight
	}
	if weightSum != 100 {
		return errors.NewValidationError("The sum of upstream weights must be 100")
	}
	if r.AuthConfig != nil {
		if err := r.AuthConfig.Validate(); err != nil {
			return err
		}
	}
	if r.FallbackConfig != nil {
		if err := r.FallbackConfig.Validate(); err != nil {
			return err
		}
	}
	for i := range r.ModelPredicates {
		if err := r.ModelPredicates[i].Validate(); err != nil {
			return err
		}
	}
	return nil
}

// AiUpstream AI上游配置
type AiUpstream struct {
	// Provider 提供商名称
	Provider string `json:"provider,omitempty"`

	// Weight 权重
	Weight int `json:"weight,omitempty"`

	// ModelMapping 模型映射
	ModelMapping map[string]string `json:"modelMapping,omitempty"`
}

// Validate 验证上游配置
func (u *AiUpstream) Validate() error {
	if u.Provider == "" {
		return errors.NewValidationError("provider cannot be null or empty")
	}
	return nil
}

// AiModelPredicate AI模型谓词
type AiModelPredicate struct {
	// MatchType 匹配类型
	MatchType string `json:"matchType,omitempty"`

	// MatchValue 匹配值
	MatchValue string `json:"matchValue,omitempty"`
}

// Validate 验证模型谓词
func (p *AiModelPredicate) Validate() error {
	if p.MatchType == "" {
		return errors.NewValidationError("matchType cannot be blank")
	}
	if p.MatchValue == "" {
		return errors.NewValidationError("matchValue cannot be blank")
	}
	return nil
}

// AiRouteFallbackConfig AI路由降级配置
type AiRouteFallbackConfig struct {
	// Enabled 是否启用
	Enabled bool `json:"enabled,omitempty"`

	// FallbackStrategy 降级策略（RANDOM/SEQUENCE）
	FallbackStrategy string `json:"fallbackStrategy,omitempty"`

	// Upstreams 降级上游列表
	Upstreams []AiUpstream `json:"upstreams,omitempty"`

	// ResponseCodes 触发降级的响应码
	ResponseCodes []string `json:"responseCodes,omitempty"`
}

// Validate 验证降级配置
func (c *AiRouteFallbackConfig) Validate() error {
	if c.Enabled && len(c.Upstreams) == 0 {
		return errors.NewValidationError("upstreams cannot be empty when fallback is enabled")
	}
	for i := range c.Upstreams {
		if err := c.Upstreams[i].Validate(); err != nil {
			return err
		}
	}
	return nil
}

// AI路由相关常量
const (
	// ModelRoutingHeader 模型路由头
	ModelRoutingHeader = "x-higress-model-routing"

	// FallbackFromHeader 降级来源头
	FallbackFromHeader = "x-higress-fallback-from"

	// AiRouteFallbackStrategyRandom 随机降级策略
	AiRouteFallbackStrategyRandom = "RANDOM"

	// AiRouteFallbackStrategySequence 顺序降级策略
	AiRouteFallbackStrategySequence = "SEQUENCE"
)
