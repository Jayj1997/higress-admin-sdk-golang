// Package model provides data models for Higress Admin SDK.
package model

import (
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/errors"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model/route"
)

// Route represents a gateway route configuration.
type Route struct {
	// Name is the route name.
	Name string `json:"name,omitempty"`

	// Version is the route version, required when updating.
	Version string `json:"version,omitempty"`

	// Domains that the route applies to.
	// If empty, the route applies to all domains.
	Domains []string `json:"domains,omitempty"`

	// Path is the path predicate.
	Path *route.RoutePredicate `json:"path,omitempty"`

	// Methods that the route applies to.
	// If empty, the route applies to all methods.
	Methods []string `json:"methods,omitempty"`

	// Headers are the header predicates.
	Headers []*route.KeyedRoutePredicate `json:"headers,omitempty"`

	// URLParams are the URL parameter predicates.
	URLParams []*route.KeyedRoutePredicate `json:"urlParams,omitempty"`

	// Services are the route upstream services.
	Services []*route.UpstreamService `json:"services,omitempty"`

	// Mock is the mock configuration (not supported yet).
	Mock *route.MockConfig `json:"mock,omitempty"`

	// Redirect is the redirect configuration (not supported yet).
	Redirect *route.RedirectConfig `json:"redirect,omitempty"`

	// RateLimit is the rate limit configuration (not supported yet).
	RateLimit *route.RateLimitConfig `json:"rateLimit,omitempty"`

	// Rewrite is the request rewrite configuration.
	Rewrite *route.RewriteConfig `json:"rewrite,omitempty"`

	// Timeout is the route timeout (not supported yet).
	Timeout string `json:"timeout,omitempty"`

	// ProxyNextUpstream is the proxy next upstream configuration.
	ProxyNextUpstream *route.ProxyNextUpstreamConfig `json:"proxyNextUpstream,omitempty"`

	// CORS is the CORS configuration.
	CORS *route.CorsConfig `json:"cors,omitempty"`

	// HeaderControl is the header control configuration.
	HeaderControl *route.HeaderControlConfig `json:"headerControl,omitempty"`

	// AuthConfig is the route auth configuration.
	AuthConfig *RouteAuthConfig `json:"authConfig,omitempty"`

	// CustomConfigs are custom configurations, e.g., custom annotations.
	CustomConfigs map[string]string `json:"customConfigs,omitempty"`

	// CustomLabels are custom labels.
	CustomLabels map[string]string `json:"customLabels,omitempty"`

	// Readonly indicates whether the route is read-only.
	Readonly *bool `json:"readonly,omitempty"`
}

// Validate validates the Route model.
func (r *Route) Validate() error {
	if r.Name == "" {
		return errors.NewValidationError("name cannot be blank")
	}
	if len(r.Services) == 0 {
		return errors.NewValidationError("services cannot be empty")
	}
	if r.Path != nil {
		if err := r.Path.Validate(); err != nil {
			return err
		}
	}
	for _, h := range r.Headers {
		if err := h.Validate(); err != nil {
			return err
		}
	}
	for _, u := range r.URLParams {
		if err := u.Validate(); err != nil {
			return err
		}
	}
	for _, s := range r.Services {
		if err := s.Validate(); err != nil {
			return err
		}
	}
	if r.AuthConfig != nil {
		if err := r.AuthConfig.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// RouteAuthConfig represents route authentication configuration.
type RouteAuthConfig struct {
	// Enabled indicates whether authentication is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// AllowedConsumers are the allowed consumer names.
	// If empty, all consumers are allowed.
	AllowedConsumers []string `json:"allowedConsumers,omitempty"`

	// AllowedCredentialTypes are the allowed credential types.
	AllowedCredentialTypes []string `json:"allowedCredentialTypes,omitempty"`
}

// Validate validates the RouteAuthConfig model.
func (c *RouteAuthConfig) Validate() error {
	return nil
}
