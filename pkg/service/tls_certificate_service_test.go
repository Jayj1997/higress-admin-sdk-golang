// Package service provides business services for the SDK
package service

import (
	"testing"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTlsCertificateService_Interface tests that the implementation satisfies the interface
func TestTlsCertificateService_Interface(t *testing.T) {
	var _ TlsCertificateService = (*TlsCertificateServiceImpl)(nil)
}

// TestTlsCertificateService_New tests creating a new TlsCertificateService
func TestTlsCertificateService_New(t *testing.T) {
	svc := NewTlsCertificateService(nil, nil)
	require.NotNil(t, svc)
}

// TestTlsCertificateModel tests the TlsCertificate model
func TestTlsCertificateModel(t *testing.T) {
	cert := model.TlsCertificate{
		Name:    "test-cert",
		Domains: []string{"example.com", "api.example.com"},
		Cert:    "-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----",
		Key:     "-----BEGIN PRIVATE KEY-----\nMIIE...\n-----END PRIVATE KEY-----",
	}

	assert.Equal(t, "test-cert", cert.Name)
	assert.Len(t, cert.Domains, 2)
	assert.Contains(t, cert.Domains, "example.com")
	assert.NotEmpty(t, cert.Cert)
	assert.NotEmpty(t, cert.Key)
}

// TestTlsCertificateValidation tests TLS certificate validation
func TestTlsCertificateValidation(t *testing.T) {
	tests := []struct {
		name        string
		cert        *model.TlsCertificate
		expectError bool
	}{
		{
			name: "valid certificate",
			cert: &model.TlsCertificate{
				Name: "valid-cert",
				Cert: "-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----",
				Key:  "-----BEGIN PRIVATE KEY-----\nMIIE...\n-----END PRIVATE KEY-----",
			},
			expectError: false,
		},
		{
			name: "missing name",
			cert: &model.TlsCertificate{
				Cert: "-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----",
				Key:  "-----BEGIN PRIVATE KEY-----\nMIIE...\n-----END PRIVATE KEY-----",
			},
			expectError: true,
		},
		{
			name: "missing cert",
			cert: &model.TlsCertificate{
				Name: "no-cert",
				Key:  "-----BEGIN PRIVATE KEY-----\nMIIE...\n-----END PRIVATE KEY-----",
			},
			expectError: true,
		},
		{
			name: "missing key",
			cert: &model.TlsCertificate{
				Name: "no-key",
				Cert: "-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----",
			},
			expectError: true,
		},
		{
			name: "certificate with multiple domains",
			cert: &model.TlsCertificate{
				Name:    "multi-domain-cert",
				Domains: []string{"example.com", "api.example.com", "*.example.com"},
				Cert:    "-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----",
				Key:     "-----BEGIN PRIVATE KEY-----\nMIIE...\n-----END PRIVATE KEY-----",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cert.Validate()
			if tt.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

// TestTlsCertificateWithWildcard tests TLS certificate with wildcard domain
func TestTlsCertificateWithWildcard(t *testing.T) {
	cert := model.TlsCertificate{
		Name:    "wildcard-cert",
		Domains: []string{"*.example.com"},
		Cert:    "-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----",
		Key:     "-----BEGIN PRIVATE KEY-----\nMIIE...\n-----END PRIVATE KEY-----",
	}

	assert.Equal(t, "wildcard-cert", cert.Name)
	assert.Contains(t, cert.Domains[0], "*")
}

// TestTlsCertificateValidity tests TLS certificate validity period
func TestTlsCertificateValidity(t *testing.T) {
	cert := model.TlsCertificate{
		Name:      "test-cert",
		Domains:   []string{"example.com"},
		Cert:      "-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----",
		Key:       "-----BEGIN PRIVATE KEY-----\nMIIE...\n-----END PRIVATE KEY-----",
		ValidFrom: "2024-01-01T00:00:00Z",
		ValidTo:   "2025-01-01T00:00:00Z",
	}

	assert.Equal(t, "2024-01-01T00:00:00Z", cert.ValidFrom)
	assert.Equal(t, "2025-01-01T00:00:00Z", cert.ValidTo)
}

// TestTlsCertificatePagination tests pagination logic for TLS certificates
func TestTlsCertificatePagination(t *testing.T) {
	certs := make([]model.TlsCertificate, 15)
	for i := 0; i < 15; i++ {
		certs[i] = model.TlsCertificate{
			Name:    "cert-" + string(rune('a'+i)),
			Domains: []string{"example.com"},
		}
	}

	total := len(certs)
	pageNum := 1
	pageSize := 10
	start := (pageNum - 1) * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}
	pagedData := certs[start:end]
	result := model.NewPaginatedResult(pagedData, total, pageNum, pageSize)

	assert.Equal(t, 10, len(result.Data))
	assert.Equal(t, 15, result.Total)
	assert.Equal(t, 2, result.TotalPages)
}
