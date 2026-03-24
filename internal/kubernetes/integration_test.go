//go:build integration
// +build integration

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

// TestIngressCRUD tests Ingress CRUD operations (requires a running Kubernetes cluster)
// 测试 Ingress CRUD 操作（需要运行中的 Kubernetes 集群）
// 调用方式: go test -v -tags=integration -run TestIngressCRUD ./internal/kubernetes/
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
		t.Fatalf("CreateIngress failed: %v", err)
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
// 调用方式: go test -v -tags=integration -run TestSecretCRUD ./internal/kubernetes/
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
		t.Fatalf("CreateSecret failed: %v", err)
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
// 调用方式: go test -v -tags=integration -run TestConfigMapCRUD ./internal/kubernetes/
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
		t.Fatalf("CreateConfigMap failed: %v", err)
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
// 调用方式: go test -v -tags=integration -run TestWasmPluginCRUD ./internal/kubernetes/
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
		t.Fatalf("CreateWasmPlugin failed: %v", err)
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
// 调用方式: go test -v -tags=integration -run TestMcpBridgeCRUD ./internal/kubernetes/
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
		t.Fatalf("CreateMcpBridge failed: %v", err)
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
// 调用方式: go test -v -tags=integration -run TestEnvoyFilterCRUD ./internal/kubernetes/
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
		t.Fatalf("CreateEnvoyFilter failed: %v", err)
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
