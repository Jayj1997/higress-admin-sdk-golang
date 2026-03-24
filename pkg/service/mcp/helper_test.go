// Package mcp provides MCP server related services
package mcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMcpServerHelper_New tests creating a new McpServerHelper
func TestMcpServerHelper_New(t *testing.T) {
	helper := NewMcpServerHelper()
	assert.NotNil(t, helper)
}

// TestMcpServerName2RouteName tests McpServerName2RouteName function
func TestMcpServerName2RouteName(t *testing.T) {
	helper := NewMcpServerHelper()

	tests := []struct {
		name          string
		mcpServerName string
		expected      string
	}{
		{
			name:          "Simple name",
			mcpServerName: "my-server",
			expected:      "mcp-server-my-server-internal",
		},
		{
			name:          "Already prefixed name",
			mcpServerName: "mcp-server-my-server-internal",
			expected:      "mcp-server-my-server-internal",
		},
		{
			name:          "Empty name",
			mcpServerName: "",
			expected:      "mcp-server--internal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helper.McpServerName2RouteName(tt.mcpServerName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestRouteName2McpServerName tests RouteName2McpServerName function
func TestRouteName2McpServerName(t *testing.T) {
	helper := NewMcpServerHelper()

	tests := []struct {
		name      string
		routeName string
		expected  string
	}{
		{
			name:      "Full prefixed route name",
			routeName: "mcp-server-my-server-internal",
			expected:  "my-server",
		},
		{
			name:      "Non-prefixed route name",
			routeName: "my-route",
			expected:  "my-route",
		},
		{
			name:      "Empty route name",
			routeName: "",
			expected:  "",
		},
		{
			name:      "Only prefix without suffix",
			routeName: "mcp-server-my-server",
			expected:  "my-server",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helper.RouteName2McpServerName(tt.routeName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestMcpServerConstants tests the constants
func TestMcpServerConstants(t *testing.T) {
	assert.Equal(t, "mcp-server-", McpServerRoutePrefix)
	assert.Equal(t, "-internal", InternalResourceNameSuffix)
}

// TestMcpServerNameRoundTrip tests round trip conversion
func TestMcpServerNameRoundTrip(t *testing.T) {
	helper := NewMcpServerHelper()

	originalNames := []string{
		"my-server",
		"test-mcp-server",
		"server123",
	}

	for _, original := range originalNames {
		t.Run("Round trip for "+original, func(t *testing.T) {
			routeName := helper.McpServerName2RouteName(original)
			result := helper.RouteName2McpServerName(routeName)
			assert.Equal(t, original, result)
		})
	}
}
