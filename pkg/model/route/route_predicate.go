// Package route provides route-related models for Higress Admin SDK.
package route

import (
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/errors"
)

// MatchType constants for route predicates.
const (
	MatchTypeExact  = "exact"
	MatchTypePrefix = "prefix"
	MatchTypeRegex  = "regex"
)

// RoutePredicate represents a route path predicate.
type RoutePredicate struct {
	// MatchType is the type of path matching.
	// Valid values: "exact", "prefix", "regex"
	MatchType string `json:"matchType,omitempty"`

	// Path is the path value to match.
	Path string `json:"path,omitempty"`

	// CaseSensitive indicates whether the path matching is case sensitive.
	CaseSensitive *bool `json:"caseSensitive,omitempty"`
}

// Validate validates the RoutePredicate model.
func (r *RoutePredicate) Validate() error {
	if r.MatchType == "" {
		return errors.NewValidationError("matchType is required")
	}
	if r.Path == "" {
		return errors.NewValidationError("path is required")
	}
	return nil
}

// KeyedRoutePredicate represents a route predicate with a key (for headers, URL params, etc.).
type KeyedRoutePredicate struct {
	// Key is the name of the header or URL parameter.
	Key string `json:"key,omitempty"`

	// MatchType is the type of matching.
	// Valid values: "exact", "prefix", "regex"
	MatchType string `json:"matchType,omitempty"`

	// Value is the value to match.
	Value string `json:"value,omitempty"`

	// CaseSensitive indicates whether the matching is case sensitive.
	CaseSensitive *bool `json:"caseSensitive,omitempty"`
}

// Validate validates the KeyedRoutePredicate model.
func (k *KeyedRoutePredicate) Validate() error {
	if k.Key == "" {
		return errors.NewValidationError("key is required")
	}
	if k.MatchType == "" {
		return errors.NewValidationError("matchType is required")
	}
	return nil
}

// UpstreamService represents an upstream service for routing.
type UpstreamService struct {
	// Name is the service name.
	Name string `json:"name,omitempty"`

	// Namespace is the service namespace.
	Namespace string `json:"namespace,omitempty"`

	// Port is the service port.
	Port int `json:"port,omitempty"`

	// Weight is the traffic weight for this service (for load balancing).
	Weight *int `json:"weight,omitempty"`

	// Version is the service version (for canary deployments).
	Version string `json:"version,omitempty"`
}

// Validate validates the UpstreamService model.
func (u *UpstreamService) Validate() error {
	if u.Name == "" {
		return errors.NewValidationError("service name is required")
	}
	return nil
}

// RewriteConfig represents URL rewrite configuration.
type RewriteConfig struct {
	// Path is the path rewrite value.
	// Supports variables like $1, $2 for regex capture groups.
	Path string `json:"path,omitempty"`

	// Host is the host rewrite value.
	Host string `json:"host,omitempty"`
}

// CorsConfig represents CORS configuration.
type CorsConfig struct {
	// AllowOrigins is the list of allowed origins.
	AllowOrigins []string `json:"allowOrigins,omitempty"`

	// AllowMethods is the list of allowed HTTP methods.
	AllowMethods []string `json:"allowMethods,omitempty"`

	// AllowHeaders is the list of allowed request headers.
	AllowHeaders []string `json:"allowHeaders,omitempty"`

	// ExposeHeaders is the list of exposed response headers.
	ExposeHeaders []string `json:"exposeHeaders,omitempty"`

	// AllowCredentials indicates whether credentials are allowed.
	AllowCredentials *bool `json:"allowCredentials,omitempty"`

	// MaxAge is the max age for preflight cache (in seconds).
	MaxAge *int `json:"maxAge,omitempty"`
}

// HeaderControlConfig represents header control configuration.
type HeaderControlConfig struct {
	// Request headers to add/set.
	RequestAddHeaders map[string]string `json:"requestAddHeaders,omitempty"`

	// Request headers to remove.
	RequestRemoveHeaders []string `json:"requestRemoveHeaders,omitempty"`

	// Response headers to add/set.
	ResponseAddHeaders map[string]string `json:"responseAddHeaders,omitempty"`

	// Response headers to remove.
	ResponseRemoveHeaders []string `json:"responseRemoveHeaders,omitempty"`
}

// ProxyNextUpstreamConfig represents proxy next upstream configuration.
type ProxyNextUpstreamConfig struct {
	// RetryOn specifies the conditions for retrying.
	RetryOn string `json:"retryOn,omitempty"`

	// NumRetries is the number of retry attempts.
	NumRetries *int `json:"numRetries,omitempty"`

	// Timeout is the timeout for each retry attempt.
	Timeout string `json:"timeout,omitempty"`
}

// MockConfig represents mock configuration (not supported yet).
type MockConfig struct {
	// StatusCode is the mock response status code.
	StatusCode *int `json:"statusCode,omitempty"`

	// Body is the mock response body.
	Body string `json:"body,omitempty"`

	// Headers are the mock response headers.
	Headers map[string]string `json:"headers,omitempty"`
}

// RedirectConfig represents redirect configuration (not supported yet).
type RedirectConfig struct {
	// StatusCode is the redirect status code (301, 302, 307, 308).
	StatusCode *int `json:"statusCode,omitempty"`

	// Host is the redirect host.
	Host string `json:"host,omitempty"`

	// Path is the redirect path.
	Path string `json:"path,omitempty"`

	// Scheme is the redirect scheme (http, https).
	Scheme string `json:"scheme,omitempty"`
}

// RateLimitConfig represents rate limit configuration (not supported yet).
type RateLimitConfig struct {
	// LimitBy is the key to limit by (ip, header, etc.).
	LimitBy string `json:"limitBy,omitempty"`

	// Limit is the rate limit value.
	Limit *int `json:"limit,omitempty"`

	// Window is the time window for rate limiting.
	Window string `json:"window,omitempty"`
}
