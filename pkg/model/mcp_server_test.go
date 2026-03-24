// Package model provides data models for the SDK
package model

import (
	"testing"

	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model/route"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMcpServerTypeEnum_Constants(t *testing.T) {
	assert.Equal(t, McpServerTypeEnum("OPEN_API"), McpServerTypeOpenApi)
	assert.Equal(t, McpServerTypeEnum("DATABASE"), McpServerTypeDatabase)
	assert.Equal(t, McpServerTypeEnum("DIRECT_ROUTE"), McpServerTypeDirectRoute)
}

func TestMcpServerDBTypeEnum_Constants(t *testing.T) {
	assert.Equal(t, McpServerDBTypeEnum("MYSQL"), McpServerDBTypeMysql)
	assert.Equal(t, McpServerDBTypeEnum("POSTGRESQL"), McpServerDBTypePostgresql)
	assert.Equal(t, McpServerDBTypeEnum("SQLITE"), McpServerDBTypeSqlite)
	assert.Equal(t, McpServerDBTypeEnum("CLICKHOUSE"), McpServerDBTypeClickhouse)
}

func TestMcpServer_Structure(t *testing.T) {
	server := &McpServer{
		Name:        "test-mcp-server",
		Description: "Test MCP Server",
		Domains:     []string{"example.com", "api.example.com"},
		Services: []route.UpstreamService{
			{Name: "upstream-service", Namespace: "default", Port: 8080},
		},
		Type:              McpServerTypeOpenApi,
		RawConfigurations: "key: value",
	}

	assert.Equal(t, "test-mcp-server", server.Name)
	assert.Equal(t, "Test MCP Server", server.Description)
	assert.Len(t, server.Domains, 2)
	assert.Len(t, server.Services, 1)
	assert.Equal(t, McpServerTypeOpenApi, server.Type)
	assert.Equal(t, "key: value", server.RawConfigurations)
}

func TestMcpServer_DatabaseType(t *testing.T) {
	server := &McpServer{
		Name:   "db-mcp-server",
		Type:   McpServerTypeDatabase,
		DBType: McpServerDBTypeMysql,
		DBConfig: &McpServerDBConfig{
			Host:     "localhost",
			Port:     3306,
			Database: "testdb",
			Username: "root",
			Password: "password",
		},
	}

	assert.Equal(t, McpServerTypeDatabase, server.Type)
	assert.Equal(t, McpServerDBTypeMysql, server.DBType)
	require.NotNil(t, server.DBConfig)
	assert.Equal(t, "localhost", server.DBConfig.Host)
	assert.Equal(t, 3306, server.DBConfig.Port)
	assert.Equal(t, "testdb", server.DBConfig.Database)
	assert.Equal(t, "root", server.DBConfig.Username)
	assert.Equal(t, "password", server.DBConfig.Password)
}

func TestMcpServer_DirectRouteType(t *testing.T) {
	server := &McpServer{
		Name: "direct-route-mcp-server",
		Type: McpServerTypeDirectRoute,
		DirectRouteConfig: &McpServerDirectRouteConfig{
			UpstreamProtocol:   "http",
			DownstreamProtocol: "https",
		},
	}

	assert.Equal(t, McpServerTypeDirectRoute, server.Type)
	require.NotNil(t, server.DirectRouteConfig)
	assert.Equal(t, "http", server.DirectRouteConfig.UpstreamProtocol)
	assert.Equal(t, "https", server.DirectRouteConfig.DownstreamProtocol)
}

func TestMcpServer_WithConsumerAuthInfo(t *testing.T) {
	server := &McpServer{
		Name: "auth-mcp-server",
		Type: McpServerTypeOpenApi,
		ConsumerAuthInfo: &ConsumerAuthInfo{
			Enable:           true,
			Type:             "key-auth",
			AllowedConsumers: []string{"consumer1", "consumer2"},
		},
	}

	require.NotNil(t, server.ConsumerAuthInfo)
	assert.True(t, server.ConsumerAuthInfo.Enable)
	assert.Equal(t, "key-auth", server.ConsumerAuthInfo.Type)
	assert.Len(t, server.ConsumerAuthInfo.AllowedConsumers, 2)
}

func TestConsumerAuthInfo_Structure(t *testing.T) {
	authInfo := &ConsumerAuthInfo{
		Enable:           true,
		Type:             "key-auth",
		AllowedConsumers: []string{"consumer1", "consumer2", "consumer3"},
	}

	assert.True(t, authInfo.Enable)
	assert.Equal(t, "key-auth", authInfo.Type)
	assert.Len(t, authInfo.AllowedConsumers, 3)
}

func TestConsumerAuthInfo_Disabled(t *testing.T) {
	authInfo := &ConsumerAuthInfo{
		Enable: false,
	}

	assert.False(t, authInfo.Enable)
}

func TestMcpServerDBConfig_Structure(t *testing.T) {
	config := &McpServerDBConfig{
		Host:     "db.example.com",
		Port:     5432,
		Database: "production",
		Username: "admin",
		Password: "secret",
	}

	assert.Equal(t, "db.example.com", config.Host)
	assert.Equal(t, 5432, config.Port)
	assert.Equal(t, "production", config.Database)
	assert.Equal(t, "admin", config.Username)
	assert.Equal(t, "secret", config.Password)
}

func TestMcpServerDirectRouteConfig_Structure(t *testing.T) {
	config := &McpServerDirectRouteConfig{
		UpstreamProtocol:   "grpc",
		DownstreamProtocol: "http",
	}

	assert.Equal(t, "grpc", config.UpstreamProtocol)
	assert.Equal(t, "http", config.DownstreamProtocol)
}

func TestMcpServer_AllTypes(t *testing.T) {
	types := []McpServerTypeEnum{
		McpServerTypeOpenApi,
		McpServerTypeDatabase,
		McpServerTypeDirectRoute,
	}

	for _, serverType := range types {
		server := &McpServer{
			Name: "test-server",
			Type: serverType,
		}
		assert.Equal(t, serverType, server.Type)
	}
}

func TestMcpServer_AllDBTypes(t *testing.T) {
	dbTypes := []McpServerDBTypeEnum{
		McpServerDBTypeMysql,
		McpServerDBTypePostgresql,
		McpServerDBTypeSqlite,
		McpServerDBTypeClickhouse,
	}

	for _, dbType := range dbTypes {
		server := &McpServer{
			Name:   "test-server",
			Type:   McpServerTypeDatabase,
			DBType: dbType,
		}
		assert.Equal(t, dbType, server.DBType)
	}
}

func TestMcpServer_EmptyFields(t *testing.T) {
	server := &McpServer{}

	assert.Empty(t, server.Name)
	assert.Empty(t, server.Description)
	assert.Nil(t, server.Domains)
	assert.Nil(t, server.Services)
	assert.Empty(t, server.Type)
	assert.Nil(t, server.ConsumerAuthInfo)
	assert.Empty(t, server.RawConfigurations)
	assert.Nil(t, server.DBConfig)
	assert.Empty(t, server.DBType)
	assert.Nil(t, server.DirectRouteConfig)
}

func TestMcpServer_MultipleServices(t *testing.T) {
	weight1, weight2 := 50, 50
	server := &McpServer{
		Name: "multi-service-mcp",
		Type: McpServerTypeOpenApi,
		Services: []route.UpstreamService{
			{Name: "service1", Namespace: "ns1", Port: 8080, Weight: &weight1},
			{Name: "service2", Namespace: "ns2", Port: 8081, Weight: &weight2},
		},
	}

	assert.Len(t, server.Services, 2)
	assert.Equal(t, "service1", server.Services[0].Name)
	assert.Equal(t, "service2", server.Services[1].Name)
}

func TestMcpServer_MultipleDomains(t *testing.T) {
	server := &McpServer{
		Name: "multi-domain-mcp",
		Type: McpServerTypeOpenApi,
		Domains: []string{
			"api1.example.com",
			"api2.example.com",
			"api3.example.com",
		},
	}

	assert.Len(t, server.Domains, 3)
}
