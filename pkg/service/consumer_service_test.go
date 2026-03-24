// Package service provides business services for the SDK
package service

import (
	"testing"

	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
	"github.com/stretchr/testify/assert"
)

// TestConsumerService_Interface tests that ConsumerService interface is defined
func TestConsumerService_Interface(t *testing.T) {
	// This test verifies the interface exists and has the expected methods
	// The actual implementation will be tested with integration tests
	var _ ConsumerService = (ConsumerService)(nil)
}

// TestConsumerModel tests the Consumer model
func TestConsumerModel(t *testing.T) {
	consumer := model.Consumer{
		Name: "test-consumer",
	}

	assert.Equal(t, "test-consumer", consumer.Name)
}

// TestConsumerWithCredentials tests consumer with credentials
func TestConsumerWithCredentials(t *testing.T) {
	cred := model.NewKeyAuthCredential("BEARER", "", []string{"test-api-key"})
	consumer := model.Consumer{
		Name:        "auth-consumer",
		Credentials: []model.Credential{cred},
	}

	assert.Equal(t, "auth-consumer", consumer.Name)
	assert.Len(t, consumer.Credentials, 1)
	assert.Equal(t, model.CredentialTypeKeyAuth, consumer.Credentials[0].GetType())
}

// TestKeyAuthCredential tests the KeyAuthCredential
func TestKeyAuthCredential(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		key     string
		values  []string
		wantErr bool
	}{
		{
			name:    "BEARER source without key",
			source:  "BEARER",
			key:     "",
			values:  []string{"test-value"},
			wantErr: false,
		},
		{
			name:    "HEADER source with key",
			source:  "HEADER",
			key:     "X-API-Key",
			values:  []string{"test-value"},
			wantErr: false,
		},
		{
			name:    "QUERY source with key",
			source:  "QUERY",
			key:     "api_key",
			values:  []string{"test-value"},
			wantErr: false,
		},
		{
			name:    "Empty source",
			source:  "",
			key:     "",
			values:  []string{"test-value"},
			wantErr: true,
		},
		{
			name:    "Invalid source",
			source:  "INVALID",
			key:     "",
			values:  []string{"test-value"},
			wantErr: true,
		},
		{
			name:    "HEADER source without key",
			source:  "HEADER",
			key:     "",
			values:  []string{"test-value"},
			wantErr: true,
		},
		{
			name:    "Empty values for create",
			source:  "BEARER",
			key:     "",
			values:  []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cred := model.NewKeyAuthCredential(tt.source, tt.key, tt.values)
			err := cred.Validate(false) // forUpdate = false
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestKeyAuthCredentialValidateForUpdate tests validation for update
func TestKeyAuthCredentialValidateForUpdate(t *testing.T) {
	// For update, empty values are allowed
	cred := model.NewKeyAuthCredential("BEARER", "", []string{})
	err := cred.Validate(true) // forUpdate = true
	assert.NoError(t, err)
}

// TestKeyAuthCredentialGetType tests GetType method
func TestKeyAuthCredentialGetType(t *testing.T) {
	cred := model.NewKeyAuthCredential("BEARER", "", []string{"test"})
	assert.Equal(t, model.CredentialTypeKeyAuth, cred.GetType())
}

// TestKeyAuthCredentialSource tests KeyAuthCredentialSource
func TestKeyAuthCredentialSource(t *testing.T) {
	tests := []struct {
		name       string
		source     string
		wantSource model.KeyAuthCredentialSource
		wantKeyReq bool
	}{
		{
			name:       "BEARER source",
			source:     "BEARER",
			wantSource: model.KeyAuthCredentialSourceBearer,
			wantKeyReq: false,
		},
		{
			name:       "HEADER source",
			source:     "HEADER",
			wantSource: model.KeyAuthCredentialSourceHeader,
			wantKeyReq: true,
		},
		{
			name:       "QUERY source",
			source:     "QUERY",
			wantSource: model.KeyAuthCredentialSourceQuery,
			wantKeyReq: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed := model.ParseKeyAuthCredentialSource(tt.source)
			assert.Equal(t, tt.wantSource, parsed)
			assert.Equal(t, tt.wantKeyReq, parsed.IsKeyRequired())
		})
	}
}

// TestConsumerValidate tests Consumer validation
func TestConsumerValidate(t *testing.T) {
	tests := []struct {
		name     string
		consumer model.Consumer
		wantErr  bool
	}{
		{
			name: "Valid consumer",
			consumer: model.Consumer{
				Name:        "test-consumer",
				Credentials: []model.Credential{model.NewKeyAuthCredential("BEARER", "", []string{"test"})},
			},
			wantErr: false,
		},
		{
			name: "Empty name",
			consumer: model.Consumer{
				Name:        "",
				Credentials: []model.Credential{model.NewKeyAuthCredential("BEARER", "", []string{"test"})},
			},
			wantErr: true,
		},
		{
			name: "Empty credentials",
			consumer: model.Consumer{
				Name:        "test-consumer",
				Credentials: []model.Credential{},
			},
			wantErr: true,
		},
		{
			name: "Nil credentials",
			consumer: model.Consumer{
				Name:        "test-consumer",
				Credentials: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.consumer.Validate(false)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestAllowListModel tests the AllowList model
func TestAllowListModel(t *testing.T) {
	authEnabled := true
	allowList := model.AllowList{
		Targets:         map[model.WasmPluginInstanceScope]string{model.WasmPluginInstanceScopeGlobal: "global"},
		AuthEnabled:     &authEnabled,
		CredentialTypes: []string{"key-auth"},
		ConsumerNames:   []string{"consumer1", "consumer2"},
	}

	assert.Len(t, allowList.ConsumerNames, 2)
	assert.Contains(t, allowList.ConsumerNames, "consumer1")
	assert.True(t, *allowList.AuthEnabled)
}

// TestNewAllowList tests NewAllowList function
func TestNewAllowList(t *testing.T) {
	allowList := model.NewAllowList()
	assert.NotNil(t, allowList)
	assert.NotNil(t, allowList.Targets)
	assert.Empty(t, allowList.CredentialTypes)
	assert.Empty(t, allowList.ConsumerNames)
}

// TestForTarget tests ForTarget function
func TestForTarget(t *testing.T) {
	allowList := model.ForTarget(model.WasmPluginInstanceScopeGlobal, "global-target")
	assert.NotNil(t, allowList)
	assert.Equal(t, "global-target", allowList.Targets[model.WasmPluginInstanceScopeGlobal])
}

// TestAllowListOperationConstants tests the AllowListOperation constants
func TestAllowListOperationConstants(t *testing.T) {
	assert.Equal(t, model.AllowListOperation("ADD"), model.AllowListOperationAdd)
	assert.Equal(t, model.AllowListOperation("REMOVE"), model.AllowListOperationRemove)
	assert.Equal(t, model.AllowListOperation("REPLACE"), model.AllowListOperationReplace)
	assert.Equal(t, model.AllowListOperation("TOGGLE_ONLY"), model.AllowListOperationToggleOnly)
}
