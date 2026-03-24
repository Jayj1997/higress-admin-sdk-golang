// Package constant provides constants for the SDK
package constant

// MCP相关常量
const (
	// MCP服务器路径前缀
	McpServerPathPre = "/mcp-server/"

	// MCP服务器路由前缀
	McpServerRoutePrefix = "mcp-server-"

	// MCP临时目录
	McpTempDir = "/tmp/mcp"

	// MCP工作目录
	McpWorkDir = "/etc/higress/mcp"

	// OpenAPI到MCP服务器脚本路径
	OpenApiToMcpServerScriptPath = "/usr/local/bin/openapi-to-mcpserver.sh"
)

// ConfigMap名称常量
const (
	// HigressConfig Higress配置ConfigMap名称
	HigressConfig = "higress-config"
)

// MCP服务器标签常量
const (
	// LabelResourceMcpServerTypeKey MCP服务器类型标签键
	LabelResourceMcpServerTypeKey = "mcp-server-type"

	// LabelMcpServerBizTypeValue MCP服务器业务类型标签值
	LabelMcpServerBizTypeValue = "mcp-server"
)

// MCP服务器注解常量
const (
	// AnnotationResourceDescriptionKey 资源描述注解键
	AnnotationResourceDescriptionKey = "description"

	// AnnotationResourceMcpServerKey MCP服务器注解键
	AnnotationResourceMcpServerKey = "mcp-server"

	// AnnotationResourceMcpServerMatchRuleDomainsKey MCP服务器匹配规则域名注解键
	AnnotationResourceMcpServerMatchRuleDomainsKey = "mcp-server-match-rule-domains"

	// AnnotationResourceMcpServerMatchRuleTypeKey MCP服务器匹配规则类型注解键
	AnnotationResourceMcpServerMatchRuleTypeKey = "mcp-server-match-rule-type"

	// AnnotationResourceMcpServerMatchRuleValueKey MCP服务器匹配规则值注解键
	AnnotationResourceMcpServerMatchRuleValueKey = "mcp-server-match-rule-value"
)

// ConfigMap数据键常量
const (
	// McpConfigKey MCP配置键
	McpConfigKey = "higress"

	// McpServerKey MCP服务器键
	McpServerKey = "mcpServer"

	// MatchListKey 匹配列表键
	MatchListKey = "match_list"

	// ServersKey 服务器列表键
	ServersKey = "servers"
)

// 匹配规则常量
const (
	// MatchRulePathKey 匹配规则路径键
	MatchRulePathKey = "match_rule_path"

	// MatchRuleDomainKey 匹配规则域名键
	MatchRuleDomainKey = "match_rule_domain"

	// MatchRuleTypeKey 匹配规则类型键
	MatchRuleTypeKey = "match_rule_type"

	// ServerNameKey 服务器名称键
	ServerNameKey = "name"
)
