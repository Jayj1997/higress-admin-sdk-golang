// Package model provides data models for Higress Admin SDK.
package model

// ProxyServer represents a proxy server configuration.
type ProxyServer struct {
	// Name is the proxy server name.
	Name string `json:"name,omitempty"`

	// Host is the proxy server host.
	Host string `json:"host,omitempty"`

	// Port is the proxy server port.
	Port int `json:"port,omitempty"`

	// Protocol is the proxy protocol.
	Protocol string `json:"protocol,omitempty"`

	// Version is the resource version.
	Version string `json:"version,omitempty"`
}
