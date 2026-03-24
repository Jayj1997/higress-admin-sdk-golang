// Package model provides data models for the SDK
package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConsumer_Validate(t *testing.T) {
	tests := []struct {
		name        string
		consumer    *Consumer
		forUpdate   bool
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid consumer with bearer credential",
			consumer: &Consumer{
				Name: "test-consumer",
				Credentials: []Credential{
					NewKeyAuthCredential("BEARER", "", []string{"test-api-key"}),
				},
			},
			forUpdate:   false,
			expectError: false,
		},
		{
			name: "valid consumer with header credential",
			consumer: &Consumer{
				Name: "test-consumer",
				Credentials: []Credential{
					NewKeyAuthCredential("HEADER", "X-API-Key", []string{"test-api-key"}),
				},
			},
			forUpdate:   false,
			expectError: false,
		},
		{
			name: "valid consumer with query credential",
			consumer: &Consumer{
				Name: "test-consumer",
				Credentials: []Credential{
					NewKeyAuthCredential("QUERY", "api_key", []string{"test-api-key"}),
				},
			},
			forUpdate:   false,
			expectError: false,
		},
		{
			name: "valid consumer with multiple credentials",
			consumer: &Consumer{
				Name: "test-consumer",
				Credentials: []Credential{
					NewKeyAuthCredential("BEARER", "", []string{"bearer-key"}),
					NewKeyAuthCredential("HEADER", "X-API-Key", []string{"header-key"}),
				},
			},
			forUpdate:   false,
			expectError: false,
		},
		{
			name: "missing name",
			consumer: &Consumer{
				Credentials: []Credential{
					NewKeyAuthCredential("BEARER", "", []string{"test-api-key"}),
				},
			},
			forUpdate:   false,
			expectError: true,
			errorMsg:    "name cannot be blank",
		},
		{
			name: "empty credentials",
			consumer: &Consumer{
				Name:        "test-consumer",
				Credentials: []Credential{},
			},
			forUpdate:   false,
			expectError: true,
			errorMsg:    "credentials cannot be empty",
		},
		{
			name: "nil credentials",
			consumer: &Consumer{
				Name:        "test-consumer",
				Credentials: nil,
			},
			forUpdate:   false,
			expectError: true,
			errorMsg:    "credentials cannot be empty",
		},
		{
			name: "invalid credential",
			consumer: &Consumer{
				Name: "test-consumer",
				Credentials: []Credential{
					NewKeyAuthCredential("", "", []string{"test-api-key"}),
				},
			},
			forUpdate:   false,
			expectError: true,
			errorMsg:    "source cannot be blank",
		},
		{
			name: "for update with empty values",
			consumer: &Consumer{
				Name: "test-consumer",
				Credentials: []Credential{
					NewKeyAuthCredential("BEARER", "", []string{}),
				},
			},
			forUpdate:   true,
			expectError: false,
		},
		{
			name: "not for update with empty values",
			consumer: &Consumer{
				Name: "test-consumer",
				Credentials: []Credential{
					NewKeyAuthCredential("BEARER", "", []string{}),
				},
			},
			forUpdate:   false,
			expectError: true,
			errorMsg:    "values cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.consumer.Validate(tt.forUpdate)
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestKeyAuthCredential_GetType(t *testing.T) {
	cred := NewKeyAuthCredential("BEARER", "", []string{"test-key"})
	assert.Equal(t, CredentialTypeKeyAuth, cred.GetType())
}

func TestKeyAuthCredential_Validate(t *testing.T) {
	tests := []struct {
		name        string
		credential  *KeyAuthCredential
		forUpdate   bool
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid bearer credential",
			credential:  NewKeyAuthCredential("BEARER", "", []string{"test-api-key"}),
			forUpdate:   false,
			expectError: false,
		},
		{
			name:        "valid header credential",
			credential:  NewKeyAuthCredential("HEADER", "X-API-Key", []string{"test-api-key"}),
			forUpdate:   false,
			expectError: false,
		},
		{
			name:        "valid query credential",
			credential:  NewKeyAuthCredential("QUERY", "api_key", []string{"test-api-key"}),
			forUpdate:   false,
			expectError: false,
		},
		{
			name:        "valid credential with multiple values",
			credential:  NewKeyAuthCredential("BEARER", "", []string{"key1", "key2", "key3"}),
			forUpdate:   false,
			expectError: false,
		},
		{
			name:        "empty source",
			credential:  NewKeyAuthCredential("", "", []string{"test-api-key"}),
			forUpdate:   false,
			expectError: true,
			errorMsg:    "source cannot be blank",
		},
		{
			name:        "unknown source",
			credential:  NewKeyAuthCredential("UNKNOWN", "", []string{"test-api-key"}),
			forUpdate:   false,
			expectError: true,
			errorMsg:    "unknown source value: UNKNOWN",
		},
		{
			name:        "header without key",
			credential:  NewKeyAuthCredential("HEADER", "", []string{"test-api-key"}),
			forUpdate:   false,
			expectError: true,
			errorMsg:    "key cannot be blank",
		},
		{
			name:        "query without key",
			credential:  NewKeyAuthCredential("QUERY", "", []string{"test-api-key"}),
			forUpdate:   false,
			expectError: true,
			errorMsg:    "key cannot be blank",
		},
		{
			name:        "bearer with key - key is ignored but valid",
			credential:  NewKeyAuthCredential("BEARER", "X-API-Key", []string{"test-api-key"}),
			forUpdate:   false,
			expectError: false,
		},
		{
			name:        "empty values for update",
			credential:  NewKeyAuthCredential("BEARER", "", []string{}),
			forUpdate:   true,
			expectError: false,
		},
		{
			name:        "empty values not for update",
			credential:  NewKeyAuthCredential("BEARER", "", []string{}),
			forUpdate:   false,
			expectError: true,
			errorMsg:    "values cannot be empty",
		},
		{
			name:        "nil values not for update",
			credential:  &KeyAuthCredential{Source: "BEARER", Key: "", Values: nil},
			forUpdate:   false,
			expectError: true,
			errorMsg:    "values cannot be empty",
		},
		{
			name:        "nil values for update",
			credential:  &KeyAuthCredential{Source: "BEARER", Key: "", Values: nil},
			forUpdate:   true,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.credential.Validate(tt.forUpdate)
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestKeyAuthCredentialSource_IsKeyRequired(t *testing.T) {
	tests := []struct {
		source         KeyAuthCredentialSource
		expectedResult bool
	}{
		{KeyAuthCredentialSourceBearer, false},
		{KeyAuthCredentialSourceHeader, true},
		{KeyAuthCredentialSourceQuery, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.source), func(t *testing.T) {
			assert.Equal(t, tt.expectedResult, tt.source.IsKeyRequired())
		})
	}
}

func TestParseKeyAuthCredentialSource(t *testing.T) {
	tests := []struct {
		input    string
		expected KeyAuthCredentialSource
	}{
		{"BEARER", KeyAuthCredentialSourceBearer},
		{"HEADER", KeyAuthCredentialSourceHeader},
		{"QUERY", KeyAuthCredentialSourceQuery},
		{"UNKNOWN", ""},
		{"", ""},
		{"bearer", ""}, // case sensitive
		{"header", ""}, // case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ParseKeyAuthCredentialSource(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestKeyAuthCredentialSource_Constants(t *testing.T) {
	assert.Equal(t, KeyAuthCredentialSource("BEARER"), KeyAuthCredentialSourceBearer)
	assert.Equal(t, KeyAuthCredentialSource("HEADER"), KeyAuthCredentialSourceHeader)
	assert.Equal(t, KeyAuthCredentialSource("QUERY"), KeyAuthCredentialSourceQuery)
}

func TestCredentialType_Constant(t *testing.T) {
	assert.Equal(t, "key-auth", CredentialTypeKeyAuth)
}

func TestNewKeyAuthCredential(t *testing.T) {
	values := []string{"key1", "key2"}
	cred := NewKeyAuthCredential("HEADER", "X-API-Key", values)

	require.NotNil(t, cred)
	assert.Equal(t, "HEADER", cred.Source)
	assert.Equal(t, "X-API-Key", cred.Key)
	assert.Equal(t, values, cred.Values)
}

// mockCredential is a mock implementation of Credential for testing
type mockCredential struct {
	getType     string
	validateErr error
}

func (m *mockCredential) GetType() string {
	return m.getType
}

func (m *mockCredential) Validate(forUpdate bool) error {
	return m.validateErr
}

func TestConsumer_Validate_WithMockCredential(t *testing.T) {
	t.Run("with mock credential that passes validation", func(t *testing.T) {
		consumer := &Consumer{
			Name: "test-consumer",
			Credentials: []Credential{
				&mockCredential{getType: "mock", validateErr: nil},
			},
		}
		err := consumer.Validate(false)
		require.NoError(t, err)
	})

	t.Run("with mock credential that fails validation", func(t *testing.T) {
		consumer := &Consumer{
			Name: "test-consumer",
			Credentials: []Credential{
				&mockCredential{getType: "mock", validateErr: assert.AnError},
			},
		}
		err := consumer.Validate(false)
		require.Error(t, err)
	})
}
