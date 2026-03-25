// Package config provides configuration for Higress Admin SDK.
package config

import (
	"testing"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/constant"
	"github.com/stretchr/testify/assert"
)

func TestNewHigressServiceConfig_Defaults(t *testing.T) {
	cfg := NewHigressServiceConfig()

	// Verify default values
	assert.Equal(t, constant.NSDefault, cfg.ControllerNamespace)
	assert.Equal(t, constant.ControllerIngressClassNameDefault, cfg.ControllerWatchedIngressClassName)
	assert.Equal(t, constant.ControllerServiceNameDefault, cfg.ControllerServiceName)
	assert.Equal(t, constant.ControllerServiceHostDefault, cfg.ControllerServiceHost)
	assert.Equal(t, constant.ControllerServicePortDefault, cfg.ControllerServicePort)
	assert.Equal(t, constant.ControllerJwtPolicyDefault, cfg.ControllerJwtPolicy)
	assert.Equal(t, constant.ServiceListSupportRegistryDefault, cfg.ServiceListSupportRegistry)
	assert.Equal(t, constant.ClusterDomainSuffixDefault, cfg.ClusterDomainSuffix)
	assert.NotNil(t, cfg.WasmPluginServiceConfig)
	assert.Empty(t, cfg.KubeConfigPath)
	assert.Empty(t, cfg.KubeConfigContent)
	assert.Empty(t, cfg.ControllerWatchedNamespace)
	assert.Empty(t, cfg.ControllerAccessToken)
}

func TestNewHigressServiceConfig_WithOptions(t *testing.T) {
	wasmCfg := &WasmPluginServiceConfig{
		ImageRegistry:         "registry.example.com",
		ImagePullPolicy:       "Always",
		CustomImageUrlPattern: "custom/${name}:${version}",
		ImagePullSecret:       "my-secret",
	}

	cfg := NewHigressServiceConfig(
		WithKubeConfigPath("/path/to/kubeconfig"),
		WithKubeConfigContent("kubeconfig-content"),
		WithControllerNamespace("custom-namespace"),
		WithControllerWatchedNamespace("watched-ns"),
		WithControllerWatchedIngressClassName("custom-ingress"),
		WithControllerServiceName("custom-controller"),
		WithControllerServiceHost("custom-host.example.com"),
		WithControllerServicePort(8080),
		WithControllerJwtPolicy("third-party-jwt"),
		WithControllerAccessToken("my-token"),
		WithWasmPluginServiceConfig(wasmCfg),
		WithServiceListSupportRegistry(false),
		WithClusterDomainSuffix("custom.local"),
	)

	// Verify all options are applied
	assert.Equal(t, "/path/to/kubeconfig", cfg.KubeConfigPath)
	assert.Equal(t, "kubeconfig-content", cfg.KubeConfigContent)
	assert.Equal(t, "custom-namespace", cfg.ControllerNamespace)
	assert.Equal(t, "watched-ns", cfg.ControllerWatchedNamespace)
	assert.Equal(t, "custom-ingress", cfg.ControllerWatchedIngressClassName)
	assert.Equal(t, "custom-controller", cfg.ControllerServiceName)
	assert.Equal(t, "custom-host.example.com", cfg.ControllerServiceHost)
	assert.Equal(t, 8080, cfg.ControllerServicePort)
	assert.Equal(t, "third-party-jwt", cfg.ControllerJwtPolicy)
	assert.Equal(t, "my-token", cfg.ControllerAccessToken)
	assert.Equal(t, wasmCfg, cfg.WasmPluginServiceConfig)
	assert.False(t, cfg.ServiceListSupportRegistry)
	assert.Equal(t, "custom.local", cfg.ClusterDomainSuffix)
}

func TestWithKubeConfigPath(t *testing.T) {
	cfg := NewHigressServiceConfig(WithKubeConfigPath("/test/path"))
	assert.Equal(t, "/test/path", cfg.KubeConfigPath)
}

func TestWithKubeConfigContent(t *testing.T) {
	cfg := NewHigressServiceConfig(WithKubeConfigContent("test-content"))
	assert.Equal(t, "test-content", cfg.KubeConfigContent)
}

func TestWithControllerNamespace(t *testing.T) {
	cfg := NewHigressServiceConfig(WithControllerNamespace("test-ns"))
	assert.Equal(t, "test-ns", cfg.ControllerNamespace)
}

func TestWithControllerWatchedNamespace(t *testing.T) {
	cfg := NewHigressServiceConfig(WithControllerWatchedNamespace("watched-ns"))
	assert.Equal(t, "watched-ns", cfg.ControllerWatchedNamespace)
}

func TestWithControllerWatchedIngressClassName(t *testing.T) {
	cfg := NewHigressServiceConfig(WithControllerWatchedIngressClassName("test-class"))
	assert.Equal(t, "test-class", cfg.ControllerWatchedIngressClassName)
}

func TestWithControllerServiceName(t *testing.T) {
	cfg := NewHigressServiceConfig(WithControllerServiceName("test-service"))
	assert.Equal(t, "test-service", cfg.ControllerServiceName)
}

func TestWithControllerServiceHost(t *testing.T) {
	cfg := NewHigressServiceConfig(WithControllerServiceHost("test.host.com"))
	assert.Equal(t, "test.host.com", cfg.ControllerServiceHost)
}

func TestWithControllerServicePort(t *testing.T) {
	cfg := NewHigressServiceConfig(WithControllerServicePort(9999))
	assert.Equal(t, 9999, cfg.ControllerServicePort)
}

func TestWithControllerJwtPolicy(t *testing.T) {
	cfg := NewHigressServiceConfig(WithControllerJwtPolicy("test-policy"))
	assert.Equal(t, "test-policy", cfg.ControllerJwtPolicy)
}

func TestWithControllerAccessToken(t *testing.T) {
	cfg := NewHigressServiceConfig(WithControllerAccessToken("test-token"))
	assert.Equal(t, "test-token", cfg.ControllerAccessToken)
}

func TestWithWasmPluginServiceConfig(t *testing.T) {
	wasmCfg := &WasmPluginServiceConfig{
		ImageRegistry:   "test-registry",
		ImagePullPolicy: "IfNotPresent",
	}
	cfg := NewHigressServiceConfig(WithWasmPluginServiceConfig(wasmCfg))
	assert.Equal(t, wasmCfg, cfg.WasmPluginServiceConfig)
	assert.Equal(t, "test-registry", cfg.WasmPluginServiceConfig.ImageRegistry)
}

func TestWithServiceListSupportRegistry(t *testing.T) {
	cfg := NewHigressServiceConfig(WithServiceListSupportRegistry(false))
	assert.False(t, cfg.ServiceListSupportRegistry)

	cfg2 := NewHigressServiceConfig(WithServiceListSupportRegistry(true))
	assert.True(t, cfg2.ServiceListSupportRegistry)
}

func TestWithClusterDomainSuffix(t *testing.T) {
	cfg := NewHigressServiceConfig(WithClusterDomainSuffix("test.local"))
	assert.Equal(t, "test.local", cfg.ClusterDomainSuffix)
}

func TestHigressServiceConfig_Getters(t *testing.T) {
	wasmCfg := &WasmPluginServiceConfig{
		ImageRegistry:   "getter-registry",
		ImagePullPolicy: "Always",
	}

	cfg := &HigressServiceConfig{
		KubeConfigPath:                    "/getter/path",
		KubeConfigContent:                 "getter-content",
		ControllerNamespace:               "getter-ns",
		ControllerWatchedNamespace:        "getter-watched",
		ControllerWatchedIngressClassName: "getter-class",
		ControllerServiceName:             "getter-service",
		ControllerServiceHost:             "getter.host.com",
		ControllerServicePort:             12345,
		ControllerJwtPolicy:               "getter-policy",
		ControllerAccessToken:             "getter-token",
		WasmPluginServiceConfig:           wasmCfg,
		ServiceListSupportRegistry:        true,
		ClusterDomainSuffix:               "getter.local",
	}

	// Test all getters
	assert.Equal(t, "/getter/path", cfg.GetKubeConfigPath())
	assert.Equal(t, "getter-content", cfg.GetKubeConfigContent())
	assert.Equal(t, "getter-ns", cfg.GetControllerNamespace())
	assert.Equal(t, "getter-watched", cfg.GetControllerWatchedNamespace())
	assert.Equal(t, "getter-class", cfg.GetControllerWatchedIngressClassName())
	assert.Equal(t, "getter-service", cfg.GetControllerServiceName())
	assert.Equal(t, "getter.host.com", cfg.GetControllerServiceHost())
	assert.Equal(t, 12345, cfg.GetControllerServicePort())
	assert.Equal(t, "getter-policy", cfg.GetControllerJwtPolicy())
	assert.Equal(t, "getter-token", cfg.GetControllerAccessToken())
	assert.Equal(t, wasmCfg, cfg.GetWasmPluginServiceConfig())
	assert.True(t, cfg.GetServiceListSupportRegistry())
	assert.Equal(t, "getter.local", cfg.GetClusterDomainSuffix())
}

func TestWasmPluginServiceConfig_Fields(t *testing.T) {
	wasmCfg := &WasmPluginServiceConfig{
		ImageRegistry:         "wasm-registry.example.com",
		ImagePullPolicy:       "IfNotPresent",
		CustomImageUrlPattern: "custom/${name}:${version}",
		ImagePullSecret:       "pull-secret",
	}

	assert.Equal(t, "wasm-registry.example.com", wasmCfg.ImageRegistry)
	assert.Equal(t, "IfNotPresent", wasmCfg.ImagePullPolicy)
	assert.Equal(t, "custom/${name}:${version}", wasmCfg.CustomImageUrlPattern)
	assert.Equal(t, "pull-secret", wasmCfg.ImagePullSecret)
}

func TestNewHigressServiceConfig_MultipleOptions(t *testing.T) {
	// Test that multiple options can be combined
	cfg := NewHigressServiceConfig(
		WithKubeConfigPath("/path1"),
		WithKubeConfigPath("/path2"), // Second call should override
		WithControllerNamespace("ns1"),
		WithControllerServicePort(1111),
		WithControllerServicePort(2222), // Second call should override
	)

	assert.Equal(t, "/path2", cfg.KubeConfigPath) // Last one wins
	assert.Equal(t, "ns1", cfg.ControllerNamespace)
	assert.Equal(t, 2222, cfg.ControllerServicePort) // Last one wins
}

func TestNewHigressServiceConfig_EmptyOptions(t *testing.T) {
	// Test with no options
	cfg := NewHigressServiceConfig()
	assert.NotNil(t, cfg)
	assert.Equal(t, constant.NSDefault, cfg.ControllerNamespace)
}

func TestHigressServiceConfig_NilWasmPluginConfig(t *testing.T) {
	cfg := &HigressServiceConfig{
		WasmPluginServiceConfig: nil,
	}

	assert.Nil(t, cfg.GetWasmPluginServiceConfig())
}

func TestHigressServiceConfig_EmptyStrings(t *testing.T) {
	cfg := &HigressServiceConfig{
		KubeConfigPath:                    "",
		KubeConfigContent:                 "",
		ControllerNamespace:               "",
		ControllerWatchedNamespace:        "",
		ControllerWatchedIngressClassName: "",
		ControllerServiceName:             "",
		ControllerServiceHost:             "",
		ControllerJwtPolicy:               "",
		ControllerAccessToken:             "",
		ClusterDomainSuffix:               "",
	}

	assert.Empty(t, cfg.GetKubeConfigPath())
	assert.Empty(t, cfg.GetKubeConfigContent())
	assert.Empty(t, cfg.GetControllerNamespace())
	assert.Empty(t, cfg.GetControllerWatchedNamespace())
	assert.Empty(t, cfg.GetControllerWatchedIngressClassName())
	assert.Empty(t, cfg.GetControllerServiceName())
	assert.Empty(t, cfg.GetControllerServiceHost())
	assert.Empty(t, cfg.GetControllerJwtPolicy())
	assert.Empty(t, cfg.GetControllerAccessToken())
	assert.Empty(t, cfg.GetClusterDomainSuffix())
}

func TestHigressServiceConfig_ZeroPort(t *testing.T) {
	cfg := &HigressServiceConfig{
		ControllerServicePort: 0,
	}

	assert.Equal(t, 0, cfg.GetControllerServicePort())
}

func TestOption_FunctionPattern(t *testing.T) {
	// Test that Option type works as a function
	var opt Option = WithControllerNamespace("test-ns")
	assert.NotNil(t, opt)

	cfg := &HigressServiceConfig{}
	opt(cfg)
	assert.Equal(t, "test-ns", cfg.ControllerNamespace)
}
