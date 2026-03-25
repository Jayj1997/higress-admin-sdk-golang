// Package model provides data models for Higress Admin SDK.
package model

import "github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/errors"

// EnableHTTPS constants for domain HTTPS configuration.
const (
	// EnableHTTPSOff disables HTTPS for the domain
	EnableHTTPSOff = "off"
	// EnableHTTPSOn enables HTTPS for the domain
	EnableHTTPSOn = "on"
	// EnableHTTPSForce forces HTTPS redirect for the domain
	EnableHTTPSForce = "force"
)

// Domain represents a gateway domain configuration.
//
// Domain is used to configure the domain names that the gateway will handle,
// including HTTPS settings and certificate bindings.
type Domain struct {
	// Name is the domain name (e.g., "example.com", "*.example.com")
	Name string `json:"name,omitempty"`

	// Version is the domain version, required when updating.
	// This is typically used for optimistic concurrency control.
	Version string `json:"version,omitempty"`

	// EnableHTTPS is the HTTPS configuration.
	// Valid values: "off", "on", "force"
	// - off: HTTP only
	// - on: HTTP and HTTPS both enabled
	// - force: HTTPS only, HTTP requests will be redirected to HTTPS
	EnableHTTPS string `json:"enableHttps,omitempty"`

	// CertIdentifier is the certificate name/identifier for HTTPS.
	// Required when EnableHTTPS is "on" or "force".
	CertIdentifier string `json:"certIdentifier,omitempty"`
}

// Validate validates the Domain model.
func (d *Domain) Validate() error {
	if d.Name == "" {
		return errors.NewValidationError("domain name is required")
	}
	return nil
}

// IsHTTPS returns true if HTTPS is enabled for the domain.
func (d *Domain) IsHTTPS() bool {
	return d.EnableHTTPS == EnableHTTPSOn || d.EnableHTTPS == EnableHTTPSForce
}

// IsForceHTTPS returns true if HTTPS is forced for the domain.
func (d *Domain) IsForceHTTPS() bool {
	return d.EnableHTTPS == EnableHTTPSForce
}
