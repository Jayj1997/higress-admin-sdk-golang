// Package model provides data models for the SDK
package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLlmProvider_Validate(t *testing.T) {
	tests := []struct {
		name        string
		provider    *LlmProvider
		forUpdate   bool
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid provider with minimal fields",
			provider: &LlmProvider{
				Name: "test-provider",
				Type: LlmProviderTypeOpenai,
			},
			forUpdate:   false,
			expectError: false,
		},
		{
			name: "valid provider with all fields",
			provider: &LlmProvider{
				Name:     "test-provider",
				Type:     LlmProviderTypeQwen,
				Protocol: LlmProviderProtocolOpenaiV1,
				Tokens:   []string{"token1", "token2"},
				TokenFailoverConfig: &TokenFailoverConfig{
					Enabled:             true,
					FailureThreshold:    3,
					SuccessThreshold:    2,
					HealthCheckInterval: 30,
					HealthCheckTimeout:  10,
					HealthCheckModel:    "gpt-3.5-turbo",
				},
			},
			forUpdate:   false,
			expectError: false,
		},
		{
			name: "missing name",
			provider: &LlmProvider{
				Type: LlmProviderTypeOpenai,
			},
			forUpdate:   false,
			expectError: true,
			errorMsg:    "name cannot be blank",
		},
		{
			name: "name with slash",
			provider: &LlmProvider{
				Name: "test/provider",
				Type: LlmProviderTypeOpenai,
			},
			forUpdate:   false,
			expectError: true,
			errorMsg:    "slashes (/) are not allowed in name",
		},
		{
			name: "missing type",
			provider: &LlmProvider{
				Name: "test-provider",
			},
			forUpdate:   false,
			expectError: true,
			errorMsg:    "type cannot be blank",
		},
		{
			name: "empty protocol defaults to openai/v1",
			provider: &LlmProvider{
				Name:     "test-provider",
				Type:     LlmProviderTypeOpenai,
				Protocol: "",
			},
			forUpdate:   false,
			expectError: false,
		},
		{
			name: "invalid protocol",
			provider: &LlmProvider{
				Name:     "test-provider",
				Type:     LlmProviderTypeOpenai,
				Protocol: "invalid-protocol",
			},
			forUpdate:   false,
			expectError: true,
			errorMsg:    "Unknown protocol: invalid-protocol",
		},
		{
			name: "valid original protocol",
			provider: &LlmProvider{
				Name:     "test-provider",
				Type:     LlmProviderTypeOllama,
				Protocol: LlmProviderProtocolOriginal,
			},
			forUpdate:   false,
			expectError: false,
		},
		{
			name: "valid openai/v1 protocol",
			provider: &LlmProvider{
				Name:     "test-provider",
				Type:     LlmProviderTypeOpenai,
				Protocol: LlmProviderProtocolOpenaiV1,
			},
			forUpdate:   false,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.provider.Validate(tt.forUpdate)
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestLlmProvider_Validate_SetsDefaultProtocol(t *testing.T) {
	provider := &LlmProvider{
		Name: "test-provider",
		Type: LlmProviderTypeOpenai,
	}

	err := provider.Validate(false)
	require.NoError(t, err)
	assert.Equal(t, LlmProviderProtocolOpenaiV1, provider.Protocol)
}

func TestLlmProviderEndpoint_Validate(t *testing.T) {
	tests := []struct {
		name                string
		endpoint            *LlmProviderEndpoint
		expectError         bool
		errorMsg            string
		checkDefault        bool
		expectedProtocol    string
		expectedContextPath string
	}{
		{
			name: "valid endpoint with all fields",
			endpoint: &LlmProviderEndpoint{
				Protocol:    "https",
				Address:     "api.openai.com",
				Port:        443,
				ContextPath: "/v1",
			},
			expectError: false,
		},
		{
			name: "valid endpoint with minimal fields",
			endpoint: &LlmProviderEndpoint{
				Address: "api.openai.com",
				Port:    443,
			},
			expectError:         false,
			checkDefault:        true,
			expectedProtocol:    "https",
			expectedContextPath: "/",
		},
		{
			name: "missing address",
			endpoint: &LlmProviderEndpoint{
				Port: 443,
			},
			expectError: true,
			errorMsg:    "endpoint address cannot be empty",
		},
		{
			name: "missing port",
			endpoint: &LlmProviderEndpoint{
				Address: "api.openai.com",
			},
			expectError: true,
			errorMsg:    "endpoint port must be positive",
		},
		{
			name: "zero port",
			endpoint: &LlmProviderEndpoint{
				Address: "api.openai.com",
				Port:    0,
			},
			expectError: true,
			errorMsg:    "endpoint port must be positive",
		},
		{
			name: "negative port",
			endpoint: &LlmProviderEndpoint{
				Address: "api.openai.com",
				Port:    -1,
			},
			expectError: true,
			errorMsg:    "endpoint port must be positive",
		},
		{
			name: "http protocol",
			endpoint: &LlmProviderEndpoint{
				Protocol: "http",
				Address:  "localhost",
				Port:     8080,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.endpoint.Validate()
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
				if tt.checkDefault {
					assert.Equal(t, tt.expectedProtocol, tt.endpoint.Protocol)
					assert.Equal(t, tt.expectedContextPath, tt.endpoint.ContextPath)
				}
			}
		})
	}
}

func TestTokenFailoverConfig_Structure(t *testing.T) {
	config := &TokenFailoverConfig{
		Enabled:             true,
		FailureThreshold:    5,
		SuccessThreshold:    3,
		HealthCheckInterval: 60,
		HealthCheckTimeout:  15,
		HealthCheckModel:    "gpt-4",
	}

	assert.True(t, config.Enabled)
	assert.Equal(t, 5, config.FailureThreshold)
	assert.Equal(t, 3, config.SuccessThreshold)
	assert.Equal(t, 60, config.HealthCheckInterval)
	assert.Equal(t, 15, config.HealthCheckTimeout)
	assert.Equal(t, "gpt-4", config.HealthCheckModel)
}

func TestLlmProviderType_Constants(t *testing.T) {
	// Test that all provider type constants are defined
	providerTypes := []string{
		LlmProviderTypeQwen,
		LlmProviderTypeOpenai,
		LlmProviderTypeMoonshot,
		LlmProviderTypeAzure,
		LlmProviderTypeAi360,
		LlmProviderTypeGithub,
		LlmProviderTypeGroq,
		LlmProviderTypeBaichuan,
		LlmProviderTypeYi,
		LlmProviderTypeDeepSeek,
		LlmProviderTypeZhipuai,
		LlmProviderTypeOllama,
		LlmProviderTypeClaude,
		LlmProviderTypeBaidu,
		LlmProviderTypeHunyuan,
		LlmProviderTypeStepfun,
		LlmProviderTypeMinimax,
		LlmProviderTypeCloudflare,
		LlmProviderTypeSpark,
		LlmProviderTypeGemini,
		LlmProviderTypeDeepl,
		LlmProviderTypeMistral,
		LlmProviderTypeCohere,
		LlmProviderTypeDoubao,
		LlmProviderTypeCoze,
		LlmProviderTypeTogetherAi,
		LlmProviderTypeBedrock,
		LlmProviderTypeVertex,
		LlmProviderTypeOpenrouter,
		LlmProviderTypeGrok,
	}

	for _, pt := range providerTypes {
		assert.NotEmpty(t, pt)
	}
}

func TestLlmProviderProtocol_Constants(t *testing.T) {
	assert.Equal(t, "openai/v1", LlmProviderProtocolOpenaiV1)
	assert.Equal(t, "original", LlmProviderProtocolOriginal)
}

func TestIsValidLlmProviderProtocol(t *testing.T) {
	tests := []struct {
		protocol string
		expected bool
	}{
		{LlmProviderProtocolOpenaiV1, true},
		{LlmProviderProtocolOriginal, true},
		{"invalid", false},
		{"", false},
		{"openai", false},
		{"OPENAI/V1", false}, // case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.protocol, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsValidLlmProviderProtocol(tt.protocol))
		})
	}
}

func TestLlmProviderProtocolFromValue(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"empty string defaults to openai/v1", "", LlmProviderProtocolOpenaiV1},
		{"valid openai/v1", LlmProviderProtocolOpenaiV1, LlmProviderProtocolOpenaiV1},
		{"valid original", LlmProviderProtocolOriginal, LlmProviderProtocolOriginal},
		{"invalid protocol returns empty", "invalid", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LlmProviderProtocolFromValue(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLlmProvider_WithTokens(t *testing.T) {
	provider := &LlmProvider{
		Name:   "test-provider",
		Type:   LlmProviderTypeOpenai,
		Tokens: []string{"sk-test-1", "sk-test-2", "sk-test-3"},
	}

	err := provider.Validate(false)
	require.NoError(t, err)
	assert.Len(t, provider.Tokens, 3)
}

func TestLlmProvider_WithRawConfigs(t *testing.T) {
	provider := &LlmProvider{
		Name: "test-provider",
		Type: LlmProviderTypeOpenai,
		RawConfigs: map[string]interface{}{
			"customConfig": "value",
			"timeout":      30,
		},
	}

	err := provider.Validate(false)
	require.NoError(t, err)
	assert.Len(t, provider.RawConfigs, 2)
}
