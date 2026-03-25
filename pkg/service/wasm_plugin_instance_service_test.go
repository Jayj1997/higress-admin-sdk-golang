// Package service provides business services for the SDK
package service

import (
	"testing"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model"
	"github.com/stretchr/testify/assert"
)

// TestWasmPluginInstanceService_Interface2 tests that WasmPluginInstanceService interface is defined
func TestWasmPluginInstanceService_Interface2(t *testing.T) {
	// This test verifies the interface exists and has the expected methods
	var _ WasmPluginInstanceService = (WasmPluginInstanceService)(nil)
}

// TestWasmPluginInstanceModel2 tests the WasmPluginInstance model
func TestWasmPluginInstanceModel2(t *testing.T) {
	enabled := true
	instance := model.WasmPluginInstance{
		PluginName: "test-plugin",
		Scope:      model.WasmPluginInstanceScopeGlobal,
		Target:     "global",
		Enabled:    &enabled,
	}

	assert.Equal(t, "test-plugin", instance.PluginName)
	assert.Equal(t, model.WasmPluginInstanceScopeGlobal, instance.Scope)
	assert.Equal(t, "global", instance.Target)
	assert.True(t, *instance.Enabled)
}

// TestWasmPluginInstanceWithConfigurations tests WasmPluginInstance with configurations
func TestWasmPluginInstanceWithConfigurations(t *testing.T) {
	config := map[string]interface{}{
		"key": "value",
		"nested": map[string]interface{}{
			"nestedKey": "nestedValue",
		},
	}

	instance := model.WasmPluginInstance{
		PluginName:     "configured-plugin",
		Scope:          model.WasmPluginInstanceScopeDomain,
		Target:         "example.com",
		Configurations: config,
	}

	assert.NotNil(t, instance.Configurations)
	assert.Equal(t, "value", instance.Configurations["key"])
}

// TestWasmPluginInstanceScope2 tests WasmPluginInstanceScope constants
func TestWasmPluginInstanceScope2(t *testing.T) {
	tests := []struct {
		name  string
		scope model.WasmPluginInstanceScope
	}{
		{"Global scope", model.WasmPluginInstanceScopeGlobal},
		{"Domain scope", model.WasmPluginInstanceScopeDomain},
		{"Service scope", model.WasmPluginInstanceScopeService},
		{"Route scope", model.WasmPluginInstanceScopeRoute},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, string(tt.scope))
		})
	}
}

// TestWasmPluginInstanceScopePriority2 tests scope priority
func TestWasmPluginInstanceScopePriority2(t *testing.T) {
	tests := []struct {
		name             string
		scope            model.WasmPluginInstanceScope
		expectedPriority int
	}{
		{"Global has lowest priority", model.WasmPluginInstanceScopeGlobal, 0},
		{"Domain has second priority", model.WasmPluginInstanceScopeDomain, 10},
		{"Route has third priority", model.WasmPluginInstanceScopeRoute, 100},
		{"Service has highest priority", model.WasmPluginInstanceScopeService, 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedPriority, tt.scope.Priority())
		})
	}
}

// TestWasmPluginInstanceValidate2 tests WasmPluginInstance validation
func TestWasmPluginInstanceValidate2(t *testing.T) {
	tests := []struct {
		name     string
		instance model.WasmPluginInstance
		wantErr  bool
	}{
		{
			name: "Valid instance with scope",
			instance: model.WasmPluginInstance{
				PluginName: "test-plugin",
				Scope:      model.WasmPluginInstanceScopeGlobal,
				Target:     "global",
			},
			wantErr: false,
		},
		{
			name: "Valid instance with targets",
			instance: model.WasmPluginInstance{
				PluginName: "test-plugin",
				Targets: map[model.WasmPluginInstanceScope]string{
					model.WasmPluginInstanceScopeDomain: "example.com",
				},
			},
			wantErr: false,
		},
		{
			name: "Empty plugin name",
			instance: model.WasmPluginInstance{
				PluginName: "",
				Scope:      model.WasmPluginInstanceScopeGlobal,
				Target:     "global",
			},
			wantErr: true,
		},
		{
			name: "Missing scope and targets",
			instance: model.WasmPluginInstance{
				PluginName: "test-plugin",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.instance.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestWasmPluginInstanceSyncDeprecatedFields2 tests SyncDeprecatedFields method
func TestWasmPluginInstanceSyncDeprecatedFields2(t *testing.T) {
	instance := model.WasmPluginInstance{
		PluginName: "test-plugin",
		Scope:      model.WasmPluginInstanceScopeDomain,
		Target:     "example.com",
	}

	// SyncDeprecatedFields should not panic
	instance.SyncDeprecatedFields()

	// After sync, Targets should be populated
	assert.NotNil(t, instance.Targets)
	assert.Equal(t, "example.com", instance.Targets[model.WasmPluginInstanceScopeDomain])
}

// TestWasmPluginInstanceInternal2 tests internal flag
func TestWasmPluginInstanceInternal2(t *testing.T) {
	internalTrue := true
	internalFalse := false

	instanceWithInternal := model.WasmPluginInstance{
		PluginName: "internal-plugin",
		Scope:      model.WasmPluginInstanceScopeGlobal,
		Target:     "global",
		Internal:   &internalTrue,
	}

	instanceWithoutInternal := model.WasmPluginInstance{
		PluginName: "external-plugin",
		Scope:      model.WasmPluginInstanceScopeGlobal,
		Target:     "global",
		Internal:   &internalFalse,
	}

	assert.True(t, *instanceWithInternal.Internal)
	assert.False(t, *instanceWithoutInternal.Internal)
}

// TestWasmPluginInstanceTargets tests Targets field
func TestWasmPluginInstanceTargets(t *testing.T) {
	targets := map[model.WasmPluginInstanceScope]string{
		model.WasmPluginInstanceScopeDomain:  "example.com",
		model.WasmPluginInstanceScopeService: "my-service",
	}

	instance := model.WasmPluginInstance{
		PluginName: "multi-target-plugin",
		Scope:      model.WasmPluginInstanceScopeDomain,
		Target:     "example.com",
		Targets:    targets,
	}

	assert.NotNil(t, instance.Targets)
	assert.Equal(t, "example.com", instance.Targets[model.WasmPluginInstanceScopeDomain])
	assert.Equal(t, "my-service", instance.Targets[model.WasmPluginInstanceScopeService])
}

// TestWasmPluginInstanceSortByScope2 tests sorting instances by scope
func TestWasmPluginInstanceSortByScope2(t *testing.T) {
	instances := []model.WasmPluginInstance{
		{PluginName: "route-plugin", Scope: model.WasmPluginInstanceScopeRoute, Target: "route1"},
		{PluginName: "global-plugin", Scope: model.WasmPluginInstanceScopeGlobal, Target: "global"},
		{PluginName: "domain-plugin", Scope: model.WasmPluginInstanceScopeDomain, Target: "example.com"},
		{PluginName: "service-plugin", Scope: model.WasmPluginInstanceScopeService, Target: "service1"},
	}

	// Sort by scope priority (ascending)
	for i := 0; i < len(instances)-1; i++ {
		for j := i + 1; j < len(instances); j++ {
			if instances[i].Scope.Priority() > instances[j].Scope.Priority() {
				instances[i], instances[j] = instances[j], instances[i]
			}
		}
	}

	// After sorting, global should be first (lowest priority = 0)
	// Order: Global(0) < Domain(10) < Route(100) < Service(1000)
	assert.Equal(t, model.WasmPluginInstanceScopeGlobal, instances[0].Scope)
	assert.Equal(t, model.WasmPluginInstanceScopeDomain, instances[1].Scope)
	assert.Equal(t, model.WasmPluginInstanceScopeRoute, instances[2].Scope)
	assert.Equal(t, model.WasmPluginInstanceScopeService, instances[3].Scope)
}

// TestWasmPluginInstanceID tests ID field
func TestWasmPluginInstanceID(t *testing.T) {
	instance := model.WasmPluginInstance{
		ID:         "instance-123",
		PluginName: "test-plugin",
		Scope:      model.WasmPluginInstanceScopeGlobal,
		Target:     "global",
	}

	assert.Equal(t, "instance-123", instance.ID)
}

// TestWasmPluginInstancePluginVersion tests PluginVersion field
func TestWasmPluginInstancePluginVersion(t *testing.T) {
	instance := model.WasmPluginInstance{
		PluginName:    "test-plugin",
		PluginVersion: "1.0.0",
		Scope:         model.WasmPluginInstanceScopeGlobal,
		Target:        "global",
	}

	assert.Equal(t, "1.0.0", instance.PluginVersion)
}

// TestWasmPluginInstanceHasScopedTarget tests HasScopedTarget method
func TestWasmPluginInstanceHasScopedTarget(t *testing.T) {
	instance := model.WasmPluginInstance{
		PluginName: "test-plugin",
		Scope:      model.WasmPluginInstanceScopeDomain,
		Target:     "example.com",
		Targets: map[model.WasmPluginInstanceScope]string{
			model.WasmPluginInstanceScopeDomain:  "example.com",
			model.WasmPluginInstanceScopeService: "my-service",
		},
	}

	// Test with Targets map
	assert.True(t, instance.HasScopedTarget(model.WasmPluginInstanceScopeDomain, "example.com"))
	assert.True(t, instance.HasScopedTarget(model.WasmPluginInstanceScopeService, "my-service"))
	assert.True(t, instance.HasScopedTarget(model.WasmPluginInstanceScopeDomain, "")) // empty target matches any
	assert.False(t, instance.HasScopedTarget(model.WasmPluginInstanceScopeRoute, "route1"))
}
