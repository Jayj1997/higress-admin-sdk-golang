// Package client provides the main client for Higress Admin SDK.
package client

import (
	"testing"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/config"
)

func TestNewHigressServiceProvider(t *testing.T) {
	// Test with empty config (should fail because no Kubernetes cluster is available)
	cfg := config.NewHigressServiceConfig()

	_, err := NewHigressServiceProvider(cfg)
	if err == nil {
		t.Log("Provider created successfully (cluster may be available)")
	} else {
		t.Logf("Expected error when no cluster available: %v", err)
	}
}

func TestHigressServiceProviderInterface(t *testing.T) {
	// Test that HigressServiceProviderImpl implements HigressServiceProvider interface
	var _ HigressServiceProvider = (*HigressServiceProviderImpl)(nil)
}

func TestConfigOptions(t *testing.T) {
	tests := []struct {
		name     string
		options  []config.Option
		validate func(t *testing.T, cfg *config.HigressServiceConfig)
	}{
		{
			name: "WithKubeConfigPath",
			options: []config.Option{
				config.WithKubeConfigPath("/path/to/kubeconfig"),
			},
			validate: func(t *testing.T, cfg *config.HigressServiceConfig) {
				if cfg.KubeConfigPath != "/path/to/kubeconfig" {
					t.Errorf("KubeConfigPath mismatch: got %s, want /path/to/kubeconfig", cfg.KubeConfigPath)
				}
			},
		},
		{
			name: "WithControllerNamespace",
			options: []config.Option{
				config.WithControllerNamespace("custom-namespace"),
			},
			validate: func(t *testing.T, cfg *config.HigressServiceConfig) {
				if cfg.ControllerNamespace != "custom-namespace" {
					t.Errorf("ControllerNamespace mismatch: got %s, want custom-namespace", cfg.ControllerNamespace)
				}
			},
		},
		{
			name: "WithServiceListSupportRegistry_true",
			options: []config.Option{
				config.WithServiceListSupportRegistry(true),
			},
			validate: func(t *testing.T, cfg *config.HigressServiceConfig) {
				if !cfg.ServiceListSupportRegistry {
					t.Errorf("ServiceListSupportRegistry should be true")
				}
			},
		},
		{
			name: "WithServiceListSupportRegistry_false",
			options: []config.Option{
				config.WithServiceListSupportRegistry(false),
			},
			validate: func(t *testing.T, cfg *config.HigressServiceConfig) {
				if cfg.ServiceListSupportRegistry {
					t.Errorf("ServiceListSupportRegistry should be false")
				}
			},
		},
		{
			name: "Multiple options",
			options: []config.Option{
				config.WithKubeConfigPath("/path/to/kubeconfig"),
				config.WithControllerNamespace("custom-namespace"),
				config.WithServiceListSupportRegistry(false),
			},
			validate: func(t *testing.T, cfg *config.HigressServiceConfig) {
				if cfg.KubeConfigPath != "/path/to/kubeconfig" {
					t.Errorf("KubeConfigPath mismatch: got %s", cfg.KubeConfigPath)
				}
				if cfg.ControllerNamespace != "custom-namespace" {
					t.Errorf("ControllerNamespace mismatch: got %s", cfg.ControllerNamespace)
				}
				if cfg.ServiceListSupportRegistry {
					t.Errorf("ServiceListSupportRegistry should be false")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.NewHigressServiceConfig(tt.options...)
			tt.validate(t, cfg)
		})
	}
}
