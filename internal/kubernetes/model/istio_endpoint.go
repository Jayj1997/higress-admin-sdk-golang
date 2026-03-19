// Package model contains internal Kubernetes models
package model

// IstioEndpoint represents an Istio endpoint
type IstioEndpoint struct {
	// Address is the endpoint address
	Address string `json:"address,omitempty"`
	// Port is the endpoint port
	Port uint32 `json:"port,omitempty"`
	// ServicePortName is the service port name
	ServicePortName string `json:"servicePortName,omitempty"`
	// ServicePort is the service port
	ServicePort uint32 `json:"servicePort,omitempty"`
	// ServiceName is the service name
	ServiceName string `json:"serviceName,omitempty"`
	// ServiceNamespace is the service namespace
	ServiceNamespace string `json:"serviceNamespace,omitempty"`
	// Labels are the endpoint labels
	Labels map[string]string `json:"labels,omitempty"`
	// UID is the endpoint UID
	UID string `json:"uid,omitempty"`
	// ClusterID is the cluster ID
	ClusterID string `json:"clusterId,omitempty"`
	// Network is the network
	Network string `json:"network,omitempty"`
	// Locality is the locality
	Locality string `json:"locality,omitempty"`
	// EndpointPort is the endpoint port
	EndpointPort uint32 `json:"endpointPort,omitempty"`
}

// IstioEndpointShard represents an Istio endpoint shard
type IstioEndpointShard struct {
	// Hostname is the hostname
	Hostname string `json:"hostname,omitempty"`
	// ClusterID is the cluster ID
	ClusterID string `json:"clusterId,omitempty"`
	// Ports is the ports map
	Ports map[string]*Port `json:"ports,omitempty"`
	// Endpoints is the endpoints
	Endpoints []*IstioEndpoint `json:"endpoints,omitempty"`
	// Service is the service
	Service string `json:"service,omitempty"`
	// Namespace is the namespace
	Namespace string `json:"namespace,omitempty"`
}

// Port represents a port definition
type Port struct {
	// Name is the port name
	Name string `json:"name,omitempty"`
	// Port is the port number
	Port uint32 `json:"port,omitempty"`
	// Protocol is the port protocol
	Protocol string `json:"protocol,omitempty"`
}

// RegistryzService represents a service from the registry
type RegistryzService struct {
	// Hostname is the service hostname
	Hostname string `json:"hostname,omitempty"`
	// ClusterVIPs is the cluster VIPs
	ClusterVIPs map[string][]string `json:"clusterVIPs,omitempty"`
	// Ports is the service ports
	Ports []*Port `json:"ports,omitempty"`
	// Attributes is the service attributes
	Attributes *RegistryzServiceAttributes `json:"attributes,omitempty"`
}

// RegistryzServiceAttributes represents service attributes
type RegistryzServiceAttributes struct {
	// Name is the service name
	Name string `json:"name,omitempty"`
	// Namespace is the service namespace
	Namespace string `json:"namespace,omitempty"`
	// Labels are the service labels
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations are the service annotations
	Annotations map[string]string `json:"annotations,omitempty"`
	// ServiceRegistry is the service registry
	ServiceRegistry string `json:"serviceRegistry,omitempty"`
	// ClusterID is the cluster ID
	ClusterID string `json:"clusterId,omitempty"`
	// Locality is the locality
	Locality string `json:"locality,omitempty"`
	// MeshExternal indicates if the service is external to the mesh
	MeshExternal bool `json:"meshExternal,omitempty"`
}
