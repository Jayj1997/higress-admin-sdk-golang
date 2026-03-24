// Package mcp provides MCP server related services
package mcp

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
)

// McpServerDBConfigDsnConverter 数据库配置DSN转换器
type McpServerDBConfigDsnConverter struct{}

// NewMcpServerDBConfigDsnConverter 创建数据库配置DSN转换器
func NewMcpServerDBConfigDsnConverter() *McpServerDBConfigDsnConverter {
	return &McpServerDBConfigDsnConverter{}
}

// ConvertToDsn 将数据库配置转换为DSN字符串
func (c *McpServerDBConfigDsnConverter) ConvertToDsn(config *model.McpServerDBConfig, dbType model.McpServerDBTypeEnum) (string, error) {
	if config == nil {
		return "", fmt.Errorf("database config is nil")
	}

	switch dbType {
	case model.McpServerDBTypeMysql:
		return c.convertToMySQLDsn(config), nil
	case model.McpServerDBTypePostgresql:
		return c.convertToPostgreSQLDsn(config), nil
	case model.McpServerDBTypeSqlite:
		return c.convertToSQLiteDsn(config), nil
	case model.McpServerDBTypeClickhouse:
		return c.convertToClickHouseDsn(config), nil
	default:
		return "", fmt.Errorf("unsupported database type: %s", dbType)
	}
}

// ConvertFromDsn 从DSN字符串解析数据库配置
func (c *McpServerDBConfigDsnConverter) ConvertFromDsn(dsn string, dbType model.McpServerDBTypeEnum) (*model.McpServerDBConfig, error) {
	if dsn == "" {
		return nil, fmt.Errorf("dsn is empty")
	}

	switch dbType {
	case model.McpServerDBTypeMysql:
		return c.convertFromMySQLDsn(dsn)
	case model.McpServerDBTypePostgresql:
		return c.convertFromPostgreSQLDsn(dsn)
	case model.McpServerDBTypeSqlite:
		return c.convertFromSQLiteDsn(dsn)
	case model.McpServerDBTypeClickhouse:
		return c.convertFromClickHouseDsn(dsn)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}

// convertToMySQLDsn 转换为MySQL DSN
// 格式: username:password@tcp(host:port)/database
func (c *McpServerDBConfigDsnConverter) convertToMySQLDsn(config *model.McpServerDBConfig) string {
	port := config.Port
	if port == 0 {
		port = 3306
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.Username,
		config.Password,
		config.Host,
		port,
		config.Database,
	)
}

// convertFromMySQLDsn 从MySQL DSN解析
func (c *McpServerDBConfigDsnConverter) convertFromMySQLDsn(dsn string) (*model.McpServerDBConfig, error) {
	// 格式: username:password@tcp(host:port)/database
	config := &model.McpServerDBConfig{}

	// 分离用户名密码和连接信息
	atIndex := strings.Index(dsn, "@")
	if atIndex == -1 {
		return nil, fmt.Errorf("invalid MySQL DSN format: missing @")
	}

	userPass := dsn[:atIndex]
	connPart := dsn[atIndex+1:]

	// 解析用户名密码
	colonIndex := strings.Index(userPass, ":")
	if colonIndex == -1 {
		config.Username = userPass
	} else {
		config.Username = userPass[:colonIndex]
		config.Password = userPass[colonIndex+1:]
	}

	// 解析连接信息
	// 格式: tcp(host:port)/database
	if !strings.HasPrefix(connPart, "tcp(") {
		return nil, fmt.Errorf("invalid MySQL DSN format: missing tcp()")
	}

	tcpEnd := strings.Index(connPart, ")")
	if tcpEnd == -1 {
		return nil, fmt.Errorf("invalid MySQL DSN format: missing )")
	}

	hostPort := connPart[4:tcpEnd]
	databasePart := connPart[tcpEnd+1:]

	// 解析主机端口
	lastColon := strings.LastIndex(hostPort, ":")
	if lastColon == -1 {
		config.Host = hostPort
	} else {
		config.Host = hostPort[:lastColon]
		port, err := strconv.Atoi(hostPort[lastColon+1:])
		if err == nil {
			config.Port = port
		}
	}

	// 解析数据库名
	if strings.HasPrefix(databasePart, "/") {
		config.Database = databasePart[1:]
	}

	return config, nil
}

// convertToPostgreSQLDsn 转换为PostgreSQL DSN
// 格式: host=host port=port user=username password=password dbname=database
func (c *McpServerDBConfigDsnConverter) convertToPostgreSQLDsn(config *model.McpServerDBConfig) string {
	port := config.Port
	if port == 0 {
		port = 5432
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		config.Host,
		port,
		config.Username,
		config.Password,
		config.Database,
	)
}

// convertFromPostgreSQLDsn 从PostgreSQL DSN解析
func (c *McpServerDBConfigDsnConverter) convertFromPostgreSQLDsn(dsn string) (*model.McpServerDBConfig, error) {
	config := &model.McpServerDBConfig{}

	pairs := strings.Split(dsn, " ")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := kv[0]
		value := kv[1]

		switch key {
		case "host":
			config.Host = value
		case "port":
			port, err := strconv.Atoi(value)
			if err == nil {
				config.Port = port
			}
		case "user":
			config.Username = value
		case "password":
			config.Password = value
		case "dbname":
			config.Database = value
		}
	}

	return config, nil
}

// convertToSQLiteDsn 转换为SQLite DSN
// 格式: file:path
func (c *McpServerDBConfigDsnConverter) convertToSQLiteDsn(config *model.McpServerDBConfig) string {
	return fmt.Sprintf("file:%s", config.Database)
}

// convertFromSQLiteDsn 从SQLite DSN解析
func (c *McpServerDBConfigDsnConverter) convertFromSQLiteDsn(dsn string) (*model.McpServerDBConfig, error) {
	config := &model.McpServerDBConfig{}

	if strings.HasPrefix(dsn, "file:") {
		config.Database = dsn[5:]
	} else {
		config.Database = dsn
	}

	return config, nil
}

// convertToClickHouseDsn 转换为ClickHouse DSN
// 格式: tcp://host:port?database=database&username=username&password=password
func (c *McpServerDBConfigDsnConverter) convertToClickHouseDsn(config *model.McpServerDBConfig) string {
	port := config.Port
	if port == 0 {
		port = 9000
	}
	return fmt.Sprintf("tcp://%s:%d?database=%s&username=%s&password=%s",
		config.Host,
		port,
		url.QueryEscape(config.Database),
		url.QueryEscape(config.Username),
		url.QueryEscape(config.Password),
	)
}

// convertFromClickHouseDsn 从ClickHouse DSN解析
func (c *McpServerDBConfigDsnConverter) convertFromClickHouseDsn(dsn string) (*model.McpServerDBConfig, error) {
	config := &model.McpServerDBConfig{}

	u, err := url.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ClickHouse DSN: %w", err)
	}

	config.Host = u.Hostname()
	if u.Port() != "" {
		port, err := strconv.Atoi(u.Port())
		if err == nil {
			config.Port = port
		}
	}

	query := u.Query()
	config.Database = query.Get("database")
	config.Username = query.Get("username")
	config.Password = query.Get("password")

	return config, nil
}
