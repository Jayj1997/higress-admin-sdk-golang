// Package model provides data models for the SDK
package model

import (
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model/route"
)

// McpServer MCP服务器配置
type McpServer struct {
	// ID MCP服务器ID（已废弃）
	ID string `json:"id,omitempty"`

	// Name MCP服务器名称
	Name string `json:"name,omitempty"`

	// Description MCP服务器描述
	Description string `json:"description,omitempty"`

	// Domains MCP服务器适用的域名列表
	Domains []string `json:"domains,omitempty"`

	// Services MCP服务器上游服务列表
	Services []route.UpstreamService `json:"services,omitempty"`

	// Type MCP服务器类型
	Type McpServerTypeEnum `json:"type,omitempty"`

	// ConsumerAuthInfo 消费者认证信息
	ConsumerAuthInfo *ConsumerAuthInfo `json:"consumerAuthInfo,omitempty"`

	// RawConfigurations YAML格式的原始配置
	RawConfigurations string `json:"rawConfigurations,omitempty"`

	// DBConfig 数据库配置，对于DATABASE类型服务器是必需的
	DBConfig *McpServerDBConfig `json:"dbConfig,omitempty"`

	// DBType 数据库类型
	DBType McpServerDBTypeEnum `json:"dbType,omitempty"`

	// DirectRouteConfig 直连路由配置
	DirectRouteConfig *McpServerDirectRouteConfig `json:"directRouteConfig,omitempty"`
}

// McpServerTypeEnum MCP服务器类型
type McpServerTypeEnum string

const (
	McpServerTypeOpenApi     McpServerTypeEnum = "OPEN_API"
	McpServerTypeDatabase    McpServerTypeEnum = "DATABASE"
	McpServerTypeDirectRoute McpServerTypeEnum = "DIRECT_ROUTE"
)

// McpServerDBTypeEnum MCP服务器数据库类型
type McpServerDBTypeEnum string

const (
	McpServerDBTypeMysql      McpServerDBTypeEnum = "MYSQL"
	McpServerDBTypePostgresql McpServerDBTypeEnum = "POSTGRESQL"
	McpServerDBTypeSqlite     McpServerDBTypeEnum = "SQLITE"
	McpServerDBTypeClickhouse McpServerDBTypeEnum = "CLICKHOUSE"
)

// ConsumerAuthInfo 消费者认证信息
type ConsumerAuthInfo struct {
	// Enable 是否启用认证
	Enable bool `json:"enable,omitempty"`

	// Type 凭证类型
	Type string `json:"type,omitempty"`

	// AllowedConsumers 允许的消费者列表
	AllowedConsumers []string `json:"allowedConsumers,omitempty"`
}

// McpServerDBConfig MCP服务器数据库配置
type McpServerDBConfig struct {
	// Host 数据库主机
	Host string `json:"host,omitempty"`

	// Port 数据库端口
	Port int `json:"port,omitempty"`

	// Database 数据库名称
	Database string `json:"database,omitempty"`

	// Username 数据库用户名
	Username string `json:"username,omitempty"`

	// Password 数据库密码
	Password string `json:"password,omitempty"`
}

// McpServerDirectRouteConfig MCP服务器直连路由配置
type McpServerDirectRouteConfig struct {
	// UpstreamProtocol 上游协议
	UpstreamProtocol string `json:"upstreamProtocol,omitempty"`

	// DownstreamProtocol 下游协议
	DownstreamProtocol string `json:"downstreamProtocol,omitempty"`
}

// McpServerPageQuery MCP服务器分页查询
type McpServerPageQuery struct {
	CommonPageQuery
	// McpServerName MCP服务器名称（模糊匹配）
	McpServerName string `json:"mcpServerName,omitempty"`
	// Type MCP服务器类型
	Type string `json:"type,omitempty"`
}

// McpServerConsumer MCP服务器消费者
type McpServerConsumer struct {
	// ConsumerName 消费者名称
	ConsumerName string `json:"consumerName,omitempty"`
}

// McpServerConsumerDetail MCP服务器消费者详情
type McpServerConsumerDetail struct {
	McpServerConsumer
	// McpServerName MCP服务器名称
	McpServerName string `json:"mcpServerName,omitempty"`
	// Consumer 消费者详情
	Consumer *Consumer `json:"consumer,omitempty"`
}

// McpServerConsumers MCP服务器消费者列表
type McpServerConsumers struct {
	// McpServerName MCP服务器名称
	McpServerName string `json:"mcpServerName,omitempty"`
	// Consumers 消费者名称列表
	Consumers []string `json:"consumers,omitempty"`
}

// McpServerConsumersPageQuery MCP服务器消费者分页查询
type McpServerConsumersPageQuery struct {
	CommonPageQuery
	// McpServerName MCP服务器名称
	McpServerName string `json:"mcpServerName,omitempty"`
	// ConsumerName 消费者名称（模糊匹配）
	ConsumerName string `json:"consumerName,omitempty"`
}

// McpServerConfigMap MCP服务器ConfigMap配置
type McpServerConfigMap struct {
	// Servers 服务器列表
	Servers []McpServerConfigMapServer `json:"servers,omitempty"`
	// MatchList 匹配规则列表
	MatchList []McpServerConfigMapMatchList `json:"match_list,omitempty"`
}

// McpServerConfigMapServer MCP服务器ConfigMap服务器配置
type McpServerConfigMapServer struct {
	// Name 服务器名称
	Name string `json:"name,omitempty"`
	// Config 服务器配置
	Config map[string]interface{} `json:"config,omitempty"`
}

// McpServerConfigMapMatchList MCP服务器ConfigMap匹配规则
type McpServerConfigMapMatchList struct {
	// MatchRulePath 匹配规则路径
	MatchRulePath string `json:"match_rule_path,omitempty"`
	// MatchRuleDomain 匹配规则域名
	MatchRuleDomain string `json:"match_rule_domain,omitempty"`
	// MatchRuleType 匹配规则类型
	MatchRuleType string `json:"match_rule_type,omitempty"`
}

// SwaggerContent Swagger内容
type SwaggerContent struct {
	// Content Swagger文件内容
	Content string `json:"content,omitempty"`
}

// ParseMcpServerTypeEnum 从字符串解析MCP服务器类型
func ParseMcpServerTypeEnum(name string) McpServerTypeEnum {
	switch name {
	case "OPEN_API", "open_api":
		return McpServerTypeOpenApi
	case "DATABASE", "database":
		return McpServerTypeDatabase
	case "DIRECT_ROUTE", "direct_route":
		return McpServerTypeDirectRoute
	default:
		return ""
	}
}

// Value 获取MCP服务器类型的值
func (t McpServerTypeEnum) Value() string {
	switch t {
	case McpServerTypeOpenApi:
		return "open_api"
	case McpServerTypeDatabase:
		return "database"
	case McpServerTypeDirectRoute:
		return "direct_route"
	default:
		return ""
	}
}

// ParseMcpServerDBTypeEnum 从字符串解析数据库类型
func ParseMcpServerDBTypeEnum(name string) McpServerDBTypeEnum {
	switch name {
	case "MYSQL", "mysql":
		return McpServerDBTypeMysql
	case "POSTGRESQL", "postgresql", "postgres":
		return McpServerDBTypePostgresql
	case "SQLITE", "sqlite":
		return McpServerDBTypeSqlite
	case "CLICKHOUSE", "clickhouse":
		return McpServerDBTypeClickhouse
	default:
		return ""
	}
}

// Value 获取数据库类型的值
func (t McpServerDBTypeEnum) Value() string {
	switch t {
	case McpServerDBTypeMysql:
		return "mysql"
	case McpServerDBTypePostgresql:
		return "postgres"
	case McpServerDBTypeSqlite:
		return "sqlite"
	case McpServerDBTypeClickhouse:
		return "clickhouse"
	default:
		return ""
	}
}
