// Package service provides business services for the SDK
package service

import (
	"testing"

	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
	"github.com/stretchr/testify/assert"
)

// TestLlmProviderService_New tests creating a new LlmProviderService
func TestLlmProviderService_New(t *testing.T) {
	svc := NewLlmProviderService(nil, nil)
	assert.NotNil(t, svc)
}

// TestLlmProviderService_SetAiRouteService tests setting the AI route service
func TestLlmProviderService_SetAiRouteService(t *testing.T) {
	svc := NewLlmProviderService(nil, nil)
	// Setting nil should not panic
	svc.SetAiRouteService(nil)
}

// TestLlmProviderModel tests the LlmProvider model
func TestLlmProviderModel(t *testing.T) {
	provider := model.LlmProvider{
		Name:     "test-provider",
		Type:     "openai",
		Protocol: "openai/v1",
	}

	assert.Equal(t, "test-provider", provider.Name)
	assert.Equal(t, "openai", provider.Type)
	assert.Equal(t, "openai/v1", provider.Protocol)
}

// TestLlmProviderWithTokens tests LlmProvider with tokens
func TestLlmProviderWithTokens(t *testing.T) {
	provider := model.LlmProvider{
		Name:     "azure-provider",
		Type:     "azure",
		Protocol: "openai/v1",
		Tokens:   []string{"token1", "token2"},
	}

	assert.Equal(t, "azure-provider", provider.Name)
	assert.Equal(t, "azure", provider.Type)
	assert.Len(t, provider.Tokens, 2)
}

// TestLlmProviderWithTokenFailover tests LlmProvider with token failover config
func TestLlmProviderWithTokenFailover(t *testing.T) {
	failoverConfig := &model.TokenFailoverConfig{
		Enabled:             true,
		FailureThreshold:    3,
		SuccessThreshold:    2,
		HealthCheckInterval: 60,
		HealthCheckTimeout:  10,
		HealthCheckModel:    "gpt-3.5-turbo",
	}

	provider := model.LlmProvider{
		Name:                "test-provider",
		Type:                "openai",
		Protocol:            "openai/v1",
		TokenFailoverConfig: failoverConfig,
	}

	assert.True(t, provider.TokenFailoverConfig.Enabled)
	assert.Equal(t, 3, provider.TokenFailoverConfig.FailureThreshold)
	assert.Equal(t, 2, provider.TokenFailoverConfig.SuccessThreshold)
}

// TestLlmProviderTypes tests different LLM provider types
func TestLlmProviderTypes(t *testing.T) {
	tests := []struct {
		name         string
		providerType string
	}{
		{"OpenAI provider", model.LlmProviderTypeOpenai},
		{"Azure provider", model.LlmProviderTypeAzure},
		{"Claude provider", model.LlmProviderTypeClaude},
		{"Qwen provider", model.LlmProviderTypeQwen},
		{"Moonshot provider", model.LlmProviderTypeMoonshot},
		{"Yi provider", model.LlmProviderTypeYi},
		{"Baichuan provider", model.LlmProviderTypeBaichuan},
		{"DeepSeek provider", model.LlmProviderTypeDeepSeek},
		{"Zhipuai provider", model.LlmProviderTypeZhipuai},
		{"Ollama provider", model.LlmProviderTypeOllama},
		{"Baidu provider", model.LlmProviderTypeBaidu},
		{"Hunyuan provider", model.LlmProviderTypeHunyuan},
		{"Stepfun provider", model.LlmProviderTypeStepfun},
		{"Minimax provider", model.LlmProviderTypeMinimax},
		{"Gemini provider", model.LlmProviderTypeGemini},
		{"Mistral provider", model.LlmProviderTypeMistral},
		{"Cohere provider", model.LlmProviderTypeCohere},
		{"Doubao provider", model.LlmProviderTypeDoubao},
		{"Coze provider", model.LlmProviderTypeCoze},
		{"Bedrock provider", model.LlmProviderTypeBedrock},
		{"Vertex provider", model.LlmProviderTypeVertex},
		{"Grok provider", model.LlmProviderTypeGrok},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &model.LlmProvider{Type: tt.providerType}
			assert.Equal(t, tt.providerType, provider.Type)
		})
	}
}

// TestLlmProviderProtocols tests different LLM provider protocols
func TestLlmProviderProtocols(t *testing.T) {
	tests := []struct {
		name     string
		protocol string
		valid    bool
	}{
		{"OpenAI v1 protocol", model.LlmProviderProtocolOpenaiV1, true},
		{"Original protocol", model.LlmProviderProtocolOriginal, true},
		{"Invalid protocol", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.valid, model.IsValidLlmProviderProtocol(tt.protocol))
		})
	}
}

// TestLlmProviderProtocolFromValue tests LlmProviderProtocolFromValue function
func TestLlmProviderProtocolFromValue(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"Empty value defaults to openai/v1", "", model.LlmProviderProtocolOpenaiV1},
		{"Valid openai/v1 protocol", model.LlmProviderProtocolOpenaiV1, model.LlmProviderProtocolOpenaiV1},
		{"Valid original protocol", model.LlmProviderProtocolOriginal, model.LlmProviderProtocolOriginal},
		{"Invalid protocol returns empty", "invalid", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := model.LlmProviderProtocolFromValue(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestLlmProviderValidate tests LlmProvider validation
func TestLlmProviderValidate(t *testing.T) {
	tests := []struct {
		name     string
		provider model.LlmProvider
		wantErr  bool
	}{
		{
			name: "Valid provider",
			provider: model.LlmProvider{
				Name:     "test-provider",
				Type:     "openai",
				Protocol: "openai/v1",
			},
			wantErr: false,
		},
		{
			name: "Empty name",
			provider: model.LlmProvider{
				Name:     "",
				Type:     "openai",
				Protocol: "openai/v1",
			},
			wantErr: true,
		},
		{
			name: "Name with slashes",
			provider: model.LlmProvider{
				Name:     "test/provider",
				Type:     "openai",
				Protocol: "openai/v1",
			},
			wantErr: true,
		},
		{
			name: "Empty type",
			provider: model.LlmProvider{
				Name:     "test-provider",
				Type:     "",
				Protocol: "openai/v1",
			},
			wantErr: true,
		},
		{
			name: "Empty protocol defaults to openai/v1",
			provider: model.LlmProvider{
				Name: "test-provider",
				Type: "openai",
			},
			wantErr: false,
		},
		{
			name: "Invalid protocol",
			provider: model.LlmProvider{
				Name:     "test-provider",
				Type:     "openai",
				Protocol: "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.provider.Validate(false)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestLlmProviderEndpoint tests LlmProviderEndpoint model
func TestLlmProviderEndpoint(t *testing.T) {
	endpoint := model.LlmProviderEndpoint{
		Protocol:    "https",
		Address:     "api.openai.com",
		Port:        443,
		ContextPath: "/v1",
	}

	assert.Equal(t, "https", endpoint.Protocol)
	assert.Equal(t, "api.openai.com", endpoint.Address)
	assert.Equal(t, 443, endpoint.Port)
	assert.Equal(t, "/v1", endpoint.ContextPath)
}

// TestLlmProviderEndpointValidate tests LlmProviderEndpoint validation
func TestLlmProviderEndpointValidate(t *testing.T) {
	tests := []struct {
		name     string
		endpoint model.LlmProviderEndpoint
		wantErr  bool
	}{
		{
			name: "Valid endpoint",
			endpoint: model.LlmProviderEndpoint{
				Address: "api.openai.com",
				Port:    443,
			},
			wantErr: false,
		},
		{
			name: "Empty address",
			endpoint: model.LlmProviderEndpoint{
				Address: "",
				Port:    443,
			},
			wantErr: true,
		},
		{
			name: "Invalid port",
			endpoint: model.LlmProviderEndpoint{
				Address: "api.openai.com",
				Port:    0,
			},
			wantErr: true,
		},
		{
			name: "Negative port",
			endpoint: model.LlmProviderEndpoint{
				Address: "api.openai.com",
				Port:    -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.endpoint.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestLlmProviderEndpointDefaults tests default values for endpoint
func TestLlmProviderEndpointDefaults(t *testing.T) {
	endpoint := model.LlmProviderEndpoint{
		Address: "api.openai.com",
		Port:    443,
	}

	err := endpoint.Validate()
	assert.NoError(t, err)
	assert.Equal(t, "https", endpoint.Protocol) // default
	assert.Equal(t, "/", endpoint.ContextPath)  // default
}

// TestTokenFailoverConfig tests TokenFailoverConfig model
func TestTokenFailoverConfig(t *testing.T) {
	config := model.TokenFailoverConfig{
		Enabled:             true,
		FailureThreshold:    5,
		SuccessThreshold:    3,
		HealthCheckInterval: 120,
		HealthCheckTimeout:  30,
		HealthCheckModel:    "gpt-4",
	}

	assert.True(t, config.Enabled)
	assert.Equal(t, 5, config.FailureThreshold)
	assert.Equal(t, 3, config.SuccessThreshold)
	assert.Equal(t, 120, config.HealthCheckInterval)
	assert.Equal(t, 30, config.HealthCheckTimeout)
	assert.Equal(t, "gpt-4", config.HealthCheckModel)
}
