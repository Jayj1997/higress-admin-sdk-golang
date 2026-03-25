// Package kubernetes provides Kubernetes client functionality
package kubernetes

import (
	"os"
	"testing"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/internal/kubernetes/crd/istio"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/internal/kubernetes/crd/mcp"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/internal/kubernetes/crd/wasm"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/config"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/constant"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// skipIfNoKubernetes skips the test if no Kubernetes cluster is available
// 如果没有可用的 Kubernetes 集群，则跳过测试
func skipIfNoKubernetes(t *testing.T) {
	t.Helper()
	// Check if we're in a Kubernetes cluster or have kubeconfig available
	// 检查是否在 Kubernetes 集群中或有可用的 kubeconfig
	if isInCluster() {
		return
	}

	// Check for kubeconfig file
	// 检查 kubeconfig 文件
	home := os.Getenv("HOME")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home != "" {
		kubeConfigPath := home + "/.kube/config"
		if _, err := os.Stat(kubeConfigPath); err == nil {
			return
		}
	}

	// Check for KUBECONFIG environment variable
	// 检查 KUBECONFIG 环境变量
	kubeConfig := os.Getenv("KUBECONFIG")
	if kubeConfig != "" {
		if _, err := os.Stat(kubeConfig); err == nil {
			return
		}
	}

	t.Skip("Skipping test: no Kubernetes cluster available (not in-cluster and no kubeconfig found)")
}

// TestNewKubernetesClientService tests the creation of KubernetesClientService
// 测试 KubernetesClientService 的创建
// 调用方式: go test -v -run TestNewKubernetesClientService ./internal/kubernetes/
func TestNewKubernetesClientService(t *testing.T) {
	tests := []struct {
		name       string
		config     *config.HigressServiceConfig
		wantErr    bool
		needK8s    bool // Whether this test case needs a real Kubernetes cluster
		skipReason string
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
			wantErr:    false,
			needK8s:    true, // This test case needs a real Kubernetes cluster
			skipReason: "no Kubernetes cluster available",
		},
		{
			name: "missing controller namespace",
			config: &config.HigressServiceConfig{
				ControllerServiceHost: "localhost",
				ControllerServicePort: 8080,
				ControllerJwtPolicy:   constant.JwtPolicyFirstPartyJwt,
			},
			wantErr: true,
			needK8s: false, // This test case only validates config, doesn't need K8s
		},
		{
			name: "missing controller service host",
			config: &config.HigressServiceConfig{
				ControllerServicePort: 8080,
				ControllerNamespace:   "higress-system",
				ControllerJwtPolicy:   constant.JwtPolicyFirstPartyJwt,
			},
			wantErr: true,
			needK8s: false,
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
			needK8s: false,
		},
		{
			name: "missing jwt policy",
			config: &config.HigressServiceConfig{
				ControllerServiceHost: "localhost",
				ControllerServicePort: 8080,
				ControllerNamespace:   "higress-system",
			},
			wantErr: true,
			needK8s: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip test cases that need Kubernetes if no cluster is available
			// 如果没有可用的集群，跳过需要 Kubernetes 的测试用例
			if tt.needK8s {
				skipIfNoKubernetes(t)
			}

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
	// Skip if no Kubernetes cluster is available
	// 如果没有可用的 Kubernetes 集群，则跳过
	skipIfNoKubernetes(t)

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
	// Skip if no Kubernetes cluster is available
	// 如果没有可用的 Kubernetes 集群，则跳过
	skipIfNoKubernetes(t)

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

// Helper function to create int64 pointer
func ptrInt64(v int64) *int64 {
	return &v
}
