// Package kubernetes provides Kubernetes client functionality
package kubernetes

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/internal/kubernetes/crd/wasm"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/constant"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/errors"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model/route"
	jsoniter "github.com/json-iterator/go"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Annotation key constants
const (
	AnnotationKeyPrefix = "higress.io/"

	// Route related annotations
	AnnotationKeyUseRegex                 = "higress.io/use-regex"
	AnnotationKeyDestination              = "higress.io/destination"
	AnnotationKeySSLRedirect              = "higress.io/ssl-redirect"
	AnnotationKeyRewriteEnabled           = "higress.io/enable-rewrite"
	AnnotationKeyRewritePath              = "higress.io/rewrite-path"
	AnnotationKeyRewriteTarget            = "higress.io/rewrite-target"
	AnnotationKeyUpstreamVhost            = "higress.io/upstream-vhost"
	AnnotationKeyProxyNextUpstream        = "higress.io/proxy-next-upstream"
	AnnotationKeyProxyNextUpstreamTries   = "higress.io/proxy-next-upstream-tries"
	AnnotationKeyProxyNextUpstreamTimeout = "higress.io/proxy-next-upstream-timeout"
	AnnotationKeyHeaderControlEnabled     = "higress.io/enable-header-control"
	AnnotationKeyRequestHeaderAdd         = "higress.io/request-header-control-add"
	AnnotationKeyRequestHeaderUpdate      = "higress.io/request-header-control-update"
	AnnotationKeyRequestHeaderRemove      = "higress.io/request-header-control-remove"
	AnnotationKeyResponseHeaderAdd        = "higress.io/response-header-control-add"
	AnnotationKeyResponseHeaderUpdate     = "higress.io/response-header-control-update"
	AnnotationKeyResponseHeaderRemove     = "higress.io/response-header-control-remove"
	AnnotationKeyCorsEnabled              = "higress.io/enable-cors"
	AnnotationKeyCorsAllowOrigin          = "higress.io/cors-allow-origin"
	AnnotationKeyCorsAllowMethods         = "higress.io/cors-allow-methods"
	AnnotationKeyCorsAllowHeaders         = "higress.io/cors-allow-headers"
	AnnotationKeyCorsExposeHeaders        = "higress.io/cors-expose-headers"
	AnnotationKeyCorsAllowCredentials     = "higress.io/cors-allow-credentials"
	AnnotationKeyCorsMaxAge               = "higress.io/cors-max-age"
	AnnotationKeyMethod                   = "higress.io/match-method"
	AnnotationKeyIgnorePathCase           = "higress.io/ignore-path-case"
	AnnotationKeyComment                  = "higress.io/comment"

	// Query and header match
	AnnotationKeyQueryMatchKeyword  = "-match-query-"
	AnnotationKeyQueryMatchFormat   = "higress.io/%s" + AnnotationKeyQueryMatchKeyword + "%s"
	AnnotationKeyHeaderMatchKeyword = "-match-header-"
	AnnotationKeyHeaderMatchFormat  = "higress.io/%s" + AnnotationKeyHeaderMatchKeyword + "%s"

	// WasmPlugin annotations
	AnnotationKeyWasmPluginTitle       = "higress.io/wasm-plugin-title"
	AnnotationKeyWasmPluginDescription = "higress.io/wasm-plugin-description"
	AnnotationKeyWasmPluginIcon        = "higress.io/wasm-plugin-icon"
)

// Label key constants
const (
	LabelKeyDomainPrefix         = "higress.io/domain_"
	LabelKeyDomainValue          = "true"
	LabelKeyConfigMapType        = "higress.io/config-map-type"
	LabelKeyConfigMapTypeDomain  = "domain"
	LabelKeyConfigMapTypeAiRoute = "ai-route"
	LabelKeyInternal             = "higress.io/internal"
	LabelKeyWasmPluginName       = "higress.io/wasm-plugin-name"
	LabelKeyWasmPluginVersion    = "higress.io/wasm-plugin-version"
	LabelKeyWasmPluginBuiltIn    = "higress.io/wasm-plugin-built-in"
	LabelKeyWasmPluginCategory   = "higress.io/wasm-plugin-category"
)

// ConfigMap data keys
const (
	ConfigMapKeyDomain      = "domain"
	ConfigMapKeyCert        = "cert"
	ConfigMapKeyEnableHTTPS = "enableHttps"
	ConfigMapKeyData        = "data"
)

// Secret constants
const (
	SecretTypeTLS     = "kubernetes.io/tls"
	SecretTLSCrtField = "tls.crt"
	SecretTLSKeyField = "tls.key"
)

// Ingress path type constants
const (
	IngressPathTypeExact                  = "Exact"
	IngressPathTypePrefix                 = "Prefix"
	IngressPathTypeImplementationSpecific = "ImplementationSpecific"
)

// Default values
const (
	DefaultWeight = 100
	YAMLSeparator = "---\n"
)

var jsonFast = jsoniter.ConfigCompatibleWithStandardLibrary

// KubernetesModelConverter converts between Kubernetes resources and business models.
type KubernetesModelConverter struct {
	clientService *KubernetesClientService
}

// NewKubernetesModelConverter creates a new KubernetesModelConverter.
func NewKubernetesModelConverter(clientService *KubernetesClientService) *KubernetesModelConverter {
	return &KubernetesModelConverter{
		clientService: clientService,
	}
}

// =============================================================================
// Ingress <-> Route Conversion
// =============================================================================

// IngressToRoute converts a V1Ingress to a Route model.
func (c *KubernetesModelConverter) IngressToRoute(ingress *networkingv1.Ingress) (*model.Route, error) {
	if ingress == nil {
		return nil, errors.NewValidationError("ingress cannot be nil")
	}

	r := &model.Route{}

	// Fill metadata
	c.fillRouteMetadata(r, ingress)

	// Fill route info
	if err := c.fillRouteInfo(r, ingress); err != nil {
		return nil, err
	}

	// Fill custom configs
	c.fillCustomConfigs(r, ingress)

	// Fill custom labels
	c.fillCustomLabels(r, ingress)

	// Set readonly flag
	readonly := !c.IsIngressSupported(ingress)
	r.Readonly = &readonly

	return r, nil
}

// RouteToIngress converts a Route model to a V1Ingress.
func (c *KubernetesModelConverter) RouteToIngress(r *model.Route) (*networkingv1.Ingress, error) {
	if r == nil {
		return nil, errors.NewValidationError("route cannot be nil")
	}

	ingress := &networkingv1.Ingress{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "networking.k8s.io/v1",
			Kind:       "Ingress",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: r.Name,
		},
	}

	// Fill metadata
	c.fillIngressMetadata(ingress, r)

	// Fill spec
	if err := c.fillIngressSpec(ingress, r); err != nil {
		return nil, err
	}

	// Fill CORS
	c.fillIngressCors(ingress, r)

	// Fill annotations
	if err := c.fillIngressAnnotations(ingress, r); err != nil {
		return nil, err
	}

	// Fill labels
	c.fillIngressLabels(ingress, r)

	return ingress, nil
}

// IngressesToRoutes converts a list of V1Ingress to a list of Route models.
func (c *KubernetesModelConverter) IngressesToRoutes(ingresses []networkingv1.Ingress) ([]model.Route, error) {
	routes := make([]model.Route, 0, len(ingresses))
	for i := range ingresses {
		r, err := c.IngressToRoute(&ingresses[i])
		if err != nil {
			return nil, err
		}
		routes = append(routes, *r)
	}
	return routes, nil
}

// IsIngressSupported checks if an ingress is supported by the converter.
func (c *KubernetesModelConverter) IsIngressSupported(ingress *networkingv1.Ingress) bool {
	if ingress == nil || ingress.ObjectMeta.Name == "" {
		return false
	}

	spec := ingress.Spec
	rules := spec.Rules

	// Only support single rule
	if len(rules) == 0 || len(rules) > 1 {
		return false
	}

	rule := rules[0]
	if rule.HTTP == nil {
		return false
	}

	// Only support single path
	if len(rule.HTTP.Paths) == 0 || len(rule.HTTP.Paths) > 1 {
		return false
	}

	path := rule.HTTP.Paths[0]
	if path.PathType == nil {
		return false
	}

	pathType := string(*path.PathType)
	if pathType != IngressPathTypeExact && pathType != IngressPathTypePrefix {
		return false
	}

	return true
}

// =============================================================================
// ConfigMap <-> ServiceSource Conversion
// =============================================================================

// ConfigMapToServiceSource converts a V1ConfigMap to a ServiceSource model.
func (c *KubernetesModelConverter) ConfigMapToServiceSource(cm *corev1.ConfigMap) (*model.ServiceSource, error) {
	if cm == nil {
		return nil, errors.NewValidationError("configmap cannot be nil")
	}

	source := &model.ServiceSource{
		Name: cm.Name,
	}

	if cm.ObjectMeta.ResourceVersion != "" {
		source.Version = cm.ObjectMeta.ResourceVersion
	}

	data := cm.Data
	if data == nil {
		return source, nil
	}

	// Parse service source data from configmap
	if domain, ok := data["domain"]; ok {
		source.Domain = domain
	}
	if portStr, ok := data["port"]; ok {
		if port, err := strconv.Atoi(portStr); err == nil {
			source.Port = &port
		}
	}
	if ns, ok := data["namespace"]; ok {
		source.Namespace = ns
	}
	if group, ok := data["group"]; ok {
		source.Group = group
	}
	if sourceType, ok := data["type"]; ok {
		source.Type = sourceType
	}

	return source, nil
}

// ServiceSourceToConfigMap converts a ServiceSource model to a V1ConfigMap.
func (c *KubernetesModelConverter) ServiceSourceToConfigMap(source *model.ServiceSource) (*corev1.ConfigMap, error) {
	if source == nil {
		return nil, errors.NewValidationError("service source cannot be nil")
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:            source.Name,
			ResourceVersion: source.Version,
			Labels: map[string]string{
				constant.LabelResourceDefinerKey: constant.LabelResourceDefinerValue,
			},
		},
		Data: make(map[string]string),
	}

	if source.Domain != "" {
		cm.Data["domain"] = source.Domain
	}
	if source.Port != nil {
		cm.Data["port"] = strconv.Itoa(*source.Port)
	}
	if source.Namespace != "" {
		cm.Data["namespace"] = source.Namespace
	}
	if source.Group != "" {
		cm.Data["group"] = source.Group
	}
	if source.Type != "" {
		cm.Data["type"] = source.Type
	}

	return cm, nil
}

// =============================================================================
// Secret <-> TlsCertificate Conversion
// =============================================================================

// SecretToTlsCertificate converts a V1Secret to a TlsCertificate model.
func (c *KubernetesModelConverter) SecretToTlsCertificate(secret *corev1.Secret) (*model.TlsCertificate, error) {
	if secret == nil {
		return nil, errors.NewValidationError("secret cannot be nil")
	}

	cert := &model.TlsCertificate{
		Name: secret.Name,
	}

	if secret.ObjectMeta.ResourceVersion != "" {
		cert.Version = secret.ObjectMeta.ResourceVersion
	}

	// Extract TLS data
	data := secret.Data
	if data != nil {
		if crtData, ok := data[SecretTLSCrtField]; ok {
			cert.Cert = string(crtData)
		}
		if keyData, ok := data[SecretTLSKeyField]; ok {
			cert.Key = string(keyData)
		}
	}

	// Extract domains from labels
	if secret.Labels != nil {
		domains := make([]string, 0)
		for key := range secret.Labels {
			if strings.HasPrefix(key, LabelKeyDomainPrefix) {
				domain := strings.TrimPrefix(key, LabelKeyDomainPrefix)
				domains = append(domains, domain)
			}
		}
		cert.Domains = domains
	}

	// Fill certificate details (validity, etc.)
	c.fillTlsCertificateDetails(cert)

	return cert, nil
}

// TlsCertificateToSecret converts a TlsCertificate model to a V1Secret.
func (c *KubernetesModelConverter) TlsCertificateToSecret(cert *model.TlsCertificate) (*corev1.Secret, error) {
	if cert == nil {
		return nil, errors.NewValidationError("certificate cannot be nil")
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:            cert.Name,
			ResourceVersion: cert.Version,
			Labels: map[string]string{
				constant.LabelResourceDefinerKey: constant.LabelResourceDefinerValue,
			},
		},
		Type: corev1.SecretType(SecretTypeTLS),
		Data: make(map[string][]byte),
	}

	if cert.Cert != "" {
		secret.Data[SecretTLSCrtField] = []byte(cert.Cert)
	}
	if cert.Key != "" {
		secret.Data[SecretTLSKeyField] = []byte(cert.Key)
	}

	// Add domain labels
	for _, domain := range cert.Domains {
		normalizedDomain := normalizeDomainName(domain)
		labelKey := LabelKeyDomainPrefix + normalizedDomain
		if secret.Labels == nil {
			secret.Labels = make(map[string]string)
		}
		secret.Labels[labelKey] = LabelKeyDomainValue
	}

	return secret, nil
}

// =============================================================================
// ConfigMap <-> Domain Conversion
// =============================================================================

// ConfigMapToDomain converts a V1ConfigMap to a Domain model.
func (c *KubernetesModelConverter) ConfigMapToDomain(cm *corev1.ConfigMap) (*model.Domain, error) {
	if cm == nil {
		return nil, errors.NewValidationError("configmap cannot be nil")
	}

	domain := &model.Domain{}

	if cm.ObjectMeta.ResourceVersion != "" {
		domain.Version = cm.ObjectMeta.ResourceVersion
	}

	data := cm.Data
	if data == nil {
		return domain, nil
	}

	if name, ok := data[ConfigMapKeyDomain]; ok {
		domain.Name = name
	}
	if certID, ok := data[ConfigMapKeyCert]; ok {
		domain.CertIdentifier = certID
	}
	if enableHTTPS, ok := data[ConfigMapKeyEnableHTTPS]; ok {
		domain.EnableHTTPS = enableHTTPS
	}

	return domain, nil
}

// DomainToConfigMap converts a Domain model to a V1ConfigMap.
func (c *KubernetesModelConverter) DomainToConfigMap(domain *model.Domain) (*corev1.ConfigMap, error) {
	if domain == nil {
		return nil, errors.NewValidationError("domain cannot be nil")
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:            c.DomainNameToConfigMapName(domain.Name),
			ResourceVersion: domain.Version,
			Labels: map[string]string{
				LabelKeyConfigMapType: LabelKeyConfigMapTypeDomain,
			},
		},
		Data: map[string]string{
			ConfigMapKeyDomain:      domain.Name,
			ConfigMapKeyCert:        domain.CertIdentifier,
			ConfigMapKeyEnableHTTPS: domain.EnableHTTPS,
		},
	}

	return cm, nil
}

// DomainNameToConfigMapName converts a domain name to a ConfigMap name.
func (c *KubernetesModelConverter) DomainNameToConfigMapName(domainName string) string {
	return "domain-" + normalizeDomainName(domainName)
}

// =============================================================================
// WasmPlugin CRD <-> Model Conversion
// =============================================================================

// WasmPluginCRDToModel converts a V1alpha1WasmPlugin CRD to a WasmPlugin model.
func (c *KubernetesModelConverter) WasmPluginCRDToModel(plugin *wasm.V1alpha1WasmPlugin) (*model.WasmPlugin, error) {
	if plugin == nil {
		return nil, errors.NewValidationError("plugin cannot be nil")
	}

	result := &model.WasmPlugin{}

	metadata := plugin.Metadata
	if metadata != nil {
		// Extract name and version from labels
		if metadata.Labels != nil {
			if name, ok := metadata.Labels[LabelKeyWasmPluginName]; ok {
				result.Name = name
			}
			if version, ok := metadata.Labels[LabelKeyWasmPluginVersion]; ok {
				result.Version = version
			}
			if category, ok := metadata.Labels[LabelKeyWasmPluginCategory]; ok {
				result.Category = category
			}
			if builtIn, ok := metadata.Labels[LabelKeyWasmPluginBuiltIn]; ok {
				val := strings.ToLower(builtIn) == "true"
				result.BuiltIn = &val
			}
		}

		// Extract annotations
		if metadata.Annotations != nil {
			if title, ok := metadata.Annotations[AnnotationKeyWasmPluginTitle]; ok {
				result.Title = title
			}
			if desc, ok := metadata.Annotations[AnnotationKeyWasmPluginDescription]; ok {
				result.Description = desc
			}
			if icon, ok := metadata.Annotations[AnnotationKeyWasmPluginIcon]; ok {
				result.Icon = icon
			}
		}

		result.Version = metadata.ResourceVersion
	}

	// Extract spec
	spec := plugin.Spec
	if spec != nil {
		result.Phase = string(spec.Phase)
		if spec.Priority != nil {
			// Convert int64 to int
			priority := int(*spec.Priority)
			result.Priority = &priority
		}
		if spec.Url != "" {
			result.ImageURL = spec.Url
		}
	}

	return result, nil
}

// ModelToWasmPluginCRD converts a WasmPlugin model to a V1alpha1WasmPlugin CRD.
func (c *KubernetesModelConverter) ModelToWasmPluginCRD(plugin *model.WasmPlugin) (*wasm.V1alpha1WasmPlugin, error) {
	if plugin == nil {
		return nil, errors.NewValidationError("plugin cannot be nil")
	}

	cr := wasm.NewV1alpha1WasmPlugin()
	cr.Metadata.Name = plugin.Name + "-" + plugin.Version

	if cr.Metadata.Labels == nil {
		cr.Metadata.Labels = make(map[string]string)
	}
	cr.Metadata.Labels[LabelKeyWasmPluginName] = plugin.Name
	cr.Metadata.Labels[LabelKeyWasmPluginVersion] = plugin.Version

	if plugin.BuiltIn != nil {
		cr.Metadata.Labels[LabelKeyWasmPluginBuiltIn] = strconv.FormatBool(*plugin.BuiltIn)
	}
	if plugin.Category != "" {
		cr.Metadata.Labels[LabelKeyWasmPluginCategory] = plugin.Category
	}

	// Add annotations
	if cr.Metadata.Annotations == nil {
		cr.Metadata.Annotations = make(map[string]string)
	}
	if plugin.Title != "" {
		cr.Metadata.Annotations[AnnotationKeyWasmPluginTitle] = plugin.Title
	}
	if plugin.Description != "" {
		cr.Metadata.Annotations[AnnotationKeyWasmPluginDescription] = plugin.Description
	}
	if plugin.Icon != "" {
		cr.Metadata.Annotations[AnnotationKeyWasmPluginIcon] = plugin.Icon
	}

	// Set spec
	cr.Spec.Phase = wasm.PluginPhase(plugin.Phase)
	if plugin.Priority != nil {
		priority := int64(*plugin.Priority)
		cr.Spec.Priority = &priority
	}
	if plugin.ImageURL != "" {
		cr.Spec.Url = plugin.ImageURL
	}

	return cr, nil
}

// =============================================================================
// WasmPluginInstance Conversion
// =============================================================================

// GetWasmPluginInstancesFromCR extracts all WasmPluginInstance from a V1alpha1WasmPlugin CRD.
func (c *KubernetesModelConverter) GetWasmPluginInstancesFromCR(plugin *wasm.V1alpha1WasmPlugin) ([]model.WasmPluginInstance, error) {
	if plugin == nil {
		return nil, errors.NewValidationError("plugin cannot be nil")
	}

	metadata := plugin.Metadata
	if metadata == nil || metadata.Labels == nil {
		return nil, nil
	}

	name := metadata.Labels[LabelKeyWasmPluginName]
	version := metadata.Labels[LabelKeyWasmPluginVersion]
	if name == "" || version == "" {
		return nil, nil
	}

	spec := plugin.Spec
	if spec == nil {
		return nil, nil
	}

	instances := make([]model.WasmPluginInstance, 0)

	// Handle default config (global scope)
	// DefaultConfigDisable is a bool value, not pointer
	enabled := !spec.DefaultConfigDisable
	if spec.DefaultConfig != nil {
		configs, ok := spec.DefaultConfig.(map[string]interface{})
		if !ok {
			configs = make(map[string]interface{})
		}
		instance := model.WasmPluginInstance{
			PluginName:    name,
			PluginVersion: version,
			Targets: map[model.WasmPluginInstanceScope]string{
				model.WasmPluginInstanceScopeGlobal: "",
			},
			Enabled:        &enabled,
			Configurations: configs,
		}
		instances = append(instances, instance)
	}

	// Handle match rules
	if len(spec.MatchRules) > 0 {
		for _, rule := range spec.MatchRules {
			ruleEnabled := rule.Enable

			// Build targets from match rule
			targets := make(map[model.WasmPluginInstanceScope]string)
			if rule.Domain != "" {
				targets[model.WasmPluginInstanceScopeDomain] = rule.Domain
			}
			if rule.Ingress != "" {
				targets[model.WasmPluginInstanceScopeRoute] = rule.Ingress
			}
			if rule.Service != "" {
				targets[model.WasmPluginInstanceScopeService] = rule.Service
			}

			configs, ok := rule.Config.(map[string]interface{})
			if !ok {
				configs = make(map[string]interface{})
			}

			instance := model.WasmPluginInstance{
				PluginName:     name,
				PluginVersion:  version,
				Targets:        targets,
				Enabled:        &ruleEnabled,
				Configurations: configs,
			}
			instances = append(instances, instance)
		}
	}

	return instances, nil
}

// SetWasmPluginInstanceToCR sets a WasmPluginInstance to a V1alpha1WasmPlugin CRD.
func (c *KubernetesModelConverter) SetWasmPluginInstanceToCR(cr *wasm.V1alpha1WasmPlugin, instance *model.WasmPluginInstance) error {
	if cr == nil {
		return errors.NewValidationError("cr cannot be nil")
	}
	if instance == nil {
		return errors.NewValidationError("instance cannot be nil")
	}

	if cr.Spec == nil {
		cr.Spec = &wasm.V1alpha1WasmPluginSpec{}
	}

	spec := cr.Spec
	enabled := true
	if instance.Enabled != nil {
		enabled = *instance.Enabled
	}

	// Handle global scope
	if len(instance.Targets) == 1 {
		if _, ok := instance.Targets[model.WasmPluginInstanceScopeGlobal]; ok {
			spec.DefaultConfigDisable = !enabled
			spec.DefaultConfig = instance.Configurations
			return nil
		}
	}

	// Handle other scopes - create or update match rule
	matchRule := &wasm.MatchRule{
		Enable: enabled,
		Config: instance.Configurations,
	}

	for scope, target := range instance.Targets {
		switch scope {
		case model.WasmPluginInstanceScopeDomain:
			matchRule.Domain = target
		case model.WasmPluginInstanceScopeRoute:
			matchRule.Ingress = target
		case model.WasmPluginInstanceScopeService:
			matchRule.Service = target
		}
	}

	// Add or update match rule
	if spec.MatchRules == nil {
		spec.MatchRules = make([]*wasm.MatchRule, 0)
	}

	// Check if rule already exists and update, otherwise append
	found := false
	for i, existing := range spec.MatchRules {
		if matchRuleKeyEquals(existing, matchRule) {
			spec.MatchRules[i] = matchRule
			found = true
			break
		}
	}
	if !found {
		spec.MatchRules = append(spec.MatchRules, matchRule)
	}

	// Sort match rules
	sortWasmPluginMatchRules(spec.MatchRules)

	return nil
}

// =============================================================================
// Helper Methods
// =============================================================================

// fillRouteMetadata fills route metadata from ingress.
func (c *KubernetesModelConverter) fillRouteMetadata(r *model.Route, ingress *networkingv1.Ingress) {
	if ingress == nil {
		return
	}
	r.Name = ingress.Name
	if ingress.ObjectMeta.ResourceVersion != "" {
		r.Version = ingress.ObjectMeta.ResourceVersion
	}
}

// fillRouteInfo fills route information from ingress.
func (c *KubernetesModelConverter) fillRouteInfo(r *model.Route, ingress *networkingv1.Ingress) error {
	if ingress == nil {
		return nil
	}

	spec := ingress.Spec
	rules := spec.Rules
	if len(rules) == 0 {
		return nil
	}

	r.Domains = make([]string, 0)

	// Process first rule
	rule := rules[0]
	if rule.HTTP != nil && len(rule.HTTP.Paths) > 0 {
		path := rule.HTTP.Paths[0]
		if err := c.fillPathRoute(r, ingress, &path); err != nil {
			return err
		}
	}

	// Get domains
	if rule.Host != "" {
		r.Domains = append(r.Domains, rule.Host)
	} else if len(spec.TLS) > 0 {
		for _, tls := range spec.TLS {
			if len(tls.Hosts) > 0 {
				r.Domains = append(r.Domains, tls.Hosts...)
			}
		}
	}

	// Process annotations
	if ingress.ObjectMeta.Annotations != nil {
		annotations := ingress.ObjectMeta.Annotations
		c.fillRewriteConfig(r, annotations)
		c.fillProxyNextUpstreamConfig(r, annotations)
		c.fillHeaderAndQueryConfig(r, annotations)
		c.fillMethodConfig(r, annotations)
		c.fillHeaderControlConfig(r, annotations)
		c.fillRouteCors(r, &ingress.ObjectMeta)
	}

	return nil
}

// fillPathRoute fills path route information.
func (c *KubernetesModelConverter) fillPathRoute(r *model.Route, ingress *networkingv1.Ingress, path *networkingv1.HTTPIngressPath) error {
	// Fill path predicate
	c.fillPathPredicates(r, ingress, path)

	// Fill route destinations
	c.fillRouteDestinations(r, ingress, path)

	return nil
}

// fillPathPredicates fills path predicates.
func (c *KubernetesModelConverter) fillPathPredicates(r *model.Route, ingress *networkingv1.Ingress, path *networkingv1.HTTPIngressPath) {
	if path == nil {
		return
	}

	r.Path = &route.RoutePredicate{
		Path: path.Path,
	}

	if path.PathType != nil {
		switch string(*path.PathType) {
		case IngressPathTypeExact:
			r.Path.MatchType = route.MatchTypeExact
		case IngressPathTypePrefix:
			// Check for regex
			if ingress != nil && ingress.ObjectMeta.Annotations != nil {
				if useRegex, ok := ingress.ObjectMeta.Annotations[AnnotationKeyUseRegex]; ok && useRegex == "true" {
					r.Path.MatchType = route.MatchTypeRegex
				} else {
					r.Path.MatchType = route.MatchTypePrefix
				}
			} else {
				r.Path.MatchType = route.MatchTypePrefix
			}
		}
	}

	// Check ignore path case
	if ingress != nil && ingress.ObjectMeta.Annotations != nil {
		if ignoreCase, ok := ingress.ObjectMeta.Annotations[AnnotationKeyIgnorePathCase]; ok {
			caseSensitive := !strings.EqualFold(ignoreCase, "true")
			r.Path.CaseSensitive = &caseSensitive
		}
	}
}

// fillRouteDestinations fills route destination services.
func (c *KubernetesModelConverter) fillRouteDestinations(r *model.Route, ingress *networkingv1.Ingress, path *networkingv1.HTTPIngressPath) {
	if ingress == nil || ingress.ObjectMeta.Annotations == nil {
		return
	}

	destination := ingress.ObjectMeta.Annotations[AnnotationKeyDestination]
	if destination == "" {
		return
	}

	services := make([]*route.UpstreamService, 0)
	lines := strings.Split(destination, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		svc := c.buildDestination(line)
		if svc != nil {
			services = append(services, svc)
		}
	}

	r.Services = services
}

// buildDestination builds an UpstreamService from a destination string.
func (c *KubernetesModelConverter) buildDestination(config string) *route.UpstreamService {
	fields := strings.Fields(config)
	if len(fields) == 0 {
		return nil
	}

	weight := DefaultWeight
	addrIndex := 0

	// Check if first field is weight (ends with %)
	if strings.HasSuffix(fields[0], "%") {
		weightStr := strings.TrimSuffix(fields[0], "%")
		if w, err := strconv.Atoi(weightStr); err == nil {
			weight = w
		}
		addrIndex = 1
	}

	if len(fields) <= addrIndex {
		return nil
	}

	address := fields[addrIndex]

	// Parse address (format: namespace/service-name:port)
	parts := strings.Split(address, "/")
	if len(parts) != 2 {
		return nil
	}

	namespace := parts[0]
	serviceAndPort := parts[1]

	// Split service name and port
	serviceParts := strings.Split(serviceAndPort, ":")
	if len(serviceParts) != 2 {
		return nil
	}

	serviceName := serviceParts[0]
	port, err := strconv.Atoi(serviceParts[1])
	if err != nil {
		return nil
	}

	return &route.UpstreamService{
		Name:      serviceName,
		Namespace: namespace,
		Port:      port,
		Weight:    &weight,
	}
}

// fillRewriteConfig fills rewrite configuration.
func (c *KubernetesModelConverter) fillRewriteConfig(r *model.Route, annotations map[string]string) {
	if annotations == nil {
		return
	}

	rewritePath := annotations[AnnotationKeyRewritePath]
	rewriteTarget := annotations[AnnotationKeyRewriteTarget]

	if rewritePath != "" || rewriteTarget != "" {
		r.Rewrite = &route.RewriteConfig{}
		if rewritePath != "" {
			r.Rewrite.Path = rewritePath
		}
		if rewriteTarget != "" {
			r.Rewrite.Path = rewriteTarget
		}
	}
}

// fillProxyNextUpstreamConfig fills proxy next upstream configuration.
func (c *KubernetesModelConverter) fillProxyNextUpstreamConfig(r *model.Route, annotations map[string]string) {
	if annotations == nil {
		return
	}

	retryOn := annotations[AnnotationKeyProxyNextUpstream]
	triesStr := annotations[AnnotationKeyProxyNextUpstreamTries]
	timeoutStr := annotations[AnnotationKeyProxyNextUpstreamTimeout]

	if retryOn != "" || triesStr != "" || timeoutStr != "" {
		r.ProxyNextUpstream = &route.ProxyNextUpstreamConfig{
			RetryOn: retryOn,
		}
		if triesStr != "" {
			if tries, err := strconv.Atoi(triesStr); err == nil {
				r.ProxyNextUpstream.NumRetries = &tries
			}
		}
		if timeoutStr != "" {
			r.ProxyNextUpstream.Timeout = timeoutStr
		}
	}
}

// fillHeaderAndQueryConfig fills header and query configuration.
func (c *KubernetesModelConverter) fillHeaderAndQueryConfig(r *model.Route, annotations map[string]string) {
	if annotations == nil {
		return
	}

	r.Headers = make([]*route.KeyedRoutePredicate, 0)
	r.URLParams = make([]*route.KeyedRoutePredicate, 0)

	for key, value := range annotations {
		// Check for header match
		if strings.Contains(key, AnnotationKeyHeaderMatchKeyword) {
			parts := strings.Split(key, AnnotationKeyHeaderMatchKeyword)
			if len(parts) == 2 {
				headerKey := parts[0]
				headerName := parts[1]
				_ = headerKey // prefix, usually "higress.io"
				r.Headers = append(r.Headers, &route.KeyedRoutePredicate{
					Key:       headerName,
					Value:     value,
					MatchType: route.MatchTypeExact,
				})
			}
		}

		// Check for query match
		if strings.Contains(key, AnnotationKeyQueryMatchKeyword) {
			parts := strings.Split(key, AnnotationKeyQueryMatchKeyword)
			if len(parts) == 2 {
				queryKey := parts[0]
				queryName := parts[1]
				_ = queryKey // prefix
				r.URLParams = append(r.URLParams, &route.KeyedRoutePredicate{
					Key:       queryName,
					Value:     value,
					MatchType: route.MatchTypeExact,
				})
			}
		}
	}
}

// fillMethodConfig fills method configuration.
func (c *KubernetesModelConverter) fillMethodConfig(r *model.Route, annotations map[string]string) {
	if annotations == nil {
		return
	}

	if method, ok := annotations[AnnotationKeyMethod]; ok {
		r.Methods = []string{method}
	}
}

// fillHeaderControlConfig fills header control configuration.
func (c *KubernetesModelConverter) fillHeaderControlConfig(r *model.Route, annotations map[string]string) {
	if annotations == nil {
		return
	}

	hasHeaderControl := false
	config := &route.HeaderControlConfig{}

	if add, ok := annotations[AnnotationKeyRequestHeaderAdd]; ok {
		config.RequestAddHeaders = parseHeaderConfig(add)
		hasHeaderControl = true
	}
	if remove, ok := annotations[AnnotationKeyRequestHeaderRemove]; ok {
		config.RequestRemoveHeaders = strings.Split(remove, ",")
		hasHeaderControl = true
	}
	if add, ok := annotations[AnnotationKeyResponseHeaderAdd]; ok {
		config.ResponseAddHeaders = parseHeaderConfig(add)
		hasHeaderControl = true
	}
	if remove, ok := annotations[AnnotationKeyResponseHeaderRemove]; ok {
		config.ResponseRemoveHeaders = strings.Split(remove, ",")
		hasHeaderControl = true
	}

	if hasHeaderControl {
		r.HeaderControl = config
	}
}

// fillRouteCors fills CORS configuration.
func (c *KubernetesModelConverter) fillRouteCors(r *model.Route, metadata *metav1.ObjectMeta) {
	if metadata == nil || metadata.Annotations == nil {
		return
	}

	annotations := metadata.Annotations
	cors := &route.CorsConfig{}

	if maxAge, ok := annotations[AnnotationKeyCorsMaxAge]; ok {
		if age, err := strconv.Atoi(maxAge); err == nil {
			cors.MaxAge = &age
		}
	}
	if credentials, ok := annotations[AnnotationKeyCorsAllowCredentials]; ok {
		allowCreds := strings.EqualFold(credentials, "true")
		cors.AllowCredentials = &allowCreds
	}
	if origins, ok := annotations[AnnotationKeyCorsAllowOrigin]; ok {
		cors.AllowOrigins = strings.Split(origins, ",")
	}
	if methods, ok := annotations[AnnotationKeyCorsAllowMethods]; ok {
		cors.AllowMethods = strings.Split(methods, ",")
	}
	if headers, ok := annotations[AnnotationKeyCorsAllowHeaders]; ok {
		cors.AllowHeaders = strings.Split(headers, ",")
	}
	if exposeHeaders, ok := annotations[AnnotationKeyCorsExposeHeaders]; ok {
		cors.ExposeHeaders = strings.Split(exposeHeaders, ",")
	}

	r.CORS = cors
}

// fillCustomConfigs fills custom configurations.
func (c *KubernetesModelConverter) fillCustomConfigs(r *model.Route, ingress *networkingv1.Ingress) {
	if ingress == nil || ingress.ObjectMeta.Annotations == nil {
		return
	}

	customConfigs := make(map[string]string)
	for key, value := range ingress.ObjectMeta.Annotations {
		if c.isCustomAnnotation(key) {
			customConfigs[key] = value
		}
	}

	if len(customConfigs) > 0 {
		r.CustomConfigs = customConfigs
	}
}

// fillCustomLabels fills custom labels.
func (c *KubernetesModelConverter) fillCustomLabels(r *model.Route, ingress *networkingv1.Ingress) {
	if ingress == nil || ingress.ObjectMeta.Labels == nil {
		return
	}

	customLabels := make(map[string]string)
	for key, value := range ingress.ObjectMeta.Labels {
		if c.isCustomLabel(key) {
			customLabels[key] = value
		}
	}

	if len(customLabels) > 0 {
		r.CustomLabels = customLabels
	}
}

// fillIngressMetadata fills ingress metadata from route.
func (c *KubernetesModelConverter) fillIngressMetadata(ingress *networkingv1.Ingress, r *model.Route) {
	if ingress == nil || r == nil {
		return
	}

	ingress.Name = r.Name
	if r.Version != "" {
		ingress.ResourceVersion = r.Version
	}

	// Set default labels
	if ingress.Labels == nil {
		ingress.Labels = make(map[string]string)
	}
	ingress.Labels[constant.LabelResourceDefinerKey] = constant.LabelResourceDefinerValue
}

// fillIngressSpec fills ingress spec from route.
func (c *KubernetesModelConverter) fillIngressSpec(ingress *networkingv1.Ingress, r *model.Route) error {
	if ingress == nil || r == nil {
		return nil
	}

	// Build rule
	rule := networkingv1.IngressRule{}
	if len(r.Domains) > 0 {
		rule.Host = r.Domains[0]
	}

	// Build HTTP rule
	rule.HTTP = &networkingv1.HTTPIngressRuleValue{
		Paths: make([]networkingv1.HTTPIngressPath, 0),
	}

	// Build path
	path := networkingv1.HTTPIngressPath{}
	if r.Path != nil {
		path.Path = r.Path.Path
		switch r.Path.MatchType {
		case route.MatchTypeExact:
			pathType := networkingv1.PathTypeExact
			path.PathType = &pathType
		case route.MatchTypeRegex:
			pathType := networkingv1.PathTypePrefix
			path.PathType = &pathType
			// Add use-regex annotation
			if ingress.Annotations == nil {
				ingress.Annotations = make(map[string]string)
			}
			ingress.Annotations[AnnotationKeyUseRegex] = "true"
		default:
			pathType := networkingv1.PathTypePrefix
			path.PathType = &pathType
		}
	}

	// Build backend
	if len(r.Services) > 0 {
		// Build destination annotation
		var destBuilder strings.Builder
		for i, svc := range r.Services {
			if i > 0 {
				destBuilder.WriteString("\n")
			}
			if svc.Weight != nil && *svc.Weight != DefaultWeight {
				destBuilder.WriteString(strconv.Itoa(*svc.Weight))
				destBuilder.WriteString("% ")
			}
			destBuilder.WriteString(svc.Namespace)
			destBuilder.WriteString("/")
			destBuilder.WriteString(svc.Name)
			destBuilder.WriteString(":")
			destBuilder.WriteString(strconv.Itoa(svc.Port))
		}
		if ingress.Annotations == nil {
			ingress.Annotations = make(map[string]string)
		}
		ingress.Annotations[AnnotationKeyDestination] = destBuilder.String()
	}

	rule.HTTP.Paths = append(rule.HTTP.Paths, path)
	ingress.Spec.Rules = append(ingress.Spec.Rules, rule)

	return nil
}

// fillIngressCors fills ingress CORS configuration from route.
func (c *KubernetesModelConverter) fillIngressCors(ingress *networkingv1.Ingress, r *model.Route) {
	if ingress == nil || r == nil || r.CORS == nil {
		return
	}

	cors := r.CORS
	if ingress.Annotations == nil {
		ingress.Annotations = make(map[string]string)
	}

	if cors.AllowCredentials != nil {
		ingress.Annotations[AnnotationKeyCorsAllowCredentials] = strconv.FormatBool(*cors.AllowCredentials)
	}
	if cors.MaxAge != nil {
		ingress.Annotations[AnnotationKeyCorsMaxAge] = strconv.Itoa(*cors.MaxAge)
	}
	if len(cors.AllowOrigins) > 0 {
		ingress.Annotations[AnnotationKeyCorsAllowOrigin] = strings.Join(cors.AllowOrigins, ",")
	}
	if len(cors.AllowMethods) > 0 {
		ingress.Annotations[AnnotationKeyCorsAllowMethods] = strings.Join(cors.AllowMethods, ",")
	}
	if len(cors.AllowHeaders) > 0 {
		ingress.Annotations[AnnotationKeyCorsAllowHeaders] = strings.Join(cors.AllowHeaders, ",")
	}
	if len(cors.ExposeHeaders) > 0 {
		ingress.Annotations[AnnotationKeyCorsExposeHeaders] = strings.Join(cors.ExposeHeaders, ",")
	}
}

// fillIngressAnnotations fills ingress annotations from route.
func (c *KubernetesModelConverter) fillIngressAnnotations(ingress *networkingv1.Ingress, r *model.Route) error {
	if ingress == nil || r == nil {
		return nil
	}

	if len(r.CustomConfigs) == 0 {
		return nil
	}

	if ingress.Annotations == nil {
		ingress.Annotations = make(map[string]string)
	}

	for key, value := range r.CustomConfigs {
		if !c.isCustomAnnotation(key) {
			return errors.NewValidationError(fmt.Sprintf("annotation [%s] is already supported by Console", key))
		}
		ingress.Annotations[key] = value
	}

	return nil
}

// fillIngressLabels fills ingress labels from route.
func (c *KubernetesModelConverter) fillIngressLabels(ingress *networkingv1.Ingress, r *model.Route) {
	if ingress == nil || r == nil {
		return
	}

	if len(r.CustomLabels) == 0 {
		return
	}

	if ingress.Labels == nil {
		ingress.Labels = make(map[string]string)
	}

	for key, value := range r.CustomLabels {
		ingress.Labels[key] = value
	}
}

// fillTlsCertificateDetails fills TLS certificate details.
func (c *KubernetesModelConverter) fillTlsCertificateDetails(cert *model.TlsCertificate) {
	// This would typically parse the certificate to extract validity dates
	// For now, we leave it as a placeholder
	// In a real implementation, you would use crypto/x509 to parse the cert
}

// isCustomAnnotation checks if an annotation key is a custom annotation.
func (c *KubernetesModelConverter) isCustomAnnotation(key string) bool {
	// Check if it's a known annotation
	knownAnnotations := map[string]bool{
		AnnotationKeyUseRegex:                 true,
		AnnotationKeyDestination:              true,
		AnnotationKeySSLRedirect:              true,
		AnnotationKeyRewriteEnabled:           true,
		AnnotationKeyRewritePath:              true,
		AnnotationKeyRewriteTarget:            true,
		AnnotationKeyUpstreamVhost:            true,
		AnnotationKeyProxyNextUpstream:        true,
		AnnotationKeyProxyNextUpstreamTries:   true,
		AnnotationKeyProxyNextUpstreamTimeout: true,
		AnnotationKeyHeaderControlEnabled:     true,
		AnnotationKeyRequestHeaderAdd:         true,
		AnnotationKeyRequestHeaderUpdate:      true,
		AnnotationKeyRequestHeaderRemove:      true,
		AnnotationKeyResponseHeaderAdd:        true,
		AnnotationKeyResponseHeaderUpdate:     true,
		AnnotationKeyResponseHeaderRemove:     true,
		AnnotationKeyCorsEnabled:              true,
		AnnotationKeyCorsAllowOrigin:          true,
		AnnotationKeyCorsAllowMethods:         true,
		AnnotationKeyCorsAllowHeaders:         true,
		AnnotationKeyCorsExposeHeaders:        true,
		AnnotationKeyCorsAllowCredentials:     true,
		AnnotationKeyCorsMaxAge:               true,
		AnnotationKeyMethod:                   true,
		AnnotationKeyIgnorePathCase:           true,
		AnnotationKeyComment:                  true,
	}

	// Check for header/query match patterns
	if strings.Contains(key, AnnotationKeyHeaderMatchKeyword) ||
		strings.Contains(key, AnnotationKeyQueryMatchKeyword) {
		return false
	}

	return !knownAnnotations[key] && strings.HasPrefix(key, AnnotationKeyPrefix)
}

// isCustomLabel checks if a label key is a custom label.
func (c *KubernetesModelConverter) isCustomLabel(key string) bool {
	// Known label prefixes
	knownLabelPrefixes := []string{
		constant.LabelKeyPrefix,
		"kubernetes.io/",
		"app.kubernetes.io/",
	}

	for _, prefix := range knownLabelPrefixes {
		if strings.HasPrefix(key, prefix) {
			return false
		}
	}

	return true
}

// =============================================================================
// Utility Functions
// =============================================================================

// normalizeDomainName normalizes a domain name for use in labels.
func normalizeDomainName(domain string) string {
	// Replace dots and wildcards with dashes
	domain = strings.ReplaceAll(domain, ".", "-")
	domain = strings.ReplaceAll(domain, "*", "wildcard")
	domain = strings.ToLower(domain)
	return domain
}

// parseHeaderConfig parses header configuration string.
func parseHeaderConfig(config string) map[string]string {
	result := make(map[string]string)
	if config == "" {
		return result
	}

	// Format: "header1:value1,header2:value2"
	pairs := strings.Split(config, ",")
	for _, pair := range pairs {
		parts := strings.SplitN(pair, ":", 2)
		if len(parts) == 2 {
			result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	return result
}

// matchRuleKeyEquals checks if two match rules have the same key.
func matchRuleKeyEquals(r1, r2 *wasm.MatchRule) bool {
	if r1 == nil || r2 == nil {
		return false
	}

	// Compare domains
	if r1.Domain != r2.Domain {
		return false
	}

	// Compare ingresses
	if r1.Ingress != r2.Ingress {
		return false
	}

	// Compare services
	if r1.Service != r2.Service {
		return false
	}

	return true
}

// sortWasmPluginMatchRules sorts match rules by priority.
func sortWasmPluginMatchRules(rules []*wasm.MatchRule) {
	if len(rules) == 0 {
		return
	}

	sort.Slice(rules, func(i, j int) bool {
		return compareMatchRules(rules[i], rules[j]) < 0
	})
}

// compareMatchRules compares two match rules for sorting.
func compareMatchRules(r1, r2 *wasm.MatchRule) int {
	hasDomain1 := r1.Domain != ""
	hasDomain2 := r2.Domain != ""
	hasIngress1 := r1.Ingress != ""
	hasIngress2 := r2.Ingress != ""
	hasService1 := r1.Service != ""
	hasService2 := r2.Service != ""

	empty1 := !hasDomain1 && !hasIngress1 && !hasService1
	empty2 := !hasDomain2 && !hasIngress2 && !hasService2

	if empty1 && empty2 {
		return 0
	}
	if empty1 != empty2 {
		if empty1 {
			return 1
		}
		return -1
	}

	// Service rules come first
	if hasService1 != hasService2 {
		if hasService1 {
			return -1
		}
		return 1
	}

	if hasService1 {
		return strings.Compare(r1.Service, r2.Service)
	}

	// Ingress rules come next
	if hasIngress1 != hasIngress2 {
		if hasIngress1 {
			return -1
		}
		return 1
	}

	if hasIngress1 {
		return strings.Compare(r1.Ingress, r2.Ingress)
	}

	// Domain rules come last
	if hasDomain1 != hasDomain2 {
		if hasDomain1 {
			return 1
		}
		return -1
	}

	if hasDomain1 {
		return strings.Compare(r1.Domain, r2.Domain)
	}

	return 0
}

// Base64Encode encodes a string to base64.
func Base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// Base64Decode decodes a base64 string.
func Base64Decode(data string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// ToYAML converts an object to YAML string.
func ToYAML(obj interface{}) (string, error) {
	data, err := jsonFast.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ParseURL parses a URL string.
func ParseURL(rawURL string) (*url.URL, error) {
	return url.Parse(rawURL)
}

// CompileRegex compiles a regex pattern.
func CompileRegex(pattern string) (*regexp.Regexp, error) {
	return regexp.Compile(pattern)
}

// MapToInterface converts a map[string]string to map[string]interface{}.
func MapToInterface(m map[string]string) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		result[k] = v
	}
	return result
}

// InterfaceToMap converts a map[string]interface{} to map[string]string.
func InterfaceToMap(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		if s, ok := v.(string); ok {
			result[k] = s
		} else if b, err := jsonFast.Marshal(v); err == nil {
			result[k] = string(b)
		}
	}
	return result
}
