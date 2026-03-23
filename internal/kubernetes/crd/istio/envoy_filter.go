// Package istio contains Istio CRD types for Kubernetes
package istio

// V1alpha3EnvoyFilter represents an Istio EnvoyFilter CRD.
//
// Note: APIVersion and Kind are defined directly in this struct (not via embedded TypeMeta)
// because jsoniter library doesn't properly support the `json:",inline"` tag. When using
// jsoniter to marshal a struct with embedded TypeMeta, the apiVersion and kind fields are
// nested under a "TypeMeta" key instead of being at the root level, which causes Kubernetes
// API server to reject the request with "Object 'Kind' is missing" error.
//
// +k8s:deepcopy-gen=true
type V1alpha3EnvoyFilter struct {
	// APIVersion is the API version (format: group/version)
	APIVersion string `json:"apiVersion,omitempty"`
	// Kind is the resource kind
	Kind string `json:"kind,omitempty"`
	// Metadata contains object metadata
	Metadata *V1ObjectMeta `json:"metadata,omitempty"`
	// Spec contains the EnvoyFilter specification
	Spec *V1alpha3EnvoyFilterSpec `json:"spec,omitempty"`
}

// Constants for EnvoyFilter
const (
	EnvoyFilterAPIGroup   = "networking.istio.io"
	EnvoyFilterAPIVersion = "networking.istio.io/v1alpha3"
	EnvoyFilterVersion    = "v1alpha3" // Version only, for GVR
	EnvoyFilterKind       = "EnvoyFilter"
	EnvoyFilterPlural     = "envoyfilters"
)

// NewV1alpha3EnvoyFilter creates a new EnvoyFilter
func NewV1alpha3EnvoyFilter() *V1alpha3EnvoyFilter {
	return &V1alpha3EnvoyFilter{
		APIVersion: EnvoyFilterAPIVersion,
		Kind:       EnvoyFilterKind,
		Metadata:   &V1ObjectMeta{},
		Spec:       &V1alpha3EnvoyFilterSpec{},
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

// V1alpha3EnvoyFilterSpec represents the EnvoyFilter specification
type V1alpha3EnvoyFilterSpec struct {
	// WorkloadSelector selects the workloads to apply the filter to
	WorkloadSelector *V1alpha3WorkloadSelector `json:"workloadSelector,omitempty"`
	// ConfigPatches are the configuration patches to apply
	ConfigPatches []*V1alpha3EnvoyConfigObjectPatch `json:"configPatches,omitempty"`
}

// V1alpha3WorkloadSelector represents a workload selector
type V1alpha3WorkloadSelector struct {
	// Labels are the labels to match
	Labels map[string]string `json:"labels,omitempty"`
}

// V1alpha3EnvoyConfigObjectPatch represents a config object patch
type V1alpha3EnvoyConfigObjectPatch struct {
	// ApplyTo specifies where to apply the patch
	ApplyTo ApplyTo `json:"applyTo,omitempty"`
	// Match specifies the match criteria
	Match *V1alpha3EnvoyConfigObjectMatch `json:"match,omitempty"`
	// Patch specifies the patch to apply
	Patch *V1alpha3Patch `json:"patch,omitempty"`
}

// ApplyTo represents what to apply the patch to
type ApplyTo string

const (
	// ApplyToInvalid represents an invalid apply target
	ApplyToInvalid ApplyTo = "INVALID"
	// ApplyToListener represents a listener
	ApplyToListener ApplyTo = "LISTENER"
	// ApplyToFilterChain represents a filter chain
	ApplyToFilterChain ApplyTo = "FILTER_CHAIN"
	// ApplyToNetworkFilter represents a network filter
	ApplyToNetworkFilter ApplyTo = "NETWORK_FILTER"
	// ApplyToHTTPFilter represents an HTTP filter
	ApplyToHTTPFilter ApplyTo = "HTTP_FILTER"
)

// V1alpha3EnvoyConfigObjectMatch represents match criteria for config objects
type V1alpha3EnvoyConfigObjectMatch struct {
	// Context specifies the context for matching
	Context PatchContext `json:"context,omitempty"`
	// Proxy specifies proxy match criteria
	Proxy *V1alpha3ProxyMatch `json:"proxy,omitempty"`
	// Listener specifies listener match criteria
	Listener *V1alpha3ListenerMatch `json:"listener,omitempty"`
	// RouteConfiguration specifies route configuration match criteria
	RouteConfiguration *V1alpha3RouteConfigurationMatch `json:"routeConfiguration,omitempty"`
	// Cluster specifies cluster match criteria
	Cluster *V1alpha3ClusterMatch `json:"cluster,omitempty"`
}

// PatchContext represents the patch context
type PatchContext string

const (
	// PatchContextAny matches any context
	PatchContextAny PatchContext = "ANY"
	// PatchContextInbound matches inbound context
	PatchContextInbound PatchContext = "INBOUND"
	// PatchContextOutbound matches outbound context
	PatchContextOutbound PatchContext = "OUTBOUND"
	// PatchContextGateway matches gateway context
	PatchContextGateway PatchContext = "GATEWAY"
)

// V1alpha3ProxyMatch represents proxy match criteria
type V1alpha3ProxyMatch struct {
	// ProxyVersion is the proxy version to match
	ProxyVersion string `json:"proxyVersion,omitempty"`
	// Metadata is the metadata to match
	Metadata map[string]string `json:"metadata,omitempty"`
}

// V1alpha3ListenerMatch represents listener match criteria
type V1alpha3ListenerMatch struct {
	// PortNumber is the port number to match
	PortNumber uint32 `json:"portNumber,omitempty"`
	// PortName is the port name to match
	PortName string `json:"portName,omitempty"`
	// FilterChain specifies filter chain match criteria
	FilterChain *V1alpha3FilterChainMatch `json:"filterChain,omitempty"`
}

// V1alpha3FilterChainMatch represents filter chain match criteria
type V1alpha3FilterChainMatch struct {
	// Filter specifies filter match criteria
	Filter *V1alpha3FilterMatch `json:"filter,omitempty"`
}

// V1alpha3FilterMatch represents filter match criteria
type V1alpha3FilterMatch struct {
	// Name is the filter name to match
	Name string `json:"name,omitempty"`
	// SubFilter specifies sub-filter match criteria
	SubFilter *V1alpha3SubFilterMatch `json:"subFilter,omitempty"`
}

// V1alpha3SubFilterMatch represents sub-filter match criteria
type V1alpha3SubFilterMatch struct {
	// Name is the sub-filter name to match
	Name string `json:"name,omitempty"`
}

// V1alpha3RouteConfigurationMatch represents route configuration match criteria
type V1alpha3RouteConfigurationMatch struct {
	// PortNumber is the port number to match
	PortNumber uint32 `json:"portNumber,omitempty"`
	// PortName is the port name to match
	PortName string `json:"portName,omitempty"`
	// Gateway is the gateway to match
	Gateway string `json:"gateway,omitempty"`
	// Vhost specifies virtual host match criteria
	Vhost *V1alpha3VirtualHostMatch `json:"vhost,omitempty"`
	// Name is the route configuration name
	Name string `json:"name,omitempty"`
}

// V1alpha3VirtualHostMatch represents virtual host match criteria
type V1alpha3VirtualHostMatch struct {
	// Name is the virtual host name
	Name string `json:"name,omitempty"`
	// Route specifies route match criteria
	Route *V1alpha3RouteMatch `json:"route,omitempty"`
}

// V1alpha3RouteMatch represents route match criteria
type V1alpha3RouteMatch struct {
	// Name is the route name
	Name string `json:"name,omitempty"`
}

// V1alpha3ClusterMatch represents cluster match criteria
type V1alpha3ClusterMatch struct {
	// PortNumber is the port number to match
	PortNumber uint32 `json:"portNumber,omitempty"`
	// Service is the service to match
	Service string `json:"service,omitempty"`
	// Subset is the subset to match
	Subset string `json:"subset,omitempty"`
	// Name is the cluster name
	Name string `json:"name,omitempty"`
}

// V1alpha3Patch represents a patch to apply
type V1alpha3Patch struct {
	// Operation is the patch operation
	Operation Operation `json:"operation,omitempty"`
	// Value is the patch value
	Value interface{} `json:"value,omitempty"`
}

// Operation represents a patch operation
type Operation string

const (
	// OperationInvalid represents an invalid operation
	OperationInvalid Operation = "INVALID"
	// OperationMerge represents a merge operation
	OperationMerge Operation = "MERGE"
	// OperationAdd represents an add operation
	OperationAdd Operation = "ADD"
	// OperationRemove represents a remove operation
	OperationRemove Operation = "REMOVE"
	// OperationInsertBefore represents an insert before operation
	OperationInsertBefore Operation = "INSERT_BEFORE"
	// OperationInsertAfter represents an insert after operation
	OperationInsertAfter Operation = "INSERT_AFTER"
	// OperationReplace represents a replace operation
	OperationReplace Operation = "REPLACE"
	// OperationCopy represents a copy operation
	OperationCopy Operation = "COPY"
	// OperationMove represents a move operation
	OperationMove Operation = "MOVE"
)

// V1alpha3EnvoyFilterList represents a list of EnvoyFilters
type V1alpha3EnvoyFilterList struct {
	// APIVersion is the API version
	APIVersion string `json:"apiVersion,omitempty"`
	// Kind is the resource kind
	Kind string `json:"kind,omitempty"`
	// Metadata contains list metadata
	Metadata *V1ListMeta `json:"metadata,omitempty"`
	// Items is the list of EnvoyFilters
	Items []*V1alpha3EnvoyFilter `json:"items,omitempty"`
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
