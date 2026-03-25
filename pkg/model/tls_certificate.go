// Package model provides data models for Higress Admin SDK.
package model

import (
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/errors"
)

// TlsCertificate represents a TLS/SSL certificate.
type TlsCertificate struct {
	// Name is the certificate name/identifier.
	Name string `json:"name,omitempty"`

	// Version is the certificate version, required when updating.
	Version string `json:"version,omitempty"`

	// Domains are the domain names this certificate is valid for.
	Domains []string `json:"domains,omitempty"`

	// Cert is the certificate content (PEM format).
	Cert string `json:"cert,omitempty"`

	// Key is the private key content (PEM format).
	Key string `json:"key,omitempty"`

	// Namespace is the namespace where the certificate is stored.
	Namespace string `json:"namespace,omitempty"`

	// ValidFrom is the certificate validity start time.
	ValidFrom string `json:"validFrom,omitempty"`

	// ValidTo is the certificate validity end time.
	ValidTo string `json:"validTo,omitempty"`
}

// Validate validates the TlsCertificate model.
func (t *TlsCertificate) Validate() error {
	if t.Name == "" {
		return errors.NewValidationError("name is required")
	}
	if t.Cert == "" {
		return errors.NewValidationError("cert is required")
	}
	if t.Key == "" {
		return errors.NewValidationError("key is required")
	}
	return nil
}
