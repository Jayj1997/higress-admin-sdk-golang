// Package model provides data models for the SDK
package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAllowListOperation_Constants(t *testing.T) {
	// Test that operation constants are defined correctly
	assert.Equal(t, AllowListOperation("ADD"), AllowListOperationAdd)
	assert.Equal(t, AllowListOperation("REMOVE"), AllowListOperationRemove)
	assert.Equal(t, AllowListOperation("REPLACE"), AllowListOperationReplace)
	assert.Equal(t, AllowListOperation("TOGGLE_ONLY"), AllowListOperationToggleOnly)
}

func TestNewAllowList(t *testing.T) {
	allowList := NewAllowList()

	require.NotNil(t, allowList)
	assert.NotNil(t, allowList.Targets)
	assert.Empty(t, allowList.Targets)
	assert.NotNil(t, allowList.CredentialTypes)
	assert.Empty(t, allowList.CredentialTypes)
	assert.NotNil(t, allowList.ConsumerNames)
	assert.Empty(t, allowList.ConsumerNames)
	assert.Nil(t, allowList.AuthEnabled)
}

func TestForTarget(t *testing.T) {
	tests := []struct {
		name   string
		scope  WasmPluginInstanceScope
		target string
	}{
		{
			name:   "global scope",
			scope:  WasmPluginInstanceScopeGlobal,
			target: "",
		},
		{
			name:   "domain scope",
			scope:  WasmPluginInstanceScopeDomain,
			target: "example.com",
		},
		{
			name:   "route scope",
			scope:  WasmPluginInstanceScopeRoute,
			target: "my-route",
		},
		{
			name:   "service scope",
			scope:  WasmPluginInstanceScopeService,
			target: "my-service.default.svc.cluster.local",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowList := ForTarget(tt.scope, tt.target)

			require.NotNil(t, allowList)
			require.NotNil(t, allowList.Targets)
			assert.Len(t, allowList.Targets, 1)
			assert.Equal(t, tt.target, allowList.Targets[tt.scope])
		})
	}
}

func TestAllowList_Structure(t *testing.T) {
	// Test with all fields populated
	authEnabled := true
	allowList := &AllowList{
		Targets: map[WasmPluginInstanceScope]string{
			WasmPluginInstanceScopeDomain: "example.com",
			WasmPluginInstanceScopeRoute:  "my-route",
		},
		AuthEnabled:     &authEnabled,
		CredentialTypes: []string{"key-auth", "basic-auth"},
		ConsumerNames:   []string{"consumer1", "consumer2", "consumer3"},
	}

	assert.Len(t, allowList.Targets, 2)
	assert.Equal(t, "example.com", allowList.Targets[WasmPluginInstanceScopeDomain])
	assert.Equal(t, "my-route", allowList.Targets[WasmPluginInstanceScopeRoute])
	assert.True(t, *allowList.AuthEnabled)
	assert.Len(t, allowList.CredentialTypes, 2)
	assert.Len(t, allowList.ConsumerNames, 3)
}

func TestAllowList_AuthEnabled(t *testing.T) {
	t.Run("auth enabled true", func(t *testing.T) {
		authEnabled := true
		allowList := &AllowList{
			AuthEnabled: &authEnabled,
		}
		assert.True(t, *allowList.AuthEnabled)
	})

	t.Run("auth enabled false", func(t *testing.T) {
		authEnabled := false
		allowList := &AllowList{
			AuthEnabled: &authEnabled,
		}
		assert.False(t, *allowList.AuthEnabled)
	})

	t.Run("auth enabled nil", func(t *testing.T) {
		allowList := &AllowList{}
		assert.Nil(t, allowList.AuthEnabled)
	})
}

func TestAllowList_MultipleTargets(t *testing.T) {
	// Test that multiple targets can be set
	allowList := NewAllowList()
	allowList.Targets[WasmPluginInstanceScopeGlobal] = ""
	allowList.Targets[WasmPluginInstanceScopeDomain] = "api.example.com"
	allowList.Targets[WasmPluginInstanceScopeRoute] = "api-route"

	assert.Len(t, allowList.Targets, 3)
	assert.Equal(t, "", allowList.Targets[WasmPluginInstanceScopeGlobal])
	assert.Equal(t, "api.example.com", allowList.Targets[WasmPluginInstanceScopeDomain])
	assert.Equal(t, "api-route", allowList.Targets[WasmPluginInstanceScopeRoute])
}

func TestAllowList_EmptyLists(t *testing.T) {
	// Test that empty lists are valid
	allowList := &AllowList{
		Targets:         map[WasmPluginInstanceScope]string{},
		CredentialTypes: []string{},
		ConsumerNames:   []string{},
	}

	assert.Empty(t, allowList.Targets)
	assert.Empty(t, allowList.CredentialTypes)
	assert.Empty(t, allowList.ConsumerNames)
}

func TestAllowList_NilLists(t *testing.T) {
	// Test that nil lists are handled correctly
	allowList := &AllowList{}

	assert.Nil(t, allowList.Targets)
	assert.Nil(t, allowList.CredentialTypes)
	assert.Nil(t, allowList.ConsumerNames)
}

func TestAllowListOperation_String(t *testing.T) {
	// Test that AllowListOperation can be used as a string
	var op AllowListOperation = "ADD"
	assert.Equal(t, "ADD", string(op))

	op = AllowListOperationRemove
	assert.Equal(t, "REMOVE", string(op))
}
