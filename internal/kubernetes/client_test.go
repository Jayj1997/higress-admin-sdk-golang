// Package kubernetes provides Kubernetes client functionality
package kubernetes

import (
	"context"
	"testing"

	"github.com/Jayj1997/higress-admin-sdk-golang/internal/kubernetes/crd/istio"
	"github.com/Jayj1997/higress-admin-sdk-golang/internal/kubernetes/crd/mcp"
	"github.com/Jayj1997/higress-admin-sdk-golang/internal/kubernetes/crd/wasm"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/config"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/constant"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestNewKubernetesClientService tests the creation of KubernetesClientService
// 测试 KubernetesClientService 的创建
// 调用方式: go test -v -run TestNewKubernetesClientService ./internal/kubernetes/
func TestNewKubernetesClientService(t *testing.T) {
	tests := []struct {
		name    string
		config  *config.HigressServiceConfig
		wantErr bool
	}{
		{
			name: "valid config with kubeconfig path",
			config: &config.HigressServiceConfig{
				ControllerServiceHost:      "localhost",
				ControllerServicePort:      8080,
				ControllerNamespace:        "higress-system",
				ControllerJwtPolicy:        constant.JwtPolicyFirstPartyJwt,
				ControllerWatchedNamespace: "",
			},
			wantErr: false,
		},
		{
			name: "missing controller namespace",
			config: &config.HigressServiceConfig{
				ControllerServiceHost: "localhost",
				ControllerServicePort: 8080,
				ControllerJwtPolicy:   constant.JwtPolicyFirstPartyJwt,
			},
			wantErr: true,
		},
		{
			name: "missing controller service host",
			config: &config.HigressServiceConfig{
				ControllerServicePort: 8080,
				ControllerNamespace:   "higress-system",
				ControllerJwtPolicy:   constant.JwtPolicyFirstPartyJwt,
			},
			wantErr: true,
		},
		{
			name: "invalid controller service port",
			config: &config.HigressServiceConfig{
				ControllerServiceHost: "localhost",
				ControllerServicePort: 0,
				ControllerNamespace:   "higress-system",
				ControllerJwtPolicy:   constant.JwtPolicyFirstPartyJwt,
			},
			wantErr: true,
		},
		{
			name: "missing jwt policy",
			config: &config.HigressServiceConfig{
				ControllerServiceHost: "localhost",
				ControllerServicePort: 8080,
				ControllerNamespace:   "higress-system",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewKubernetesClientService(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewKubernetesClientService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestIsNamespaceProtected tests the IsNamespaceProtected function
// 测试 IsNamespaceProtected 函数
// 调用方式: go test -v -run TestIsNamespaceProtected ./internal/kubernetes/
func TestIsNamespaceProtected(t *testing.T) {
	cfg := &config.HigressServiceConfig{
		ControllerServiceHost: "localhost",
		ControllerServicePort: 8080,
		ControllerNamespace:   "higress-system",
		ControllerJwtPolicy:   constant.JwtPolicyFirstPartyJwt,
	}

	svc, err := NewKubernetesClientService(cfg)
	if err != nil {
		t.Fatalf("Failed to create KubernetesClientService: %v", err)
	}

	tests := []struct {
		name      string
		namespace string
		want      bool
	}{
		{
			name:      "kube-system namespace is protected",
			namespace: constant.KubeSystemNamespace,
			want:      true,
		},
		{
			name:      "controller namespace is protected",
			namespace: "higress-system",
			want:      true,
		},
		{
			name:      "other namespace is not protected",
			namespace: "default",
			want:      false,
		},
		{
			name:      "empty namespace is not protected",
			namespace: "",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := svc.IsNamespaceProtected(tt.namespace); got != tt.want {
				t.Errorf("IsNamespaceProtected() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestRenderDefaultMetadata tests the renderDefaultMetadata function
// 测试 renderDefaultMetadata 函数
// 调用方式: go test -v -run TestRenderDefaultMetadata ./internal/kubernetes/
func TestRenderDefaultMetadata(t *testing.T) {
	cfg := &config.HigressServiceConfig{
		ControllerServiceHost: "localhost",
		ControllerServicePort: 8080,
		ControllerNamespace:   "higress-system",
		ControllerJwtPolicy:   constant.JwtPolicyFirstPartyJwt,
	}

	svc, err := NewKubernetesClientService(cfg)
	if err != nil {
		t.Fatalf("Failed to create KubernetesClientService: %v", err)
	}

	t.Run("ingress without labels", func(t *testing.T) {
		ingress := &networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-ingress",
			},
		}
		svc.renderDefaultMetadata(ingress)

		labels := ingress.GetLabels()
		if labels == nil {
			t.Error("Expected labels to be set")
			return
		}
		if labels[constant.LabelResourceDefinerKey] != constant.LabelResourceDefinerValue {
			t.Errorf("Expected label %s=%s, got %s", constant.LabelResourceDefinerKey, constant.LabelResourceDefinerValue, labels[constant.LabelResourceDefinerKey])
		}
	})

	t.Run("ingress with existing labels", func(t *testing.T) {
		ingress := &networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-ingress",
				Labels: map[string]string{
					"app": "test",
				},
			},
		}
		svc.renderDefaultMetadata(ingress)

		labels := ingress.GetLabels()
		if labels == nil {
			t.Error("Expected labels to be set")
			return
		}
		if labels["app"] != "test" {
			t.Error("Expected existing label to be preserved")
		}
		if labels[constant.LabelResourceDefinerKey] != constant.LabelResourceDefinerValue {
			t.Errorf("Expected label %s=%s", constant.LabelResourceDefinerKey, constant.LabelResourceDefinerValue)
		}
	})
}

// TestBuildLabelSelector tests the buildLabelSelector function
// 测试 buildLabelSelector 函数
// 调用方式: go test -v -run TestBuildLabelSelector ./internal/kubernetes/
func TestBuildLabelSelector(t *testing.T) {
	tests := []struct {
		key   string
		value string
		want  string
	}{
		{
			key:   "app",
			value: "nginx",
			want:  "app=nginx",
		},
		{
			key:   constant.LabelResourceDefinerKey,
			value: constant.LabelResourceDefinerValue,
			want:  "higress.io/resource-definer=higress-console",
		},
		{
			key:   "environment",
			value: "production",
			want:  "environment=production",
		},
	}

	for _, tt := range tests {
		t.Run(tt.key+"="+tt.value, func(t *testing.T) {
			if got := buildLabelSelector(tt.key, tt.value); got != tt.want {
				t.Errorf("buildLabelSelector() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsInCluster tests the isInCluster function
// 测试 isInCluster 函数
// 调用方式: go test -v -run TestIsInCluster ./internal/kubernetes/
func TestIsInCluster(t *testing.T) {
	// This test checks if the function runs without error
	// The result depends on the environment
	result := isInCluster()
	t.Logf("isInCluster() = %v", result)
	// We don't assert a specific value since it depends on the environment
}

// TestWasmPluginConversion tests the WasmPlugin conversion functions
// 测试 WasmPlugin 转换函数
// 调用方式: go test -v -run TestWasmPluginConversion ./internal/kubernetes/
func TestWasmPluginConversion(t *testing.T) {
	plugin := wasm.NewV1alpha1WasmPlugin()
	plugin.Metadata.Name = "test-plugin"
	plugin.Metadata.Namespace = "higress-system"
	plugin.Spec.Url = "oci://test-image:v1"
	plugin.Spec.Phase = wasm.PluginPhaseAuthn
	plugin.Spec.Priority = ptrInt64(100)

	// Test conversion to unstructured
	unstructuredObj := wasmPluginToUnstructured(plugin)
	if unstructuredObj == nil {
		t.Fatal("Expected unstructured object to be created")
	}

	// Test conversion back to WasmPlugin
	convertedPlugin := unstructuredToWasmPlugin(unstructuredObj)
	if convertedPlugin == nil {
		t.Fatal("Expected WasmPlugin to be created")
	}

	if convertedPlugin.Metadata.Name != plugin.Metadata.Name {
		t.Errorf("Expected name %s, got %s", plugin.Metadata.Name, convertedPlugin.Metadata.Name)
	}

	if convertedPlugin.Spec.Url != plugin.Spec.Url {
		t.Errorf("Expected URL %s, got %s", plugin.Spec.Url, convertedPlugin.Spec.Url)
	}
}

// TestMcpBridgeConversion tests the McpBridge conversion functions
// 测试 McpBridge 转换函数
// 调用方式: go test -v -run TestMcpBridgeConversion ./internal/kubernetes/
func TestMcpBridgeConversion(t *testing.T) {
	bridge := mcp.NewV1McpBridge()
	bridge.Metadata.Name = "test-bridge"
	bridge.Metadata.Namespace = "higress-system"
	bridge.Spec.Registries = []*mcp.V1RegistryConfig{
		{
			Name:           "test-registry",
			Type:           "nacos",
			Domain:         "nacos.default.svc.cluster.local",
			Port:           8848,
			NacosNamespace: "public",
		},
	}

	// Test conversion to unstructured
	unstructuredObj := mcpBridgeToUnstructured(bridge)
	if unstructuredObj == nil {
		t.Fatal("Expected unstructured object to be created")
	}

	// Test conversion back to McpBridge
	convertedBridge := unstructuredToMcpBridge(unstructuredObj)
	if convertedBridge == nil {
		t.Fatal("Expected McpBridge to be created")
	}

	if convertedBridge.Metadata.Name != bridge.Metadata.Name {
		t.Errorf("Expected name %s, got %s", bridge.Metadata.Name, convertedBridge.Metadata.Name)
	}

	if len(convertedBridge.Spec.Registries) != 1 {
		t.Errorf("Expected 1 registry, got %d", len(convertedBridge.Spec.Registries))
	}
}

// TestEnvoyFilterConversion tests the EnvoyFilter conversion functions
// 测试 EnvoyFilter 转换函数
// 调用方式: go test -v -run TestEnvoyFilterConversion ./internal/kubernetes/
func TestEnvoyFilterConversion(t *testing.T) {
	filter := istio.NewV1alpha3EnvoyFilter()
	filter.Metadata.Name = "test-filter"
	filter.Metadata.Namespace = "higress-system"
	filter.Spec.ConfigPatches = []*istio.V1alpha3EnvoyConfigObjectPatch{
		{
			ApplyTo: istio.ApplyToHTTPFilter,
			Patch: &istio.V1alpha3Patch{
				Operation: "INSERT_BEFORE",
				Value:     map[string]interface{}{"name": "test-filter"},
			},
		},
	}

	// Test conversion to unstructured
	unstructuredObj := envoyFilterToUnstructured(filter)
	if unstructuredObj == nil {
		t.Fatal("Expected unstructured object to be created")
	}

	// Test conversion back to EnvoyFilter
	convertedFilter := unstructuredToEnvoyFilter(unstructuredObj)
	if convertedFilter == nil {
		t.Fatal("Expected EnvoyFilter to be created")
	}

	if convertedFilter.Metadata.Name != filter.Metadata.Name {
		t.Errorf("Expected name %s, got %s", filter.Metadata.Name, convertedFilter.Metadata.Name)
	}
}

// TestValidateConfig tests the validateConfig function
// 测试 validateConfig 函数
// 调用方式: go test -v -run TestValidateConfig ./internal/kubernetes/
func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *config.HigressServiceConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: &config.HigressServiceConfig{
				ControllerServiceHost: "localhost",
				ControllerServicePort: 8080,
				ControllerNamespace:   "higress-system",
				ControllerJwtPolicy:   constant.JwtPolicyFirstPartyJwt,
			},
			wantErr: false,
		},
		{
			name: "missing controller service host",
			config: &config.HigressServiceConfig{
				ControllerServicePort: 8080,
				ControllerNamespace:   "higress-system",
				ControllerJwtPolicy:   constant.JwtPolicyFirstPartyJwt,
			},
			wantErr: true,
		},
		{
			name: "missing controller namespace",
			config: &config.HigressServiceConfig{
				ControllerServiceHost: "localhost",
				ControllerServicePort: 8080,
				ControllerJwtPolicy:   constant.JwtPolicyFirstPartyJwt,
			},
			wantErr: true,
		},
		{
			name: "invalid port - zero",
			config: &config.HigressServiceConfig{
				ControllerServiceHost: "localhost",
				ControllerServicePort: 0,
				ControllerNamespace:   "higress-system",
				ControllerJwtPolicy:   constant.JwtPolicyFirstPartyJwt,
			},
			wantErr: true,
		},
		{
			name: "invalid port - too large",
			config: &config.HigressServiceConfig{
				ControllerServiceHost: "localhost",
				ControllerServicePort: 70000,
				ControllerNamespace:   "higress-system",
				ControllerJwtPolicy:   constant.JwtPolicyFirstPartyJwt,
			},
			wantErr: true,
		},
		{
			name: "missing jwt policy",
			config: &config.HigressServiceConfig{
				ControllerServiceHost: "localhost",
				ControllerServicePort: 8080,
				ControllerNamespace:   "higress-system",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestIngressCRUD tests Ingress CRUD operations (requires a running Kubernetes cluster)
// 测试 Ingress CRUD 操作（需要运行中的 Kubernetes 集群）
// 调用方式: go test -v -run TestIngressCRUD ./internal/kubernetes/
// 注意: 此测试需要真实的 Kubernetes 集群环境
func TestIngressCRUD(t *testing.T) {
	cfg := &config.HigressServiceConfig{
		ControllerServiceHost: "localhost",
		ControllerServicePort: 8080,
		ControllerNamespace:   "higress-system",
		ControllerJwtPolicy:   constant.JwtPolicyFirstPartyJwt,
	}

	svc, err := NewKubernetesClientService(cfg)
	if err != nil {
		t.Fatalf("Failed to create KubernetesClientService: %v", err)
	}

	ctx := context.Background()

	// Create test ingress
	pathType := networkingv1.PathTypePrefix
	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "test-ingress-sdk",
			Labels: map[string]string{},
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: "test.example.com",
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: "test-service",
											Port: networkingv1.ServiceBackendPort{
												Number: 80,
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
	}

	// Test Create
	created, err := svc.CreateIngress(ctx, ingress)
	if err != nil {
		t.Logf("CreateIngress failed (expected if no cluster): %v", err)
		return
	}
	t.Logf("Created ingress: %s", created.Name)

	// Test Get
	got, err := svc.GetIngress(ctx, "test-ingress-sdk")
	if err != nil {
		t.Errorf("GetIngress failed: %v", err)
	}
	if got == nil {
		t.Error("Expected ingress to be found")
	}

	// Test List
	list, err := svc.ListIngresses(ctx)
	if err != nil {
		t.Errorf("ListIngresses failed: %v", err)
	}
	t.Logf("Found %d ingresses", len(list))

	// Test Update
	if got != nil {
		got.Labels["updated"] = "true"
		updated, err := svc.UpdateIngress(ctx, got)
		if err != nil {
			t.Errorf("UpdateIngress failed: %v", err)
		}
		t.Logf("Updated ingress: %s", updated.Name)
	}

	// Test Delete
	err = svc.DeleteIngress(ctx, "test-ingress-sdk")
	if err != nil {
		t.Errorf("DeleteIngress failed: %v", err)
	}
	t.Log("Deleted ingress")
}

// TestSecretCRUD tests Secret CRUD operations (requires a running Kubernetes cluster)
// 测试 Secret CRUD 操作（需要运行中的 Kubernetes 集群）
// 调用方式: go test -v -run TestSecretCRUD ./internal/kubernetes/
// 注意: 此测试需要真实的 Kubernetes 集群环境
func TestSecretCRUD(t *testing.T) {
	cfg := &config.HigressServiceConfig{
		ControllerServiceHost: "localhost",
		ControllerServicePort: 8080,
		ControllerNamespace:   "higress-system",
		ControllerJwtPolicy:   constant.JwtPolicyFirstPartyJwt,
	}

	svc, err := NewKubernetesClientService(cfg)
	if err != nil {
		t.Fatalf("Failed to create KubernetesClientService: %v", err)
	}

	ctx := context.Background()

	// Create test secret
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "test-secret-sdk",
			Labels: map[string]string{},
		},
		Type: corev1.SecretTypeTLS,
		Data: map[string][]byte{
			corev1.TLSCertKey:       []byte("test-cert"),
			corev1.TLSPrivateKeyKey: []byte("test-key"),
		},
	}

	// Test Create
	created, err := svc.CreateSecret(ctx, secret)
	if err != nil {
		t.Logf("CreateSecret failed (expected if no cluster): %v", err)
		return
	}
	t.Logf("Created secret: %s", created.Name)

	// Test Get
	got, err := svc.GetSecret(ctx, "test-secret-sdk")
	if err != nil {
		t.Errorf("GetSecret failed: %v", err)
	}
	if got == nil {
		t.Error("Expected secret to be found")
	}

	// Test List
	list, err := svc.ListSecrets(ctx, string(corev1.SecretTypeTLS))
	if err != nil {
		t.Errorf("ListSecrets failed: %v", err)
	}
	t.Logf("Found %d secrets", len(list))

	// Test Delete
	err = svc.DeleteSecret(ctx, "test-secret-sdk")
	if err != nil {
		t.Errorf("DeleteSecret failed: %v", err)
	}
	t.Log("Deleted secret")
}

// TestConfigMapCRUD tests ConfigMap CRUD operations (requires a running Kubernetes cluster)
// 测试 ConfigMap CRUD 操作（需要运行中的 Kubernetes 集群）
// 调用方式: go test -v -run TestConfigMapCRUD ./internal/kubernetes/
// 注意: 此测试需要真实的 Kubernetes 集群环境
func TestConfigMapCRUD(t *testing.T) {
	cfg := &config.HigressServiceConfig{
		ControllerServiceHost: "localhost",
		ControllerServicePort: 8080,
		ControllerNamespace:   "higress-system",
		ControllerJwtPolicy:   constant.JwtPolicyFirstPartyJwt,
	}

	svc, err := NewKubernetesClientService(cfg)
	if err != nil {
		t.Fatalf("Failed to create KubernetesClientService: %v", err)
	}

	ctx := context.Background()

	// Create test configmap
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "test-configmap-sdk",
			Labels: map[string]string{},
		},
		Data: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	// Test Create
	created, err := svc.CreateConfigMap(ctx, cm)
	if err != nil {
		t.Logf("CreateConfigMap failed (expected if no cluster): %v", err)
		return
	}
	t.Logf("Created configmap: %s", created.Name)

	// Test Get
	got, err := svc.GetConfigMap(ctx, "test-configmap-sdk")
	if err != nil {
		t.Errorf("GetConfigMap failed: %v", err)
	}
	if got == nil {
		t.Error("Expected configmap to be found")
	}

	// Test List
	list, err := svc.ListConfigMaps(ctx, nil)
	if err != nil {
		t.Errorf("ListConfigMaps failed: %v", err)
	}
	t.Logf("Found %d configmaps", len(list))

	// Test Delete
	err = svc.DeleteConfigMap(ctx, "test-configmap-sdk")
	if err != nil {
		t.Errorf("DeleteConfigMap failed: %v", err)
	}
	t.Log("Deleted configmap")
}

// TestWasmPluginCRUD tests WasmPlugin CRUD operations (requires a running Kubernetes cluster)
// 测试 WasmPlugin CRUD 操作（需要运行中的 Kubernetes 集群）
// 调用方式: go test -v -run TestWasmPluginCRUD ./internal/kubernetes/
// 注意: 此测试需要真实的 Kubernetes 集群环境
func TestWasmPluginCRUD(t *testing.T) {
	cfg := &config.HigressServiceConfig{
		ControllerServiceHost: "localhost",
		ControllerServicePort: 8080,
		ControllerNamespace:   "higress-system",
		ControllerJwtPolicy:   constant.JwtPolicyFirstPartyJwt,
	}

	svc, err := NewKubernetesClientService(cfg)
	if err != nil {
		t.Fatalf("Failed to create KubernetesClientService: %v", err)
	}

	ctx := context.Background()

	// Create test WasmPlugin
	plugin := wasm.NewV1alpha1WasmPlugin()
	plugin.Metadata.Name = "test-wasm-plugin-sdk"
	plugin.Metadata.Namespace = "higress-system"
	plugin.Spec.Url = "oci://test-plugin:v1"
	plugin.Spec.Phase = wasm.PluginPhaseAuthn
	plugin.Spec.Priority = ptrInt64(100)

	// Test Create
	created, err := svc.CreateWasmPlugin(ctx, plugin)
	if err != nil {
		t.Logf("CreateWasmPlugin failed (expected if no cluster): %v", err)
		return
	}
	t.Logf("Created WasmPlugin: %s", created.Metadata.Name)

	// Test Get
	got, err := svc.GetWasmPlugin(ctx, "test-wasm-plugin-sdk")
	if err != nil {
		t.Errorf("GetWasmPlugin failed: %v", err)
	}
	if got == nil {
		t.Error("Expected WasmPlugin to be found")
	}

	// Test List
	list, err := svc.ListWasmPlugins(ctx, "", "", nil)
	if err != nil {
		t.Errorf("ListWasmPlugins failed: %v", err)
	}
	t.Logf("Found %d WasmPlugins", len(list))

	// Test Delete
	err = svc.DeleteWasmPlugin(ctx, "test-wasm-plugin-sdk")
	if err != nil {
		t.Errorf("DeleteWasmPlugin failed: %v", err)
	}
	t.Log("Deleted WasmPlugin")
}

// TestMcpBridgeCRUD tests McpBridge CRUD operations (requires a running Kubernetes cluster)
// 测试 McpBridge CRUD 操作（需要运行中的 Kubernetes 集群）
// 调用方式: go test -v -run TestMcpBridgeCRUD ./internal/kubernetes/
// 注意: 此测试需要真实的 Kubernetes 集群环境
func TestMcpBridgeCRUD(t *testing.T) {
	cfg := &config.HigressServiceConfig{
		ControllerServiceHost: "localhost",
		ControllerServicePort: 8080,
		ControllerNamespace:   "higress-system",
		ControllerJwtPolicy:   constant.JwtPolicyFirstPartyJwt,
	}

	svc, err := NewKubernetesClientService(cfg)
	if err != nil {
		t.Fatalf("Failed to create KubernetesClientService: %v", err)
	}

	ctx := context.Background()

	// Create test McpBridge
	bridge := mcp.NewV1McpBridge()
	bridge.Metadata.Name = "test-mcp-bridge-sdk"
	bridge.Metadata.Namespace = "higress-system"
	bridge.Spec.Registries = []*mcp.V1RegistryConfig{
		{
			Name:           "test-nacos",
			Type:           "nacos",
			Domain:         "nacos.default.svc.cluster.local",
			Port:           8848,
			NacosNamespace: "public",
		},
	}

	// Test Create
	created, err := svc.CreateMcpBridge(ctx, bridge)
	if err != nil {
		t.Logf("CreateMcpBridge failed (expected if no cluster): %v", err)
		return
	}
	t.Logf("Created McpBridge: %s", created.Metadata.Name)

	// Test Get
	got, err := svc.GetMcpBridge(ctx, "test-mcp-bridge-sdk")
	if err != nil {
		t.Errorf("GetMcpBridge failed: %v", err)
	}
	if got == nil {
		t.Error("Expected McpBridge to be found")
	}

	// Test List
	list, err := svc.ListMcpBridges(ctx)
	if err != nil {
		t.Errorf("ListMcpBridges failed: %v", err)
	}
	t.Logf("Found %d McpBridges", len(list))

	// Test Delete
	err = svc.DeleteMcpBridge(ctx, "test-mcp-bridge-sdk")
	if err != nil {
		t.Errorf("DeleteMcpBridge failed: %v", err)
	}
	t.Log("Deleted McpBridge")
}

// TestEnvoyFilterCRUD tests EnvoyFilter CRUD operations (requires a running Kubernetes cluster)
// 测试 EnvoyFilter CRUD 操作（需要运行中的 Kubernetes 集群）
// 调用方式: go test -v -run TestEnvoyFilterCRUD ./internal/kubernetes/
// 注意: 此测试需要真实的 Kubernetes 集群环境
func TestEnvoyFilterCRUD(t *testing.T) {
	cfg := &config.HigressServiceConfig{
		ControllerServiceHost: "localhost",
		ControllerServicePort: 8080,
		ControllerNamespace:   "higress-system",
		ControllerJwtPolicy:   constant.JwtPolicyFirstPartyJwt,
	}

	svc, err := NewKubernetesClientService(cfg)
	if err != nil {
		t.Fatalf("Failed to create KubernetesClientService: %v", err)
	}

	ctx := context.Background()

	// Create test EnvoyFilter
	filter := istio.NewV1alpha3EnvoyFilter()
	filter.Metadata.Name = "test-envoy-filter-sdk"
	filter.Metadata.Namespace = "higress-system"
	filter.Spec.ConfigPatches = []*istio.V1alpha3EnvoyConfigObjectPatch{
		{
			ApplyTo: istio.ApplyToHTTPFilter,
			Patch: &istio.V1alpha3Patch{
				Operation: "INSERT_BEFORE",
				Value:     map[string]interface{}{"name": "test-filter"},
			},
		},
	}

	// Test Create
	created, err := svc.CreateEnvoyFilter(ctx, filter)
	if err != nil {
		t.Logf("CreateEnvoyFilter failed (expected if no cluster): %v", err)
		return
	}
	t.Logf("Created EnvoyFilter: %s", created.Metadata.Name)

	// Test Get
	got, err := svc.GetEnvoyFilter(ctx, "test-envoy-filter-sdk")
	if err != nil {
		t.Errorf("GetEnvoyFilter failed: %v", err)
	}
	if got == nil {
		t.Error("Expected EnvoyFilter to be found")
	}

	// Test Delete
	err = svc.DeleteEnvoyFilter(ctx, "test-envoy-filter-sdk")
	if err != nil {
		t.Errorf("DeleteEnvoyFilter failed: %v", err)
	}
	t.Log("Deleted EnvoyFilter")
}

// Helper function to create int64 pointer
func ptrInt64(v int64) *int64 {
	return &v
}
