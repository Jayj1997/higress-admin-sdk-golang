// Package service provides business services for the SDK
package service

import (
	"testing"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDomainService_Interface tests that the implementation satisfies the interface
func TestDomainService_Interface(t *testing.T) {
	var _ DomainService = (*DomainServiceImpl)(nil)
}

// TestDomainService_New tests creating a new DomainService
func TestDomainService_New(t *testing.T) {
	// Create service with nil dependencies (for testing without cluster)
	svc := NewDomainService(nil, nil, nil, nil)
	require.NotNil(t, svc)
}

// TestDomainModel tests the Domain model
func TestDomainModel(t *testing.T) {
	domain := model.Domain{
		Name:        "example.com",
		EnableHTTPS: model.EnableHTTPSOn,
	}

	assert.Equal(t, "example.com", domain.Name)
	assert.Equal(t, model.EnableHTTPSOn, domain.EnableHTTPS)
	assert.True(t, domain.IsHTTPS())
	assert.False(t, domain.IsForceHTTPS())
}

// TestDomainModel_Validate tests domain validation
func TestDomainModel_Validate(t *testing.T) {
	tests := []struct {
		name        string
		domain      *model.Domain
		expectError bool
	}{
		{
			name: "valid domain",
			domain: &model.Domain{
				Name: "example.com",
			},
			expectError: false,
		},
		{
			name: "empty domain name",
			domain: &model.Domain{
				Name: "",
			},
			expectError: true,
		},
		{
			name: "domain with subdomain",
			domain: &model.Domain{
				Name: "api.example.com",
			},
			expectError: false,
		},
		{
			name: "domain with wildcard",
			domain: &model.Domain{
				Name: "*.example.com",
			},
			expectError: false,
		},
		{
			name: "domain with HTTPS enabled",
			domain: &model.Domain{
				Name:        "secure.example.com",
				EnableHTTPS: model.EnableHTTPSOn,
			},
			expectError: false,
		},
		{
			name: "domain with forced HTTPS",
			domain: &model.Domain{
				Name:        "forced.example.com",
				EnableHTTPS: model.EnableHTTPSForce,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.domain.Validate()
			if tt.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

// TestDomainModel_HTTPSMethods tests HTTPS-related methods
func TestDomainModel_HTTPSMethods(t *testing.T) {
	tests := []struct {
		name         string
		enableHTTPS  string
		isHTTPS      bool
		isForceHTTPS bool
	}{
		{
			name:         "HTTPS off",
			enableHTTPS:  model.EnableHTTPSOff,
			isHTTPS:      false,
			isForceHTTPS: false,
		},
		{
			name:         "HTTPS on",
			enableHTTPS:  model.EnableHTTPSOn,
			isHTTPS:      true,
			isForceHTTPS: false,
		},
		{
			name:         "HTTPS force",
			enableHTTPS:  model.EnableHTTPSForce,
			isHTTPS:      true,
			isForceHTTPS: true,
		},
		{
			name:         "empty HTTPS setting",
			enableHTTPS:  "",
			isHTTPS:      false,
			isForceHTTPS: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domain := &model.Domain{
				Name:        "test.com",
				EnableHTTPS: tt.enableHTTPS,
			}
			assert.Equal(t, tt.isHTTPS, domain.IsHTTPS())
			assert.Equal(t, tt.isForceHTTPS, domain.IsForceHTTPS())
		})
	}
}

// TestDomainPagination tests pagination logic for domains
func TestDomainPagination(t *testing.T) {
	domains := make([]model.Domain, 25)
	for i := 0; i < 25; i++ {
		domains[i] = model.Domain{
			Name: "domain-" + string(rune('a'+i)),
		}
	}

	// Test first page
	total := len(domains)
	pageNum := 1
	pageSize := 10
	start := (pageNum - 1) * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}
	pagedData := domains[start:end]
	result := model.NewPaginatedResult(pagedData, total, pageNum, pageSize)

	assert.Equal(t, 10, len(result.Data))
	assert.Equal(t, 25, result.Total)
	assert.Equal(t, 1, result.PageNum)
	assert.Equal(t, 10, result.PageSize)
	assert.Equal(t, 3, result.TotalPages)

	// Test second page
	pageNum = 2
	start = (pageNum - 1) * pageSize
	end = start + pageSize
	if end > total {
		end = total
	}
	pagedData = domains[start:end]
	result = model.NewPaginatedResult(pagedData, total, pageNum, pageSize)

	assert.Equal(t, 10, len(result.Data))
	assert.Equal(t, 2, result.PageNum)

	// Test last page
	pageNum = 3
	start = (pageNum - 1) * pageSize
	end = start + pageSize
	if end > total {
		end = total
	}
	pagedData = domains[start:end]
	result = model.NewPaginatedResult(pagedData, total, pageNum, pageSize)

	assert.Equal(t, 5, len(result.Data))
	assert.Equal(t, 3, result.PageNum)
}

// TestCommonPageQueryForDomain tests CommonPageQuery for domain listing
func TestCommonPageQueryForDomain(t *testing.T) {
	// Test default values
	query := &model.CommonPageQuery{}
	assert.Equal(t, 0, query.GetOffset())
	assert.Equal(t, 10, query.GetPageSize())

	// Test custom values
	query = &model.CommonPageQuery{
		PageNum:  2,
		PageSize: 20,
	}
	assert.Equal(t, 20, query.GetOffset())
	assert.Equal(t, 20, query.GetPageSize())

	// Test max page size
	query = &model.CommonPageQuery{
		PageNum:  1,
		PageSize: 200,
	}
	assert.Equal(t, 100, query.GetPageSize())
}

// TestDomainEnableHTTPSConstants tests the HTTPS constants
func TestDomainEnableHTTPSConstants(t *testing.T) {
	assert.Equal(t, "off", model.EnableHTTPSOff)
	assert.Equal(t, "on", model.EnableHTTPSOn)
	assert.Equal(t, "force", model.EnableHTTPSForce)
}
