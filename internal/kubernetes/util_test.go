// Package kubernetes provides Kubernetes client functionality
package kubernetes

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
)

// TestGetIngressHosts tests the GetIngressHosts function
// 测试 GetIngressHosts 函数
// 调用方式: go test -v -run TestGetIngressHosts ./internal/kubernetes/
func TestGetIngressHosts(t *testing.T) {
	util := NewKubernetesUtil()

	tests := []struct {
		name     string
		ingress  *networkingv1.Ingress
		expected []string
	}{
		{
			name:     "nil ingress",
			ingress:  nil,
			expected: nil,
		},
		{
			name: "no rules",
			ingress: &networkingv1.Ingress{
				Spec: networkingv1.IngressSpec{},
			},
			expected: nil,
		},
		{
			name: "single host",
			ingress: &networkingv1.Ingress{
				Spec: networkingv1.IngressSpec{
					Rules: []networkingv1.IngressRule{
						{Host: "example.com"},
					},
				},
			},
			expected: []string{"example.com"},
		},
		{
			name: "multiple hosts",
			ingress: &networkingv1.Ingress{
				Spec: networkingv1.IngressSpec{
					Rules: []networkingv1.IngressRule{
						{Host: "example.com"},
						{Host: "api.example.com"},
						{Host: "www.example.com"},
					},
				},
			},
			expected: []string{"example.com", "api.example.com", "www.example.com"},
		},
		{
			name: "empty host",
			ingress: &networkingv1.Ingress{
				Spec: networkingv1.IngressSpec{
					Rules: []networkingv1.IngressRule{
						{Host: ""},
						{Host: "example.com"},
					},
				},
			},
			expected: []string{"example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.GetIngressHosts(tt.ingress)
			if len(result) != len(tt.expected) {
				t.Errorf("GetIngressHosts() = %v, want %v", result, tt.expected)
				return
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("GetIngressHosts()[%d] = %v, want %v", i, v, tt.expected[i])
				}
			}
		})
	}
}

// TestGetIngressPaths tests the GetIngressPaths function
// 测试 GetIngressPaths 函数
// 调用方式: go test -v -run TestGetIngressPaths ./internal/kubernetes/
func TestGetIngressPaths(t *testing.T) {
	util := NewKubernetesUtil()

	tests := []struct {
		name     string
		ingress  *networkingv1.Ingress
		expected []string
	}{
		{
			name:     "nil ingress",
			ingress:  nil,
			expected: nil,
		},
		{
			name: "no rules",
			ingress: &networkingv1.Ingress{
				Spec: networkingv1.IngressSpec{},
			},
			expected: nil,
		},
		{
			name: "single path",
			ingress: &networkingv1.Ingress{
				Spec: networkingv1.IngressSpec{
					Rules: []networkingv1.IngressRule{
						{
							IngressRuleValue: networkingv1.IngressRuleValue{
								HTTP: &networkingv1.HTTPIngressRuleValue{
									Paths: []networkingv1.HTTPIngressPath{
										{Path: "/"},
									},
								},
							},
						},
					},
				},
			},
			expected: []string{"/"},
		},
		{
			name: "multiple paths",
			ingress: &networkingv1.Ingress{
				Spec: networkingv1.IngressSpec{
					Rules: []networkingv1.IngressRule{
						{
							IngressRuleValue: networkingv1.IngressRuleValue{
								HTTP: &networkingv1.HTTPIngressRuleValue{
									Paths: []networkingv1.HTTPIngressPath{
										{Path: "/"},
										{Path: "/api"},
										{Path: "/health"},
									},
								},
							},
						},
					},
				},
			},
			expected: []string{"/", "/api", "/health"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.GetIngressPaths(tt.ingress)
			if len(result) != len(tt.expected) {
				t.Errorf("GetIngressPaths() = %v, want %v", result, tt.expected)
				return
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("GetIngressPaths()[%d] = %v, want %v", i, v, tt.expected[i])
				}
			}
		})
	}
}

// TestGetIngressBackend tests the GetIngressBackend function
// 测试 GetIngressBackend 函数
// 调用方式: go test -v -run TestGetIngressBackend ./internal/kubernetes/
func TestGetIngressBackend(t *testing.T) {
	util := NewKubernetesUtil()

	tests := []struct {
		name            string
		ingress         *networkingv1.Ingress
		expectedService string
		expectedPort    int32
	}{
		{
			name:            "nil ingress",
			ingress:         nil,
			expectedService: "",
			expectedPort:    0,
		},
		{
			name: "default backend",
			ingress: &networkingv1.Ingress{
				Spec: networkingv1.IngressSpec{
					DefaultBackend: &networkingv1.IngressBackend{
						Service: &networkingv1.IngressServiceBackend{
							Name: "default-service",
							Port: networkingv1.ServiceBackendPort{Number: 80},
						},
					},
				},
			},
			expectedService: "default-service",
			expectedPort:    80,
		},
		{
			name: "rule backend",
			ingress: &networkingv1.Ingress{
				Spec: networkingv1.IngressSpec{
					Rules: []networkingv1.IngressRule{
						{
							IngressRuleValue: networkingv1.IngressRuleValue{
								HTTP: &networkingv1.HTTPIngressRuleValue{
									Paths: []networkingv1.HTTPIngressPath{
										{
											Backend: networkingv1.IngressBackend{
												Service: &networkingv1.IngressServiceBackend{
													Name: "rule-service",
													Port: networkingv1.ServiceBackendPort{Number: 8080},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedService: "rule-service",
			expectedPort:    8080,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, port := util.GetIngressBackend(tt.ingress)
			if service != tt.expectedService {
				t.Errorf("GetIngressBackend() service = %v, want %v", service, tt.expectedService)
			}
			if port != tt.expectedPort {
				t.Errorf("GetIngressBackend() port = %v, want %v", port, tt.expectedPort)
			}
		})
	}
}

// TestGetSecretType tests the GetSecretType function
// 测试 GetSecretType 函数
// 调用方式: go test -v -run TestGetSecretType ./internal/kubernetes/
func TestGetSecretType(t *testing.T) {
	util := NewKubernetesUtil()

	tests := []struct {
		name     string
		secret   *corev1.Secret
		expected string
	}{
		{
			name:     "nil secret",
			secret:   nil,
			expected: "",
		},
		{
			name: "TLS secret",
			secret: &corev1.Secret{
				Type: corev1.SecretTypeTLS,
			},
			expected: "kubernetes.io/tls",
		},
		{
			name: "Opaque secret",
			secret: &corev1.Secret{
				Type: corev1.SecretTypeOpaque,
			},
			expected: "Opaque",
		},
		{
			name: "Docker config secret",
			secret: &corev1.Secret{
				Type: corev1.SecretTypeDockerConfigJson,
			},
			expected: "kubernetes.io/dockerconfigjson",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.GetSecretType(tt.secret)
			if result != tt.expected {
				t.Errorf("GetSecretType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestIsTLSSecret tests the IsTLSSecret function
// 测试 IsTLSSecret 函数
// 调用方式: go test -v -run TestIsTLSSecret ./internal/kubernetes/
func TestIsTLSSecret(t *testing.T) {
	util := NewKubernetesUtil()

	tests := []struct {
		name     string
		secret   *corev1.Secret
		expected bool
	}{
		{
			name:     "nil secret",
			secret:   nil,
			expected: false,
		},
		{
			name: "TLS secret",
			secret: &corev1.Secret{
				Type: corev1.SecretTypeTLS,
			},
			expected: true,
		},
		{
			name: "Opaque secret",
			secret: &corev1.Secret{
				Type: corev1.SecretTypeOpaque,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.IsTLSSecret(tt.secret)
			if result != tt.expected {
				t.Errorf("IsTLSSecret() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestBuildIngressName tests the BuildIngressName function
// 测试 BuildIngressName 函数
// 调用方式: go test -v -run TestBuildIngressName ./internal/kubernetes/
func TestBuildIngressName(t *testing.T) {
	util := NewKubernetesUtil()

	tests := []struct {
		domain   string
		path     string
		expected string
	}{
		{
			domain:   "example.com",
			path:     "",
			expected: "example.com",
		},
		{
			domain:   "example.com",
			path:     "/",
			expected: "example.com",
		},
		{
			domain:   "example.com",
			path:     "/api",
			expected: "example.com-api",
		},
		{
			domain:   "example.com",
			path:     "/api/v1/users",
			expected: "example.com-api-v1-users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.domain+tt.path, func(t *testing.T) {
			result := util.BuildIngressName(tt.domain, tt.path)
			if result != tt.expected {
				t.Errorf("BuildIngressName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestBuildSecretName tests the BuildSecretName function
// 测试 BuildSecretName 函数
// 调用方式: go test -v -run TestBuildSecretName ./internal/kubernetes/
func TestBuildSecretName(t *testing.T) {
	util := NewKubernetesUtil()

	tests := []struct {
		domain   string
		expected string
	}{
		{
			domain:   "example.com",
			expected: "example-com",
		},
		{
			domain:   "api.example.com",
			expected: "api-example-com",
		},
		{
			domain:   "www.example.com",
			expected: "www-example-com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.domain, func(t *testing.T) {
			result := util.BuildSecretName(tt.domain)
			if result != tt.expected {
				t.Errorf("BuildSecretName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestMergeLabels tests the MergeLabels function
// 测试 MergeLabels 函数
// 调用方式: go test -v -run TestMergeLabels ./internal/kubernetes/
func TestMergeLabels(t *testing.T) {
	util := NewKubernetesUtil()

	tests := []struct {
		name      string
		existing  map[string]string
		newLabels map[string]string
		expected  map[string]string
	}{
		{
			name:      "nil existing",
			existing:  nil,
			newLabels: map[string]string{"app": "test"},
			expected:  map[string]string{"app": "test"},
		},
		{
			name:      "empty existing",
			existing:  map[string]string{},
			newLabels: map[string]string{"app": "test"},
			expected:  map[string]string{"app": "test"},
		},
		{
			name:      "merge with existing",
			existing:  map[string]string{"env": "prod"},
			newLabels: map[string]string{"app": "test"},
			expected:  map[string]string{"env": "prod", "app": "test"},
		},
		{
			name:      "overwrite existing",
			existing:  map[string]string{"app": "old"},
			newLabels: map[string]string{"app": "new"},
			expected:  map[string]string{"app": "new"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.MergeLabels(tt.existing, tt.newLabels)
			if len(result) != len(tt.expected) {
				t.Errorf("MergeLabels() = %v, want %v", result, tt.expected)
				return
			}
			for k, v := range tt.expected {
				if result[k] != v {
					t.Errorf("MergeLabels()[%s] = %v, want %v", k, result[k], v)
				}
			}
		})
	}
}

// TestCreateObjectMeta tests the CreateObjectMeta function
// 测试 CreateObjectMeta 函数
// 调用方式: go test -v -run TestCreateObjectMeta ./internal/kubernetes/
func TestCreateObjectMeta(t *testing.T) {
	util := NewKubernetesUtil()

	meta := util.CreateObjectMeta("test-name", "test-namespace")

	if meta.Name != "test-name" {
		t.Errorf("Expected name 'test-name', got '%s'", meta.Name)
	}
	if meta.Namespace != "test-namespace" {
		t.Errorf("Expected namespace 'test-namespace', got '%s'", meta.Namespace)
	}
}

// TestCreateObjectMetaWithLabels tests the CreateObjectMetaWithLabels function
// 测试 CreateObjectMetaWithLabels 函数
// 调用方式: go test -v -run TestCreateObjectMetaWithLabels ./internal/kubernetes/
func TestCreateObjectMetaWithLabels(t *testing.T) {
	util := NewKubernetesUtil()

	labels := map[string]string{"app": "test"}
	meta := util.CreateObjectMetaWithLabels("test-name", "test-namespace", labels)

	if meta.Name != "test-name" {
		t.Errorf("Expected name 'test-name', got '%s'", meta.Name)
	}
	if meta.Namespace != "test-namespace" {
		t.Errorf("Expected namespace 'test-namespace', got '%s'", meta.Namespace)
	}
	if meta.Labels["app"] != "test" {
		t.Errorf("Expected label app='test', got '%s'", meta.Labels["app"])
	}
}

// TestCreateTLSSecret tests the CreateTLSSecret function
// 测试 CreateTLSSecret 函数
// 调用方式: go test -v -run TestCreateTLSSecret ./internal/kubernetes/
func TestCreateTLSSecret(t *testing.T) {
	util := NewKubernetesUtil()

	secret := util.CreateTLSSecret("test-secret", "test-namespace", "test-cert", "test-key")

	if secret.Name != "test-secret" {
		t.Errorf("Expected name 'test-secret', got '%s'", secret.Name)
	}
	if secret.Namespace != "test-namespace" {
		t.Errorf("Expected namespace 'test-namespace', got '%s'", secret.Namespace)
	}
	if secret.Type != corev1.SecretTypeTLS {
		t.Errorf("Expected type 'kubernetes.io/tls', got '%s'", secret.Type)
	}
	if string(secret.Data[corev1.TLSCertKey]) != "test-cert" {
		t.Errorf("Expected cert 'test-cert', got '%s'", secret.Data[corev1.TLSCertKey])
	}
	if string(secret.Data[corev1.TLSPrivateKeyKey]) != "test-key" {
		t.Errorf("Expected key 'test-key', got '%s'", secret.Data[corev1.TLSPrivateKeyKey])
	}
}

// TestCreateOpaqueSecret tests the CreateOpaqueSecret function
// 测试 CreateOpaqueSecret 函数
// 调用方式: go test -v -run TestCreateOpaqueSecret ./internal/kubernetes/
func TestCreateOpaqueSecret(t *testing.T) {
	util := NewKubernetesUtil()

	data := map[string]string{"username": "admin", "password": "secret"}
	secret := util.CreateOpaqueSecret("test-secret", "test-namespace", data)

	if secret.Name != "test-secret" {
		t.Errorf("Expected name 'test-secret', got '%s'", secret.Name)
	}
	if secret.Namespace != "test-namespace" {
		t.Errorf("Expected namespace 'test-namespace', got '%s'", secret.Namespace)
	}
	if secret.Type != corev1.SecretTypeOpaque {
		t.Errorf("Expected type 'Opaque', got '%s'", secret.Type)
	}
	if string(secret.Data["username"]) != "admin" {
		t.Errorf("Expected username 'admin', got '%s'", secret.Data["username"])
	}
	if string(secret.Data["password"]) != "secret" {
		t.Errorf("Expected password 'secret', got '%s'", secret.Data["password"])
	}
}

// TestCreateConfigMap tests the CreateConfigMap function
// 测试 CreateConfigMap 函数
// 调用方式: go test -v -run TestCreateConfigMap ./internal/kubernetes/
func TestCreateConfigMap(t *testing.T) {
	util := NewKubernetesUtil()

	data := map[string]string{"key1": "value1", "key2": "value2"}
	cm := util.CreateConfigMap("test-cm", "test-namespace", data)

	if cm.Name != "test-cm" {
		t.Errorf("Expected name 'test-cm', got '%s'", cm.Name)
	}
	if cm.Namespace != "test-namespace" {
		t.Errorf("Expected namespace 'test-namespace', got '%s'", cm.Namespace)
	}
	if cm.Data["key1"] != "value1" {
		t.Errorf("Expected key1='value1', got '%s'", cm.Data["key1"])
	}
	if cm.Data["key2"] != "value2" {
		t.Errorf("Expected key2='value2', got '%s'", cm.Data["key2"])
	}
}

// TestCreateIngressBackend tests the CreateIngressBackend function
// 测试 CreateIngressBackend 函数
// 调用方式: go test -v -run TestCreateIngressBackend ./internal/kubernetes/
func TestCreateIngressBackend(t *testing.T) {
	util := NewKubernetesUtil()

	backend := util.CreateIngressBackend("test-service", 8080)

	if backend.Service.Name != "test-service" {
		t.Errorf("Expected service name 'test-service', got '%s'", backend.Service.Name)
	}
	if backend.Service.Port.Number != 8080 {
		t.Errorf("Expected port 8080, got %d", backend.Service.Port.Number)
	}
}

// TestCreateHTTPIngressPath tests the CreateHTTPIngressPath function
// 测试 CreateHTTPIngressPath 函数
// 调用方式: go test -v -run TestCreateHTTPIngressPath ./internal/kubernetes/
func TestCreateHTTPIngressPath(t *testing.T) {
	util := NewKubernetesUtil()

	backend := util.CreateIngressBackend("test-service", 80)

	tests := []struct {
		name         string
		path         string
		pathType     string
		expectedType networkingv1.PathType
	}{
		{
			name:         "prefix path type",
			path:         "/api",
			pathType:     "Prefix",
			expectedType: networkingv1.PathTypePrefix,
		},
		{
			name:         "exact path type",
			path:         "/health",
			pathType:     "Exact",
			expectedType: networkingv1.PathTypeExact,
		},
		{
			name:         "implementation specific path type",
			path:         "/",
			pathType:     "",
			expectedType: networkingv1.PathTypeImplementationSpecific,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.CreateHTTPIngressPath(tt.path, tt.pathType, backend)

			if result.Path != tt.path {
				t.Errorf("Expected path '%s', got '%s'", tt.path, result.Path)
			}
			if *result.PathType != tt.expectedType {
				t.Errorf("Expected path type '%s', got '%s'", tt.expectedType, *result.PathType)
			}
		})
	}
}

// TestCreateIngress tests the CreateIngress function
// 测试 CreateIngress 函数
// 调用方式: go test -v -run TestCreateIngress ./internal/kubernetes/
func TestCreateIngress(t *testing.T) {
	util := NewKubernetesUtil()

	backend := util.CreateIngressBackend("test-service", 80)
	path := util.CreateHTTPIngressPath("/", "Prefix", backend)
	httpRule := util.CreateHTTPIngressRuleValue([]networkingv1.HTTPIngressPath{path})
	rule := util.CreateIngressRule("example.com", httpRule)
	spec := util.CreateIngressSpec([]networkingv1.IngressRule{rule}, nil, nil)
	ingress := util.CreateIngress("test-ingress", "test-namespace", spec)

	if ingress.Name != "test-ingress" {
		t.Errorf("Expected name 'test-ingress', got '%s'", ingress.Name)
	}
	if ingress.Namespace != "test-namespace" {
		t.Errorf("Expected namespace 'test-namespace', got '%s'", ingress.Namespace)
	}
	if len(ingress.Spec.Rules) != 1 {
		t.Errorf("Expected 1 rule, got %d", len(ingress.Spec.Rules))
	}
	if ingress.Spec.Rules[0].Host != "example.com" {
		t.Errorf("Expected host 'example.com', got '%s'", ingress.Spec.Rules[0].Host)
	}
}

// TestIsIngressReady tests the IsIngressReady function
// 测试 IsIngressReady 函数
// 调用方式: go test -v -run TestIsIngressReady ./internal/kubernetes/
func TestIsIngressReady(t *testing.T) {
	util := NewKubernetesUtil()

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
			name: "no load balancer",
			ingress: &networkingv1.Ingress{
				Status: networkingv1.IngressStatus{},
			},
			expected: false,
		},
		{
			name: "with hostname",
			ingress: &networkingv1.Ingress{
				Status: networkingv1.IngressStatus{
					LoadBalancer: networkingv1.IngressLoadBalancerStatus{
						Ingress: []networkingv1.IngressLoadBalancerIngress{
							{Hostname: "example.com"},
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "with IP",
			ingress: &networkingv1.Ingress{
				Status: networkingv1.IngressStatus{
					LoadBalancer: networkingv1.IngressLoadBalancerStatus{
						Ingress: []networkingv1.IngressLoadBalancerIngress{
							{IP: "192.168.1.1"},
						},
					},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.IsIngressReady(tt.ingress)
			if result != tt.expected {
				t.Errorf("IsIngressReady() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestGetIngressLoadBalancer tests the GetIngressLoadBalancer function
// 测试 GetIngressLoadBalancer 函数
// 调用方式: go test -v -run TestGetIngressLoadBalancer ./internal/kubernetes/
func TestGetIngressLoadBalancer(t *testing.T) {
	util := NewKubernetesUtil()

	tests := []struct {
		name     string
		ingress  *networkingv1.Ingress
		expected string
	}{
		{
			name:     "nil ingress",
			ingress:  nil,
			expected: "",
		},
		{
			name: "no load balancer",
			ingress: &networkingv1.Ingress{
				Status: networkingv1.IngressStatus{},
			},
			expected: "",
		},
		{
			name: "with hostname",
			ingress: &networkingv1.Ingress{
				Status: networkingv1.IngressStatus{
					LoadBalancer: networkingv1.IngressLoadBalancerStatus{
						Ingress: []networkingv1.IngressLoadBalancerIngress{
							{Hostname: "example.com"},
						},
					},
				},
			},
			expected: "example.com",
		},
		{
			name: "with IP",
			ingress: &networkingv1.Ingress{
				Status: networkingv1.IngressStatus{
					LoadBalancer: networkingv1.IngressLoadBalancerStatus{
						Ingress: []networkingv1.IngressLoadBalancerIngress{
							{IP: "192.168.1.1"},
						},
					},
				},
			},
			expected: "192.168.1.1",
		},
		{
			name: "hostname takes precedence",
			ingress: &networkingv1.Ingress{
				Status: networkingv1.IngressStatus{
					LoadBalancer: networkingv1.IngressLoadBalancerStatus{
						Ingress: []networkingv1.IngressLoadBalancerIngress{
							{Hostname: "example.com", IP: "192.168.1.1"},
						},
					},
				},
			},
			expected: "example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.GetIngressLoadBalancer(tt.ingress)
			if result != tt.expected {
				t.Errorf("GetIngressLoadBalancer() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestGetTLSCertificate tests the GetTLSCertificate function
// 测试 GetTLSCertificate 函数
// 调用方式: go test -v -run TestGetTLSCertificate ./internal/kubernetes/
func TestGetTLSCertificate(t *testing.T) {
	util := NewKubernetesUtil()

	tests := []struct {
		name     string
		secret   *corev1.Secret
		expected string
	}{
		{
			name:     "nil secret",
			secret:   nil,
			expected: "",
		},
		{
			name: "no data",
			secret: &corev1.Secret{
				Data: nil,
			},
			expected: "",
		},
		{
			name: "with certificate",
			secret: &corev1.Secret{
				Data: map[string][]byte{
					corev1.TLSCertKey: []byte("test-cert"),
				},
			},
			expected: "test-cert",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.GetTLSCertificate(tt.secret)
			if result != tt.expected {
				t.Errorf("GetTLSCertificate() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestGetTLSPrivateKey tests the GetTLSPrivateKey function
// 测试 GetTLSPrivateKey 函数
// 调用方式: go test -v -run TestGetTLSPrivateKey ./internal/kubernetes/
func TestGetTLSPrivateKey(t *testing.T) {
	util := NewKubernetesUtil()

	tests := []struct {
		name     string
		secret   *corev1.Secret
		expected string
	}{
		{
			name:     "nil secret",
			secret:   nil,
			expected: "",
		},
		{
			name: "no data",
			secret: &corev1.Secret{
				Data: nil,
			},
			expected: "",
		},
		{
			name: "with private key",
			secret: &corev1.Secret{
				Data: map[string][]byte{
					corev1.TLSPrivateKeyKey: []byte("test-key"),
				},
			},
			expected: "test-key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.GetTLSPrivateKey(tt.secret)
			if result != tt.expected {
				t.Errorf("GetTLSPrivateKey() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestParseIngressClass tests the ParseIngressClass function
// 测试 ParseIngressClass 函数
// 调用方式: go test -v -run TestParseIngressClass ./internal/kubernetes/
func TestParseIngressClass(t *testing.T) {
	util := NewKubernetesUtil()

	ingressClassName := "nginx"

	tests := []struct {
		name     string
		ingress  *networkingv1.Ingress
		expected string
	}{
		{
			name:     "nil ingress",
			ingress:  nil,
			expected: "",
		},
		{
			name: "no ingress class",
			ingress: &networkingv1.Ingress{
				Spec: networkingv1.IngressSpec{},
			},
			expected: "",
		},
		{
			name: "with ingress class",
			ingress: &networkingv1.Ingress{
				Spec: networkingv1.IngressSpec{
					IngressClassName: &ingressClassName,
				},
			},
			expected: "nginx",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.ParseIngressClass(tt.ingress)
			if result != tt.expected {
				t.Errorf("ParseIngressClass() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestGetSecretData tests the GetSecretData function
// 测试 GetSecretData 函数
// 调用方式: go test -v -run TestGetSecretData ./internal/kubernetes/
func TestGetSecretData(t *testing.T) {
	util := NewKubernetesUtil()

	tests := []struct {
		name     string
		secret   *corev1.Secret
		expected map[string]string
	}{
		{
			name:     "nil secret",
			secret:   nil,
			expected: nil,
		},
		{
			name: "no data",
			secret: &corev1.Secret{
				Data: nil,
			},
			expected: nil,
		},
		{
			name: "with data",
			secret: &corev1.Secret{
				Data: map[string][]byte{
					"key1": []byte("value1"),
					"key2": []byte("value2"),
				},
			},
			expected: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.GetSecretData(tt.secret)
			if len(result) != len(tt.expected) {
				t.Errorf("GetSecretData() = %v, want %v", result, tt.expected)
				return
			}
			for k, v := range tt.expected {
				if result[k] != v {
					t.Errorf("GetSecretData()[%s] = %v, want %v", k, result[k], v)
				}
			}
		})
	}
}

// TestGetConfigMapData tests the GetConfigMapData function
// 测试 GetConfigMapData 函数
// 调用方式: go test -v -run TestGetConfigMapData ./internal/kubernetes/
func TestGetConfigMapData(t *testing.T) {
	util := NewKubernetesUtil()

	tests := []struct {
		name     string
		cm       *corev1.ConfigMap
		expected map[string]string
	}{
		{
			name:     "nil configmap",
			cm:       nil,
			expected: nil,
		},
		{
			name: "no data",
			cm: &corev1.ConfigMap{
				Data: nil,
			},
			expected: nil,
		},
		{
			name: "with data",
			cm: &corev1.ConfigMap{
				Data: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
			expected: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.GetConfigMapData(tt.cm)
			if len(result) != len(tt.expected) {
				t.Errorf("GetConfigMapData() = %v, want %v", result, tt.expected)
				return
			}
			for k, v := range tt.expected {
				if result[k] != v {
					t.Errorf("GetConfigMapData()[%s] = %v, want %v", k, result[k], v)
				}
			}
		})
	}
}

// TestCreateObjectMetaFull tests the CreateObjectMetaFull function
// 测试 CreateObjectMetaFull 函数
// 调用方式: go test -v -run TestCreateObjectMetaFull ./internal/kubernetes/
func TestCreateObjectMetaFull(t *testing.T) {
	util := NewKubernetesUtil()

	labels := map[string]string{"app": "test"}
	annotations := map[string]string{"description": "test annotation"}

	meta := util.CreateObjectMetaFull("test-name", "test-namespace", labels, annotations)

	if meta.Name != "test-name" {
		t.Errorf("Expected name 'test-name', got '%s'", meta.Name)
	}
	if meta.Namespace != "test-namespace" {
		t.Errorf("Expected namespace 'test-namespace', got '%s'", meta.Namespace)
	}
	if meta.Labels["app"] != "test" {
		t.Errorf("Expected label app='test', got '%s'", meta.Labels["app"])
	}
	if meta.Annotations["description"] != "test annotation" {
		t.Errorf("Expected annotation description='test annotation', got '%s'", meta.Annotations["description"])
	}
}

// TestMergeAnnotations tests the MergeAnnotations function
// 测试 MergeAnnotations 函数
// 调用方式: go test -v -run TestMergeAnnotations ./internal/kubernetes/
func TestMergeAnnotations(t *testing.T) {
	util := NewKubernetesUtil()

	tests := []struct {
		name           string
		existing       map[string]string
		newAnnotations map[string]string
		expected       map[string]string
	}{
		{
			name:           "nil existing",
			existing:       nil,
			newAnnotations: map[string]string{"key": "value"},
			expected:       map[string]string{"key": "value"},
		},
		{
			name:           "merge with existing",
			existing:       map[string]string{"existing": "value"},
			newAnnotations: map[string]string{"new": "value"},
			expected:       map[string]string{"existing": "value", "new": "value"},
		},
		{
			name:           "overwrite existing",
			existing:       map[string]string{"key": "old"},
			newAnnotations: map[string]string{"key": "new"},
			expected:       map[string]string{"key": "new"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.MergeAnnotations(tt.existing, tt.newAnnotations)
			if len(result) != len(tt.expected) {
				t.Errorf("MergeAnnotations() = %v, want %v", result, tt.expected)
				return
			}
			for k, v := range tt.expected {
				if result[k] != v {
					t.Errorf("MergeAnnotations()[%s] = %v, want %v", k, result[k], v)
				}
			}
		})
	}
}
