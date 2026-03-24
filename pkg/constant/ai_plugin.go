// Package constant provides constants for Higress Admin SDK.
package constant

// AI代理插件相关常量
const (
	// BuiltInPluginAiProxy AI代理插件名称
	BuiltInPluginAiProxy = "ai-proxy"

	// BuiltInPluginModelRouter 模型路由插件名称
	BuiltInPluginModelRouter = "model-router"

	// BuiltInPluginModelMapper 模型映射插件名称
	BuiltInPluginModelMapper = "model-mapper"

	// BuiltInPluginAiStatistics AI统计插件名称
	BuiltInPluginAiStatistics = "ai-statistics"
)

// AI代理插件配置键
const (
	// AiProxyConfigProviders 提供商列表配置键
	AiProxyConfigProviders = "providers"

	// AiProxyConfigProviderId 提供商ID配置键
	AiProxyConfigProviderId = "id"

	// AiProxyConfigProviderType 提供商类型配置键
	AiProxyConfigProviderType = "type"

	// AiProxyConfigProtocol 协议配置键
	AiProxyConfigProtocol = "protocol"

	// AiProxyConfigProviderApiTokens API Token列表配置键
	AiProxyConfigProviderApiTokens = "apiTokens"

	// AiProxyConfigFailover 故障转移配置键
	AiProxyConfigFailover = "failover"

	// AiProxyConfigFailoverEnabled 故障转移启用配置键
	AiProxyConfigFailoverEnabled = "enabled"

	// AiProxyConfigFailoverFailureThreshold 故障阈值配置键
	AiProxyConfigFailoverFailureThreshold = "failureThreshold"

	// AiProxyConfigFailoverSuccessThreshold 成功阈值配置键
	AiProxyConfigFailoverSuccessThreshold = "successThreshold"

	// AiProxyConfigFailoverHealthCheckInterval 健康检查间隔配置键
	AiProxyConfigFailoverHealthCheckInterval = "healthCheckInterval"

	// AiProxyConfigFailoverHealthCheckTimeout 健康检查超时配置键
	AiProxyConfigFailoverHealthCheckTimeout = "healthCheckTimeout"

	// AiProxyConfigFailoverHealthCheckModel 健康检查模型配置键
	AiProxyConfigFailoverHealthCheckModel = "healthCheckModel"

	// AiProxyConfigRetryOnFailure 失败重试配置键
	AiProxyConfigRetryOnFailure = "retryOnFailure"

	// AiProxyConfigRetryEnabled 重试启用配置键
	AiProxyConfigRetryEnabled = "enabled"

	// AiProxyConfigActiveProviderId 活跃提供商ID配置键
	AiProxyConfigActiveProviderId = "activeProviderId"
)

// 模型路由插件配置键
const (
	// ModelRouterConfigModelToHeader 模型到头的映射配置键
	ModelRouterConfigModelToHeader = "modelToHeader"
)

// 模型映射插件配置键
const (
	// ModelMapperConfigModelMapping 模型映射配置键
	ModelMapperConfigModelMapping = "modelMapping"
)

// AI统计插件配置键
const (
	// AiStatisticsConfigAttributes 属性配置键
	AiStatisticsConfigAttributes = "attributes"

	// AiStatisticsConfigValueSource 值来源
	AiStatisticsConfigValueSource = "valueSource"

	// AiStatisticsConfigKey 键
	AiStatisticsConfigKey = "key"

	// AiStatisticsConfigValue 值
	AiStatisticsConfigValue = "value"

	// AiStatisticsConfigRule 规则
	AiStatisticsConfigRule = "rule"

	// AiStatisticsConfigAppend 追加模式
	AiStatisticsConfigAppend = "append"

	// AiStatisticsConfigRequestBody 请求体来源
	AiStatisticsConfigRequestBody = "request_body"

	// AiStatisticsConfigResponseBody 响应体来源
	AiStatisticsConfigResponseBody = "response_body"

	// AiStatisticsConfigResponseStreamingBody 流式响应体来源
	AiStatisticsConfigResponseStreamingBody = "response_streaming_body"
)
