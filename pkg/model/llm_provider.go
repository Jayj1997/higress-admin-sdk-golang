// Package model provides data models for the SDK
package model

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

// TokenFailoverConfig Token故障转移配置
type TokenFailoverConfig struct {
	// Enabled 是否启用
	Enabled bool `json:"enabled,omitempty"`

	// FailureThreshold 故障阈值
	FailureThreshold int `json:"failureThreshold,omitempty"`

	// RecoveryTimeout 恢复超时时间（秒）
	RecoveryTimeout int `json:"recoveryTimeout,omitempty"`
}

// LlmProviderType LLM提供商类型常量
const (
	LlmProviderTypeOpenai   = "openai"
	LlmProviderTypeAzure    = "azure"
	LlmProviderTypeQwen     = "qwen"
	LlmProviderTypeMoonshot = "moonshot"
	LlmProviderTypeYi       = "yi"
	LlmProviderTypeDeepSeek = "deepseek"
	LlmProviderTypeBaichuan = "baichuan"
	LlmProviderTypeZhipuai  = "zhipuai"
	LlmProviderTypeOllama   = "ollama"
	LlmProviderTypeBedrock  = "bedrock"
	LlmProviderTypeVertex   = "vertex"
	LlmProviderTypeCustom   = "custom"
)

// LlmProviderProtocol LLM提供商协议常量
const (
	LlmProviderProtocolDefault  = "openai"
	LlmProviderProtocolOriginal = "original"
)
