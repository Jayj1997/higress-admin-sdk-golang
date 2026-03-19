// Package model provides data models for Higress Admin SDK.
package model

// Service represents a backend service.
type Service struct {
	// Name is the service name.
	Name string `json:"name,omitempty"`

	// Namespace is the service namespace.
	Namespace string `json:"namespace,omitempty"`

	// Ports are the service ports.
	Ports []ServicePort `json:"ports,omitempty"`

	// Labels are the service labels.
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations are the service annotations.
	Annotations map[string]string `json:"annotations,omitempty"`
}

// ServicePort represents a service port.
type ServicePort struct {
	// Name is the port name.
	Name string `json:"name,omitempty"`

	// Port is the port number.
	Port int `json:"port,omitempty"`

	// Protocol is the port protocol (TCP, UDP).
	Protocol string `json:"protocol,omitempty"`

	// TargetPort is the target port on the pods.
	TargetPort int `json:"targetPort,omitempty"`
}
