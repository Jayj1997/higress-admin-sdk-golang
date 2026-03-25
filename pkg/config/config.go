// Package config provides configuration for Higress Admin SDK.
package config

import "github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/constant"

// HigressServiceConfig is the configuration for Higress service provider.
type HigressServiceConfig struct {
	// KubeConfigPath is the path to the kubeconfig file.
	KubeConfigPath string

	// KubeConfigContent is the content of kubeconfig as a string.
	KubeConfigContent string

	// ControllerNamespace is the namespace where Higress controller is deployed.
	ControllerNamespace string

	// ControllerWatchedNamespace is the namespace that the controller watches.
	// If empty, the controller watches all namespaces.
	ControllerWatchedNamespace string

	// ControllerWatchedIngressClassName is the ingress class name that the controller watches.
	ControllerWatchedIngressClassName string

	// ControllerServiceName is the name of the controller service.
	ControllerServiceName string

	// ControllerServiceHost is the host of the controller service.
	ControllerServiceHost string

	// ControllerServicePort is the port of the controller service.
	ControllerServicePort int

	// ControllerJwtPolicy is the JWT policy for controller authentication.
	ControllerJwtPolicy string

	// ControllerAccessToken is the access token for controller authentication.
	ControllerAccessToken string

	// WasmPluginServiceConfig is the configuration for WASM plugin service.
	WasmPluginServiceConfig *WasmPluginServiceConfig

	// ServiceListSupportRegistry indicates whether the service list interface supports registry.
	// If true, the service list interface will support registry and depend on the controller API.
	// If false, the service list implementation will not support registry and will interact with the API server directly.
	ServiceListSupportRegistry bool

	// ClusterDomainSuffix is the cluster domain suffix.
	ClusterDomainSuffix string
}

// WasmPluginServiceConfig is the configuration for WASM plugin service.
type WasmPluginServiceConfig struct {
	// ImageRegistry is the registry for WASM plugin images.
	ImageRegistry string

	// ImagePullPolicy is the pull policy for WASM plugin images.
	ImagePullPolicy string

	// CustomImageUrlPattern is the pattern for custom image URLs.
	// Supported placeholders: ${name}, ${version}
	CustomImageUrlPattern string

	// ImagePullSecret is the secret for pulling WASM plugin images.
	ImagePullSecret string
}

// Option is a function that modifies the configuration.
type Option func(*HigressServiceConfig)

// NewHigressServiceConfig creates a new HigressServiceConfig with the given options.
func NewHigressServiceConfig(opts ...Option) *HigressServiceConfig {
	cfg := &HigressServiceConfig{
		ControllerNamespace:               constant.NSDefault,
		ControllerWatchedIngressClassName: constant.ControllerIngressClassNameDefault,
		ControllerServiceName:             constant.ControllerServiceNameDefault,
		ControllerServiceHost:             constant.ControllerServiceHostDefault,
		ControllerServicePort:             constant.ControllerServicePortDefault,
		ControllerJwtPolicy:               constant.ControllerJwtPolicyDefault,
		WasmPluginServiceConfig:           &WasmPluginServiceConfig{},
		ServiceListSupportRegistry:        constant.ServiceListSupportRegistryDefault,
		ClusterDomainSuffix:               constant.ClusterDomainSuffixDefault,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

// WithKubeConfigPath sets the kubeconfig path.
func WithKubeConfigPath(path string) Option {
	return func(c *HigressServiceConfig) {
		c.KubeConfigPath = path
	}
}

// WithKubeConfigContent sets the kubeconfig content.
func WithKubeConfigContent(content string) Option {
	return func(c *HigressServiceConfig) {
		c.KubeConfigContent = content
	}
}

// WithControllerNamespace sets the controller namespace.
func WithControllerNamespace(ns string) Option {
	return func(c *HigressServiceConfig) {
		c.ControllerNamespace = ns
	}
}

// WithControllerWatchedNamespace sets the controller watched namespace.
func WithControllerWatchedNamespace(ns string) Option {
	return func(c *HigressServiceConfig) {
		c.ControllerWatchedNamespace = ns
	}
}

// WithControllerWatchedIngressClassName sets the controller watched ingress class name.
func WithControllerWatchedIngressClassName(name string) Option {
	return func(c *HigressServiceConfig) {
		c.ControllerWatchedIngressClassName = name
	}
}

// WithControllerServiceName sets the controller service name.
func WithControllerServiceName(name string) Option {
	return func(c *HigressServiceConfig) {
		c.ControllerServiceName = name
	}
}

// WithControllerServiceHost sets the controller service host.
func WithControllerServiceHost(host string) Option {
	return func(c *HigressServiceConfig) {
		c.ControllerServiceHost = host
	}
}

// WithControllerServicePort sets the controller service port.
func WithControllerServicePort(port int) Option {
	return func(c *HigressServiceConfig) {
		c.ControllerServicePort = port
	}
}

// WithControllerJwtPolicy sets the controller JWT policy.
func WithControllerJwtPolicy(policy string) Option {
	return func(c *HigressServiceConfig) {
		c.ControllerJwtPolicy = policy
	}
}

// WithControllerAccessToken sets the controller access token.
func WithControllerAccessToken(token string) Option {
	return func(c *HigressServiceConfig) {
		c.ControllerAccessToken = token
	}
}

// WithWasmPluginServiceConfig sets the WASM plugin service config.
func WithWasmPluginServiceConfig(wasmCfg *WasmPluginServiceConfig) Option {
	return func(c *HigressServiceConfig) {
		c.WasmPluginServiceConfig = wasmCfg
	}
}

// WithServiceListSupportRegistry sets whether service list supports registry.
func WithServiceListSupportRegistry(support bool) Option {
	return func(c *HigressServiceConfig) {
		c.ServiceListSupportRegistry = support
	}
}

// WithClusterDomainSuffix sets the cluster domain suffix.
func WithClusterDomainSuffix(suffix string) Option {
	return func(c *HigressServiceConfig) {
		c.ClusterDomainSuffix = suffix
	}
}

// Getter methods for HigressServiceConfig

// GetKubeConfigPath returns the kubeconfig path.
func (c *HigressServiceConfig) GetKubeConfigPath() string {
	return c.KubeConfigPath
}

// GetKubeConfigContent returns the kubeconfig content.
func (c *HigressServiceConfig) GetKubeConfigContent() string {
	return c.KubeConfigContent
}

// GetControllerNamespace returns the controller namespace.
func (c *HigressServiceConfig) GetControllerNamespace() string {
	return c.ControllerNamespace
}

// GetControllerWatchedNamespace returns the controller watched namespace.
func (c *HigressServiceConfig) GetControllerWatchedNamespace() string {
	return c.ControllerWatchedNamespace
}

// GetControllerWatchedIngressClassName returns the controller watched ingress class name.
func (c *HigressServiceConfig) GetControllerWatchedIngressClassName() string {
	return c.ControllerWatchedIngressClassName
}

// GetControllerServiceName returns the controller service name.
func (c *HigressServiceConfig) GetControllerServiceName() string {
	return c.ControllerServiceName
}

// GetControllerServiceHost returns the controller service host.
func (c *HigressServiceConfig) GetControllerServiceHost() string {
	return c.ControllerServiceHost
}

// GetControllerServicePort returns the controller service port.
func (c *HigressServiceConfig) GetControllerServicePort() int {
	return c.ControllerServicePort
}

// GetControllerJwtPolicy returns the controller JWT policy.
func (c *HigressServiceConfig) GetControllerJwtPolicy() string {
	return c.ControllerJwtPolicy
}

// GetControllerAccessToken returns the controller access token.
func (c *HigressServiceConfig) GetControllerAccessToken() string {
	return c.ControllerAccessToken
}

// GetWasmPluginServiceConfig returns the WASM plugin service config.
func (c *HigressServiceConfig) GetWasmPluginServiceConfig() *WasmPluginServiceConfig {
	return c.WasmPluginServiceConfig
}

// GetServiceListSupportRegistry returns whether service list supports registry.
func (c *HigressServiceConfig) GetServiceListSupportRegistry() bool {
	return c.ServiceListSupportRegistry
}

// GetClusterDomainSuffix returns the cluster domain suffix.
func (c *HigressServiceConfig) GetClusterDomainSuffix() string {
	return c.ClusterDomainSuffix
}
