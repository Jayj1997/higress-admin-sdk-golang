// Package consumer provides consumer management services
package consumer

import (
	"testing"

	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
	"github.com/stretchr/testify/assert"
)

// TestKeyAuthCredentialHandler_GetType tests the GetType method
func TestKeyAuthCredentialHandler_GetType(t *testing.T) {
	handler := NewKeyAuthCredentialHandler()
	assert.Equal(t, model.CredentialTypeKeyAuth, handler.GetType())
}

// TestKeyAuthCredentialHandler_GetPluginName tests the GetPluginName method
func TestKeyAuthCredentialHandler_GetPluginName(t *testing.T) {
	handler := NewKeyAuthCredentialHandler()
	assert.Equal(t, PluginNameKeyAuth, handler.GetPluginName())
}

// TestKeyAuthCredentialHandler_IsConsumerInUse tests the IsConsumerInUse method
func TestKeyAuthCredentialHandler_IsConsumerInUse(t *testing.T) {
	handler := NewKeyAuthCredentialHandler()

	tests := []struct {
		name          string
		consumerName  string
		instances     []*model.WasmPluginInstance
		expectedInUse bool
	}{
		{
			name:          "empty instances list",
			consumerName:  "test-consumer",
			instances:     []*model.WasmPluginInstance{},
			expectedInUse: false,
		},
		{
			name:          "nil instances list",
			consumerName:  "test-consumer",
			instances:     nil,
			expectedInUse: false,
		},
		{
			name:         "instance with nil configurations",
			consumerName: "test-consumer",
			instances: []*model.WasmPluginInstance{
				{
					PluginName:     "key-auth",
					Configurations: nil,
				},
			},
			expectedInUse: false,
		},
		{
			name:         "instance without allow config",
			consumerName: "test-consumer",
			instances: []*model.WasmPluginInstance{
				{
					PluginName: "key-auth",
					Configurations: map[string]interface{}{
						"consumers": []interface{}{},
					},
				},
			},
			expectedInUse: false,
		},
		{
			name:         "consumer in allow list",
			consumerName: "test-consumer",
			instances: []*model.WasmPluginInstance{
				{
					PluginName: "key-auth",
					Configurations: map[string]interface{}{
						"allow": []interface{}{"test-consumer", "other-consumer"},
					},
				},
			},
			expectedInUse: true,
		},
		{
			name:         "consumer not in allow list",
			consumerName: "test-consumer",
			instances: []*model.WasmPluginInstance{
				{
					PluginName: "key-auth",
					Configurations: map[string]interface{}{
						"allow": []interface{}{"other-consumer", "another-consumer"},
					},
				},
			},
			expectedInUse: false,
		},
		{
			name:         "allow config is not a slice",
			consumerName: "test-consumer",
			instances: []*model.WasmPluginInstance{
				{
					PluginName: "key-auth",
					Configurations: map[string]interface{}{
						"allow": "not-a-slice",
					},
				},
			},
			expectedInUse: false,
		},
		{
			name:         "multiple instances with consumer in second",
			consumerName: "test-consumer",
			instances: []*model.WasmPluginInstance{
				{
					PluginName: "key-auth",
					Configurations: map[string]interface{}{
						"allow": []interface{}{"other-consumer"},
					},
				},
				{
					PluginName: "key-auth",
					Configurations: map[string]interface{}{
						"allow": []interface{}{"test-consumer"},
					},
				},
			},
			expectedInUse: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.IsConsumerInUse(tt.consumerName, tt.instances)
			assert.Equal(t, tt.expectedInUse, result)
		})
	}
}

// TestKeyAuthCredentialHandler_ExtractConsumers tests the ExtractConsumers method
func TestKeyAuthCredentialHandler_ExtractConsumers(t *testing.T) {
	handler := NewKeyAuthCredentialHandler()

	tests := []struct {
		name          string
		instance      *model.WasmPluginInstance
		expectedCount int
		expectedNames []string
	}{
		{
			name:          "nil instance",
			instance:      nil,
			expectedCount: 0,
			expectedNames: nil,
		},
		{
			name:          "instance with nil configurations",
			instance:      &model.WasmPluginInstance{Configurations: nil},
			expectedCount: 0,
			expectedNames: nil,
		},
		{
			name:          "instance without consumers config",
			instance:      &model.WasmPluginInstance{Configurations: map[string]interface{}{}},
			expectedCount: 0,
			expectedNames: nil,
		},
		{
			name: "consumers config is not a slice",
			instance: &model.WasmPluginInstance{
				Configurations: map[string]interface{}{
					"consumers": "not-a-slice",
				},
			},
			expectedCount: 0,
			expectedNames: nil,
		},
		{
			name: "empty consumers list",
			instance: &model.WasmPluginInstance{
				Configurations: map[string]interface{}{
					"consumers": []interface{}{},
				},
			},
			expectedCount: 0,
			expectedNames: nil,
		},
		{
			name: "consumer item is not a map",
			instance: &model.WasmPluginInstance{
				Configurations: map[string]interface{}{
					"consumers": []interface{}{"not-a-map"},
				},
			},
			expectedCount: 0,
			expectedNames: nil,
		},
		{
			name: "consumer without name",
			instance: &model.WasmPluginInstance{
				Configurations: map[string]interface{}{
					"consumers": []interface{}{
						map[string]interface{}{
							"keys": []interface{}{"x-api-key"},
						},
					},
				},
			},
			expectedCount: 0,
			expectedNames: nil,
		},
		{
			name: "consumer with valid bearer credential",
			instance: &model.WasmPluginInstance{
				Configurations: map[string]interface{}{
					"consumers": []interface{}{
						map[string]interface{}{
							"name":      "test-consumer",
							"keys":      []interface{}{"Authorization"},
							"in_header": true,
							"credentials": []interface{}{
								"Bearer token123",
							},
						},
					},
				},
			},
			expectedCount: 1,
			expectedNames: []string{"test-consumer"},
		},
		{
			name: "consumer with valid header credential",
			instance: &model.WasmPluginInstance{
				Configurations: map[string]interface{}{
					"consumers": []interface{}{
						map[string]interface{}{
							"name":      "header-consumer",
							"keys":      []interface{}{"x-api-key"},
							"in_header": true,
							"credentials": []interface{}{
								"secret-key",
							},
						},
					},
				},
			},
			expectedCount: 1,
			expectedNames: []string{"header-consumer"},
		},
		{
			name: "consumer with valid query credential",
			instance: &model.WasmPluginInstance{
				Configurations: map[string]interface{}{
					"consumers": []interface{}{
						map[string]interface{}{
							"name":     "query-consumer",
							"keys":     []interface{}{"api_key"},
							"in_query": true,
							"credentials": []interface{}{
								"query-token",
							},
						},
					},
				},
			},
			expectedCount: 1,
			expectedNames: []string{"query-consumer"},
		},
		{
			name: "multiple consumers",
			instance: &model.WasmPluginInstance{
				Configurations: map[string]interface{}{
					"consumers": []interface{}{
						map[string]interface{}{
							"name":      "consumer1",
							"keys":      []interface{}{"x-api-key"},
							"in_header": true,
							"credentials": []interface{}{
								"key1",
							},
						},
						map[string]interface{}{
							"name":     "consumer2",
							"keys":     []interface{}{"api_key"},
							"in_query": true,
							"credentials": []interface{}{
								"key2",
							},
						},
					},
				},
			},
			expectedCount: 2,
			expectedNames: []string{"consumer1", "consumer2"},
		},
		{
			name: "consumer without in_header or in_query",
			instance: &model.WasmPluginInstance{
				Configurations: map[string]interface{}{
					"consumers": []interface{}{
						map[string]interface{}{
							"name":        "invalid-consumer",
							"keys":        []interface{}{"x-api-key"},
							"credentials": []interface{}{"key"},
						},
					},
				},
			},
			expectedCount: 0,
			expectedNames: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			consumers := handler.ExtractConsumers(tt.instance)
			assert.Len(t, consumers, tt.expectedCount)
			if tt.expectedNames != nil {
				names := make([]string, len(consumers))
				for i, c := range consumers {
					names[i] = c.Name
				}
				assert.ElementsMatch(t, tt.expectedNames, names)
			}
		})
	}
}

// TestKeyAuthCredentialHandler_InitDefaultGlobalConfigs tests the InitDefaultGlobalConfigs method
func TestKeyAuthCredentialHandler_InitDefaultGlobalConfigs(t *testing.T) {
	handler := NewKeyAuthCredentialHandler()

	tests := []struct {
		name               string
		instance           *model.WasmPluginInstance
		expectedConfigKeys []string
	}{
		{
			name: "nil configurations",
			instance: &model.WasmPluginInstance{
				Configurations: nil,
			},
			expectedConfigKeys: []string{"global_auth", "allow", "keys", "consumers"},
		},
		{
			name: "empty configurations",
			instance: &model.WasmPluginInstance{
				Configurations: map[string]interface{}{},
			},
			expectedConfigKeys: []string{"global_auth", "allow", "keys", "consumers"},
		},
		{
			name: "existing configurations preserved",
			instance: &model.WasmPluginInstance{
				Configurations: map[string]interface{}{
					"global_auth": true,
					"allow":       []string{"existing-consumer"},
				},
			},
			expectedConfigKeys: []string{"global_auth", "allow", "keys", "consumers"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler.InitDefaultGlobalConfigs(tt.instance)
			assert.NotNil(t, tt.instance.Configurations)
			for _, key := range tt.expectedConfigKeys {
				_, exists := tt.instance.Configurations[key]
				assert.True(t, exists, "Expected key %s to exist", key)
			}
		})
	}
}

// TestKeyAuthCredentialHandler_GetAllowedConsumers tests the GetAllowedConsumers method
func TestKeyAuthCredentialHandler_GetAllowedConsumers(t *testing.T) {
	handler := NewKeyAuthCredentialHandler()

	tests := []struct {
		name           string
		instance       *model.WasmPluginInstance
		expectedResult []string
	}{
		{
			name:           "nil instance",
			instance:       nil,
			expectedResult: []string{},
		},
		{
			name:           "nil configurations",
			instance:       &model.WasmPluginInstance{Configurations: nil},
			expectedResult: []string{},
		},
		{
			name:           "no allow config",
			instance:       &model.WasmPluginInstance{Configurations: map[string]interface{}{}},
			expectedResult: []string{},
		},
		{
			name: "allow config is not a slice",
			instance: &model.WasmPluginInstance{
				Configurations: map[string]interface{}{
					"allow": "not-a-slice",
				},
			},
			expectedResult: []string{},
		},
		{
			name: "empty allow list",
			instance: &model.WasmPluginInstance{
				Configurations: map[string]interface{}{
					"allow": []interface{}{},
				},
			},
			expectedResult: []string{},
		},
		{
			name: "allow list with consumers",
			instance: &model.WasmPluginInstance{
				Configurations: map[string]interface{}{
					"allow": []interface{}{"consumer1", "consumer2", "consumer3"},
				},
			},
			expectedResult: []string{"consumer1", "consumer2", "consumer3"},
		},
		{
			name: "allow list with non-string items",
			instance: &model.WasmPluginInstance{
				Configurations: map[string]interface{}{
					"allow": []interface{}{"consumer1", 123, true, "consumer2"},
				},
			},
			expectedResult: []string{"consumer1", "consumer2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.GetAllowedConsumers(tt.instance)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

// TestKeyAuthCredentialHandler_UpdateAllowList tests the UpdateAllowList method
func TestKeyAuthCredentialHandler_UpdateAllowList(t *testing.T) {
	handler := NewKeyAuthCredentialHandler()

	tests := []struct {
		name           string
		initialAllow   []interface{}
		operation      model.AllowListOperation
		consumerNames  []string
		expectedResult []string
	}{
		{
			name:           "ADD to empty list",
			initialAllow:   []interface{}{},
			operation:      model.AllowListOperationAdd,
			consumerNames:  []string{"consumer1", "consumer2"},
			expectedResult: []string{"consumer1", "consumer2"},
		},
		{
			name:           "ADD to existing list",
			initialAllow:   []interface{}{"consumer1"},
			operation:      model.AllowListOperationAdd,
			consumerNames:  []string{"consumer2", "consumer3"},
			expectedResult: []string{"consumer1", "consumer2", "consumer3"},
		},
		{
			name:           "ADD duplicate consumer",
			initialAllow:   []interface{}{"consumer1", "consumer2"},
			operation:      model.AllowListOperationAdd,
			consumerNames:  []string{"consumer2", "consumer3"},
			expectedResult: []string{"consumer1", "consumer2", "consumer3"},
		},
		{
			name:           "REMOVE from list",
			initialAllow:   []interface{}{"consumer1", "consumer2", "consumer3"},
			operation:      model.AllowListOperationRemove,
			consumerNames:  []string{"consumer2"},
			expectedResult: []string{"consumer1", "consumer3"},
		},
		{
			name:           "REMOVE non-existent consumer",
			initialAllow:   []interface{}{"consumer1", "consumer2"},
			operation:      model.AllowListOperationRemove,
			consumerNames:  []string{"consumer3"},
			expectedResult: []string{"consumer1", "consumer2"},
		},
		{
			name:           "REPLACE list",
			initialAllow:   []interface{}{"consumer1", "consumer2"},
			operation:      model.AllowListOperationReplace,
			consumerNames:  []string{"consumer3", "consumer4"},
			expectedResult: []string{"consumer3", "consumer4"},
		},
		{
			name:           "TOGGLE_ONLY with empty list",
			initialAllow:   []interface{}{},
			operation:      model.AllowListOperationToggleOnly,
			consumerNames:  nil,
			expectedResult: []string{},
		},
		{
			name:           "TOGGLE_ONLY with existing list",
			initialAllow:   []interface{}{"consumer1", "consumer2"},
			operation:      model.AllowListOperationToggleOnly,
			consumerNames:  nil,
			expectedResult: []string{"consumer1", "consumer2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instance := &model.WasmPluginInstance{
				Configurations: map[string]interface{}{
					"allow": tt.initialAllow,
				},
			}

			handler.UpdateAllowList(tt.operation, instance, tt.consumerNames)
			// GetAllowedConsumers returns []string from []interface{}, but UpdateAllowList stores []string
			// So we need to check the stored value directly
			allowRaw, exists := instance.Configurations["allow"]
			assert.True(t, exists)
			// After UpdateAllowList, the value is stored as []string
			if allowSlice, ok := allowRaw.([]string); ok {
				assert.Equal(t, tt.expectedResult, allowSlice)
			} else if allowSlice, ok := allowRaw.([]interface{}); ok {
				// Convert to []string for comparison
				result := make([]string, 0, len(allowSlice))
				for _, item := range allowSlice {
					if str, ok := item.(string); ok {
						result = append(result, str)
					}
				}
				assert.Equal(t, tt.expectedResult, result)
			} else {
				t.Errorf("Unexpected type for allow: %T", allowRaw)
			}
		})
	}
}

// TestKeyAuthCredentialHandler_DeleteConsumer tests the DeleteConsumer method
func TestKeyAuthCredentialHandler_DeleteConsumer(t *testing.T) {
	handler := NewKeyAuthCredentialHandler()

	tests := []struct {
		name            string
		instance        *model.WasmPluginInstance
		consumerName    string
		expectedDeleted bool
		expectedCount   int
	}{
		{
			name:            "nil instance",
			instance:        nil,
			consumerName:    "test-consumer",
			expectedDeleted: false,
			expectedCount:   0,
		},
		{
			name:            "nil configurations",
			instance:        &model.WasmPluginInstance{Configurations: nil},
			consumerName:    "test-consumer",
			expectedDeleted: false,
			expectedCount:   0,
		},
		{
			name: "no consumers config",
			instance: &model.WasmPluginInstance{
				Configurations: map[string]interface{}{},
			},
			consumerName:    "test-consumer",
			expectedDeleted: false,
			expectedCount:   0,
		},
		{
			name: "consumers config is not a slice",
			instance: &model.WasmPluginInstance{
				Configurations: map[string]interface{}{
					"consumers": "not-a-slice",
				},
			},
			consumerName:    "test-consumer",
			expectedDeleted: false,
			expectedCount:   0,
		},
		{
			name: "delete existing consumer",
			instance: &model.WasmPluginInstance{
				Configurations: map[string]interface{}{
					"consumers": []interface{}{
						map[string]interface{}{"name": "consumer1"},
						map[string]interface{}{"name": "consumer2"},
						map[string]interface{}{"name": "consumer3"},
					},
				},
			},
			consumerName:    "consumer2",
			expectedDeleted: true,
			expectedCount:   2,
		},
		{
			name: "delete non-existent consumer",
			instance: &model.WasmPluginInstance{
				Configurations: map[string]interface{}{
					"consumers": []interface{}{
						map[string]interface{}{"name": "consumer1"},
						map[string]interface{}{"name": "consumer2"},
					},
				},
			},
			consumerName:    "consumer3",
			expectedDeleted: false,
			expectedCount:   2,
		},
		{
			name: "consumer item is not a map",
			instance: &model.WasmPluginInstance{
				Configurations: map[string]interface{}{
					"consumers": []interface{}{
						"not-a-map",
						map[string]interface{}{"name": "consumer1"},
					},
				},
			},
			consumerName:    "consumer1",
			expectedDeleted: true,
			expectedCount:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.DeleteConsumer(tt.instance, tt.consumerName)
			assert.Equal(t, tt.expectedDeleted, result)
			if tt.instance != nil && tt.instance.Configurations != nil {
				consumers, _ := tt.instance.Configurations["consumers"].([]interface{})
				assert.Len(t, consumers, tt.expectedCount)
			}
		})
	}
}

// TestKeyAuthCredentialHandler_SaveConsumer tests the SaveConsumer method
func TestKeyAuthCredentialHandler_SaveConsumer(t *testing.T) {
	handler := NewKeyAuthCredentialHandler()

	tests := []struct {
		name          string
		instance      *model.WasmPluginInstance
		consumer      *model.Consumer
		expectedSaved bool
		expectPanic   bool
	}{
		{
			name:          "consumer with no credentials",
			instance:      &model.WasmPluginInstance{Configurations: map[string]interface{}{}},
			consumer:      &model.Consumer{Name: "test-consumer", Credentials: []model.Credential{}},
			expectedSaved: false,
			expectPanic:   false,
		},
		{
			name:          "consumer with nil credentials",
			instance:      &model.WasmPluginInstance{Configurations: map[string]interface{}{}},
			consumer:      &model.Consumer{Name: "test-consumer", Credentials: nil},
			expectedSaved: false,
			expectPanic:   false,
		},
		{
			name: "save new consumer with bearer credential",
			instance: &model.WasmPluginInstance{
				Configurations: map[string]interface{}{
					"consumers": []interface{}{},
				},
			},
			consumer: &model.Consumer{
				Name: "test-consumer",
				Credentials: []model.Credential{
					model.NewKeyAuthCredential("BEARER", "", []string{"token123"}),
				},
			},
			expectedSaved: true,
			expectPanic:   false,
		},
		{
			name: "save new consumer with header credential",
			instance: &model.WasmPluginInstance{
				Configurations: map[string]interface{}{
					"consumers": []interface{}{},
				},
			},
			consumer: &model.Consumer{
				Name: "test-consumer",
				Credentials: []model.Credential{
					model.NewKeyAuthCredential("HEADER", "x-api-key", []string{"secret-key"}),
				},
			},
			expectedSaved: true,
			expectPanic:   false,
		},
		{
			name: "save new consumer with query credential",
			instance: &model.WasmPluginInstance{
				Configurations: map[string]interface{}{
					"consumers": []interface{}{},
				},
			},
			consumer: &model.Consumer{
				Name: "test-consumer",
				Credentials: []model.Credential{
					model.NewKeyAuthCredential("QUERY", "api_key", []string{"query-token"}),
				},
			},
			expectedSaved: true,
			expectPanic:   false,
		},
		{
			name: "save consumer with non-keyauth credential",
			instance: &model.WasmPluginInstance{
				Configurations: map[string]interface{}{
					"consumers": []interface{}{},
				},
			},
			consumer: &model.Consumer{
				Name: "test-consumer",
				Credentials: []model.Credential{
					&mockCredential{},
				},
			},
			expectedSaved: false,
			expectPanic:   false,
		},
		{
			name: "update existing consumer",
			instance: &model.WasmPluginInstance{
				Configurations: map[string]interface{}{
					"consumers": []interface{}{
						map[string]interface{}{
							"name":        "test-consumer",
							"keys":        []interface{}{"x-api-key"},
							"in_header":   true,
							"credentials": []interface{}{"old-key"},
						},
					},
				},
			},
			consumer: &model.Consumer{
				Name: "test-consumer",
				Credentials: []model.Credential{
					model.NewKeyAuthCredential("HEADER", "x-api-key", []string{"new-key"}),
				},
			},
			expectedSaved: true,
			expectPanic:   false,
		},
		{
			name: "save consumer with invalid source should panic",
			instance: &model.WasmPluginInstance{
				Configurations: map[string]interface{}{
					"consumers": []interface{}{},
				},
			},
			consumer: &model.Consumer{
				Name: "test-consumer",
				Credentials: []model.Credential{
					model.NewKeyAuthCredential("invalid-source", "key", []string{"value"}),
				},
			},
			expectedSaved: false,
			expectPanic:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectPanic {
				defer func() {
					r := recover()
					assert.True(t, r != nil, "Expected panic but didn't get one")
				}()
			}

			result := handler.SaveConsumer(tt.instance, tt.consumer)
			assert.Equal(t, tt.expectedSaved, result)
		})
	}
}

// TestKeyAuthCredentialHandler_parseCredential tests the parseCredential method
func TestKeyAuthCredentialHandler_parseCredential(t *testing.T) {
	handler := NewKeyAuthCredentialHandler()

	tests := []struct {
		name           string
		consumerMap    map[string]interface{}
		expectNil      bool
		expectedSource string
		expectedKey    string
		expectedValues []string
	}{
		{
			name:        "nil map",
			consumerMap: nil,
			expectNil:   true,
		},
		{
			name:        "no keys",
			consumerMap: map[string]interface{}{},
			expectNil:   true,
		},
		{
			name: "keys is not a slice",
			consumerMap: map[string]interface{}{
				"keys": "not-a-slice",
			},
			expectNil: true,
		},
		{
			name: "empty keys slice",
			consumerMap: map[string]interface{}{
				"keys": []interface{}{},
			},
			expectNil: true,
		},
		{
			name: "keys slice with empty strings",
			consumerMap: map[string]interface{}{
				"keys": []interface{}{"", ""},
			},
			expectNil: true,
		},
		{
			name: "bearer token credential",
			consumerMap: map[string]interface{}{
				"keys":        []interface{}{"Authorization"},
				"in_header":   true,
				"credentials": []interface{}{"Bearer token123"},
			},
			expectNil:      false,
			expectedSource: "BEARER",
			expectedKey:    "",
			expectedValues: []string{"token123"},
		},
		{
			name: "header credential",
			consumerMap: map[string]interface{}{
				"keys":        []interface{}{"x-api-key"},
				"in_header":   true,
				"credentials": []interface{}{"secret-key"},
			},
			expectNil:      false,
			expectedSource: "HEADER",
			expectedKey:    "x-api-key",
			expectedValues: []string{"secret-key"},
		},
		{
			name: "query credential",
			consumerMap: map[string]interface{}{
				"keys":        []interface{}{"api_key"},
				"in_query":    true,
				"credentials": []interface{}{"query-token"},
			},
			expectNil:      false,
			expectedSource: "QUERY",
			expectedKey:    "api_key",
			expectedValues: []string{"query-token"},
		},
		{
			name: "credential with legacy single credential field",
			consumerMap: map[string]interface{}{
				"keys":       []interface{}{"x-api-key"},
				"in_header":  true,
				"credential": "legacy-key",
			},
			expectNil:      false,
			expectedSource: "HEADER",
			expectedKey:    "x-api-key",
			expectedValues: []string{"legacy-key"},
		},
		{
			name: "mixed bearer and non-bearer tokens",
			consumerMap: map[string]interface{}{
				"keys":        []interface{}{"Authorization"},
				"in_header":   true,
				"credentials": []interface{}{"Bearer token1", "Basic base64"},
			},
			expectNil:      false,
			expectedSource: "HEADER",
			expectedKey:    "Authorization",
			expectedValues: []string{"Bearer token1", "Basic base64"},
		},
		{
			name: "no in_header or in_query",
			consumerMap: map[string]interface{}{
				"keys":        []interface{}{"x-api-key"},
				"credentials": []interface{}{"some-key"},
			},
			expectNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.parseCredential(tt.consumerMap)
			if tt.expectNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedSource, result.Source)
				assert.Equal(t, tt.expectedKey, result.Key)
				assert.Equal(t, tt.expectedValues, result.Values)
			}
		})
	}
}

// TestKeyAuthCredentialHandler_hasSameCredential tests the hasSameCredential method
func TestKeyAuthCredentialHandler_hasSameCredential(t *testing.T) {
	handler := NewKeyAuthCredentialHandler()

	tests := []struct {
		name             string
		existingConsumer *model.Consumer
		credential       *model.KeyAuthCredential
		expectedResult   bool
	}{
		{
			name:             "nil credential",
			existingConsumer: &model.Consumer{Name: "test", Credentials: []model.Credential{}},
			credential:       nil,
			expectedResult:   false,
		},
		{
			name:             "nil existing consumer",
			existingConsumer: nil,
			credential:       model.NewKeyAuthCredential("HEADER", "x-api-key", []string{"key"}),
			expectedResult:   false,
		},
		{
			name: "existing consumer without keyauth credential",
			existingConsumer: &model.Consumer{
				Name:        "test",
				Credentials: []model.Credential{&mockCredential{}},
			},
			credential:     model.NewKeyAuthCredential("HEADER", "x-api-key", []string{"key"}),
			expectedResult: false,
		},
		{
			name: "same source, key, and overlapping values",
			existingConsumer: &model.Consumer{
				Name: "test",
				Credentials: []model.Credential{
					model.NewKeyAuthCredential("HEADER", "x-api-key", []string{"key1", "key2"}),
				},
			},
			credential:     model.NewKeyAuthCredential("HEADER", "x-api-key", []string{"key2", "key3"}),
			expectedResult: true,
		},
		{
			name: "same source, key, but no overlapping values",
			existingConsumer: &model.Consumer{
				Name: "test",
				Credentials: []model.Credential{
					model.NewKeyAuthCredential("HEADER", "x-api-key", []string{"key1", "key2"}),
				},
			},
			credential:     model.NewKeyAuthCredential("HEADER", "x-api-key", []string{"key3", "key4"}),
			expectedResult: false,
		},
		{
			name: "different source",
			existingConsumer: &model.Consumer{
				Name: "test",
				Credentials: []model.Credential{
					model.NewKeyAuthCredential("HEADER", "x-api-key", []string{"key1"}),
				},
			},
			credential:     model.NewKeyAuthCredential("QUERY", "x-api-key", []string{"key1"}),
			expectedResult: false,
		},
		{
			name: "different key",
			existingConsumer: &model.Consumer{
				Name: "test",
				Credentials: []model.Credential{
					model.NewKeyAuthCredential("HEADER", "x-api-key", []string{"key1"}),
				},
			},
			credential:     model.NewKeyAuthCredential("HEADER", "x-other-key", []string{"key1"}),
			expectedResult: false,
		},
		{
			name: "empty values",
			existingConsumer: &model.Consumer{
				Name: "test",
				Credentials: []model.Credential{
					model.NewKeyAuthCredential("HEADER", "x-api-key", []string{}),
				},
			},
			credential:     model.NewKeyAuthCredential("HEADER", "x-api-key", []string{"key1"}),
			expectedResult: false,
		},
		{
			name: "case insensitive source comparison",
			existingConsumer: &model.Consumer{
				Name: "test",
				Credentials: []model.Credential{
					model.NewKeyAuthCredential("HEADER", "x-api-key", []string{"key1"}),
				},
			},
			credential:     model.NewKeyAuthCredential("header", "x-api-key", []string{"key1"}),
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.hasSameCredential(tt.existingConsumer, tt.credential)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

// TestKeyAuthCredentialHandler_mergeExistedConfig tests the mergeExistedConfig method
func TestKeyAuthCredentialHandler_mergeExistedConfig(t *testing.T) {
	handler := NewKeyAuthCredentialHandler()

	tests := []struct {
		name           string
		credential     *model.KeyAuthCredential
		consumerConfig map[string]interface{}
		expectedSource string
		expectedKey    string
		expectedValues []string
	}{
		{
			name:           "merge with nil consumer config",
			credential:     model.NewKeyAuthCredential("HEADER", "x-api-key", []string{"key1"}),
			consumerConfig: nil,
			expectedSource: "HEADER",
			expectedKey:    "x-api-key",
			expectedValues: []string{"key1"},
		},
		{
			name:           "merge with empty consumer config",
			credential:     model.NewKeyAuthCredential("HEADER", "x-api-key", []string{"key1"}),
			consumerConfig: map[string]interface{}{},
			expectedSource: "HEADER",
			expectedKey:    "x-api-key",
			expectedValues: []string{"key1"},
		},
		{
			name:       "credential has all fields",
			credential: model.NewKeyAuthCredential("QUERY", "api_key", []string{"new-key"}),
			consumerConfig: map[string]interface{}{
				"name":        "test-consumer",
				"keys":        []interface{}{"x-api-key"},
				"in_header":   true,
				"credentials": []interface{}{"old-key"},
			},
			expectedSource: "QUERY",
			expectedKey:    "api_key",
			expectedValues: []string{"new-key"},
		},
		{
			name:       "credential missing source",
			credential: model.NewKeyAuthCredential("", "x-api-key", []string{"key1"}),
			consumerConfig: map[string]interface{}{
				"name":        "test-consumer",
				"keys":        []interface{}{"x-old-key"},
				"in_query":    true,
				"credentials": []interface{}{"old-key"},
			},
			expectedSource: "QUERY",
			expectedKey:    "x-api-key",
			expectedValues: []string{"key1"},
		},
		{
			name:       "credential missing key",
			credential: model.NewKeyAuthCredential("HEADER", "", []string{"key1"}),
			consumerConfig: map[string]interface{}{
				"name":        "test-consumer",
				"keys":        []interface{}{"x-old-key"},
				"in_header":   true,
				"credentials": []interface{}{"old-key"},
			},
			expectedSource: "HEADER",
			expectedKey:    "x-old-key",
			expectedValues: []string{"key1"},
		},
		{
			name:       "credential missing values",
			credential: model.NewKeyAuthCredential("HEADER", "x-api-key", []string{}),
			consumerConfig: map[string]interface{}{
				"name":        "test-consumer",
				"keys":        []interface{}{"x-old-key"},
				"in_header":   true,
				"credentials": []interface{}{"old-key"},
			},
			expectedSource: "HEADER",
			expectedKey:    "x-api-key",
			expectedValues: []string{"old-key"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.mergeExistedConfig(tt.credential, tt.consumerConfig)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedSource, result.Source)
			assert.Equal(t, tt.expectedKey, result.Key)
			assert.Equal(t, tt.expectedValues, result.Values)
		})
	}
}

// mockCredential is a mock implementation of Credential interface for testing
type mockCredential struct{}

func (m *mockCredential) GetType() string {
	return "mock"
}

func (m *mockCredential) Validate(forUpdate bool) error {
	return nil
}
