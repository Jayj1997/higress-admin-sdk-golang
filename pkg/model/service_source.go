// Package model provides data models for Higress Admin SDK.
package model

import (
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/errors"
)

// Service source type constants.
const (
	ServiceSourceTypeDNS    = "dns"
	ServiceSourceTypeStatic = "static"
)

// ServiceSource represents a service discovery source.
type ServiceSource struct {
	// Name is the service source name.
	Name string `json:"name,omitempty"`

	// Version is the service source version, required when updating.
	Version string `json:"version,omitempty"`

	// Type is the service source type.
	// Valid values: "dns", "static", and registry types like "nacos", "zookeeper", "eureka", "consul"
	Type string `json:"type,omitempty"`

	// Domain is the service source domain(s).
	// For static type, use ip:port format.
	// For dns type, use domain list.
	Domain string `json:"domain,omitempty"`

	// Port is the service port.
	Port *int `json:"port,omitempty"`

	// Namespace is the namespace for the service source.
	Namespace string `json:"namespace,omitempty"`

	// Group is the service group (for Nacos).
	Group string `json:"group,omitempty"`

	// Services is the list of services to discover.
	Services []string `json:"services,omitempty"`

	// AuthEnabled indicates whether authentication is enabled.
	AuthEnabled *bool `json:"authEnabled,omitempty"`

	// AuthType is the authentication type.
	AuthType string `json:"authType,omitempty"`

	// AuthSecretName is the name of the secret containing auth credentials.
	AuthSecretName string `json:"authSecretName,omitempty"`

	// Properties are additional properties for the service source.
	Properties map[string]string `json:"properties,omitempty"`
}

// Validate validates the ServiceSource model.
func (s *ServiceSource) Validate() error {
	if s.Name == "" {
		return errors.NewValidationError("name is required")
	}
	if s.Type == "" {
		return errors.NewValidationError("type is required")
	}
	return nil
}
