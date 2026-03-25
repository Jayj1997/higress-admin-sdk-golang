// Package model provides data models for the SDK
package model

import (
	"testing"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model/route"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAiRoute_Validate(t *testing.T) {
	tests := []struct {
		name        string
		aiRoute     *AiRoute
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid AI route",
			aiRoute: &AiRoute{
				Name: "test-route",
				Upstreams: []AiUpstream{
					{Provider: "openai", Weight: 100},
				},
			},
			expectError: false,
		},
		{
			name: "valid AI route with path predicate",
			aiRoute: &AiRoute{
				Name: "test-route",
				PathPredicate: &route.RoutePredicate{
					MatchType: route.MatchTypePrefix,
					Path:      "/api/v1",
				},
				Upstreams: []AiUpstream{
					{Provider: "openai", Weight: 50},
					{Provider: "qwen", Weight: 50},
				},
			},
			expectError: false,
		},
		{
			name: "missing name",
			aiRoute: &AiRoute{
				Upstreams: []AiUpstream{
					{Provider: "openai", Weight: 100},
				},
			},
			expectError: true,
			errorMsg:    "name cannot be blank",
		},
		{
			name: "empty upstreams",
			aiRoute: &AiRoute{
				Name:      "test-route",
				Upstreams: []AiUpstream{},
			},
			expectError: true,
			errorMsg:    "upstreams cannot be empty",
		},
		{
			name: "nil upstreams",
			aiRoute: &AiRoute{
				Name:      "test-route",
				Upstreams: nil,
			},
			expectError: true,
			errorMsg:    "upstreams cannot be empty",
		},
		{
			name: "invalid path predicate - exact match",
			aiRoute: &AiRoute{
				Name: "test-route",
				PathPredicate: &route.RoutePredicate{
					MatchType: route.MatchTypeExact,
					Path:      "/api/v1",
				},
				Upstreams: []AiUpstream{
					{Provider: "openai", Weight: 100},
				},
			},
			expectError: true,
			errorMsg:    "pathPredicate must be of type prefix",
		},
		{
			name: "invalid path predicate - regex match",
			aiRoute: &AiRoute{
				Name: "test-route",
				PathPredicate: &route.RoutePredicate{
					MatchType: route.MatchTypeRegex,
					Path:      "/api/.*",
				},
				Upstreams: []AiUpstream{
					{Provider: "openai", Weight: 100},
				},
			},
			expectError: true,
			errorMsg:    "pathPredicate must be of type prefix",
		},
		{
			name: "invalid path predicate - missing path",
			aiRoute: &AiRoute{
				Name: "test-route",
				PathPredicate: &route.RoutePredicate{
					MatchType: route.MatchTypePrefix,
				},
				Upstreams: []AiUpstream{
					{Provider: "openai", Weight: 100},
				},
			},
			expectError: true,
			errorMsg:    "path is required",
		},
		{
			name: "header predicate with model routing header",
			aiRoute: &AiRoute{
				Name: "test-route",
				HeaderPredicates: []route.KeyedRoutePredicate{
					{Key: ModelRoutingHeader, MatchType: route.MatchTypeExact, Value: "test"},
				},
				Upstreams: []AiUpstream{
					{Provider: "openai", Weight: 100},
				},
			},
			expectError: true,
			errorMsg:    "headerPredicates cannot contain the model routing header",
		},
		{
			name: "header predicate with case insensitive model routing header",
			aiRoute: &AiRoute{
				Name: "test-route",
				HeaderPredicates: []route.KeyedRoutePredicate{
					{Key: "X-HIGRESS-MODEL-ROUTING", MatchType: route.MatchTypeExact, Value: "test"},
				},
				Upstreams: []AiUpstream{
					{Provider: "openai", Weight: 100},
				},
			},
			expectError: true,
			errorMsg:    "headerPredicates cannot contain the model routing header",
		},
		{
			name: "invalid header predicate - missing key",
			aiRoute: &AiRoute{
				Name: "test-route",
				HeaderPredicates: []route.KeyedRoutePredicate{
					{MatchType: route.MatchTypeExact, Value: "test"},
				},
				Upstreams: []AiUpstream{
					{Provider: "openai", Weight: 100},
				},
			},
			expectError: true,
			errorMsg:    "key is required",
		},
		{
			name: "invalid url param predicate - missing key",
			aiRoute: &AiRoute{
				Name: "test-route",
				UrlParamPredicates: []route.KeyedRoutePredicate{
					{MatchType: route.MatchTypeExact, Value: "test"},
				},
				Upstreams: []AiUpstream{
					{Provider: "openai", Weight: 100},
				},
			},
			expectError: true,
			errorMsg:    "key is required",
		},
		{
			name: "weight sum not 100",
			aiRoute: &AiRoute{
				Name: "test-route",
				Upstreams: []AiUpstream{
					{Provider: "openai", Weight: 30},
					{Provider: "qwen", Weight: 30},
				},
			},
			expectError: true,
			errorMsg:    "The sum of upstream weights must be 100",
		},
		{
			name: "weight sum exceeds 100",
			aiRoute: &AiRoute{
				Name: "test-route",
				Upstreams: []AiUpstream{
					{Provider: "openai", Weight: 60},
					{Provider: "qwen", Weight: 50},
				},
			},
			expectError: true,
			errorMsg:    "The sum of upstream weights must be 100",
		},
		{
			name: "upstream with empty provider",
			aiRoute: &AiRoute{
				Name: "test-route",
				Upstreams: []AiUpstream{
					{Provider: "", Weight: 100},
				},
			},
			expectError: true,
			errorMsg:    "provider cannot be null or empty",
		},
		{
			name: "valid with fallback config",
			aiRoute: &AiRoute{
				Name: "test-route",
				Upstreams: []AiUpstream{
					{Provider: "openai", Weight: 100},
				},
				FallbackConfig: &AiRouteFallbackConfig{
					Enabled: false,
				},
			},
			expectError: false,
		},
		{
			name: "fallback config enabled without upstreams",
			aiRoute: &AiRoute{
				Name: "test-route",
				Upstreams: []AiUpstream{
					{Provider: "openai", Weight: 100},
				},
				FallbackConfig: &AiRouteFallbackConfig{
					Enabled:   true,
					Upstreams: []AiUpstream{},
				},
			},
			expectError: true,
			errorMsg:    "upstreams cannot be empty when fallback is enabled",
		},
		{
			name: "valid fallback config with upstreams",
			aiRoute: &AiRoute{
				Name: "test-route",
				Upstreams: []AiUpstream{
					{Provider: "openai", Weight: 100},
				},
				FallbackConfig: &AiRouteFallbackConfig{
					Enabled: true,
					Upstreams: []AiUpstream{
						{Provider: "fallback-provider", Weight: 100},
					},
				},
			},
			expectError: false,
		},
		{
			name: "invalid model predicate - missing matchType",
			aiRoute: &AiRoute{
				Name: "test-route",
				Upstreams: []AiUpstream{
					{Provider: "openai", Weight: 100},
				},
				ModelPredicates: []AiModelPredicate{
					{MatchValue: "gpt-4"},
				},
			},
			expectError: true,
			errorMsg:    "matchType cannot be blank",
		},
		{
			name: "invalid model predicate - missing matchValue",
			aiRoute: &AiRoute{
				Name: "test-route",
				Upstreams: []AiUpstream{
					{Provider: "openai", Weight: 100},
				},
				ModelPredicates: []AiModelPredicate{
					{MatchType: "exact"},
				},
			},
			expectError: true,
			errorMsg:    "matchValue cannot be blank",
		},
		{
			name: "valid with all optional fields",
			aiRoute: &AiRoute{
				Name:    "test-route",
				Version: "v1",
				Domains: []string{"example.com", "api.example.com"},
				PathPredicate: &route.RoutePredicate{
					MatchType: route.MatchTypePrefix,
					Path:      "/api",
				},
				HeaderPredicates: []route.KeyedRoutePredicate{
					{Key: "X-Custom-Header", MatchType: route.MatchTypeExact, Value: "value"},
				},
				UrlParamPredicates: []route.KeyedRoutePredicate{
					{Key: "version", MatchType: route.MatchTypeExact, Value: "v1"},
				},
				Upstreams: []AiUpstream{
					{Provider: "openai", Weight: 70, ModelMapping: map[string]string{"gpt-4": "gpt-4-turbo"}},
					{Provider: "qwen", Weight: 30, ModelMapping: map[string]string{"qwen-turbo": "qwen-plus"}},
				},
				ModelPredicates: []AiModelPredicate{
					{MatchType: "exact", MatchValue: "gpt-4"},
				},
				CustomConfigs: map[string]string{"key1": "value1"},
				CustomLabels:  map[string]string{"label1": "value1"},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.aiRoute.Validate()
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAiUpstream_Validate(t *testing.T) {
	tests := []struct {
		name        string
		upstream    *AiUpstream
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid upstream",
			upstream: &AiUpstream{
				Provider: "openai",
				Weight:   100,
			},
			expectError: false,
		},
		{
			name: "valid upstream with model mapping",
			upstream: &AiUpstream{
				Provider:     "openai",
				Weight:       50,
				ModelMapping: map[string]string{"gpt-4": "gpt-4-turbo"},
			},
			expectError: false,
		},
		{
			name: "empty provider",
			upstream: &AiUpstream{
				Provider: "",
				Weight:   100,
			},
			expectError: true,
			errorMsg:    "provider cannot be null or empty",
		},
		{
			name: "zero weight is valid",
			upstream: &AiUpstream{
				Provider: "openai",
				Weight:   0,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.upstream.Validate()
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAiModelPredicate_Validate(t *testing.T) {
	tests := []struct {
		name        string
		predicate   *AiModelPredicate
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid predicate",
			predicate: &AiModelPredicate{
				MatchType:  "exact",
				MatchValue: "gpt-4",
			},
			expectError: false,
		},
		{
			name: "valid predicate with prefix match",
			predicate: &AiModelPredicate{
				MatchType:  "prefix",
				MatchValue: "gpt",
			},
			expectError: false,
		},
		{
			name: "missing matchType",
			predicate: &AiModelPredicate{
				MatchValue: "gpt-4",
			},
			expectError: true,
			errorMsg:    "matchType cannot be blank",
		},
		{
			name: "missing matchValue",
			predicate: &AiModelPredicate{
				MatchType: "exact",
			},
			expectError: true,
			errorMsg:    "matchValue cannot be blank",
		},
		{
			name:        "empty predicate",
			predicate:   &AiModelPredicate{},
			expectError: true,
			errorMsg:    "matchType cannot be blank",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.predicate.Validate()
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAiRouteFallbackConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      *AiRouteFallbackConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "disabled fallback is valid",
			config: &AiRouteFallbackConfig{
				Enabled: false,
			},
			expectError: false,
		},
		{
			name: "enabled fallback with upstreams",
			config: &AiRouteFallbackConfig{
				Enabled: true,
				Upstreams: []AiUpstream{
					{Provider: "fallback", Weight: 100},
				},
				FallbackStrategy: AiRouteFallbackStrategyRandom,
				ResponseCodes:    []string{"500", "503"},
			},
			expectError: false,
		},
		{
			name: "enabled fallback without upstreams",
			config: &AiRouteFallbackConfig{
				Enabled:   true,
				Upstreams: []AiUpstream{},
			},
			expectError: true,
			errorMsg:    "upstreams cannot be empty when fallback is enabled",
		},
		{
			name: "enabled fallback with nil upstreams",
			config: &AiRouteFallbackConfig{
				Enabled:   true,
				Upstreams: nil,
			},
			expectError: true,
			errorMsg:    "upstreams cannot be empty when fallback is enabled",
		},
		{
			name: "enabled fallback with invalid upstream",
			config: &AiRouteFallbackConfig{
				Enabled: true,
				Upstreams: []AiUpstream{
					{Provider: "", Weight: 100},
				},
			},
			expectError: true,
			errorMsg:    "provider cannot be null or empty",
		},
		{
			name: "sequence fallback strategy",
			config: &AiRouteFallbackConfig{
				Enabled: true,
				Upstreams: []AiUpstream{
					{Provider: "fallback1", Weight: 50},
					{Provider: "fallback2", Weight: 50},
				},
				FallbackStrategy: AiRouteFallbackStrategySequence,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAiRoute_Constants(t *testing.T) {
	// Test that constants are defined correctly
	assert.Equal(t, "x-higress-model-routing", ModelRoutingHeader)
	assert.Equal(t, "x-higress-fallback-from", FallbackFromHeader)
	assert.Equal(t, "RANDOM", AiRouteFallbackStrategyRandom)
	assert.Equal(t, "SEQUENCE", AiRouteFallbackStrategySequence)
}
