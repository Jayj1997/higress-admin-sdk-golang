// Package kubernetes provides Kubernetes client functionality
package kubernetes

import (
	"testing"

	"github.com/Jayj1997/higress-admin-sdk-golang/internal/kubernetes/crd/wasm"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/constant"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model/route"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestIngressToRoute tests the IngressToRoute conversion
func TestIngressToRoute(t *testing.T) {
	converter := &KubernetesModelConverter{}

	tests := []struct {
		name     string
		ingress  *networkingv1.Ingress
		wantErr  bool
		wantName string
	}{
		{
			name:    "nil ingress",
			ingress: nil,
			wantErr: true,
		},
		{
			name: "valid ingress",
			ingress: &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "test-route",
					ResourceVersion: "12345",
					Annotations: map[string]string{
						AnnotationKeyDestination: "default/my-service:8080",
					},
					Labels: map[string]string{
						constant.LabelResourceDefinerKey: constant.LabelResourceDefinerValue,
					},
				},
				Spec: networkingv1.IngressSpec{
					Rules: []networkingv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networkingv1.IngressRuleValue{
								HTTP: &networkingv1.HTTPIngressRuleValue{
									Paths: []networkingv1.HTTPIngressPath{
										{
											Path:     "/api",
											PathType: pathTypePtr(networkingv1.PathTypePrefix),
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr:  false,
			wantName: "test-route",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.IngressToRoute(tt.ingress)
			if (err != nil) != tt.wantErr {
				t.Errorf("IngressToRoute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Name != tt.wantName {
				t.Errorf("IngressToRoute() name = %v, want %v", got.Name, tt.wantName)
			}
		})
	}
}

// TestRouteToIngress tests the RouteToIngress conversion
func TestRouteToIngress(t *testing.T) {
	converter := &KubernetesModelConverter{}

	tests := []struct {
		name    string
		route   *model.Route
		wantErr bool
	}{
		{
			name:    "nil route",
			route:   nil,
			wantErr: true,
		},
		{
			name: "valid route",
			route: &model.Route{
				Name:    "test-route",
				Version: "12345",
				Domains: []string{"example.com"},
				Path: &route.RoutePredicate{
					Path:      "/api",
					MatchType: route.MatchTypePrefix,
				},
				Services: []*route.UpstreamService{
					{
						Name:      "my-service",
						Namespace: "default",
						Port:      8080,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.RouteToIngress(tt.route)
			if (err != nil) != tt.wantErr {
				t.Errorf("RouteToIngress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Name != tt.route.Name {
					t.Errorf("RouteToIngress() name = %v, want %v", got.Name, tt.route.Name)
				}
			}
		})
	}
}

// TestConfigMapToServiceSource tests the ConfigMapToServiceSource conversion
func TestConfigMapToServiceSource(t *testing.T) {
	converter := &KubernetesModelConverter{}

	tests := []struct {
		name      string
		configMap *corev1.ConfigMap
		wantErr   bool
		wantName  string
	}{
		{
			name:      "nil configmap",
			configMap: nil,
			wantErr:   true,
		},
		{
			name: "valid configmap",
			configMap: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "test-source",
					ResourceVersion: "12345",
				},
				Data: map[string]string{
					"domain": "nacos.example.com",
					"port":   "8848",
					"type":   "nacos",
				},
			},
			wantErr:  false,
			wantName: "test-source",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.ConfigMapToServiceSource(tt.configMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigMapToServiceSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Name != tt.wantName {
				t.Errorf("ConfigMapToServiceSource() name = %v, want %v", got.Name, tt.wantName)
			}
		})
	}
}

// TestServiceSourceToConfigMap tests the ServiceSourceToConfigMap conversion
func TestServiceSourceToConfigMap(t *testing.T) {
	converter := &KubernetesModelConverter{}

	port := 8848
	tests := []struct {
		name    string
		source  *model.ServiceSource
		wantErr bool
	}{
		{
			name:    "nil source",
			source:  nil,
			wantErr: true,
		},
		{
			name: "valid source",
			source: &model.ServiceSource{
				Name:    "test-source",
				Version: "12345",
				Type:    "nacos",
				Domain:  "nacos.example.com",
				Port:    &port,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.ServiceSourceToConfigMap(tt.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceSourceToConfigMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Name != tt.source.Name {
				t.Errorf("ServiceSourceToConfigMap() name = %v, want %v", got.Name, tt.source.Name)
			}
		})
	}
}

// TestSecretToTlsCertificate tests the SecretToTlsCertificate conversion
func TestSecretToTlsCertificate(t *testing.T) {
	converter := &KubernetesModelConverter{}

	tests := []struct {
		name     string
		secret   *corev1.Secret
		wantErr  bool
		wantName string
	}{
		{
			name:    "nil secret",
			secret:  nil,
			wantErr: true,
		},
		{
			name: "valid secret",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "test-cert",
					ResourceVersion: "12345",
					Labels: map[string]string{
						LabelKeyDomainPrefix + "example-com": "true",
					},
				},
				Type: corev1.SecretTypeTLS,
				Data: map[string][]byte{
					SecretTLSCrtField: []byte("-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----"),
					SecretTLSKeyField: []byte("-----BEGIN PRIVATE KEY-----\ntest\n-----END PRIVATE KEY-----"),
				},
			},
			wantErr:  false,
			wantName: "test-cert",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.SecretToTlsCertificate(tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("SecretToTlsCertificate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Name != tt.wantName {
				t.Errorf("SecretToTlsCertificate() name = %v, want %v", got.Name, tt.wantName)
			}
		})
	}
}

// TestTlsCertificateToSecret tests the TlsCertificateToSecret conversion
func TestTlsCertificateToSecret(t *testing.T) {
	converter := &KubernetesModelConverter{}

	tests := []struct {
		name    string
		cert    *model.TlsCertificate
		wantErr bool
	}{
		{
			name:    "nil cert",
			cert:    nil,
			wantErr: true,
		},
		{
			name: "valid cert",
			cert: &model.TlsCertificate{
				Name:    "test-cert",
				Version: "12345",
				Domains: []string{"example.com"},
				Cert:    "-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----",
				Key:     "-----BEGIN PRIVATE KEY-----\ntest\n-----END PRIVATE KEY-----",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.TlsCertificateToSecret(tt.cert)
			if (err != nil) != tt.wantErr {
				t.Errorf("TlsCertificateToSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Name != tt.cert.Name {
				t.Errorf("TlsCertificateToSecret() name = %v, want %v", got.Name, tt.cert.Name)
			}
		})
	}
}

// TestConfigMapToDomain tests the ConfigMapToDomain conversion
func TestConfigMapToDomain(t *testing.T) {
	converter := &KubernetesModelConverter{}

	tests := []struct {
		name      string
		configMap *corev1.ConfigMap
		wantErr   bool
		wantName  string
	}{
		{
			name:      "nil configmap",
			configMap: nil,
			wantErr:   true,
		},
		{
			name: "valid configmap",
			configMap: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "domain-example-com",
					ResourceVersion: "12345",
				},
				Data: map[string]string{
					ConfigMapKeyDomain:      "example.com",
					ConfigMapKeyCert:        "my-cert",
					ConfigMapKeyEnableHTTPS: "on",
				},
			},
			wantErr:  false,
			wantName: "example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.ConfigMapToDomain(tt.configMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigMapToDomain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Name != tt.wantName {
				t.Errorf("ConfigMapToDomain() name = %v, want %v", got.Name, tt.wantName)
			}
		})
	}
}

// TestDomainToConfigMap tests the DomainToConfigMap conversion
func TestDomainToConfigMap(t *testing.T) {
	converter := &KubernetesModelConverter{}

	tests := []struct {
		name    string
		domain  *model.Domain
		wantErr bool
	}{
		{
			name:    "nil domain",
			domain:  nil,
			wantErr: true,
		},
		{
			name: "valid domain",
			domain: &model.Domain{
				Name:           "example.com",
				Version:        "12345",
				EnableHTTPS:    "on",
				CertIdentifier: "my-cert",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.DomainToConfigMap(tt.domain)
			if (err != nil) != tt.wantErr {
				t.Errorf("DomainToConfigMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				expectedName := converter.DomainNameToConfigMapName(tt.domain.Name)
				if got.Name != expectedName {
					t.Errorf("DomainToConfigMap() name = %v, want %v", got.Name, expectedName)
				}
			}
		})
	}
}

// TestWasmPluginCRDToModel tests the WasmPluginCRDToModel conversion
func TestWasmPluginCRDToModel(t *testing.T) {
	converter := &KubernetesModelConverter{}

	tests := []struct {
		name    string
		plugin  *wasm.V1alpha1WasmPlugin
		wantErr bool
	}{
		{
			name:    "nil plugin",
			plugin:  nil,
			wantErr: true,
		},
		{
			name: "valid plugin",
			plugin: &wasm.V1alpha1WasmPlugin{
				APIVersion: wasm.WasmPluginAPIVersion,
				Kind:       wasm.WasmPluginKind,
				Metadata: &wasm.V1ObjectMeta{
					Name:            "my-plugin-1.0.0",
					ResourceVersion: "12345",
					Labels: map[string]string{
						LabelKeyWasmPluginName:     "my-plugin",
						LabelKeyWasmPluginVersion:  "1.0.0",
						LabelKeyWasmPluginCategory: "auth",
						LabelKeyWasmPluginBuiltIn:  "true",
					},
					Annotations: map[string]string{
						AnnotationKeyWasmPluginTitle:       "My Plugin",
						AnnotationKeyWasmPluginDescription: "A test plugin",
						AnnotationKeyWasmPluginIcon:        "https://example.com/icon.png",
					},
				},
				Spec: &wasm.V1alpha1WasmPluginSpec{
					Phase: wasm.PluginPhaseAuthn,
					Url:   "oci://registry.example.com/my-plugin:1.0.0",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.WasmPluginCRDToModel(tt.plugin)
			if (err != nil) != tt.wantErr {
				t.Errorf("WasmPluginCRDToModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.plugin != nil {
				if got.Name != tt.plugin.Metadata.Labels[LabelKeyWasmPluginName] {
					t.Errorf("WasmPluginCRDToModel() name = %v, want %v", got.Name, tt.plugin.Metadata.Labels[LabelKeyWasmPluginName])
				}
			}
		})
	}
}

// TestModelToWasmPluginCRD tests the ModelToWasmPluginCRD conversion
func TestModelToWasmPluginCRD(t *testing.T) {
	converter := &KubernetesModelConverter{}
	builtIn := true
	priority := 100

	tests := []struct {
		name    string
		plugin  *model.WasmPlugin
		wantErr bool
	}{
		{
			name:    "nil plugin",
			plugin:  nil,
			wantErr: true,
		},
		{
			name: "valid plugin",
			plugin: &model.WasmPlugin{
				Name:        "my-plugin",
				Version:     "1.0.0",
				Category:    "auth",
				Title:       "My Plugin",
				Description: "A test plugin",
				Icon:        "https://example.com/icon.png",
				BuiltIn:     &builtIn,
				Phase:       "AUTHN",
				Priority:    &priority,
				ImageURL:    "oci://registry.example.com/my-plugin:1.0.0",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.ModelToWasmPluginCRD(tt.plugin)
			if (err != nil) != tt.wantErr {
				t.Errorf("ModelToWasmPluginCRD() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.plugin != nil {
				expectedName := tt.plugin.Name + "-" + tt.plugin.Version
				if got.Metadata.Name != expectedName {
					t.Errorf("ModelToWasmPluginCRD() name = %v, want %v", got.Metadata.Name, expectedName)
				}
			}
		})
	}
}

// TestGetWasmPluginInstancesFromCR tests the GetWasmPluginInstancesFromCR conversion
func TestGetWasmPluginInstancesFromCR(t *testing.T) {
	converter := &KubernetesModelConverter{}

	tests := []struct {
		name      string
		plugin    *wasm.V1alpha1WasmPlugin
		wantErr   bool
		wantCount int
	}{
		{
			name:    "nil plugin",
			plugin:  nil,
			wantErr: true,
		},
		{
			name: "plugin with global config",
			plugin: &wasm.V1alpha1WasmPlugin{
				Metadata: &wasm.V1ObjectMeta{
					Labels: map[string]string{
						LabelKeyWasmPluginName:    "my-plugin",
						LabelKeyWasmPluginVersion: "1.0.0",
					},
				},
				Spec: &wasm.V1alpha1WasmPluginSpec{
					DefaultConfigDisable: false,
					DefaultConfig: map[string]interface{}{
						"key": "value",
					},
				},
			},
			wantErr:   false,
			wantCount: 1,
		},
		{
			name: "plugin with match rules",
			plugin: &wasm.V1alpha1WasmPlugin{
				Metadata: &wasm.V1ObjectMeta{
					Labels: map[string]string{
						LabelKeyWasmPluginName:    "my-plugin",
						LabelKeyWasmPluginVersion: "1.0.0",
					},
				},
				Spec: &wasm.V1alpha1WasmPluginSpec{
					MatchRules: []*wasm.MatchRule{
						{
							Domain: "example.com",
							Enable: true,
							Config: map[string]interface{}{
								"key": "value",
							},
						},
					},
				},
			},
			wantErr:   false,
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.GetWasmPluginInstancesFromCR(tt.plugin)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWasmPluginInstancesFromCR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != tt.wantCount {
				t.Errorf("GetWasmPluginInstancesFromCR() count = %v, want %v", len(got), tt.wantCount)
			}
		})
	}
}

// TestSetWasmPluginInstanceToCR tests the SetWasmPluginInstanceToCR conversion
func TestSetWasmPluginInstanceToCR(t *testing.T) {
	converter := &KubernetesModelConverter{}
	enabled := true

	tests := []struct {
		name     string
		cr       *wasm.V1alpha1WasmPlugin
		instance *model.WasmPluginInstance
		wantErr  bool
	}{
		{
			name:     "nil cr",
			cr:       nil,
			instance: &model.WasmPluginInstance{},
			wantErr:  true,
		},
		{
			name:     "nil instance",
			cr:       wasm.NewV1alpha1WasmPlugin(),
			instance: nil,
			wantErr:  true,
		},
		{
			name: "global instance",
			cr:   wasm.NewV1alpha1WasmPlugin(),
			instance: &model.WasmPluginInstance{
				Targets: map[model.WasmPluginInstanceScope]string{
					model.WasmPluginInstanceScopeGlobal: "",
				},
				Enabled:        &enabled,
				Configurations: map[string]interface{}{"key": "value"},
			},
			wantErr: false,
		},
		{
			name: "domain instance",
			cr:   wasm.NewV1alpha1WasmPlugin(),
			instance: &model.WasmPluginInstance{
				Targets: map[model.WasmPluginInstanceScope]string{
					model.WasmPluginInstanceScopeDomain: "example.com",
				},
				Enabled:        &enabled,
				Configurations: map[string]interface{}{"key": "value"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := converter.SetWasmPluginInstanceToCR(tt.cr, tt.instance)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetWasmPluginInstanceToCR() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestIsIngressSupported tests the IsIngressSupported method
func TestIsIngressSupported(t *testing.T) {
	converter := &KubernetesModelConverter{}

	tests := []struct {
		name     string
		ingress  *networkingv1.Ingress
		expected bool
	}{
		{
			name:     "nil ingress",
			ingress:  nil,
			expected: false,
		},
		{
			name: "no rules",
			ingress: &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec:       networkingv1.IngressSpec{},
			},
			expected: false,
		},
		{
			name: "multiple rules",
			ingress: &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: networkingv1.IngressSpec{
					Rules: []networkingv1.IngressRule{
						{Host: "a.com"},
						{Host: "b.com"},
					},
				},
			},
			expected: false,
		},
		{
			name: "valid single rule single path",
			ingress: &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: networkingv1.IngressSpec{
					Rules: []networkingv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networkingv1.IngressRuleValue{
								HTTP: &networkingv1.HTTPIngressRuleValue{
									Paths: []networkingv1.HTTPIngressPath{
										{
											Path:     "/api",
											PathType: pathTypePtr(networkingv1.PathTypePrefix),
										},
									},
								},
							},
						},
					},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := converter.IsIngressSupported(tt.ingress)
			if got != tt.expected {
				t.Errorf("IsIngressSupported() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestNormalizeDomainName tests the normalizeDomainName function
func TestNormalizeDomainName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"example.com", "example-com"},
		{"*.example.com", "wildcard-example-com"},
		{"TEST.COM", "test-com"},
		{"sub.domain.example.com", "sub-domain-example-com"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := normalizeDomainName(tt.input)
			if got != tt.expected {
				t.Errorf("normalizeDomainName() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestDomainNameToConfigMapName tests the DomainNameToConfigMapName method
func TestDomainNameToConfigMapName(t *testing.T) {
	converter := &KubernetesModelConverter{}

	tests := []struct {
		domain   string
		expected string
	}{
		{"example.com", "domain-example-com"},
		{"*.example.com", "domain-wildcard-example-com"},
	}

	for _, tt := range tests {
		t.Run(tt.domain, func(t *testing.T) {
			got := converter.DomainNameToConfigMapName(tt.domain)
			if got != tt.expected {
				t.Errorf("DomainNameToConfigMapName() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestParseHeaderConfig tests the parseHeaderConfig function
func TestParseHeaderConfig(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]string
	}{
		{"", map[string]string{}},
		{"X-Custom:value1", map[string]string{"X-Custom": "value1"}},
		{"H1:v1,H2:v2", map[string]string{"H1": "v1", "H2": "v2"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := parseHeaderConfig(tt.input)
			for k, v := range tt.expected {
				if got[k] != v {
					t.Errorf("parseHeaderConfig()[%v] = %v, want %v", k, got[k], v)
				}
			}
		})
	}
}

// TestMatchRuleKeyEquals tests the matchRuleKeyEquals function
func TestMatchRuleKeyEquals(t *testing.T) {
	tests := []struct {
		name     string
		r1       *wasm.MatchRule
		r2       *wasm.MatchRule
		expected bool
	}{
		{
			name:     "both nil",
			r1:       nil,
			r2:       nil,
			expected: false,
		},
		{
			name:     "one nil",
			r1:       &wasm.MatchRule{},
			r2:       nil,
			expected: false,
		},
		{
			name:     "equal rules",
			r1:       &wasm.MatchRule{Domain: "example.com"},
			r2:       &wasm.MatchRule{Domain: "example.com"},
			expected: true,
		},
		{
			name:     "different domain",
			r1:       &wasm.MatchRule{Domain: "a.com"},
			r2:       &wasm.MatchRule{Domain: "b.com"},
			expected: false,
		},
		{
			name:     "different ingress",
			r1:       &wasm.MatchRule{Ingress: "route-a"},
			r2:       &wasm.MatchRule{Ingress: "route-b"},
			expected: false,
		},
		{
			name:     "different service",
			r1:       &wasm.MatchRule{Service: "svc-a"},
			r2:       &wasm.MatchRule{Service: "svc-b"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchRuleKeyEquals(tt.r1, tt.r2)
			if got != tt.expected {
				t.Errorf("matchRuleKeyEquals() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestCompareMatchRules tests the compareMatchRules function
func TestCompareMatchRules(t *testing.T) {
	tests := []struct {
		name     string
		r1       *wasm.MatchRule
		r2       *wasm.MatchRule
		expected int
	}{
		{
			name:     "both empty",
			r1:       &wasm.MatchRule{},
			r2:       &wasm.MatchRule{},
			expected: 0,
		},
		{
			name:     "r1 empty",
			r1:       &wasm.MatchRule{},
			r2:       &wasm.MatchRule{Service: "svc"},
			expected: 1,
		},
		{
			name:     "r2 empty",
			r1:       &wasm.MatchRule{Service: "svc"},
			r2:       &wasm.MatchRule{},
			expected: -1,
		},
		{
			name:     "service comes first",
			r1:       &wasm.MatchRule{Service: "svc"},
			r2:       &wasm.MatchRule{Ingress: "route"},
			expected: -1,
		},
		{
			name:     "ingress comes before domain",
			r1:       &wasm.MatchRule{Ingress: "route"},
			r2:       &wasm.MatchRule{Domain: "example.com"},
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := compareMatchRules(tt.r1, tt.r2)
			if got != tt.expected {
				t.Errorf("compareMatchRules() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Helper function to create a pointer to PathType
func pathTypePtr(pt networkingv1.PathType) *networkingv1.PathType {
	return &pt
}
