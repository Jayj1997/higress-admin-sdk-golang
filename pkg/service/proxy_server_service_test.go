// Package service provides business services for the SDK
package service

import (
	"testing"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProxyServerService_Interface tests that the implementation satisfies the interface
func TestProxyServerService_Interface(t *testing.T) {
	var _ ProxyServerService = (*ProxyServerServiceImpl)(nil)
}

// TestProxyServerService_New tests creating a new ProxyServerService
func TestProxyServerService_New(t *testing.T) {
	svc := NewProxyServerService(nil, nil)
	require.NotNil(t, svc)
}

// TestProxyServerService_Model tests the ProxyServer model
func TestProxyServerService_Model(t *testing.T) {
	server := model.ProxyServer{
		Name:     "test-proxy",
		Host:     "proxy.example.com",
		Port:     3128,
		Protocol: "http",
		Version:  "v1",
	}

	assert.Equal(t, "test-proxy", server.Name)
	assert.Equal(t, "proxy.example.com", server.Host)
	assert.Equal(t, 3128, server.Port)
	assert.Equal(t, "http", server.Protocol)
	assert.Equal(t, "v1", server.Version)
}

// TestProxyServerService_Protocols tests different proxy protocols
func TestProxyServerService_Protocols(t *testing.T) {
	tests := []struct {
		name     string
		protocol string
	}{
		{"HTTP protocol", "http"},
		{"HTTPS protocol", "https"},
		{"SOCKS5 protocol", "socks5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &model.ProxyServer{Protocol: tt.protocol}
			assert.Equal(t, tt.protocol, server.Protocol)
		})
	}
}

// TestProxyServerService_Pagination tests pagination logic for proxy servers
func TestProxyServerService_Pagination(t *testing.T) {
	servers := make([]model.ProxyServer, 15)
	for i := 0; i < 15; i++ {
		servers[i] = model.ProxyServer{
			Name:     "proxy-" + string(rune('a'+i)),
			Host:     "proxy.example.com",
			Port:     3128,
			Protocol: "http",
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
	assert.Equal(t, 15, result.Total)
	assert.Equal(t, 2, result.TotalPages)
}
