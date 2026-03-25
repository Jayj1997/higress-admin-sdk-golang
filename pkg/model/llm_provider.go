// Package model provides data models for the SDK
package model

import (
	"strings"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/errors"
)

// LlmProvider LLM服务提供商
type LlmProvider struct {
	// Name 提供商名称
	Name string `json:"name,omitempty"`

	// Type 提供商类型
	Type string `json:"type,omitempty"`

	// Protocol 提供商协议
	Protocol string `json:"protocol,omitempty"`

	// ProxyName 代理服务器名称
	ProxyName string `json:"proxyName,omitempty"`

	// Tokens 用于请求提供商的Token列表
	Tokens []string `json:"tokens,omitempty"`

	// TokenFailoverConfig Token故障转移配置
	TokenFailoverConfig *TokenFailoverConfig `json:"tokenFailoverConfig,omitempty"`

	// RawConfigs ai-proxy插件使用的原始配置键值对
	RawConfigs map[string]interface{} `json:"rawConfigs,omitempty"`
}

// Validate 验证提供商配置
func (p *LlmProvider) Validate(forUpdate bool) error {
	if p.Name == "" {
		return errors.NewValidationError("name cannot be blank")
	}
	if strings.Contains(p.Name, "/") {
		return errors.NewValidationError("slashes (/) are not allowed in name")
	}
	if p.Type == "" {
		return errors.NewValidationError("type cannot be blank")
	}
	if p.Protocol == "" {
		p.Protocol = LlmProviderProtocolOpenaiV1
	} else if !IsValidLlmProviderProtocol(p.Protocol) {
		return errors.NewValidationError("Unknown protocol: " + p.Protocol)
	}
	return nil
}

// TokenFailoverConfig Token故障转移配置
type TokenFailoverConfig struct {
	// Enabled 是否启用
	Enabled bool `json:"enabled,omitempty"`

	// FailureThreshold 故障阈值
	FailureThreshold int `json:"failureThreshold,omitempty"`

	// SuccessThreshold 成功阈值
	SuccessThreshold int `json:"successThreshold,omitempty"`

	// HealthCheckInterval 健康检查间隔（秒）
	HealthCheckInterval int `json:"healthCheckInterval,omitempty"`

	// HealthCheckTimeout 健康检查超时（秒）
	HealthCheckTimeout int `json:"healthCheckTimeout,omitempty"`

	// HealthCheckModel 健康检查使用的模型
	HealthCheckModel string `json:"healthCheckModel,omitempty"`
}

// LlmProviderEndpoint LLM提供商端点
type LlmProviderEndpoint struct {
	// Protocol 协议（http/https）
	Protocol string `json:"protocol,omitempty"`

	// Address 地址（域名或IP）
	Address string `json:"address,omitempty"`

	// Port 端口
	Port int `json:"port,omitempty"`

	// ContextPath 上下文路径
	ContextPath string `json:"contextPath,omitempty"`
}

// Validate 验证端点配置
func (e *LlmProviderEndpoint) Validate() error {
	if e.Address == "" {
		return errors.NewValidationError("endpoint address cannot be empty")
	}
	if e.Port <= 0 {
		return errors.NewValidationError("endpoint port must be positive")
	}
	if e.Protocol == "" {
		e.Protocol = "https"
	}
	if e.ContextPath == "" {
		e.ContextPath = "/"
	}
	return nil
}

// LlmProviderType LLM提供商类型常量
const (
	LlmProviderTypeQwen       = "qwen"
	LlmProviderTypeOpenai     = "openai"
	LlmProviderTypeMoonshot   = "moonshot"
	LlmProviderTypeAzure      = "azure"
	LlmProviderTypeAi360      = "ai360"
	LlmProviderTypeGithub     = "github"
	LlmProviderTypeGroq       = "groq"
	LlmProviderTypeBaichuan   = "baichuan"
	LlmProviderTypeYi         = "yi"
	LlmProviderTypeDeepSeek   = "deepseek"
	LlmProviderTypeZhipuai    = "zhipuai"
	LlmProviderTypeOllama     = "ollama"
	LlmProviderTypeClaude     = "claude"
	LlmProviderTypeBaidu      = "baidu"
	LlmProviderTypeHunyuan    = "hunyuan"
	LlmProviderTypeStepfun    = "stepfun"
	LlmProviderTypeMinimax    = "minimax"
	LlmProviderTypeCloudflare = "cloudflare"
	LlmProviderTypeSpark      = "spark"
	LlmProviderTypeGemini     = "gemini"
	LlmProviderTypeDeepl      = "deepl"
	LlmProviderTypeMistral    = "mistral"
	LlmProviderTypeCohere     = "cohere"
	LlmProviderTypeDoubao     = "doubao"
	LlmProviderTypeCoze       = "coze"
	LlmProviderTypeTogetherAi = "together-ai"
	LlmProviderTypeBedrock    = "bedrock"
	LlmProviderTypeVertex     = "vertex"
	LlmProviderTypeOpenrouter = "openrouter"
	LlmProviderTypeGrok       = "grok"
)

// LlmProviderProtocol LLM提供商协议常量
const (
	LlmProviderProtocolOpenaiV1 = "openai/v1"
	LlmProviderProtocolOriginal = "original"
)

// IsValidLlmProviderProtocol 检查协议是否有效
func IsValidLlmProviderProtocol(protocol string) bool {
	switch protocol {
	case LlmProviderProtocolOpenaiV1, LlmProviderProtocolOriginal:
		return true
	default:
		return false
	}
}

// LlmProviderProtocolFromValue 从字符串值获取协议
func LlmProviderProtocolFromValue(value string) string {
	if value == "" {
		return LlmProviderProtocolOpenaiV1
	}
	if IsValidLlmProviderProtocol(value) {
		return value
	}
	return ""
}
