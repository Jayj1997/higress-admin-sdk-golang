// Package service provides business services for the SDK
package service

import (
	"testing"

	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestServiceSourceService_Interface tests that the implementation satisfies the interface
func TestServiceSourceService_Interface(t *testing.T) {
	var _ ServiceSourceService = (*ServiceSourceServiceImpl)(nil)
}

// TestServiceSourceService_New tests creating a new ServiceSourceService
func TestServiceSourceService_New(t *testing.T) {
	svc := NewServiceSourceService(nil, nil)
	require.NotNil(t, svc)
}

// TestServiceSourceModel tests the ServiceSource model
func TestServiceSourceModel(t *testing.T) {
	port := 8848
	ss := model.ServiceSource{
		Name:      "test-registry",
		Type:      "nacos",
		Domain:    "nacos.example.com",
		Port:      &port,
		Namespace: "public",
		Group:     "DEFAULT_GROUP",
	}

	assert.Equal(t, "test-registry", ss.Name)
	assert.Equal(t, "nacos", ss.Type)
	assert.Equal(t, "nacos.example.com", ss.Domain)
	assert.Equal(t, 8848, *ss.Port)
	assert.Equal(t, "public", ss.Namespace)
	assert.Equal(t, "DEFAULT_GROUP", ss.Group)
}

// TestServiceSourceType tests different service source types
func TestServiceSourceType(t *testing.T) {
	tests := []struct {
		name        string
		sourceType  string
		expectValid bool
	}{
		{"nacos type", "nacos", true},
		{"zookeeper type", "zookeeper", true},
		{"consul type", "consul", true},
		{"eureka type", "eureka", true},
		{"static type", model.ServiceSourceTypeStatic, true},
		{"dns type", model.ServiceSourceTypeDNS, true},
		{"empty type", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := &model.ServiceSource{Type: tt.sourceType}
			if tt.expectValid {
				assert.NotEmpty(t, ss.Type)
			} else {
				assert.Empty(t, ss.Type)
			}
		})
	}
}

// TestServiceSourceValidation tests service source validation
func TestServiceSourceValidation(t *testing.T) {
	tests := []struct {
		name        string
		source      *model.ServiceSource
		expectError bool
	}{
		{
			name: "valid service source",
			source: &model.ServiceSource{
				Name: "valid-source",
				Type: "nacos",
			},
			expectError: false,
		},
		{
			name: "missing name",
			source: &model.ServiceSource{
				Type: "nacos",
			},
			expectError: true,
		},
		{
			name: "missing type",
			source: &model.ServiceSource{
				Name: "no-type",
			},
			expectError: true,
		},
		{
			name: "DNS type service source",
			source: &model.ServiceSource{
				Name: "dns-source",
				Type: model.ServiceSourceTypeDNS,
				Domain: "api.example.com",
			},
			expectError: false,
		},
		{
			name: "static type service source",
			source: &model.ServiceSource{
				Name:   "static-source",
				Type:   model.ServiceSourceTypeStatic,
				Domain: "192.168.1.1:8080",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.source.Validate()
			if tt.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

// TestServiceSourceWithAuth tests service source with authentication
func TestServiceSourceWithAuth(t *testing.T) {
	authEnabled := true
	ss := model.ServiceSource{
		Name:          "auth-source",
		Type:          "nacos",
		AuthEnabled:   &authEnabled,
		AuthType:      "basic",
		AuthSecretName: "nacos-credentials",
	}

	assert.True(t, *ss.AuthEnabled)
	assert.Equal(t, "basic", ss.AuthType)
	assert.Equal(t, "nacos-credentials", ss.AuthSecretName)
}

// TestServiceSourceWithProperties tests service source with properties
func TestServiceSourceWithProperties(t *testing.T) {
	ss := model.ServiceSource{
		Name: "props-source",
		Type: "nacos",
		Properties: map[string]string{
			"serverAddr": "nacos.example.com:8848",
			"namespace":  "dev",
		},
	}

	assert.Equal(t, "nacos.example.com:8848", ss.Properties["serverAddr"])
	assert.Equal(t, "dev", ss.Properties["namespace"])
}

// TestServiceSourcePagination tests pagination logic for service sources
func TestServiceSourcePagination(t *testing.T) {
	sources := make([]model.ServiceSource, 15)
	for i := 0; i < 15; i++ {
		sources[i] = model.ServiceSource{
			Name: "source-" + string(rune('a'+i)),
			Type: "nacos",
		}
	}

	total := len(sources)
	pageNum := 1
	pageSize := 10
	start := (pageNum - 1) * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}
	pagedData := sources[start:end]
	result := model.NewPaginatedResult(pagedData, total, pageNum, pageSize)

	assert.Equal(t, 10, len(result.Data))
	assert.Equal(t, 15, result.Total)
	assert.Equal(t, 2, result.TotalPages)
}

// TestServiceSourceTypeConstants tests the service source type constants
func TestServiceSourceTypeConstants(t *testing.T) {
	assert.Equal(t, "dns", model.ServiceSourceTypeDNS)
	assert.Equal(t, "static", model.ServiceSourceTypeStatic)
}
