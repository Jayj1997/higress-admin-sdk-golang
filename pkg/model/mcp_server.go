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
	// Enabled 是否启用认证
	Enabled bool `json:"enabled,omitempty"`

	// Consumers 允许的消费者列表
	Consumers []string `json:"consumers,omitempty"`
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
}

// McpServerConsumer MCP服务器消费者
type McpServerConsumer struct {
	// ConsumerName 消费者名称
	ConsumerName string `json:"consumerName,omitempty"`
}

// McpServerConsumerDetail MCP服务器消费者详情
type McpServerConsumerDetail struct {
	McpServerConsumer
	// Consumer 消费者详情
	Consumer *Consumer `json:"consumer,omitempty"`
}
