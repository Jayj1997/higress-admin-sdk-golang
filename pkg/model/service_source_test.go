// Package model provides data models for the SDK
package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceSource_Validate(t *testing.T) {
	tests := []struct {
		name        string
		source      *ServiceSource
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid service source",
			source: &ServiceSource{
				Name: "my-service-source",
				Type: "dns",
			},
			expectError: false,
		},
		{
			name: "valid with all fields",
			source: &ServiceSource{
				Name:      "nacos-source",
				Version:   "v1",
				Type:      "nacos",
				Domain:    "nacos.example.com",
				Port:      intPtrSS(8848),
				Namespace: "public",
				Group:     "DEFAULT_GROUP",
				Services:  []string{"service1", "service2"},
			},
			expectError: false,
		},
		{
			name: "missing name",
			source: &ServiceSource{
				Type: "dns",
			},
			expectError: true,
			errorMsg:    "name is required",
		},
		{
			name: "missing type",
			source: &ServiceSource{
				Name: "my-service-source",
			},
			expectError: true,
			errorMsg:    "type is required",
		},
		{
			name: "empty name",
			source: &ServiceSource{
				Name: "",
				Type: "dns",
			},
			expectError: true,
			errorMsg:    "name is required",
		},
		{
			name: "empty type",
			source: &ServiceSource{
				Name: "my-service-source",
				Type: "",
			},
			expectError: true,
			errorMsg:    "type is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.source.Validate()
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestServiceSource_Structure(t *testing.T) {
	authEnabled := true
	port := 8848
	source := &ServiceSource{
		Name:           "nacos-registry",
		Version:        "v1",
		Type:           "nacos",
		Domain:         "nacos.example.com:8848",
		Port:           &port,
		Namespace:      "production",
		Group:          "DEFAULT_GROUP",
		Services:       []string{"user-service", "order-service"},
		AuthEnabled:    &authEnabled,
		AuthType:       "username_password",
		AuthSecretName: "nacos-credentials",
		Properties: map[string]string{
			"namespaceId": "production",
			"groupName":   "DEFAULT_GROUP",
		},
	}

	assert.Equal(t, "nacos-registry", source.Name)
	assert.Equal(t, "v1", source.Version)
	assert.Equal(t, "nacos", source.Type)
	assert.Equal(t, "nacos.example.com:8848", source.Domain)
	require.NotNil(t, source.Port)
	assert.Equal(t, 8848, *source.Port)
	assert.Equal(t, "production", source.Namespace)
	assert.Equal(t, "DEFAULT_GROUP", source.Group)
	assert.Len(t, source.Services, 2)
	require.NotNil(t, source.AuthEnabled)
	assert.True(t, *source.AuthEnabled)
	assert.Equal(t, "username_password", source.AuthType)
	assert.Equal(t, "nacos-credentials", source.AuthSecretName)
	assert.Len(t, source.Properties, 2)
}

func TestServiceSource_Types(t *testing.T) {
	tests := []struct {
		name      string
		typeValue string
	}{
		{"dns type", ServiceSourceTypeDNS},
		{"static type", ServiceSourceTypeStatic},
		{"nacos type", "nacos"},
		{"zookeeper type", "zookeeper"},
		{"eureka type", "eureka"},
		{"consul type", "consul"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := &ServiceSource{
				Name: "test-source",
				Type: tt.typeValue,
			}
			assert.Equal(t, tt.typeValue, source.Type)
		})
	}
}

func TestServiceSource_DNSType(t *testing.T) {
	source := &ServiceSource{
		Name:   "dns-source",
		Type:   ServiceSourceTypeDNS,
		Domain: "api.example.com",
	}

	assert.Equal(t, "dns", source.Type)
	assert.Equal(t, "api.example.com", source.Domain)
}

func TestServiceSource_StaticType(t *testing.T) {
	source := &ServiceSource{
		Name:   "static-source",
		Type:   ServiceSourceTypeStatic,
		Domain: "192.168.1.1:8080,192.168.1.2:8080",
	}

	assert.Equal(t, "static", source.Type)
	assert.Equal(t, "192.168.1.1:8080,192.168.1.2:8080", source.Domain)
}

func TestServiceSource_NacosType(t *testing.T) {
	port := 8848
	source := &ServiceSource{
		Name:      "nacos-source",
		Type:      "nacos",
		Domain:    "nacos.example.com",
		Port:      &port,
		Namespace: "dev",
		Group:     "DEFAULT_GROUP",
		Services:  []string{"service-a", "service-b"},
	}

	assert.Equal(t, "nacos", source.Type)
	assert.Equal(t, "nacos.example.com", source.Domain)
	require.NotNil(t, source.Port)
	assert.Equal(t, 8848, *source.Port)
	assert.Equal(t, "dev", source.Namespace)
	assert.Equal(t, "DEFAULT_GROUP", source.Group)
}

func TestServiceSource_WithAuth(t *testing.T) {
	authEnabled := true
	source := &ServiceSource{
		Name:           "secured-source",
		Type:           "nacos",
		AuthEnabled:    &authEnabled,
		AuthType:       "token",
		AuthSecretName: "nacos-token-secret",
	}

	require.NotNil(t, source.AuthEnabled)
	assert.True(t, *source.AuthEnabled)
	assert.Equal(t, "token", source.AuthType)
	assert.Equal(t, "nacos-token-secret", source.AuthSecretName)
}

func TestServiceSource_AuthDisabled(t *testing.T) {
	authEnabled := false
	source := &ServiceSource{
		Name:        "unsecured-source",
		Type:        "dns",
		AuthEnabled: &authEnabled,
	}

	require.NotNil(t, source.AuthEnabled)
	assert.False(t, *source.AuthEnabled)
}

func TestServiceSource_NilAuthEnabled(t *testing.T) {
	source := &ServiceSource{
		Name: "no-auth-source",
		Type: "dns",
	}

	assert.Nil(t, source.AuthEnabled)
}

func TestServiceSource_NilPort(t *testing.T) {
	source := &ServiceSource{
		Name: "no-port-source",
		Type: "dns",
	}

	assert.Nil(t, source.Port)
}

func TestServiceSource_Properties(t *testing.T) {
	source := &ServiceSource{
		Name: "source-with-props",
		Type: "nacos",
		Properties: map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		},
	}

	assert.Len(t, source.Properties, 3)
	assert.Equal(t, "value1", source.Properties["key1"])
	assert.Equal(t, "value2", source.Properties["key2"])
	assert.Equal(t, "value3", source.Properties["key3"])
}

func TestServiceSource_MultipleServices(t *testing.T) {
	source := &ServiceSource{
		Name:     "multi-service-source",
		Type:     "nacos",
		Services: []string{"svc1", "svc2", "svc3", "svc4", "svc5"},
	}

	assert.Len(t, source.Services, 5)
}

func TestServiceSource_EmptyFields(t *testing.T) {
	source := &ServiceSource{}

	assert.Empty(t, source.Name)
	assert.Empty(t, source.Version)
	assert.Empty(t, source.Type)
	assert.Empty(t, source.Domain)
	assert.Nil(t, source.Port)
	assert.Empty(t, source.Namespace)
	assert.Empty(t, source.Group)
	assert.Nil(t, source.Services)
	assert.Nil(t, source.AuthEnabled)
	assert.Empty(t, source.AuthType)
	assert.Empty(t, source.AuthSecretName)
	assert.Nil(t, source.Properties)
}

func TestServiceSourceType_Constants(t *testing.T) {
	assert.Equal(t, "dns", ServiceSourceTypeDNS)
	assert.Equal(t, "static", ServiceSourceTypeStatic)
}

// Helper function
func intPtrSS(i int) *int {
	return &i
}
