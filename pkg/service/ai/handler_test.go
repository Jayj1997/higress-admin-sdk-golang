// Package ai provides AI-related services for Higress Admin SDK.
package ai

import (
	"testing"

	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model/route"
	"github.com/stretchr/testify/assert"
)

// TestLlmProviderHandler_Interface tests that LlmProviderHandler interface is defined
func TestLlmProviderHandler_Interface(t *testing.T) {
	// This test verifies the interface exists
	var _ LlmProviderHandler = (LlmProviderHandler)(nil)
}

// TestHandlerRegistry tests the handler registry
func TestHandlerRegistry(t *testing.T) {
	// Create a mock handler for testing
	mockHandler := &MockLlmProviderHandler{
		providerType: "test-type",
	}

	// Register the handler
	RegisterHandler(mockHandler)

	// Verify the handler is registered
	assert.True(t, HasHandler("test-type"))
	assert.False(t, HasHandler("non-existent-type"))

	// Get the handler
	handler := GetHandler("test-type")
	assert.NotNil(t, handler)
	assert.Equal(t, "test-type", handler.GetType())

	// Get non-existent handler
	nilHandler := GetHandler("non-existent-type")
	assert.Nil(t, nilHandler)
}

// TestGetAllHandlers tests GetAllHandlers function
func TestGetAllHandlers(t *testing.T) {
	// Register a test handler
	testHandler := &MockLlmProviderHandler{
		providerType: "test-all-handlers",
	}
	RegisterHandler(testHandler)

	// Get all handlers
	allHandlers := GetAllHandlers()
	assert.NotNil(t, allHandlers)

	// Verify our test handler is in the map
	_, exists := allHandlers["test-all-handlers"]
	assert.True(t, exists)
}

// MockLlmProviderHandler is a mock implementation of LlmProviderHandler for testing
type MockLlmProviderHandler struct {
	providerType string
}

func (h *MockLlmProviderHandler) GetType() string {
	return h.providerType
}

func (h *MockLlmProviderHandler) CreateProvider() *model.LlmProvider {
	return &model.LlmProvider{
		Type: h.providerType,
	}
}

func (h *MockLlmProviderHandler) LoadConfig(provider *model.LlmProvider, configurations map[string]interface{}) bool {
	return true
}

func (h *MockLlmProviderHandler) SaveConfig(provider *model.LlmProvider, configurations map[string]interface{}) {
	// Mock implementation
}

func (h *MockLlmProviderHandler) NormalizeConfigs(configurations map[string]interface{}) {
	// Mock implementation
}

func (h *MockLlmProviderHandler) BuildServiceSource(providerName string, providerConfig map[string]interface{}) (*model.ServiceSource, error) {
	return &model.ServiceSource{
		Name: providerName,
		Type: "dns",
	}, nil
}

func (h *MockLlmProviderHandler) BuildUpstreamService(providerName string, providerConfig map[string]interface{}) (*route.UpstreamService, error) {
	return &route.UpstreamService{
		Name: providerName,
	}, nil
}

func (h *MockLlmProviderHandler) GetServiceSourceName(providerName string) string {
	return providerName + "-source"
}

func (h *MockLlmProviderHandler) GetExtraServiceSources(providerName string, providerConfig map[string]interface{}, forDelete bool) []model.ServiceSource {
	return nil
}

func (h *MockLlmProviderHandler) NeedSyncRouteAfterUpdate() bool {
	return false
}

func (h *MockLlmProviderHandler) GetProviderEndpoints(providerConfig map[string]interface{}) []model.LlmProviderEndpoint {
	return []model.LlmProviderEndpoint{
		{
			Protocol:    "https",
			Address:     "api.mock.com",
			Port:        443,
			ContextPath: "/",
		},
	}
}

// TestMockLlmProviderHandler tests the mock handler implementation
func TestMockLlmProviderHandler(t *testing.T) {
	handler := &MockLlmProviderHandler{providerType: "mock"}

	// Test GetType
	assert.Equal(t, "mock", handler.GetType())

	// Test CreateProvider
	provider := handler.CreateProvider()
	assert.NotNil(t, provider)
	assert.Equal(t, "mock", provider.Type)

	// Test LoadConfig
	config := map[string]interface{}{"key": "value"}
	loaded := handler.LoadConfig(provider, config)
	assert.True(t, loaded)

	// Test SaveConfig (should not panic)
	handler.SaveConfig(provider, config)

	// Test NormalizeConfigs (should not panic)
	handler.NormalizeConfigs(config)

	// Test BuildServiceSource
	source, err := handler.BuildServiceSource("test-provider", config)
	assert.NoError(t, err)
	assert.NotNil(t, source)
	assert.Equal(t, "test-provider", source.Name)

	// Test BuildUpstreamService
	upstream, err := handler.BuildUpstreamService("test-provider", config)
	assert.NoError(t, err)
	assert.NotNil(t, upstream)

	// Test GetServiceSourceName
	sourceName := handler.GetServiceSourceName("test-provider")
	assert.Equal(t, "test-provider-source", sourceName)

	// Test GetExtraServiceSources
	extraSources := handler.GetExtraServiceSources("test-provider", config, false)
	assert.Nil(t, extraSources)

	// Test NeedSyncRouteAfterUpdate
	needSync := handler.NeedSyncRouteAfterUpdate()
	assert.False(t, needSync)
}
