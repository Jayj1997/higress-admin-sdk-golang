// Package service provides business services for the SDK
package service

import (
	"testing"

	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
	"github.com/stretchr/testify/assert"
)

// TestMcpServerService_Interface tests that McpServerService interface is defined
func TestMcpServerService_Interface(t *testing.T) {
	// This test verifies the interface exists and has the expected methods
	var _ McpServerService = (McpServerService)(nil)
}

// TestMcpServerModel tests the McpServer model
func TestMcpServerModel(t *testing.T) {
	server := model.McpServer{
		Name:        "test-mcp-server",
		Description: "Test MCP Server",
		Type:        model.McpServerTypeOpenApi,
	}

	assert.Equal(t, "test-mcp-server", server.Name)
	assert.Equal(t, "Test MCP Server", server.Description)
	assert.Equal(t, model.McpServerTypeOpenApi, server.Type)
}

// TestMcpServerWithDomains tests McpServer with domains
func TestMcpServerWithDomains(t *testing.T) {
	server := model.McpServer{
		Name:        "test-server",
		Description: "Test server",
		Domains:     []string{"example.com", "api.example.com"},
		Type:        model.McpServerTypeOpenApi,
	}

	assert.Len(t, server.Domains, 2)
	assert.Contains(t, server.Domains, "example.com")
}

// TestMcpServerWithDBConfig tests McpServer with database config
func TestMcpServerWithDBConfig(t *testing.T) {
	server := model.McpServer{
		Name:        "db-server",
		Description: "Database MCP Server",
		Type:        model.McpServerTypeDatabase,
		DBType:      model.McpServerDBTypeMysql,
		DBConfig: &model.McpServerDBConfig{
			Host:     "localhost",
			Port:     3306,
			Database: "testdb",
			Username: "root",
			Password: "password",
		},
	}

	assert.Equal(t, model.McpServerTypeDatabase, server.Type)
	assert.Equal(t, model.McpServerDBTypeMysql, server.DBType)
	assert.NotNil(t, server.DBConfig)
	assert.Equal(t, "localhost", server.DBConfig.Host)
	assert.Equal(t, 3306, server.DBConfig.Port)
}

// TestMcpServerWithDirectRouteConfig tests McpServer with direct route config
func TestMcpServerWithDirectRouteConfig(t *testing.T) {
	server := model.McpServer{
		Name:        "direct-route-server",
		Description: "Direct Route MCP Server",
		Type:        model.McpServerTypeDirectRoute,
		DirectRouteConfig: &model.McpServerDirectRouteConfig{
			UpstreamProtocol:   "http",
			DownstreamProtocol: "https",
		},
	}

	assert.Equal(t, model.McpServerTypeDirectRoute, server.Type)
	assert.NotNil(t, server.DirectRouteConfig)
	assert.Equal(t, "http", server.DirectRouteConfig.UpstreamProtocol)
}

// TestMcpServerWithConsumerAuthInfo tests McpServer with consumer auth info
func TestMcpServerWithConsumerAuthInfo(t *testing.T) {
	server := model.McpServer{
		Name:        "auth-server",
		Description: "Server with auth",
		Type:        model.McpServerTypeOpenApi,
		ConsumerAuthInfo: &model.ConsumerAuthInfo{
			Enable:           true,
			Type:             "key-auth",
			AllowedConsumers: []string{"consumer1", "consumer2"},
		},
	}

	assert.NotNil(t, server.ConsumerAuthInfo)
	assert.True(t, server.ConsumerAuthInfo.Enable)
	assert.Equal(t, "key-auth", server.ConsumerAuthInfo.Type)
	assert.Len(t, server.ConsumerAuthInfo.AllowedConsumers, 2)
}

// TestMcpServerTypeEnum tests MCP server type enum
func TestMcpServerTypeEnum(t *testing.T) {
	tests := []struct {
		name      string
		typeEnum  model.McpServerTypeEnum
		wantValue string
	}{
		{"OPEN_API type", model.McpServerTypeOpenApi, "open_api"},
		{"DATABASE type", model.McpServerTypeDatabase, "database"},
		{"DIRECT_ROUTE type", model.McpServerTypeDirectRoute, "direct_route"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantValue, tt.typeEnum.Value())
		})
	}
}

// TestParseMcpServerTypeEnum tests parsing MCP server type enum
func TestParseMcpServerTypeEnum(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantEnum model.McpServerTypeEnum
	}{
		{"OPEN_API uppercase", "OPEN_API", model.McpServerTypeOpenApi},
		{"OPEN_API lowercase", "open_api", model.McpServerTypeOpenApi},
		{"DATABASE uppercase", "DATABASE", model.McpServerTypeDatabase},
		{"DATABASE lowercase", "database", model.McpServerTypeDatabase},
		{"DIRECT_ROUTE uppercase", "DIRECT_ROUTE", model.McpServerTypeDirectRoute},
		{"DIRECT_ROUTE lowercase", "direct_route", model.McpServerTypeDirectRoute},
		{"Invalid type", "invalid", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantEnum, model.ParseMcpServerTypeEnum(tt.input))
		})
	}
}

// TestMcpServerDBTypeEnum tests MCP server database type enum
func TestMcpServerDBTypeEnum(t *testing.T) {
	tests := []struct {
		name      string
		typeEnum  model.McpServerDBTypeEnum
		wantValue string
	}{
		{"MYSQL type", model.McpServerDBTypeMysql, "mysql"},
		{"POSTGRESQL type", model.McpServerDBTypePostgresql, "postgres"},
		{"SQLITE type", model.McpServerDBTypeSqlite, "sqlite"},
		{"CLICKHOUSE type", model.McpServerDBTypeClickhouse, "clickhouse"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantValue, tt.typeEnum.Value())
		})
	}
}

// TestParseMcpServerDBTypeEnum tests parsing MCP server database type enum
func TestParseMcpServerDBTypeEnum(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantEnum model.McpServerDBTypeEnum
	}{
		{"MYSQL uppercase", "MYSQL", model.McpServerDBTypeMysql},
		{"MYSQL lowercase", "mysql", model.McpServerDBTypeMysql},
		{"POSTGRESQL uppercase", "POSTGRESQL", model.McpServerDBTypePostgresql},
		{"POSTGRESQL lowercase", "postgresql", model.McpServerDBTypePostgresql},
		{"postgres alias", "postgres", model.McpServerDBTypePostgresql},
		{"SQLITE uppercase", "SQLITE", model.McpServerDBTypeSqlite},
		{"SQLITE lowercase", "sqlite", model.McpServerDBTypeSqlite},
		{"CLICKHOUSE uppercase", "CLICKHOUSE", model.McpServerDBTypeClickhouse},
		{"CLICKHOUSE lowercase", "clickhouse", model.McpServerDBTypeClickhouse},
		{"Invalid type", "invalid", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantEnum, model.ParseMcpServerDBTypeEnum(tt.input))
		})
	}
}

// TestMcpServerPageQuery tests the McpServerPageQuery model
func TestMcpServerPageQuery(t *testing.T) {
	query := model.McpServerPageQuery{
		CommonPageQuery: model.CommonPageQuery{
			PageNum:  1,
			PageSize: 10,
		},
		McpServerName: "test-server",
		Type:          "OPEN_API",
	}

	assert.Equal(t, 1, query.PageNum)
	assert.Equal(t, 10, query.PageSize)
	assert.Equal(t, "test-server", query.McpServerName)
	assert.Equal(t, "OPEN_API", query.Type)
}

// TestMcpServerConsumer tests the McpServerConsumer model
func TestMcpServerConsumer(t *testing.T) {
	consumer := model.McpServerConsumer{
		ConsumerName: "test-consumer",
	}

	assert.Equal(t, "test-consumer", consumer.ConsumerName)
}

// TestMcpServerConsumerDetail tests the McpServerConsumerDetail model
func TestMcpServerConsumerDetail(t *testing.T) {
	detail := model.McpServerConsumerDetail{
		McpServerConsumer: model.McpServerConsumer{
			ConsumerName: "test-consumer",
		},
		McpServerName: "test-server",
	}

	assert.Equal(t, "test-consumer", detail.ConsumerName)
	assert.Equal(t, "test-server", detail.McpServerName)
}

// TestMcpServerConsumers tests the McpServerConsumers model
func TestMcpServerConsumers(t *testing.T) {
	consumers := model.McpServerConsumers{
		McpServerName: "test-server",
		Consumers:     []string{"consumer1", "consumer2"},
	}

	assert.Equal(t, "test-server", consumers.McpServerName)
	assert.Len(t, consumers.Consumers, 2)
}

// TestMcpServerConsumersPageQuery tests the McpServerConsumersPageQuery model
func TestMcpServerConsumersPageQuery(t *testing.T) {
	query := model.McpServerConsumersPageQuery{
		CommonPageQuery: model.CommonPageQuery{
			PageNum:  1,
			PageSize: 10,
		},
		McpServerName: "test-server",
		ConsumerName:  "test-consumer",
	}

	assert.Equal(t, 1, query.PageNum)
	assert.Equal(t, 10, query.PageSize)
	assert.Equal(t, "test-server", query.McpServerName)
	assert.Equal(t, "test-consumer", query.ConsumerName)
}

// TestMcpServerConfigMap tests the McpServerConfigMap model
func TestMcpServerConfigMap(t *testing.T) {
	configMap := model.McpServerConfigMap{
		Servers: []model.McpServerConfigMapServer{
			{
				Name: "server1",
				Config: map[string]interface{}{
					"key": "value",
				},
			},
		},
		MatchList: []model.McpServerConfigMapMatchList{
			{
				MatchRulePath:   "/api/*",
				MatchRuleDomain: "example.com",
				MatchRuleType:   "prefix",
			},
		},
	}

	assert.Len(t, configMap.Servers, 1)
	assert.Equal(t, "server1", configMap.Servers[0].Name)
	assert.Len(t, configMap.MatchList, 1)
}

// TestSwaggerContent tests the SwaggerContent model
func TestSwaggerContent(t *testing.T) {
	swagger := model.SwaggerContent{
		Content: "openapi: 3.0.0\ninfo:\n  title: Test API",
	}

	assert.NotEmpty(t, swagger.Content)
	assert.Contains(t, swagger.Content, "openapi")
}

// TestConsumerAuthInfo tests the ConsumerAuthInfo model
func TestConsumerAuthInfo(t *testing.T) {
	authInfo := model.ConsumerAuthInfo{
		Enable:           true,
		Type:             "key-auth",
		AllowedConsumers: []string{"consumer1", "consumer2"},
	}

	assert.True(t, authInfo.Enable)
	assert.Equal(t, "key-auth", authInfo.Type)
	assert.Len(t, authInfo.AllowedConsumers, 2)
}

// TestMcpServerDBConfig tests the McpServerDBConfig model
func TestMcpServerDBConfig(t *testing.T) {
	config := model.McpServerDBConfig{
		Host:     "localhost",
		Port:     5432,
		Database: "testdb",
		Username: "postgres",
		Password: "secret",
	}

	assert.Equal(t, "localhost", config.Host)
	assert.Equal(t, 5432, config.Port)
	assert.Equal(t, "testdb", config.Database)
	assert.Equal(t, "postgres", config.Username)
	assert.Equal(t, "secret", config.Password)
}

// TestMcpServerDirectRouteConfig tests the McpServerDirectRouteConfig model
func TestMcpServerDirectRouteConfig(t *testing.T) {
	config := model.McpServerDirectRouteConfig{
		UpstreamProtocol:   "grpc",
		DownstreamProtocol: "http",
	}

	assert.Equal(t, "grpc", config.UpstreamProtocol)
	assert.Equal(t, "http", config.DownstreamProtocol)
}

// TestMcpServerPagination tests pagination for MCP servers
func TestMcpServerPagination(t *testing.T) {
	servers := make([]model.McpServer, 25)
	for i := 0; i < 25; i++ {
		servers[i] = model.McpServer{
			Name:        "server-" + string(rune('a'+i%26)),
			Description: "Test server",
		}
	}

	total := len(servers)
	pageNum := 1
	pageSize := 10
	start := (pageNum - 1) * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}
	pagedData := servers[start:end]
	result := model.NewPaginatedResult(pagedData, total, pageNum, pageSize)

	assert.Equal(t, 10, len(result.Data))
	assert.Equal(t, 25, result.Total)
	assert.Equal(t, 3, result.TotalPages)
}
