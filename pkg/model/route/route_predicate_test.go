// Package route provides route-related models for the SDK
package route

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoutePredicate_Validate(t *testing.T) {
	tests := []struct {
		name        string
		predicate   *RoutePredicate
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid exact match",
			predicate: &RoutePredicate{
				MatchType: MatchTypeExact,
				Path:      "/api/v1",
			},
			expectError: false,
		},
		{
			name: "valid prefix match",
			predicate: &RoutePredicate{
				MatchType: MatchTypePrefix,
				Path:      "/api",
			},
			expectError: false,
		},
		{
			name: "valid regex match",
			predicate: &RoutePredicate{
				MatchType: MatchTypeRegex,
				Path:      "/api/.*",
			},
			expectError: false,
		},
		{
			name: "missing match type",
			predicate: &RoutePredicate{
				Path: "/api",
			},
			expectError: true,
			errorMsg:    "matchType is required",
		},
		{
			name: "missing path",
			predicate: &RoutePredicate{
				MatchType: MatchTypePrefix,
			},
			expectError: true,
			errorMsg:    "path is required",
		},
		{
			name: "empty match type",
			predicate: &RoutePredicate{
				MatchType: "",
				Path:      "/api",
			},
			expectError: true,
			errorMsg:    "matchType is required",
		},
		{
			name: "empty path",
			predicate: &RoutePredicate{
				MatchType: MatchTypePrefix,
				Path:      "",
			},
			expectError: true,
			errorMsg:    "path is required",
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

func TestRoutePredicate_Structure(t *testing.T) {
	caseSensitive := true
	predicate := &RoutePredicate{
		MatchType:     MatchTypePrefix,
		Path:          "/api/v1",
		CaseSensitive: &caseSensitive,
	}

	assert.Equal(t, MatchTypePrefix, predicate.MatchType)
	assert.Equal(t, "/api/v1", predicate.Path)
	require.NotNil(t, predicate.CaseSensitive)
	assert.True(t, *predicate.CaseSensitive)
}

func TestRoutePredicate_CaseSensitive(t *testing.T) {
	t.Run("case sensitive true", func(t *testing.T) {
		caseSensitive := true
		predicate := &RoutePredicate{
			MatchType:     MatchTypeExact,
			Path:          "/API",
			CaseSensitive: &caseSensitive,
		}
		require.NotNil(t, predicate.CaseSensitive)
		assert.True(t, *predicate.CaseSensitive)
	})

	t.Run("case sensitive false", func(t *testing.T) {
		caseSensitive := false
		predicate := &RoutePredicate{
			MatchType:     MatchTypeExact,
			Path:          "/api",
			CaseSensitive: &caseSensitive,
		}
		require.NotNil(t, predicate.CaseSensitive)
		assert.False(t, *predicate.CaseSensitive)
	})

	t.Run("case sensitive nil", func(t *testing.T) {
		predicate := &RoutePredicate{
			MatchType: MatchTypeExact,
			Path:      "/api",
		}
		assert.Nil(t, predicate.CaseSensitive)
	})
}

func TestKeyedRoutePredicate_Validate(t *testing.T) {
	tests := []struct {
		name        string
		predicate   *KeyedRoutePredicate
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid predicate",
			predicate: &KeyedRoutePredicate{
				Key:       "X-Custom-Header",
				MatchType: MatchTypeExact,
				Value:     "test-value",
			},
			expectError: false,
		},
		{
			name: "missing key",
			predicate: &KeyedRoutePredicate{
				MatchType: MatchTypeExact,
				Value:     "test-value",
			},
			expectError: true,
			errorMsg:    "key is required",
		},
		{
			name: "missing match type",
			predicate: &KeyedRoutePredicate{
				Key:   "X-Custom-Header",
				Value: "test-value",
			},
			expectError: true,
			errorMsg:    "matchType is required",
		},
		{
			name: "empty key",
			predicate: &KeyedRoutePredicate{
				Key:       "",
				MatchType: MatchTypeExact,
				Value:     "test-value",
			},
			expectError: true,
			errorMsg:    "key is required",
		},
		{
			name: "empty match type",
			predicate: &KeyedRoutePredicate{
				Key:       "X-Custom-Header",
				MatchType: "",
				Value:     "test-value",
			},
			expectError: true,
			errorMsg:    "matchType is required",
		},
		{
			name: "value is optional",
			predicate: &KeyedRoutePredicate{
				Key:       "X-Custom-Header",
				MatchType: MatchTypePrefix,
				Value:     "",
			},
			expectError: false,
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

func TestKeyedRoutePredicate_Structure(t *testing.T) {
	caseSensitive := true
	predicate := &KeyedRoutePredicate{
		Key:           "X-API-Key",
		MatchType:     MatchTypeExact,
		Value:         "secret-key",
		CaseSensitive: &caseSensitive,
	}

	assert.Equal(t, "X-API-Key", predicate.Key)
	assert.Equal(t, MatchTypeExact, predicate.MatchType)
	assert.Equal(t, "secret-key", predicate.Value)
	require.NotNil(t, predicate.CaseSensitive)
	assert.True(t, *predicate.CaseSensitive)
}

func TestUpstreamService_Validate(t *testing.T) {
	tests := []struct {
		name        string
		service     *UpstreamService
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid service",
			service: &UpstreamService{
				Name:      "my-service",
				Namespace: "default",
				Port:      8080,
			},
			expectError: false,
		},
		{
			name: "valid service with weight",
			service: &UpstreamService{
				Name:      "my-service",
				Namespace: "default",
				Port:      8080,
				Weight:    intPtrUS(50),
			},
			expectError: false,
		},
		{
			name: "missing name",
			service: &UpstreamService{
				Namespace: "default",
				Port:      8080,
			},
			expectError: true,
			errorMsg:    "service name is required",
		},
		{
			name: "empty name",
			service: &UpstreamService{
				Name:      "",
				Namespace: "default",
				Port:      8080,
			},
			expectError: true,
			errorMsg:    "service name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.service.Validate()
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUpstreamService_Structure(t *testing.T) {
	weight := 80
	service := &UpstreamService{
		Name:      "backend-service",
		Namespace: "production",
		Port:      8080,
		Weight:    &weight,
		Version:   "v1",
	}

	assert.Equal(t, "backend-service", service.Name)
	assert.Equal(t, "production", service.Namespace)
	assert.Equal(t, 8080, service.Port)
	require.NotNil(t, service.Weight)
	assert.Equal(t, 80, *service.Weight)
	assert.Equal(t, "v1", service.Version)
}

func TestMatchType_Constants(t *testing.T) {
	assert.Equal(t, "exact", MatchTypeExact)
	assert.Equal(t, "prefix", MatchTypePrefix)
	assert.Equal(t, "regex", MatchTypeRegex)
}

func TestRewriteConfig_Structure(t *testing.T) {
	config := &RewriteConfig{
		Path: "/new-path/$1",
		Host: "new-host.example.com",
	}

	assert.Equal(t, "/new-path/$1", config.Path)
	assert.Equal(t, "new-host.example.com", config.Host)
}

func TestRewriteConfig_Empty(t *testing.T) {
	config := &RewriteConfig{}

	assert.Empty(t, config.Path)
	assert.Empty(t, config.Host)
}

// Helper function
func intPtrUS(i int) *int {
	return &i
}
