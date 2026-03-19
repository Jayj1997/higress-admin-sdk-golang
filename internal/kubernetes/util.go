// Package kubernetes provides Kubernetes client functionality
package kubernetes

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KubernetesUtil provides utility functions for Kubernetes resources
type KubernetesUtil struct{}

// NewKubernetesUtil creates a new KubernetesUtil
func NewKubernetesUtil() *KubernetesUtil {
	return &KubernetesUtil{}
}

// GetIngressHosts extracts hosts from ingress rules
func (u *KubernetesUtil) GetIngressHosts(ingress *networkingv1.Ingress) []string {
	if ingress == nil || ingress.Spec.Rules == nil {
		return nil
	}
	var hosts []string
	for _, rule := range ingress.Spec.Rules {
		if rule.Host != "" {
			hosts = append(hosts, rule.Host)
		}
	}
	return hosts
}

// GetIngressPaths extracts paths from ingress rules
func (u *KubernetesUtil) GetIngressPaths(ingress *networkingv1.Ingress) []string {
	if ingress == nil || ingress.Spec.Rules == nil {
		return nil
	}
	var paths []string
	for _, rule := range ingress.Spec.Rules {
		if rule.HTTP != nil {
			for _, path := range rule.HTTP.Paths {
				paths = append(paths, path.Path)
			}
		}
	}
	return paths
}

// GetIngressBackend extracts backend service name and port from ingress
func (u *KubernetesUtil) GetIngressBackend(ingress *networkingv1.Ingress) (serviceName string, servicePort int32) {
	if ingress == nil {
		return "", 0
	}

	// Check default backend
	if ingress.Spec.DefaultBackend != nil && ingress.Spec.DefaultBackend.Service != nil {
		return ingress.Spec.DefaultBackend.Service.Name, ingress.Spec.DefaultBackend.Service.Port.Number
	}

	// Check rules
	for _, rule := range ingress.Spec.Rules {
		if rule.HTTP != nil {
			for _, path := range rule.HTTP.Paths {
				if path.Backend.Service != nil {
					return path.Backend.Service.Name, path.Backend.Service.Port.Number
				}
			}
		}
	}

	return "", 0
}

// GetSecretType returns the type of a secret
func (u *KubernetesUtil) GetSecretType(secret *corev1.Secret) string {
	if secret == nil {
		return ""
	}
	return string(secret.Type)
}

// GetSecretData returns the data of a secret as a map
func (u *KubernetesUtil) GetSecretData(secret *corev1.Secret) map[string]string {
	if secret == nil || secret.Data == nil {
		return nil
	}
	data := make(map[string]string)
	for k, v := range secret.Data {
		data[k] = string(v)
	}
	return data
}

// GetConfigMapData returns the data of a configmap as a map
func (u *KubernetesUtil) GetConfigMapData(cm *corev1.ConfigMap) map[string]string {
	if cm == nil || cm.Data == nil {
		return nil
	}
	return cm.Data
}

// IsIngressReady checks if an ingress is ready
func (u *KubernetesUtil) IsIngressReady(ingress *networkingv1.Ingress) bool {
	if ingress == nil {
		return false
	}
	for _, cond := range ingress.Status.LoadBalancer.Ingress {
		if cond.Hostname != "" || cond.IP != "" {
			return true
		}
	}
	return false
}

// GetIngressLoadBalancer returns the load balancer hostname or IP
func (u *KubernetesUtil) GetIngressLoadBalancer(ingress *networkingv1.Ingress) string {
	if ingress == nil {
		return ""
	}
	for _, ing := range ingress.Status.LoadBalancer.Ingress {
		if ing.Hostname != "" {
			return ing.Hostname
		}
		if ing.IP != "" {
			return ing.IP
		}
	}
	return ""
}

// BuildIngressName builds an ingress name from domain and path
func (u *KubernetesUtil) BuildIngressName(domain, path string) string {
	if path == "" || path == "/" {
		return domain
	}
	// Replace special characters
	path = strings.TrimPrefix(path, "/")
	path = strings.ReplaceAll(path, "/", "-")
	path = strings.ReplaceAll(path, "_", "-")
	return fmt.Sprintf("%s-%s", domain, path)
}

// BuildSecretName builds a secret name from domain
func (u *KubernetesUtil) BuildSecretName(domain string) string {
	return strings.ReplaceAll(domain, ".", "-")
}

// BuildConfigMapName builds a configmap name from name and namespace
func (u *KubernetesUtil) BuildConfigMapName(name, namespace string) string {
	if namespace == "" {
		return name
	}
	return fmt.Sprintf("%s-%s", namespace, name)
}

// ParseIngressClass parses the ingress class name
func (u *KubernetesUtil) ParseIngressClass(ingress *networkingv1.Ingress) string {
	if ingress == nil {
		return ""
	}
	if ingress.Spec.IngressClassName != nil {
		return *ingress.Spec.IngressClassName
	}
	return ""
}

// GetIngressAnnotations returns the annotations of an ingress
func (u *KubernetesUtil) GetIngressAnnotations(ingress *networkingv1.Ingress) map[string]string {
	if ingress == nil || ingress.Annotations == nil {
		return nil
	}
	return ingress.Annotations
}

// GetIngressLabels returns the labels of an ingress
func (u *KubernetesUtil) GetIngressLabels(ingress *networkingv1.Ingress) map[string]string {
	if ingress == nil || ingress.Labels == nil {
		return nil
	}
	return ingress.Labels
}

// MergeLabels merges labels into existing labels
func (u *KubernetesUtil) MergeLabels(existing, newLabels map[string]string) map[string]string {
	if existing == nil {
		existing = make(map[string]string)
	}
	for k, v := range newLabels {
		existing[k] = v
	}
	return existing
}

// MergeAnnotations merges annotations into existing annotations
func (u *KubernetesUtil) MergeAnnotations(existing, newAnnotations map[string]string) map[string]string {
	if existing == nil {
		existing = make(map[string]string)
	}
	for k, v := range newAnnotations {
		existing[k] = v
	}
	return existing
}

// CreateObjectMeta creates ObjectMeta with basic fields
func (u *KubernetesUtil) CreateObjectMeta(name, namespace string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
	}
}

// CreateObjectMetaWithLabels creates ObjectMeta with labels
func (u *KubernetesUtil) CreateObjectMetaWithLabels(name, namespace string, labels map[string]string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		Labels:    labels,
	}
}

// CreateObjectMetaWithAnnotations creates ObjectMeta with annotations
func (u *KubernetesUtil) CreateObjectMetaWithAnnotations(name, namespace string, annotations map[string]string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:        name,
		Namespace:   namespace,
		Annotations: annotations,
	}
}

// CreateObjectMetaFull creates ObjectMeta with labels and annotations
func (u *KubernetesUtil) CreateObjectMetaFull(name, namespace string, labels, annotations map[string]string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:        name,
		Namespace:   namespace,
		Labels:      labels,
		Annotations: annotations,
	}
}

// IsSecretType checks if a secret is of a specific type
func (u *KubernetesUtil) IsSecretType(secret *corev1.Secret, secretType corev1.SecretType) bool {
	if secret == nil {
		return false
	}
	return secret.Type == secretType
}

// IsTLSSecret checks if a secret is a TLS secret
func (u *KubernetesUtil) IsTLSSecret(secret *corev1.Secret) bool {
	return u.IsSecretType(secret, corev1.SecretTypeTLS)
}

// IsOpaqueSecret checks if a secret is an opaque secret
func (u *KubernetesUtil) IsOpaqueSecret(secret *corev1.Secret) bool {
	return u.IsSecretType(secret, corev1.SecretTypeOpaque)
}

// IsDockerConfigSecret checks if a secret is a docker config secret
func (u *KubernetesUtil) IsDockerConfigSecret(secret *corev1.Secret) bool {
	return u.IsSecretType(secret, corev1.SecretTypeDockerConfigJson)
}

// GetTLSCertificate returns the TLS certificate from a secret
func (u *KubernetesUtil) GetTLSCertificate(secret *corev1.Secret) string {
	if secret == nil || secret.Data == nil {
		return ""
	}
	return string(secret.Data[corev1.TLSCertKey])
}

// GetTLSPrivateKey returns the TLS private key from a secret
func (u *KubernetesUtil) GetTLSPrivateKey(secret *corev1.Secret) string {
	if secret == nil || secret.Data == nil {
		return ""
	}
	return string(secret.Data[corev1.TLSPrivateKeyKey])
}

// CreateTLSSecret creates a TLS secret
func (u *KubernetesUtil) CreateTLSSecret(name, namespace, cert, key string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Type: corev1.SecretTypeTLS,
		Data: map[string][]byte{
			corev1.TLSCertKey:       []byte(cert),
			corev1.TLSPrivateKeyKey: []byte(key),
		},
	}
}

// CreateOpaqueSecret creates an opaque secret
func (u *KubernetesUtil) CreateOpaqueSecret(name, namespace string, data map[string]string) *corev1.Secret {
	secretData := make(map[string][]byte)
	for k, v := range data {
		secretData[k] = []byte(v)
	}
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Type: corev1.SecretTypeOpaque,
		Data: secretData,
	}
}

// CreateConfigMap creates a configmap
func (u *KubernetesUtil) CreateConfigMap(name, namespace string, data map[string]string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}
}

// CreateIngressBackend creates an IngressBackend
func (u *KubernetesUtil) CreateIngressBackend(serviceName string, servicePort int32) *networkingv1.IngressBackend {
	return &networkingv1.IngressBackend{
		Service: &networkingv1.IngressServiceBackend{
			Name: serviceName,
			Port: networkingv1.ServiceBackendPort{
				Number: servicePort,
			},
		},
	}
}

// CreateHTTPIngressPath creates an HTTPIngressPath
func (u *KubernetesUtil) CreateHTTPIngressPath(path, pathType string, backend *networkingv1.IngressBackend) networkingv1.HTTPIngressPath {
	pt := networkingv1.PathTypeImplementationSpecific
	switch pathType {
	case "Exact":
		pt = networkingv1.PathTypeExact
	case "Prefix":
		pt = networkingv1.PathTypePrefix
	}
	return networkingv1.HTTPIngressPath{
		Path:     path,
		PathType: &pt,
		Backend:  *backend,
	}
}

// CreateHTTPIngressRuleValue creates an HTTPIngressRuleValue
func (u *KubernetesUtil) CreateHTTPIngressRuleValue(paths []networkingv1.HTTPIngressPath) *networkingv1.HTTPIngressRuleValue {
	return &networkingv1.HTTPIngressRuleValue{
		Paths: paths,
	}
}

// CreateIngressRule creates an IngressRule
func (u *KubernetesUtil) CreateIngressRule(host string, http *networkingv1.HTTPIngressRuleValue) networkingv1.IngressRule {
	return networkingv1.IngressRule{
		Host: host,
		IngressRuleValue: networkingv1.IngressRuleValue{
			HTTP: http,
		},
	}
}

// CreateIngressSpec creates an IngressSpec
func (u *KubernetesUtil) CreateIngressSpec(rules []networkingv1.IngressRule, defaultBackend *networkingv1.IngressBackend, ingressClassName *string) *networkingv1.IngressSpec {
	return &networkingv1.IngressSpec{
		Rules:            rules,
		DefaultBackend:   defaultBackend,
		IngressClassName: ingressClassName,
	}
}

// CreateIngress creates an Ingress
func (u *KubernetesUtil) CreateIngress(name, namespace string, spec *networkingv1.IngressSpec) *networkingv1.Ingress {
	return &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: *spec,
	}
}
