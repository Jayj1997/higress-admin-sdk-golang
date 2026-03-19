package model

import (
	"encoding/json"
	"testing"

	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model/route"
	"github.com/stretchr/testify/assert"
)

// TestRouteJSONSerialization 测试 Route 模型的 JSON 序列化
// 运行命令: go test -v -run TestRouteJSONSerialization ./pkg/model/
func TestRouteJSONSerialization(t *testing.T) {
	tests := []struct {
		name     string
		route    Route
		expected string
	}{
		{
			name: "basic route",
			route: Route{
				Name: "my-route",
				Path: &route.RoutePredicate{
					MatchType: route.MatchTypePrefix,
					Path:      "/api",
				},
				Services: []*route.UpstreamService{
					{Name: "my-service", Port: 8080},
				},
			},
			expected: `{"name":"my-route","path":{"matchType":"prefix","path":"/api"},"services":[{"name":"my-service","port":8080}]}`,
		},
		{
			name: "route with domains",
			route: Route{
				Name:    "domain-route",
				Domains: []string{"example.com", "api.example.com"},
				Path: &route.RoutePredicate{
					MatchType: route.MatchTypeExact,
					Path:      "/health",
				},
				Services: []*route.UpstreamService{
					{Name: "health-service", Port: 80},
				},
			},
			expected: `{"name":"domain-route","domains":["example.com","api.example.com"],"path":{"matchType":"exact","path":"/health"},"services":[{"name":"health-service","port":80}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.route)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(data))
		})
	}
}

// TestRouteValidate 测试 Route 模型验证
// 运行命令: go test -v -run TestRouteValidate ./pkg/model/
func TestRouteValidate(t *testing.T) {
	tests := []struct {
		name    string
		route   Route
		wantErr bool
		errMsg  string
	}{
		{
			name:    "empty name should fail",
			route:   Route{Services: []*route.UpstreamService{{Name: "service"}}},
			wantErr: true,
			errMsg:  "name cannot be blank",
		},
		{
			name:    "empty services should fail",
			route:   Route{Name: "test-route"},
			wantErr: true,
			errMsg:  "services cannot be empty",
		},
		{
			name: "valid route",
			route: Route{
				Name: "test-route",
				Services: []*route.UpstreamService{
					{Name: "my-service", Port: 8080},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.route.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
