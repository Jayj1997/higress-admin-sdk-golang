// Package mcp contains MCP (Model Context Protocol) CRD types for Kubernetes
package mcp

// V1McpBridge represents an MCP Bridge CRD
type V1McpBridge struct {
	// TypeMeta contains standard Kubernetes type metadata
	TypeMeta `json:",inline"`
	// Metadata contains object metadata
	Metadata *V1ObjectMeta `json:"metadata,omitempty"`
	// Spec contains the McpBridge specification
	Spec *V1McpBridgeSpec `json:"spec,omitempty"`
}

// TypeMeta contains standard Kubernetes type metadata
type TypeMeta struct {
	// APIVersion is the API version (format: group/version)
	APIVersion string `json:"apiVersion,omitempty"`
	// Kind is the resource kind
	Kind string `json:"kind,omitempty"`
}

// Constants for McpBridge
const (
	McpBridgeAPIGroup   = "networking.higress.io"
	McpBridgeAPIVersion = "networking.higress.io/v1"
	McpBridgeKind       = "McpBridge"
	McpBridgePlural     = "mcpbridges"
)

// NewV1McpBridge creates a new McpBridge
func NewV1McpBridge() *V1McpBridge {
	return &V1McpBridge{
		TypeMeta: TypeMeta{
			APIVersion: McpBridgeAPIVersion,
			Kind:       McpBridgeKind,
		},
		Metadata: &V1ObjectMeta{},
		Spec:     &V1McpBridgeSpec{},
	}
}

// V1ObjectMeta represents object metadata
type V1ObjectMeta struct {
	// Name is the resource name
	Name string `json:"name,omitempty"`
	// Namespace is the resource namespace
	Namespace string `json:"namespace,omitempty"`
	// Labels are resource labels
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations are resource annotations
	Annotations map[string]string `json:"annotations,omitempty"`
	// ResourceVersion is the resource version
	ResourceVersion string `json:"resourceVersion,omitempty"`
	// UID is the resource UID
	UID string `json:"uid,omitempty"`
	// CreationTimestamp is the creation time
	CreationTimestamp string `json:"creationTimestamp,omitempty"`
	// Generation is the generation number
	Generation int64 `json:"generation,omitempty"`
}

// V1McpBridgeSpec represents the McpBridge specification
type V1McpBridgeSpec struct {
	// Registries is the list of registry configurations
	Registries []*V1RegistryConfig `json:"registries,omitempty"`
}

// V1RegistryConfig represents a registry configuration
type V1RegistryConfig struct {
	// Name is the registry name
	Name string `json:"name,omitempty"`
	// Type is the registry type (nacos, eureka, consul, dns, static, etc.)
	Type string `json:"type,omitempty"`
	// Domain is the service domain
	Domain string `json:"domain,omitempty"`
	// Port is the service port
	Port uint32 `json:"port,omitempty"`
	// NacosNamespace is the Nacos namespace
	NacosNamespace string `json:"nacosNamespace,omitempty"`
	// NacosGroups is the Nacos groups
	NacosGroups []string `json:"nacosGroups,omitempty"`
	// NacosUsername is the Nacos username
	NacosUsername string `json:"nacosUsername,omitempty"`
	// NacosPassword is the Nacos password
	NacosPassword string `json:"nacosPassword,omitempty"`
	// NacosAccessKey is the Nacos access key
	NacosAccessKey string `json:"nacosAccessKey,omitempty"`
	// NacosSecretKey is the Nacos secret key
	NacosSecretKey string `json:"nacosSecretKey,omitempty"`
	// ConsulNamespace is the Consul namespace
	ConsulNamespace string `json:"consulNamespace,omitempty"`
	// ConsulDatacenter is the Consul datacenter
	ConsulDatacenter string `json:"consulDatacenter,omitempty"`
	// ConsulServiceDomain is the Consul service domain
	ConsulServiceDomain string `json:"consulServiceDomain,omitempty"`
	// EurekaClientServiceUrl is the Eureka client service URL
	EurekaClientServiceUrl string `json:"eurekaClientServiceUrl,omitempty"`
	// DnsDomain is the DNS domain
	DnsDomain string `json:"dnsDomain,omitempty"`
	// DnsPort is the DNS port
	DnsPort uint32 `json:"dnsPort,omitempty"`
	// StaticAddresses is the static addresses
	StaticAddresses []string `json:"staticAddresses,omitempty"`
	// AuthSecretName is the authentication secret name
	AuthSecretName string `json:"authSecretName,omitempty"`
	// AuthSecretNamespace is the authentication secret namespace
	AuthSecretNamespace string `json:"authSecretNamespace,omitempty"`
	// Metadata contains additional registry metadata
	Metadata *V1RegistryConfigMetadata `json:"metadata,omitempty"`
	// Proxy represents proxy configuration
	Proxy *V1ProxyConfig `json:"proxy,omitempty"`
	// EnableServiceDiscovery indicates if service discovery is enabled
	EnableServiceDiscovery bool `json:"enableServiceDiscovery,omitempty"`
	// EnableMcpServer indicates if MCP server is enabled
	EnableMcpServer bool `json:"enableMcpServer,omitempty"`
	// McpServerName is the MCP server name
	McpServerName string `json:"mcpServerName,omitempty"`
	// McpServerDescription is the MCP server description
	McpServerDescription string `json:"mcpServerDescription,omitempty"`
}

// V1RegistryConfigMetadata represents registry configuration metadata
type V1RegistryConfigMetadata struct {
	// Description is the registry description
	Description string `json:"description,omitempty"`
	// Owner is the registry owner
	Owner string `json:"owner,omitempty"`
	// CustomMetadata contains custom metadata
	CustomMetadata map[string]string `json:"customMetadata,omitempty"`
}

// V1ProxyConfig represents proxy configuration
type V1ProxyConfig struct {
	// Enabled indicates if proxy is enabled
	Enabled bool `json:"enabled,omitempty"`
	// Host is the proxy host
	Host string `json:"host,omitempty"`
	// Port is the proxy port
	Port uint32 `json:"port,omitempty"`
	// Type is the proxy type (http, https, socks5)
	Type string `json:"type,omitempty"`
	// Username is the proxy username
	Username string `json:"username,omitempty"`
	// Password is the proxy password
	Password string `json:"password,omitempty"`
}

// VPort represents a port configuration
type VPort struct {
	// Number is the port number
	Number uint32 `json:"number,omitempty"`
	// Protocol is the port protocol
	Protocol string `json:"protocol,omitempty"`
	// Name is the port name
	Name string `json:"name,omitempty"`
	// TargetPort is the target port
	TargetPort uint32 `json:"targetPort,omitempty"`
}

// V1McpBridgeList represents a list of McpBridges
type V1McpBridgeList struct {
	// APIVersion is the API version
	APIVersion string `json:"apiVersion,omitempty"`
	// Kind is the resource kind
	Kind string `json:"kind,omitempty"`
	// Metadata contains list metadata
	Metadata *V1ListMeta `json:"metadata,omitempty"`
	// Items is the list of McpBridges
	Items []*V1McpBridge `json:"items,omitempty"`
}

// V1ListMeta represents list metadata
type V1ListMeta struct {
	// ResourceVersion is the resource version
	ResourceVersion string `json:"resourceVersion,omitempty"`
	// SelfLink is the self link
	SelfLink string `json:"selfLink,omitempty"`
	// Continue is the continue token
	Continue string `json:"continue,omitempty"`
	// RemainingItemCount is the remaining item count
	RemainingItemCount *int64 `json:"remainingItemCount,omitempty"`
}
