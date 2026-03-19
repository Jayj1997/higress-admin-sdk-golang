// Package wasm contains WASM Plugin CRD types for Kubernetes
package wasm

// V1alpha1WasmPlugin represents a WASM Plugin CRD
type V1alpha1WasmPlugin struct {
	// APIGroup is the API group
	APIGroup string `json:"apiGroup,omitempty"`
	// APIVersion is the API version
	APIVersion string `json:"apiVersion,omitempty"`
	// Kind is the resource kind
	Kind string `json:"kind,omitempty"`
	// Metadata contains object metadata
	Metadata *V1ObjectMeta `json:"metadata,omitempty"`
	// Spec contains the WasmPlugin specification
	Spec *V1alpha1WasmPluginSpec `json:"spec,omitempty"`
}

// Constants for WasmPlugin
const (
	WasmPluginAPIGroup   = "extensions.higress.io"
	WasmPluginAPIVersion = "v1alpha1"
	WasmPluginKind       = "WasmPlugin"
	WasmPluginPlural     = "wasmplugins"
)

// NewV1alpha1WasmPlugin creates a new WasmPlugin
func NewV1alpha1WasmPlugin() *V1alpha1WasmPlugin {
	return &V1alpha1WasmPlugin{
		APIGroup:   WasmPluginAPIGroup,
		APIVersion: WasmPluginAPIVersion,
		Kind:       WasmPluginKind,
		Metadata:   &V1ObjectMeta{},
		Spec:       &V1alpha1WasmPluginSpec{},
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

// V1alpha1WasmPluginSpec represents the WasmPlugin specification
type V1alpha1WasmPluginSpec struct {
	// PluginName is the plugin name
	PluginName string `json:"pluginName,omitempty"`
	// PluginVersion is the plugin version
	PluginVersion string `json:"pluginVersion,omitempty"`
	// Phase is the plugin phase
	Phase PluginPhase `json:"phase,omitempty"`
	// Priority is the plugin priority
	Priority *int64 `json:"priority,omitempty"`
	// Url is the plugin URL (OCI registry or HTTP URL)
	Url string `json:"url,omitempty"`
	// Sha256 is the plugin SHA256 checksum
	Sha256 string `json:"sha256,omitempty"`
	// ImagePullPolicy is the image pull policy
	ImagePullPolicy ImagePullPolicy `json:"imagePullPolicy,omitempty"`
	// ImagePullSecret is the image pull secret
	ImagePullSecret string `json:"imagePullSecret,omitempty"`
	// WasmURL is the WASM URL (deprecated, use Url)
	WasmURL string `json:"wasmUrl,omitempty"`
	// WasmName is the WASM name (deprecated, use PluginName)
	WasmName string `json:"wasmName,omitempty"`
	// WasmVersion is the WASM version (deprecated, use PluginVersion)
	WasmVersion string `json:"wasmVersion,omitempty"`
	// MatchRules are the match rules
	MatchRules []*MatchRule `json:"matchRules,omitempty"`
	// DefaultConfig is the default configuration
	DefaultConfig interface{} `json:"defaultConfig,omitempty"`
	// DefaultConfigDisable indicates if default config is disabled
	DefaultConfigDisable bool `json:"defaultConfigDisable,omitempty"`
	// FailStrategy is the fail strategy
	FailStrategy FailStrategy `json:"failStrategy,omitempty"`
	// Config is the plugin configuration
	Config interface{} `json:"config,omitempty"`
	// RawConfig is the raw configuration string
	RawConfig string `json:"rawConfig,omitempty"`
	// Info is the plugin info
	Info *PluginInfo `json:"info,omitempty"`
}

// PluginPhase represents the plugin phase
type PluginPhase string

const (
	// PluginPhaseUnspecifiedPhase represents an unspecified phase
	PluginPhaseUnspecifiedPhase PluginPhase = "UNSPECIFIED_PHASE"
	// PluginPhaseAuthn represents authentication phase
	PluginPhaseAuthn PluginPhase = "AUTHN"
	// PluginPhaseAuthz represents authorization phase
	PluginPhaseAuthz PluginPhase = "AUTHZ"
	// PluginPhaseStats represents stats phase
	PluginPhaseStats PluginPhase = "STATS"
	// PluginPhaseBeforeProxy represents before proxy phase
	PluginPhaseBeforeProxy PluginPhase = "BEFORE_PROXY"
)

// ImagePullPolicy represents the image pull policy
type ImagePullPolicy string

const (
	// ImagePullPolicyUnspecified represents an unspecified policy
	ImagePullPolicyUnspecified ImagePullPolicy = "UNSPECIFIED_POLICY"
	// ImagePullPolicyIfNotPresent represents if not present policy
	ImagePullPolicyIfNotPresent ImagePullPolicy = "IfNotPresent"
	// ImagePullPolicyAlways represents always policy
	ImagePullPolicyAlways ImagePullPolicy = "Always"
)

// FailStrategy represents the fail strategy
type FailStrategy string

const (
	// FailStrategyUnspecified represents an unspecified strategy
	FailStrategyUnspecified FailStrategy = "UNSPECIFIED_STRATEGY"
	// FailStrategyFailOpen represents fail open strategy
	FailStrategyFailOpen FailStrategy = "FAIL_OPEN"
	// FailStrategyFailClose represents fail close strategy
	FailStrategyFailClose FailStrategy = "FAIL_CLOSE"
)

// MatchRule represents a match rule
type MatchRule struct {
	// Domain is the domain to match
	Domain string `json:"domain,omitempty"`
	// Ingress is the ingress to match
	Ingress string `json:"ingress,omitempty"`
	// Service is the service to match
	Service string `json:"service,omitempty"`
	// Config is the configuration for this match rule
	Config interface{} `json:"config,omitempty"`
	// RawConfig is the raw configuration string
	RawConfig string `json:"rawConfig,omitempty"`
	// Enable indicates if this rule is enabled
	Enable bool `json:"enable,omitempty"`
}

// PluginInfo represents plugin information
type PluginInfo struct {
	// Category is the plugin category
	Category string `json:"category,omitempty"`
	// Name is the plugin name
	Name string `json:"name,omitempty"`
	// Title is the plugin title
	Title string `json:"title,omitempty"`
	// Description is the plugin description
	Description string `json:"description,omitempty"`
	// Icon is the plugin icon URL
	Icon string `json:"icon,omitempty"`
	// Version is the plugin version
	Version string `json:"version,omitempty"`
	// Author is the plugin author
	Author string `json:"author,omitempty"`
	// Keywords are the plugin keywords
	Keywords []string `json:"keywords,omitempty"`
	// Homepage is the plugin homepage URL
	Homepage string `json:"homepage,omitempty"`
	// License is the plugin license
	License string `json:"license,omitempty"`
	// Repository is the plugin repository URL
	Repository string `json:"repository,omitempty"`
	// BuiltIn indicates if the plugin is built-in
	BuiltIn bool `json:"builtIn,omitempty"`
	// Custom indicates if the plugin is custom
	Custom bool `json:"custom,omitempty"`
}

// V1alpha1WasmPluginList represents a list of WasmPlugins
type V1alpha1WasmPluginList struct {
	// APIVersion is the API version
	APIVersion string `json:"apiVersion,omitempty"`
	// Kind is the resource kind
	Kind string `json:"kind,omitempty"`
	// Metadata contains list metadata
	Metadata *V1ListMeta `json:"metadata,omitempty"`
	// Items is the list of WasmPlugins
	Items []*V1alpha1WasmPlugin `json:"items,omitempty"`
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
