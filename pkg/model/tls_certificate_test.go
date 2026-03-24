// Package model provides data models for the SDK
package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTlsCertificate_Validate(t *testing.T) {
	tests := []struct {
		name        string
		cert        *TlsCertificate
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid certificate",
			cert: &TlsCertificate{
				Name: "my-cert",
				Cert: "-----BEGIN CERTIFICATE-----\nMIID...\n-----END CERTIFICATE-----",
				Key:  "-----BEGIN PRIVATE KEY-----\nMIIE...\n-----END PRIVATE KEY-----",
			},
			expectError: false,
		},
		{
			name: "valid with all fields",
			cert: &TlsCertificate{
				Name:      "my-cert",
				Version:   "v1",
				Domains:   []string{"example.com", "www.example.com"},
				Cert:      "-----BEGIN CERTIFICATE-----\nMIID...\n-----END CERTIFICATE-----",
				Key:       "-----BEGIN PRIVATE KEY-----\nMIIE...\n-----END PRIVATE KEY-----",
				Namespace: "default",
				ValidFrom: "2024-01-01T00:00:00Z",
				ValidTo:   "2025-01-01T00:00:00Z",
			},
			expectError: false,
		},
		{
			name: "missing name",
			cert: &TlsCertificate{
				Cert: "-----BEGIN CERTIFICATE-----\nMIID...\n-----END CERTIFICATE-----",
				Key:  "-----BEGIN PRIVATE KEY-----\nMIIE...\n-----END PRIVATE KEY-----",
			},
			expectError: true,
			errorMsg:    "name is required",
		},
		{
			name: "missing cert",
			cert: &TlsCertificate{
				Name: "my-cert",
				Key:  "-----BEGIN PRIVATE KEY-----\nMIIE...\n-----END PRIVATE KEY-----",
			},
			expectError: true,
			errorMsg:    "cert is required",
		},
		{
			name: "missing key",
			cert: &TlsCertificate{
				Name: "my-cert",
				Cert: "-----BEGIN CERTIFICATE-----\nMIID...\n-----END CERTIFICATE-----",
			},
			expectError: true,
			errorMsg:    "key is required",
		},
		{
			name: "empty name",
			cert: &TlsCertificate{
				Name: "",
				Cert: "cert",
				Key:  "key",
			},
			expectError: true,
			errorMsg:    "name is required",
		},
		{
			name: "empty cert",
			cert: &TlsCertificate{
				Name: "my-cert",
				Cert: "",
				Key:  "key",
			},
			expectError: true,
			errorMsg:    "cert is required",
		},
		{
			name: "empty key",
			cert: &TlsCertificate{
				Name: "my-cert",
				Cert: "cert",
				Key:  "",
			},
			expectError: true,
			errorMsg:    "key is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cert.Validate()
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTlsCertificate_Structure(t *testing.T) {
	cert := &TlsCertificate{
		Name:      "example-com-cert",
		Version:   "v1",
		Domains:   []string{"example.com", "www.example.com", "api.example.com"},
		Cert:      "-----BEGIN CERTIFICATE-----\nMIID...\n-----END CERTIFICATE-----",
		Key:       "-----BEGIN PRIVATE KEY-----\nMIIE...\n-----END PRIVATE KEY-----",
		Namespace: "production",
		ValidFrom: "2024-01-01T00:00:00Z",
		ValidTo:   "2025-01-01T00:00:00Z",
	}

	assert.Equal(t, "example-com-cert", cert.Name)
	assert.Equal(t, "v1", cert.Version)
	assert.Len(t, cert.Domains, 3)
	assert.Contains(t, cert.Domains, "example.com")
	assert.Contains(t, cert.Domains, "www.example.com")
	assert.Contains(t, cert.Domains, "api.example.com")
	assert.NotEmpty(t, cert.Cert)
	assert.NotEmpty(t, cert.Key)
	assert.Equal(t, "production", cert.Namespace)
	assert.Equal(t, "2024-01-01T00:00:00Z", cert.ValidFrom)
	assert.Equal(t, "2025-01-01T00:00:00Z", cert.ValidTo)
}

func TestTlsCertificate_EmptyFields(t *testing.T) {
	cert := &TlsCertificate{}

	assert.Empty(t, cert.Name)
	assert.Empty(t, cert.Version)
	assert.Nil(t, cert.Domains)
	assert.Empty(t, cert.Cert)
	assert.Empty(t, cert.Key)
	assert.Empty(t, cert.Namespace)
	assert.Empty(t, cert.ValidFrom)
	assert.Empty(t, cert.ValidTo)
}

func TestTlsCertificate_SingleDomain(t *testing.T) {
	cert := &TlsCertificate{
		Name:    "single-domain-cert",
		Domains: []string{"example.com"},
		Cert:    "cert-content",
		Key:     "key-content",
	}

	assert.Len(t, cert.Domains, 1)
	assert.Equal(t, "example.com", cert.Domains[0])
}

func TestTlsCertificate_WildcardDomain(t *testing.T) {
	cert := &TlsCertificate{
		Name:    "wildcard-cert",
		Domains: []string{"*.example.com"},
		Cert:    "cert-content",
		Key:     "key-content",
	}

	assert.Len(t, cert.Domains, 1)
	assert.Equal(t, "*.example.com", cert.Domains[0])
}

func TestTlsCertificate_MultipleDomains(t *testing.T) {
	cert := &TlsCertificate{
		Name: "multi-domain-cert",
		Domains: []string{
			"example.com",
			"www.example.com",
			"api.example.com",
			"admin.example.com",
		},
		Cert: "cert-content",
		Key:  "key-content",
	}

	assert.Len(t, cert.Domains, 4)
}

func TestTlsCertificate_ValidityPeriod(t *testing.T) {
	cert := &TlsCertificate{
		Name:      "cert-with-validity",
		Cert:      "cert-content",
		Key:       "key-content",
		ValidFrom: "2024-01-01T00:00:00Z",
		ValidTo:   "2025-12-31T23:59:59Z",
	}

	assert.Equal(t, "2024-01-01T00:00:00Z", cert.ValidFrom)
	assert.Equal(t, "2025-12-31T23:59:59Z", cert.ValidTo)
}

func TestTlsCertificate_Namespace(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
	}{
		{"default namespace", "default"},
		{"production namespace", "production"},
		{"staging namespace", "staging"},
		{"custom namespace", "my-custom-ns"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cert := &TlsCertificate{
				Name:      "test-cert",
				Namespace: tt.namespace,
				Cert:      "cert",
				Key:       "key",
			}
			assert.Equal(t, tt.namespace, cert.Namespace)
		})
	}
}
