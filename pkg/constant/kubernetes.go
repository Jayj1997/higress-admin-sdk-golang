// Package constant provides constants for Higress Admin SDK.
package constant

// Kubernetes annotation keys.
const (
	// AnnotationKeyIngressClass is the ingress class annotation key.
	AnnotationKeyIngressClass = "kubernetes.io/ingress.class"

	// AnnotationKeySSLRedirect is the SSL redirect annotation key.
	AnnotationKeySSLRedirect = "nginx.ingress.kubernetes.io/ssl-redirect"

	// AnnotationKeySSLForceRedirect is the force SSL redirect annotation key.
	AnnotationKeySSLForceRedirect = "nginx.ingress.kubernetes.io/force-ssl-redirect"

	// AnnotationKeyRewriteTarget is the rewrite target annotation key.
	AnnotationKeyRewriteTarget = "nginx.ingress.kubernetes.io/rewrite-target"

	// AnnotationKeyUseRegex is the use regex annotation key.
	AnnotationKeyUseRegex = "nginx.ingress.kubernetes.io/use-regex"

	// AnnotationKeyBackendProtocol is the backend protocol annotation key.
	AnnotationKeyBackendProtocol = "nginx.ingress.kubernetes.io/backend-protocol"

	// AnnotationKeyCorsEnabled is the CORS enabled annotation key.
	AnnotationKeyCorsEnabled = "nginx.ingress.kubernetes.io/enable-cors"

	// AnnotationKeyCorsAllowOrigin is the CORS allow origin annotation key.
	AnnotationKeyCorsAllowOrigin = "nginx.ingress.kubernetes.io/cors-allow-origin"

	// AnnotationKeyCorsAllowMethods is the CORS allow methods annotation key.
	AnnotationKeyCorsAllowMethods = "nginx.ingress.kubernetes.io/cors-allow-methods"

	// AnnotationKeyCorsAllowHeaders is the CORS allow headers annotation key.
	AnnotationKeyCorsAllowHeaders = "nginx.ingress.kubernetes.io/cors-allow-headers"

	// AnnotationKeyProxyNextUpstream is the proxy next upstream annotation key.
	AnnotationKeyProxyNextUpstream = "nginx.ingress.kubernetes.io/proxy-next-upstream"

	// AnnotationKeyProxyNextUpstreamTries is the proxy next upstream tries annotation key.
	AnnotationKeyProxyNextUpstreamTries = "nginx.ingress.kubernetes.io/proxy-next-upstream-tries"

	// AnnotationKeyProxyNextUpstreamTimeout is the proxy next upstream timeout annotation key.
	AnnotationKeyProxyNextUpstreamTimeout = "nginx.ingress.kubernetes.io/proxy-next-upstream-timeout"

	// AnnotationKeyConfigurationSnippet is the configuration snippet annotation key.
	AnnotationKeyConfigurationSnippet = "nginx.ingress.kubernetes.io/configuration-snippet"

	// AnnotationKeyServerSnippet is the server snippet annotation key.
	AnnotationKeyServerSnippet = "nginx.ingress.kubernetes.io/server-snippet"
)

// Kubernetes label keys.
const (
	// LabelKeyAppName is the app name label key.
	LabelKeyAppName = "app"

	// LabelKeyAppVersion is the app version label key.
	LabelKeyAppVersion = "version"

	// LabelKeyCreatedBy is the created by label key.
	LabelKeyCreatedBy = "app.kubernetes.io/created-by"

	// LabelKeyPartOf is the part of label key.
	LabelKeyPartOf = "app.kubernetes.io/part-of"

	// LabelKeyManagedBy is the managed by label key.
	LabelKeyManagedBy = "app.kubernetes.io/managed-by"

	// LabelResourceBizTypeKey is the resource business type label key.
	LabelResourceBizTypeKey = "higress.io/resource-biz-type"

	// LabelInternalKey is the internal resource label key.
	LabelInternalKey = "higress.io/internal"
)

// Kubernetes resource names.
const (
	// ResourceNameConfigMap is the config map resource name.
	ResourceNameConfigMap = "configmaps"

	// ResourceNameSecret is the secret resource name.
	ResourceNameSecret = "secrets"

	// ResourceNameIngress is the ingress resource name.
	ResourceNameIngress = "ingresses"

	// ResourceNameService is the service resource name.
	ResourceNameService = "services"

	// ResourceNameNamespace is the namespace resource name.
	ResourceNameNamespace = "namespaces"
)

// CRD Group and Version constants.
const (
	// CRDGroupIstio is the Istio CRD group.
	CRDGroupIstio = "networking.istio.io"

	// CRDGroupWasm is the WASM plugin CRD group.
	CRDGroupWasm = "extensions.higress.io"

	// CRDGroupMcp is the MCP bridge CRD group.
	CRDGroupMcp = "networking.higress.io"

	// CRDVersionV1Alpha3 is the v1alpha3 CRD version.
	CRDVersionV1Alpha3 = "v1alpha3"

	// CRDVersionV1 is the v1 CRD version.
	CRDVersionV1 = "v1"
)

// CRD Kind constants.
const (
	CRDKindEnvoyFilter     = "EnvoyFilter"
	CRDKindWasmPlugin      = "WasmPlugin"
	CRDKindMcpBridge       = "McpBridge"
	CRDKindVirtualService  = "VirtualService"
	CRDKindDestinationRule = "DestinationRule"
)
