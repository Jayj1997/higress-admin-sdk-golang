// Package model provides data models for the SDK
package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWasmPluginInstanceScope_Priority(t *testing.T) {
	tests := []struct {
		scope            WasmPluginInstanceScope
		expectedPriority int
	}{
		{WasmPluginInstanceScopeGlobal, 0},
		{WasmPluginInstanceScopeDomain, 10},
		{WasmPluginInstanceScopeRoute, 100},
		{WasmPluginInstanceScopeService, 1000},
		{WasmPluginInstanceScope("unknown"), 0},
	}

	for _, tt := range tests {
		t.Run(string(tt.scope), func(t *testing.T) {
			assert.Equal(t, tt.expectedPriority, tt.scope.Priority())
		})
	}
}

func TestWasmPluginInstanceScope_Constants(t *testing.T) {
	assert.Equal(t, WasmPluginInstanceScope("global"), WasmPluginInstanceScopeGlobal)
	assert.Equal(t, WasmPluginInstanceScope("domain"), WasmPluginInstanceScopeDomain)
	assert.Equal(t, WasmPluginInstanceScope("route"), WasmPluginInstanceScopeRoute)
	assert.Equal(t, WasmPluginInstanceScope("service"), WasmPluginInstanceScopeService)
}

func TestWasmPluginInstance_Validate(t *testing.T) {
	tests := []struct {
		name        string
		instance    *WasmPluginInstance
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid instance with scope",
			instance: &WasmPluginInstance{
				PluginName: "basic-auth",
				Scope:      WasmPluginInstanceScopeGlobal,
			},
			expectError: false,
		},
		{
			name: "valid instance with targets",
			instance: &WasmPluginInstance{
				PluginName: "basic-auth",
				Targets: map[WasmPluginInstanceScope]string{
					WasmPluginInstanceScopeDomain: "example.com",
				},
			},
			expectError: false,
		},
		{
			name: "missing plugin name",
			instance: &WasmPluginInstance{
				Scope: WasmPluginInstanceScopeGlobal,
			},
			expectError: true,
			errorMsg:    "plugin name is required",
		},
		{
			name: "missing scope and targets",
			instance: &WasmPluginInstance{
				PluginName: "basic-auth",
			},
			expectError: true,
			errorMsg:    "scope or targets is required",
		},
		{
			name: "empty plugin name",
			instance: &WasmPluginInstance{
				PluginName: "",
				Scope:      WasmPluginInstanceScopeGlobal,
			},
			expectError: true,
			errorMsg:    "plugin name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.instance.Validate()
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestWasmPluginInstance_HasScopedTarget(t *testing.T) {
	t.Run("with targets map", func(t *testing.T) {
		instance := &WasmPluginInstance{
			PluginName: "test-plugin",
			Targets: map[WasmPluginInstanceScope]string{
				WasmPluginInstanceScopeDomain: "example.com",
				WasmPluginInstanceScopeRoute:  "my-route",
			},
		}

		assert.True(t, instance.HasScopedTarget(WasmPluginInstanceScopeDomain, "example.com"))
		assert.True(t, instance.HasScopedTarget(WasmPluginInstanceScopeDomain, ""))
		assert.False(t, instance.HasScopedTarget(WasmPluginInstanceScopeDomain, "other.com"))
		assert.True(t, instance.HasScopedTarget(WasmPluginInstanceScopeRoute, "my-route"))
		assert.False(t, instance.HasScopedTarget(WasmPluginInstanceScopeService, ""))
	})

	t.Run("with deprecated scope/target", func(t *testing.T) {
		instance := &WasmPluginInstance{
			PluginName: "test-plugin",
			Scope:      WasmPluginInstanceScopeDomain,
			Target:     "example.com",
		}

		assert.True(t, instance.HasScopedTarget(WasmPluginInstanceScopeDomain, "example.com"))
		assert.True(t, instance.HasScopedTarget(WasmPluginInstanceScopeDomain, ""))
		assert.False(t, instance.HasScopedTarget(WasmPluginInstanceScopeRoute, ""))
	})

	t.Run("with both targets and deprecated fields", func(t *testing.T) {
		instance := &WasmPluginInstance{
			PluginName: "test-plugin",
			Scope:      WasmPluginInstanceScopeDomain,
			Target:     "old-domain.com",
			Targets: map[WasmPluginInstanceScope]string{
				WasmPluginInstanceScopeDomain: "new-domain.com",
			},
		}

		// Targets map takes precedence
		assert.True(t, instance.HasScopedTarget(WasmPluginInstanceScopeDomain, "new-domain.com"))
	})
}

func TestWasmPluginInstance_SyncDeprecatedFields(t *testing.T) {
	t.Run("sync scope/target to targets", func(t *testing.T) {
		instance := &WasmPluginInstance{
			PluginName: "test-plugin",
			Scope:      WasmPluginInstanceScopeDomain,
			Target:     "example.com",
		}

		instance.SyncDeprecatedFields()

		require.NotNil(t, instance.Targets)
		assert.Equal(t, "example.com", instance.Targets[WasmPluginInstanceScopeDomain])
	})

	t.Run("sync targets to scope/target", func(t *testing.T) {
		instance := &WasmPluginInstance{
			PluginName: "test-plugin",
			Targets: map[WasmPluginInstanceScope]string{
				WasmPluginInstanceScopeRoute: "my-route",
			},
		}

		instance.SyncDeprecatedFields()

		assert.Equal(t, WasmPluginInstanceScopeRoute, instance.Scope)
		assert.Equal(t, "my-route", instance.Target)
	})

	t.Run("multiple targets does not sync to scope/target", func(t *testing.T) {
		instance := &WasmPluginInstance{
			PluginName: "test-plugin",
			Targets: map[WasmPluginInstanceScope]string{
				WasmPluginInstanceScopeDomain: "example.com",
				WasmPluginInstanceScopeRoute:  "my-route",
			},
		}

		instance.SyncDeprecatedFields()

		// Should not sync when multiple targets
		assert.Empty(t, instance.Scope)
		assert.Empty(t, instance.Target)
	})
}

func TestValidationError_Error(t *testing.T) {
	t.Run("with field", func(t *testing.T) {
		err := &ValidationError{
			Field:   "name",
			Message: "is required",
		}
		assert.Equal(t, "name: is required", err.Error())
	})

	t.Run("without field", func(t *testing.T) {
		err := &ValidationError{
			Message: "validation failed",
		}
		assert.Equal(t, "validation failed", err.Error())
	})

	t.Run("empty field", func(t *testing.T) {
		err := &ValidationError{
			Field:   "",
			Message: "something went wrong",
		}
		assert.Equal(t, "something went wrong", err.Error())
	})
}

func TestWasmPlugin_Structure(t *testing.T) {
	builtIn := true
	internal := false
	priority := 100

	plugin := &WasmPlugin{
		Name:              "basic-auth",
		Version:           "1.0.0",
		PluginVersion:     "1.0.0",
		Category:          "security",
		Title:             "Basic Auth",
		Description:       "Basic authentication plugin",
		ImageURL:          "oci://registry.example.com/plugins/basic-auth:1.0.0",
		ImageRepository:   "registry.example.com/plugins/basic-auth",
		ImageVersion:      "1.0.0",
		Icon:              "data:image/svg+xml;base64,...",
		BuiltIn:           &builtIn,
		Internal:          &internal,
		Phase:             "AUTHN",
		Priority:          &priority,
		ConfigSchema:      map[string]interface{}{"type": "object"},
		RouteConfigSchema: map[string]interface{}{"type": "object"},
		ImagePullPolicy:   "IfNotPresent",
		ImagePullSecret:   "registry-secret",
		Lang:              "zh-CN",
	}

	assert.Equal(t, "basic-auth", plugin.Name)
	assert.Equal(t, "1.0.0", plugin.Version)
	assert.Equal(t, "security", plugin.Category)
	assert.Equal(t, "Basic Auth", plugin.Title)
	require.NotNil(t, plugin.BuiltIn)
	assert.True(t, *plugin.BuiltIn)
	require.NotNil(t, plugin.Priority)
	assert.Equal(t, 100, *plugin.Priority)
}

func TestWasmPluginInstance_Structure(t *testing.T) {
	enabled := true
	internal := false

	instance := &WasmPluginInstance{
		ID:            "instance-123",
		PluginName:    "basic-auth",
		PluginVersion: "1.0.0",
		Scope:         WasmPluginInstanceScopeDomain,
		Target:        "example.com",
		Targets: map[WasmPluginInstanceScope]string{
			WasmPluginInstanceScopeDomain: "example.com",
		},
		Enabled:  &enabled,
		Internal: &internal,
		Configurations: map[string]interface{}{
			"users": []string{"admin", "user1"},
		},
	}

	assert.Equal(t, "instance-123", instance.ID)
	assert.Equal(t, "basic-auth", instance.PluginName)
	assert.Equal(t, "1.0.0", instance.PluginVersion)
	assert.Equal(t, WasmPluginInstanceScopeDomain, instance.Scope)
	assert.Equal(t, "example.com", instance.Target)
	require.NotNil(t, instance.Enabled)
	assert.True(t, *instance.Enabled)
}

func TestWasmPluginConfig_Structure(t *testing.T) {
	config := &WasmPluginConfig{
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"username": map[string]interface{}{"type": "string"},
				"password": map[string]interface{}{"type": "string"},
			},
		},
	}

	require.NotNil(t, config.Schema)
	assert.Equal(t, "object", config.Schema["type"])
}
