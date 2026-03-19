// Package kubernetes provides Kubernetes client functionality
package kubernetes

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jayj1997/higress-admin-sdk-golang/internal/kubernetes/crd/istio"
	"github.com/Jayj1997/higress-admin-sdk-golang/internal/kubernetes/crd/mcp"
	"github.com/Jayj1997/higress-admin-sdk-golang/internal/kubernetes/crd/wasm"
	"github.com/Jayj1997/higress-admin-sdk-golang/internal/kubernetes/model"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/config"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/constant"
	jsoniter "github.com/json-iterator/go"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// KubernetesClientService provides Kubernetes API operations
type KubernetesClientService struct {
	config             *config.HigressServiceConfig
	clientset          *kubernetes.Clientset
	dynamicClient      dynamic.Interface
	restConfig         *rest.Config
	httpClient         *http.Client
	inClusterMode      bool
	defaultLabels      string
	ingressV1Supported bool
}

// NewKubernetesClientService creates a new KubernetesClientService
func NewKubernetesClientService(cfg *config.HigressServiceConfig) (*KubernetesClientService, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	svc := &KubernetesClientService{
		config:        cfg,
		httpClient:    &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}},
		defaultLabels: buildLabelSelector(constant.LabelResourceDefinerKey, constant.LabelResourceDefinerValue),
	}

	var restConfig *rest.Config
	var err error

	// Determine authentication mode
	svc.inClusterMode = isInCluster() && cfg.GetKubeConfigPath() == "" && cfg.GetKubeConfigContent() == ""

	if svc.inClusterMode {
		// InCluster mode - use ServiceAccount
		restConfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to get in-cluster config: %w", err)
		}
	} else if cfg.GetKubeConfigContent() != "" {
		// KubeConfig content mode
		restConfig, err = clientcmd.RESTConfigFromKubeConfig([]byte(cfg.GetKubeConfigContent()))
		if err != nil {
			return nil, fmt.Errorf("failed to parse kubeconfig content: %w", err)
		}
	} else {
		// KubeConfig file mode
		kubeConfigPath := cfg.GetKubeConfigPath()
		if kubeConfigPath == "" {
			home := os.Getenv("HOME")
			if home == "" {
				home = os.Getenv("USERPROFILE")
			}
			kubeConfigPath = filepath.Join(home, ".kube", "config")
		}

		restConfig, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			return nil, fmt.Errorf("failed to build config from kubeconfig: %w", err)
		}
	}

	svc.restConfig = restConfig

	// Create clientset
	svc.clientset, err = kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	// Create dynamic client for CRDs
	svc.dynamicClient, err = dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	// Check API capabilities
	svc.ingressV1Supported = true

	return svc, nil
}

// IsNamespaceProtected checks if a namespace is protected
func (s *KubernetesClientService) IsNamespaceProtected(namespace string) bool {
	return namespace == constant.KubeSystemNamespace || namespace == s.config.GetControllerNamespace()
}

// ==================== Ingress Operations ====================

// ListIngresses lists all ingresses in the controller namespace
func (s *KubernetesClientService) ListIngresses(ctx context.Context) ([]networkingv1.Ingress, error) {
	list, err := s.clientset.NetworkingV1().Ingresses(s.config.GetControllerNamespace()).List(ctx, metav1.ListOptions{
		LabelSelector: s.defaultLabels,
	})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

// ListAllIngresses lists all ingresses across namespaces
func (s *KubernetesClientService) ListAllIngresses(ctx context.Context) ([]networkingv1.Ingress, error) {
	watchedNamespace := s.config.GetControllerWatchedNamespace()
	var ingresses []networkingv1.Ingress

	if watchedNamespace == "" {
		list, err := s.clientset.NetworkingV1().Ingresses("").List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		ingresses = list.Items
	} else {
		for _, ns := range []string{s.config.GetControllerNamespace(), watchedNamespace} {
			list, err := s.clientset.NetworkingV1().Ingresses(ns).List(ctx, metav1.ListOptions{})
			if err != nil {
				return nil, err
			}
			ingresses = append(ingresses, list.Items...)
		}
	}

	return ingresses, nil
}

// GetIngress gets an ingress by name
func (s *KubernetesClientService) GetIngress(ctx context.Context, name string) (*networkingv1.Ingress, error) {
	ingress, err := s.clientset.NetworkingV1().Ingresses(s.config.GetControllerNamespace()).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ingress, nil
}

// CreateIngress creates an ingress
func (s *KubernetesClientService) CreateIngress(ctx context.Context, ingress *networkingv1.Ingress) (*networkingv1.Ingress, error) {
	s.renderDefaultMetadata(ingress)
	return s.clientset.NetworkingV1().Ingresses(s.config.GetControllerNamespace()).Create(ctx, ingress, metav1.CreateOptions{})
}

// UpdateIngress updates an ingress
func (s *KubernetesClientService) UpdateIngress(ctx context.Context, ingress *networkingv1.Ingress) (*networkingv1.Ingress, error) {
	s.renderDefaultMetadata(ingress)
	return s.clientset.NetworkingV1().Ingresses(s.config.GetControllerNamespace()).Update(ctx, ingress, metav1.UpdateOptions{})
}

// DeleteIngress deletes an ingress
func (s *KubernetesClientService) DeleteIngress(ctx context.Context, name string) error {
	err := s.clientset.NetworkingV1().Ingresses(s.config.GetControllerNamespace()).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

// ==================== Secret Operations ====================

// ListSecrets lists secrets in the controller namespace
func (s *KubernetesClientService) ListSecrets(ctx context.Context, secretType string) ([]corev1.Secret, error) {
	fieldSelector := ""
	if secretType != "" {
		fieldSelector = fmt.Sprintf("type=%s", secretType)
	}
	list, err := s.clientset.CoreV1().Secrets(s.config.GetControllerNamespace()).List(ctx, metav1.ListOptions{
		LabelSelector: s.defaultLabels,
		FieldSelector: fieldSelector,
	})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

// GetSecret gets a secret by name
func (s *KubernetesClientService) GetSecret(ctx context.Context, name string) (*corev1.Secret, error) {
	secret, err := s.clientset.CoreV1().Secrets(s.config.GetControllerNamespace()).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return secret, nil
}

// CreateSecret creates a secret
func (s *KubernetesClientService) CreateSecret(ctx context.Context, secret *corev1.Secret) (*corev1.Secret, error) {
	s.renderDefaultMetadata(secret)
	return s.clientset.CoreV1().Secrets(s.config.GetControllerNamespace()).Create(ctx, secret, metav1.CreateOptions{})
}

// UpdateSecret updates a secret
func (s *KubernetesClientService) UpdateSecret(ctx context.Context, secret *corev1.Secret) (*corev1.Secret, error) {
	s.renderDefaultMetadata(secret)
	return s.clientset.CoreV1().Secrets(s.config.GetControllerNamespace()).Update(ctx, secret, metav1.UpdateOptions{})
}

// DeleteSecret deletes a secret
func (s *KubernetesClientService) DeleteSecret(ctx context.Context, name string) error {
	err := s.clientset.CoreV1().Secrets(s.config.GetControllerNamespace()).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

// ==================== ConfigMap Operations ====================

// ListConfigMaps lists configmaps in the controller namespace
func (s *KubernetesClientService) ListConfigMaps(ctx context.Context, labelSelectors map[string]string) ([]corev1.ConfigMap, error) {
	labelSelector := s.defaultLabels
	if len(labelSelectors) > 0 {
		var selectors []string
		for k, v := range labelSelectors {
			selectors = append(selectors, buildLabelSelector(k, v))
		}
		selectors = append(selectors, s.defaultLabels)
		labelSelector = strings.Join(selectors, ",")
	}
	list, err := s.clientset.CoreV1().ConfigMaps(s.config.GetControllerNamespace()).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

// GetConfigMap gets a configmap by name
func (s *KubernetesClientService) GetConfigMap(ctx context.Context, name string) (*corev1.ConfigMap, error) {
	cm, err := s.clientset.CoreV1().ConfigMaps(s.config.GetControllerNamespace()).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return cm, nil
}

// CreateConfigMap creates a configmap
func (s *KubernetesClientService) CreateConfigMap(ctx context.Context, cm *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	s.renderDefaultMetadata(cm)
	return s.clientset.CoreV1().ConfigMaps(s.config.GetControllerNamespace()).Create(ctx, cm, metav1.CreateOptions{})
}

// UpdateConfigMap updates a configmap
func (s *KubernetesClientService) UpdateConfigMap(ctx context.Context, cm *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	s.renderDefaultMetadata(cm)
	return s.clientset.CoreV1().ConfigMaps(s.config.GetControllerNamespace()).Update(ctx, cm, metav1.UpdateOptions{})
}

// DeleteConfigMap deletes a configmap
func (s *KubernetesClientService) DeleteConfigMap(ctx context.Context, name string) error {
	err := s.clientset.CoreV1().ConfigMaps(s.config.GetControllerNamespace()).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

// ==================== WasmPlugin CRD Operations ====================

var wasmPluginGVR = schema.GroupVersionResource{
	Group:    wasm.WasmPluginAPIGroup,
	Version:  wasm.WasmPluginAPIVersion,
	Resource: wasm.WasmPluginPlural,
}

// ListWasmPlugins lists WasmPlugins
func (s *KubernetesClientService) ListWasmPlugins(ctx context.Context, name, version string, builtIn *bool) ([]*wasm.V1alpha1WasmPlugin, error) {
	labelSelectors := []string{s.defaultLabels}
	if name != "" {
		labelSelectors = append(labelSelectors, buildLabelSelector(constant.LabelWasmPluginNameKey, name))
	}
	if version != "" {
		labelSelectors = append(labelSelectors, buildLabelSelector(constant.LabelWasmPluginVersionKey, version))
	}
	if builtIn != nil {
		labelSelectors = append(labelSelectors, buildLabelSelector(constant.LabelWasmPluginBuiltInKey, fmt.Sprintf("%v", *builtIn)))
	}

	list, err := s.dynamicClient.Resource(wasmPluginGVR).Namespace(s.config.GetControllerNamespace()).List(ctx, metav1.ListOptions{
		LabelSelector: strings.Join(labelSelectors, ","),
	})
	if err != nil {
		return nil, err
	}

	var plugins []*wasm.V1alpha1WasmPlugin
	for _, item := range list.Items {
		plugin := unstructuredToWasmPlugin(&item)
		if plugin != nil {
			plugins = append(plugins, plugin)
		}
	}
	return plugins, nil
}

// GetWasmPlugin gets a WasmPlugin by name
func (s *KubernetesClientService) GetWasmPlugin(ctx context.Context, name string) (*wasm.V1alpha1WasmPlugin, error) {
	obj, err := s.dynamicClient.Resource(wasmPluginGVR).Namespace(s.config.GetControllerNamespace()).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return unstructuredToWasmPlugin(obj), nil
}

// CreateWasmPlugin creates a WasmPlugin
func (s *KubernetesClientService) CreateWasmPlugin(ctx context.Context, plugin *wasm.V1alpha1WasmPlugin) (*wasm.V1alpha1WasmPlugin, error) {
	obj := wasmPluginToUnstructured(plugin)
	s.renderDefaultMetadataUnstructured(obj)
	result, err := s.dynamicClient.Resource(wasmPluginGVR).Namespace(s.config.GetControllerNamespace()).Create(ctx, obj, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return unstructuredToWasmPlugin(result), nil
}

// UpdateWasmPlugin updates a WasmPlugin
func (s *KubernetesClientService) UpdateWasmPlugin(ctx context.Context, plugin *wasm.V1alpha1WasmPlugin) (*wasm.V1alpha1WasmPlugin, error) {
	obj := wasmPluginToUnstructured(plugin)
	s.renderDefaultMetadataUnstructured(obj)
	result, err := s.dynamicClient.Resource(wasmPluginGVR).Namespace(s.config.GetControllerNamespace()).Update(ctx, obj, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return unstructuredToWasmPlugin(result), nil
}

// DeleteWasmPlugin deletes a WasmPlugin
func (s *KubernetesClientService) DeleteWasmPlugin(ctx context.Context, name string) error {
	err := s.dynamicClient.Resource(wasmPluginGVR).Namespace(s.config.GetControllerNamespace()).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

// ==================== McpBridge CRD Operations ====================

var mcpBridgeGVR = schema.GroupVersionResource{
	Group:    mcp.McpBridgeAPIGroup,
	Version:  mcp.McpBridgeAPIVersion,
	Resource: mcp.McpBridgePlural,
}

// ListMcpBridges lists McpBridges
func (s *KubernetesClientService) ListMcpBridges(ctx context.Context) ([]*mcp.V1McpBridge, error) {
	list, err := s.dynamicClient.Resource(mcpBridgeGVR).Namespace(s.config.GetControllerNamespace()).List(ctx, metav1.ListOptions{
		LabelSelector: s.defaultLabels,
	})
	if err != nil {
		return nil, err
	}

	var bridges []*mcp.V1McpBridge
	for _, item := range list.Items {
		bridge := unstructuredToMcpBridge(&item)
		if bridge != nil {
			bridges = append(bridges, bridge)
		}
	}
	return bridges, nil
}

// GetMcpBridge gets a McpBridge by name
func (s *KubernetesClientService) GetMcpBridge(ctx context.Context, name string) (*mcp.V1McpBridge, error) {
	obj, err := s.dynamicClient.Resource(mcpBridgeGVR).Namespace(s.config.GetControllerNamespace()).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return unstructuredToMcpBridge(obj), nil
}

// CreateMcpBridge creates a McpBridge
func (s *KubernetesClientService) CreateMcpBridge(ctx context.Context, bridge *mcp.V1McpBridge) (*mcp.V1McpBridge, error) {
	obj := mcpBridgeToUnstructured(bridge)
	s.renderDefaultMetadataUnstructured(obj)
	result, err := s.dynamicClient.Resource(mcpBridgeGVR).Namespace(s.config.GetControllerNamespace()).Create(ctx, obj, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return unstructuredToMcpBridge(result), nil
}

// UpdateMcpBridge updates a McpBridge
func (s *KubernetesClientService) UpdateMcpBridge(ctx context.Context, bridge *mcp.V1McpBridge) (*mcp.V1McpBridge, error) {
	obj := mcpBridgeToUnstructured(bridge)
	s.renderDefaultMetadataUnstructured(obj)
	result, err := s.dynamicClient.Resource(mcpBridgeGVR).Namespace(s.config.GetControllerNamespace()).Update(ctx, obj, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return unstructuredToMcpBridge(result), nil
}

// DeleteMcpBridge deletes a McpBridge
func (s *KubernetesClientService) DeleteMcpBridge(ctx context.Context, name string) error {
	err := s.dynamicClient.Resource(mcpBridgeGVR).Namespace(s.config.GetControllerNamespace()).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

// ==================== EnvoyFilter CRD Operations ====================

var envoyFilterGVR = schema.GroupVersionResource{
	Group:    istio.EnvoyFilterAPIGroup,
	Version:  istio.EnvoyFilterAPIVersion,
	Resource: istio.EnvoyFilterPlural,
}

// GetEnvoyFilter gets an EnvoyFilter by name
func (s *KubernetesClientService) GetEnvoyFilter(ctx context.Context, name string) (*istio.V1alpha3EnvoyFilter, error) {
	obj, err := s.dynamicClient.Resource(envoyFilterGVR).Namespace(s.config.GetControllerNamespace()).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return unstructuredToEnvoyFilter(obj), nil
}

// CreateEnvoyFilter creates an EnvoyFilter
func (s *KubernetesClientService) CreateEnvoyFilter(ctx context.Context, filter *istio.V1alpha3EnvoyFilter) (*istio.V1alpha3EnvoyFilter, error) {
	obj := envoyFilterToUnstructured(filter)
	s.renderDefaultMetadataUnstructured(obj)
	result, err := s.dynamicClient.Resource(envoyFilterGVR).Namespace(s.config.GetControllerNamespace()).Create(ctx, obj, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return unstructuredToEnvoyFilter(result), nil
}

// UpdateEnvoyFilter updates an EnvoyFilter
func (s *KubernetesClientService) UpdateEnvoyFilter(ctx context.Context, filter *istio.V1alpha3EnvoyFilter) (*istio.V1alpha3EnvoyFilter, error) {
	obj := envoyFilterToUnstructured(filter)
	s.renderDefaultMetadataUnstructured(obj)
	result, err := s.dynamicClient.Resource(envoyFilterGVR).Namespace(s.config.GetControllerNamespace()).Update(ctx, obj, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return unstructuredToEnvoyFilter(result), nil
}

// DeleteEnvoyFilter deletes an EnvoyFilter
func (s *KubernetesClientService) DeleteEnvoyFilter(ctx context.Context, name string) error {
	err := s.dynamicClient.Resource(envoyFilterGVR).Namespace(s.config.GetControllerNamespace()).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

// ==================== Controller API Operations ====================

// GatewayServiceList gets the gateway service list from the controller
func (s *KubernetesClientService) GatewayServiceList(ctx context.Context) ([]*model.RegistryzService, error) {
	req, err := s.buildControllerRequest("/debug/registryz")
	if err != nil {
		return nil, err
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get gateway service list: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var services []*model.RegistryzService
	if err := jsoniter.ConfigFastest.Unmarshal(body, &services); err != nil {
		return nil, err
	}
	return services, nil
}

// GatewayServiceEndpoint gets the gateway service endpoints from the controller
func (s *KubernetesClientService) GatewayServiceEndpoint(ctx context.Context) (map[string]map[string]*model.IstioEndpointShard, error) {
	req, err := s.buildControllerRequest("/debug/endpointShardz")
	if err != nil {
		return nil, err
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get service endpoints: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var endpoints map[string]map[string]*model.IstioEndpointShard
	if err := jsoniter.ConfigFastest.Unmarshal(body, &endpoints); err != nil {
		return nil, err
	}
	return endpoints, nil
}

// buildControllerRequest builds an HTTP request to the controller
func (s *KubernetesClientService) buildControllerRequest(path string) (*http.Request, error) {
	var host string
	if s.inClusterMode {
		host = fmt.Sprintf("%s.%s", s.config.GetControllerServiceName(), s.config.GetControllerNamespace())
	} else {
		host = s.config.GetControllerServiceHost()
	}

	url := fmt.Sprintf("http://%s:%d%s", host, s.config.GetControllerServicePort(), path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	token := s.config.GetControllerAccessToken()
	if token == "" && s.inClusterMode {
		token = s.readTokenFromFile()
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	return req, nil
}

// readTokenFromFile reads the service account token from file
func (s *KubernetesClientService) readTokenFromFile() string {
	fileName := "/var/run/secrets/access-token/token"
	if s.config.GetControllerJwtPolicy() == constant.JwtPolicyFirstPartyJwt {
		fileName = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	}
	data, err := os.ReadFile(fileName)
	if err != nil {
		return ""
	}
	return string(data)
}

// ==================== Helper Functions ====================

func (s *KubernetesClientService) renderDefaultMetadata(obj metav1.Object) {
	labels := obj.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}
	labels[constant.LabelResourceDefinerKey] = constant.LabelResourceDefinerValue
	obj.SetLabels(labels)
}

func (s *KubernetesClientService) renderDefaultMetadataUnstructured(obj *unstructured.Unstructured) {
	labels, _, _ := unstructured.NestedStringMap(obj.Object, "metadata", "labels")
	if labels == nil {
		labels = make(map[string]string)
	}
	labels[constant.LabelResourceDefinerKey] = constant.LabelResourceDefinerValue
	unstructured.SetNestedStringMap(obj.Object, labels, "metadata", "labels")
}

func buildLabelSelector(key, value string) string {
	return fmt.Sprintf("%s=%s", key, value)
}

func isInCluster() bool {
	_, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount/token")
	return !os.IsNotExist(err)
}

func validateConfig(cfg *config.HigressServiceConfig) error {
	if isInCluster() {
		if cfg.GetControllerServiceName() == "" {
			return fmt.Errorf("controllerServiceName is required")
		}
	} else {
		if cfg.GetControllerServiceHost() == "" {
			return fmt.Errorf("controllerServiceHost is required")
		}
	}
	if cfg.GetControllerNamespace() == "" {
		return fmt.Errorf("controllerNamespace is required")
	}
	if cfg.GetControllerServicePort() <= 0 || cfg.GetControllerServicePort() > 65535 {
		return fmt.Errorf("controllerServicePort is invalid")
	}
	if cfg.GetControllerJwtPolicy() == "" {
		return fmt.Errorf("controllerJwtPolicy is required")
	}
	return nil
}

// ==================== Conversion Functions ====================

func unstructuredToWasmPlugin(obj *unstructured.Unstructured) *wasm.V1alpha1WasmPlugin {
	data, err := json.Marshal(obj.Object)
	if err != nil {
		return nil
	}
	var plugin wasm.V1alpha1WasmPlugin
	if err := jsoniter.ConfigFastest.Unmarshal(data, &plugin); err != nil {
		return nil
	}
	return &plugin
}

func wasmPluginToUnstructured(plugin *wasm.V1alpha1WasmPlugin) *unstructured.Unstructured {
	data, err := jsoniter.ConfigFastest.Marshal(plugin)
	if err != nil {
		return &unstructured.Unstructured{Object: make(map[string]interface{})}
	}
	var obj map[string]interface{}
	if err := jsoniter.ConfigFastest.Unmarshal(data, &obj); err != nil {
		return &unstructured.Unstructured{Object: make(map[string]interface{})}
	}
	return &unstructured.Unstructured{Object: obj}
}

func unstructuredToMcpBridge(obj *unstructured.Unstructured) *mcp.V1McpBridge {
	data, err := json.Marshal(obj.Object)
	if err != nil {
		return nil
	}
	var bridge mcp.V1McpBridge
	if err := jsoniter.ConfigFastest.Unmarshal(data, &bridge); err != nil {
		return nil
	}
	return &bridge
}

func mcpBridgeToUnstructured(bridge *mcp.V1McpBridge) *unstructured.Unstructured {
	data, err := jsoniter.ConfigFastest.Marshal(bridge)
	if err != nil {
		return &unstructured.Unstructured{Object: make(map[string]interface{})}
	}
	var obj map[string]interface{}
	if err := jsoniter.ConfigFastest.Unmarshal(data, &obj); err != nil {
		return &unstructured.Unstructured{Object: make(map[string]interface{})}
	}
	return &unstructured.Unstructured{Object: obj}
}

func unstructuredToEnvoyFilter(obj *unstructured.Unstructured) *istio.V1alpha3EnvoyFilter {
	data, err := json.Marshal(obj.Object)
	if err != nil {
		return nil
	}
	var filter istio.V1alpha3EnvoyFilter
	if err := jsoniter.ConfigFastest.Unmarshal(data, &filter); err != nil {
		return nil
	}
	return &filter
}

func envoyFilterToUnstructured(filter *istio.V1alpha3EnvoyFilter) *unstructured.Unstructured {
	data, err := jsoniter.ConfigFastest.Marshal(filter)
	if err != nil {
		return &unstructured.Unstructured{Object: make(map[string]interface{})}
	}
	var obj map[string]interface{}
	if err := jsoniter.ConfigFastest.Unmarshal(data, &obj); err != nil {
		return &unstructured.Unstructured{Object: make(map[string]interface{})}
	}
	return &unstructured.Unstructured{Object: obj}
}
