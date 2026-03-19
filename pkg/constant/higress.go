// Package constant provides constants for Higress Admin SDK.
package constant

// Default values for Higress controller configuration.
const (
	// NSDefault is the default namespace for Higress controller.
	NSDefault = "higress-system"

	// ControllerIngressClassNameDefault is the default ingress class name.
	ControllerIngressClassNameDefault = "higress"

	// ControllerServiceNameDefault is the default controller service name.
	ControllerServiceNameDefault = "higress-controller"

	// ControllerServiceHostDefault is the default controller service host.
	ControllerServiceHostDefault = "higress-controller.higress-system.svc.cluster.local"

	// ControllerServicePortDefault is the default controller service port.
	ControllerServicePortDefault = 15014

	// ControllerJwtPolicyDefault is the default JWT policy.
	ControllerJwtPolicyDefault = "first-party-jwt"

	// ServiceListSupportRegistryDefault is the default value for service list support registry.
	ServiceListSupportRegistryDefault = true

	// ClusterDomainSuffixDefault is the default cluster domain suffix.
	ClusterDomainSuffixDefault = "cluster.local"
)

// Controller JWT policy constants.
const (
	JwtPolicyFirstParty    = "first-party-jwt"
	JwtPolicyThirdParty    = "third-party-jwt"
	JwtPolicyFirstPartyJwt = "first-party-jwt"
)

// Resource definer constants.
const (
	LabelResourceDefinerKey   = "higress.io/resource-definer"
	LabelResourceDefinerValue = "higress-console"
)

// WasmPlugin label keys.
const (
	LabelWasmPluginNameKey    = "higress.io/wasm-plugin-name"
	LabelWasmPluginVersionKey = "higress.io/wasm-plugin-version"
	LabelWasmPluginBuiltInKey = "higress.io/wasm-plugin-built-in"
)

// Kubernetes namespace constants.
const (
	KubeSystemNamespace = "kube-system"
)

// Common key constants.
const (
	// KeyDelimiter is the delimiter used for composite keys.
	KeyDelimiter = "/"

	// LabelKeyPrefix is the prefix for Higress labels.
	LabelKeyPrefix = "higress.io/"

	// AnnotationKeyPrefix is the prefix for Higress annotations.
	AnnotationKeyPrefix = "higress.io/"
)

// Label keys.
const (
	LabelKeyDomain      = LabelKeyPrefix + "domain"
	LabelKeyRoute       = LabelKeyPrefix + "route"
	LabelKeyService     = LabelKeyPrefix + "service"
	LabelKeyWasmPlugin  = LabelKeyPrefix + "wasm-plugin"
	LabelKeyInternal    = LabelKeyPrefix + "internal"
	LabelKeyBuiltIn     = LabelKeyPrefix + "built-in"
	LabelKeyCategory    = LabelKeyPrefix + "category"
	LabelKeyVersion     = LabelKeyPrefix + "version"
	LabelKeyEnabled     = LabelKeyPrefix + "enabled"
	LabelKeyScope       = LabelKeyPrefix + "scope"
	LabelKeyTarget      = LabelKeyPrefix + "target"
	LabelKeyPluginName  = LabelKeyPrefix + "plugin-name"
	LabelKeyConsumer    = LabelKeyPrefix + "consumer"
	LabelKeyMcpServer   = LabelKeyPrefix + "mcp-server"
	LabelKeyAiRoute     = LabelKeyPrefix + "ai-route"
	LabelKeyLlmProvider = LabelKeyPrefix + "llm-provider"
)

// Annotation keys.
const (
	AnnotationKeyConfig = AnnotationKeyPrefix + "config"
)
