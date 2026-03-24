// Package model provides data models for the SDK
package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_Structure(t *testing.T) {
	service := &Service{
		Name:      "my-service",
		Namespace: "default",
		Port:      8080,
		Ports: []ServicePort{
			{Name: "http", Port: 8080, Protocol: "TCP"},
			{Name: "grpc", Port: 9090, Protocol: "TCP"},
		},
		Endpoints: []string{"10.0.0.1:8080", "10.0.0.2:8080"},
		Labels: map[string]string{
			"app":  "my-app",
			"tier": "backend",
		},
		Annotations: map[string]string{
			"description": "test service",
		},
	}

	assert.Equal(t, "my-service", service.Name)
	assert.Equal(t, "default", service.Namespace)
	assert.Equal(t, 8080, service.Port)
	assert.Len(t, service.Ports, 2)
	assert.Len(t, service.Endpoints, 2)
	assert.Len(t, service.Labels, 2)
	assert.Len(t, service.Annotations, 1)
}

func TestService_EmptyFields(t *testing.T) {
	service := &Service{}

	assert.Empty(t, service.Name)
	assert.Empty(t, service.Namespace)
	assert.Equal(t, 0, service.Port)
	assert.Nil(t, service.Ports)
	assert.Nil(t, service.Endpoints)
	assert.Nil(t, service.Labels)
	assert.Nil(t, service.Annotations)
}

func TestService_SinglePort(t *testing.T) {
	service := &Service{
		Name:      "simple-service",
		Namespace: "production",
		Port:      80,
	}

	assert.Equal(t, "simple-service", service.Name)
	assert.Equal(t, "production", service.Namespace)
	assert.Equal(t, 80, service.Port)
}

func TestService_MultiplePorts(t *testing.T) {
	service := &Service{
		Name: "multi-port-service",
		Ports: []ServicePort{
			{Name: "http", Port: 80, Protocol: "TCP", TargetPort: 8080},
			{Name: "https", Port: 443, Protocol: "TCP", TargetPort: 8443},
			{Name: "metrics", Port: 9090, Protocol: "TCP", TargetPort: 9090},
		},
	}

	assert.Len(t, service.Ports, 3)
	assert.Equal(t, "http", service.Ports[0].Name)
	assert.Equal(t, 80, service.Ports[0].Port)
	assert.Equal(t, "https", service.Ports[1].Name)
	assert.Equal(t, 443, service.Ports[1].Port)
	assert.Equal(t, "metrics", service.Ports[2].Name)
	assert.Equal(t, 9090, service.Ports[2].Port)
}

func TestService_WithEndpoints(t *testing.T) {
	service := &Service{
		Name: "service-with-endpoints",
		Endpoints: []string{
			"192.168.1.1:8080",
			"192.168.1.2:8080",
			"192.168.1.3:8080",
		},
	}

	assert.Len(t, service.Endpoints, 3)
	assert.Equal(t, "192.168.1.1:8080", service.Endpoints[0])
}

func TestServicePort_Structure(t *testing.T) {
	port := ServicePort{
		Name:       "http",
		Port:       80,
		Protocol:   "TCP",
		TargetPort: 8080,
	}

	assert.Equal(t, "http", port.Name)
	assert.Equal(t, 80, port.Port)
	assert.Equal(t, "TCP", port.Protocol)
	assert.Equal(t, 8080, port.TargetPort)
}

func TestServicePort_EmptyFields(t *testing.T) {
	port := ServicePort{}

	assert.Empty(t, port.Name)
	assert.Equal(t, 0, port.Port)
	assert.Empty(t, port.Protocol)
	assert.Equal(t, 0, port.TargetPort)
}

func TestServicePort_DifferentProtocols(t *testing.T) {
	tests := []struct {
		name     string
		protocol string
	}{
		{"TCP protocol", "TCP"},
		{"UDP protocol", "UDP"},
		{"SCTP protocol", "SCTP"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			port := ServicePort{
				Name:     "test-port",
				Port:     8080,
				Protocol: tt.protocol,
			}
			assert.Equal(t, tt.protocol, port.Protocol)
		})
	}
}

func TestService_Labels(t *testing.T) {
	service := &Service{
		Name: "labeled-service",
		Labels: map[string]string{
			"app.kubernetes.io/name":    "my-app",
			"app.kubernetes.io/version": "v1.0.0",
			"environment":               "production",
		},
	}

	assert.Len(t, service.Labels, 3)
	assert.Equal(t, "my-app", service.Labels["app.kubernetes.io/name"])
	assert.Equal(t, "v1.0.0", service.Labels["app.kubernetes.io/version"])
	assert.Equal(t, "production", service.Labels["environment"])
}

func TestService_Annotations(t *testing.T) {
	service := &Service{
		Name: "annotated-service",
		Annotations: map[string]string{
			"prometheus.io/scrape": "true",
			"prometheus.io/port":   "9090",
		},
	}

	assert.Len(t, service.Annotations, 2)
	assert.Equal(t, "true", service.Annotations["prometheus.io/scrape"])
	assert.Equal(t, "9090", service.Annotations["prometheus.io/port"])
}

func TestService_FullExample(t *testing.T) {
	service := &Service{
		Name:      "api-gateway",
		Namespace: "infrastructure",
		Port:      443,
		Ports: []ServicePort{
			{Name: "https", Port: 443, Protocol: "TCP", TargetPort: 8443},
			{Name: "admin", Port: 8080, Protocol: "TCP", TargetPort: 8080},
		},
		Endpoints: []string{
			"10.0.0.1:8443",
			"10.0.0.2:8443",
		},
		Labels: map[string]string{
			"app":  "api-gateway",
			"team": "platform",
		},
		Annotations: map[string]string{
			"service.beta.kubernetes.io/aws-load-balancer-type": "nlb",
		},
	}

	// Verify all fields
	assert.Equal(t, "api-gateway", service.Name)
	assert.Equal(t, "infrastructure", service.Namespace)
	assert.Equal(t, 443, service.Port)
	assert.Len(t, service.Ports, 2)
	assert.Len(t, service.Endpoints, 2)
	assert.Len(t, service.Labels, 2)
	assert.Len(t, service.Annotations, 1)
}
