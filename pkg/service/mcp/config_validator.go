// Package mcp provides MCP server related services
package mcp

import (
	"fmt"

	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/errors"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
)

// McpServerDBConfigValidator 数据库配置验证器
type McpServerDBConfigValidator struct{}

// NewMcpServerDBConfigValidator 创建数据库配置验证器
func NewMcpServerDBConfigValidator() *McpServerDBConfigValidator {
	return &McpServerDBConfigValidator{}
}

// Validate 验证数据库配置
func (v *McpServerDBConfigValidator) Validate(config *model.McpServerDBConfig, dbType model.McpServerDBTypeEnum) error {
	if config == nil {
		return errors.NewValidationError("database config is required")
	}

	switch dbType {
	case model.McpServerDBTypeMysql:
		return v.ValidateMySQL(config)
	case model.McpServerDBTypePostgresql:
		return v.ValidatePostgreSQL(config)
	case model.McpServerDBTypeSqlite:
		return v.ValidateSQLite(config)
	case model.McpServerDBTypeClickhouse:
		return v.ValidateClickHouse(config)
	default:
		return errors.NewValidationError(fmt.Sprintf("unsupported database type: %s", dbType))
	}
}

// ValidateMySQL 验证MySQL配置
func (v *McpServerDBConfigValidator) ValidateMySQL(config *model.McpServerDBConfig) error {
	if config.Host == "" {
		return errors.NewValidationErrorWithField("host is required for MySQL", "host")
	}
	if config.Database == "" {
		return errors.NewValidationErrorWithField("database is required for MySQL", "database")
	}
	if config.Username == "" {
		return errors.NewValidationErrorWithField("username is required for MySQL", "username")
	}
	if config.Port == 0 {
		// 使用默认端口
		return nil
	}
	if config.Port < 1 || config.Port > 65535 {
		return errors.NewValidationErrorWithField("port must be between 1 and 65535", "port")
	}
	return nil
}

// ValidatePostgreSQL 验证PostgreSQL配置
func (v *McpServerDBConfigValidator) ValidatePostgreSQL(config *model.McpServerDBConfig) error {
	if config.Host == "" {
		return errors.NewValidationErrorWithField("host is required for PostgreSQL", "host")
	}
	if config.Database == "" {
		return errors.NewValidationErrorWithField("database is required for PostgreSQL", "database")
	}
	if config.Username == "" {
		return errors.NewValidationErrorWithField("username is required for PostgreSQL", "username")
	}
	if config.Port == 0 {
		// 使用默认端口
		return nil
	}
	if config.Port < 1 || config.Port > 65535 {
		return errors.NewValidationErrorWithField("port must be between 1 and 65535", "port")
	}
	return nil
}

// ValidateSQLite 验证SQLite配置
func (v *McpServerDBConfigValidator) ValidateSQLite(config *model.McpServerDBConfig) error {
	if config.Database == "" {
		return errors.NewValidationErrorWithField("database path is required for SQLite", "database")
	}
	return nil
}

// ValidateClickHouse 验证ClickHouse配置
func (v *McpServerDBConfigValidator) ValidateClickHouse(config *model.McpServerDBConfig) error {
	if config.Host == "" {
		return errors.NewValidationErrorWithField("host is required for ClickHouse", "host")
	}
	if config.Database == "" {
		return errors.NewValidationErrorWithField("database is required for ClickHouse", "database")
	}
	if config.Port == 0 {
		// 使用默认端口
		return nil
	}
	if config.Port < 1 || config.Port > 65535 {
		return errors.NewValidationErrorWithField("port must be between 1 and 65535", "port")
	}
	return nil
}

// ValidateMcpServer 验证MCP服务器配置
func (v *McpServerDBConfigValidator) ValidateMcpServer(server *model.McpServer) error {
	if server == nil {
		return errors.NewValidationError("MCP server is nil")
	}

	if server.Name == "" {
		return errors.NewValidationErrorWithField("name is required", "name")
	}

	// 根据类型验证
	switch server.Type {
	case model.McpServerTypeOpenApi:
		// OpenAPI类型需要rawConfigurations
		if server.RawConfigurations == "" {
			return errors.NewValidationErrorWithField("rawConfigurations is required for OPEN_API type", "rawConfigurations")
		}
	case model.McpServerTypeDatabase:
		// 数据库类型需要dbConfig和dbType
		if server.DBConfig == nil {
			return errors.NewValidationErrorWithField("dbConfig is required for DATABASE type", "dbConfig")
		}
		if server.DBType == "" {
			return errors.NewValidationErrorWithField("dbType is required for DATABASE type", "dbType")
		}
		if err := v.Validate(server.DBConfig, server.DBType); err != nil {
			return err
		}
	case model.McpServerTypeDirectRoute:
		// 直连路由类型需要services
		if len(server.Services) == 0 {
			return errors.NewValidationErrorWithField("services is required for DIRECT_ROUTE type", "services")
		}
	case "":
		return errors.NewValidationErrorWithField("type is required", "type")
	default:
		return errors.NewValidationErrorWithField(fmt.Sprintf("unknown type: %s", server.Type), "type")
	}

	return nil
}
