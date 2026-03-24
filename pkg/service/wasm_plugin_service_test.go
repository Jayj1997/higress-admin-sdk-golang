// Package service provides business services for the SDK
package service

import (
	"context"
	"testing"

	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/config"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWasmPluginService_Interface tests that the implementation satisfies the interface
func TestWasmPluginService_Interface(t *testing.T) {
	var _ WasmPluginService = (*WasmPluginServiceImpl)(nil)
}

// TestWasmPluginService_New tests creating a new WasmPluginService
func TestWasmPluginService_New(t *testing.T) {
	// Create service with nil kubernetes client (for testing without cluster)
	svc, err := NewWasmPluginService(
		nil, // kubernetesClient
		nil, // modelConverter
		&config.WasmPluginServiceConfig{
			ImageRegistry: "higress-registry.cn-hangzhou.cr.aliyuncs.com/plugins",
		},
	)
	require.NoError(t, err)
	require.NotNil(t, svc)
}

// TestWasmPluginService_List tests listing plugins
func TestWasmPluginService_List(t *testing.T) {
	svc, err := NewWasmPluginService(
		nil,
		nil,
		&config.WasmPluginServiceConfig{},
	)
	require.NoError(t, err)

	ctx := context.Background()
	query := &model.WasmPluginPageQuery{
		CommonPageQuery: model.CommonPageQuery{
			PageNum:  1,
			PageSize: 10,
		},
	}

	result, err := svc.List(ctx, query)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Should have at least some built-in plugins
	assert.Greater(t, result.Total, 0)
	assert.NotEmpty(t, result.Data)
}

// TestWasmPluginService_Get tests getting a specific plugin
func TestWasmPluginService_Get(t *testing.T) {
	svc, err := NewWasmPluginService(
		nil,
		nil,
		&config.WasmPluginServiceConfig{},
	)
	require.NoError(t, err)

	ctx := context.Background()

	// First list to get a valid plugin name
	list, err := svc.List(ctx, &model.WasmPluginPageQuery{
		CommonPageQuery: model.CommonPageQuery{PageNum: 1, PageSize: 1},
	})
	require.NoError(t, err)
	require.NotEmpty(t, list.Data, "Should have at least one plugin")

	pluginName := list.Data[0].Name

	// Get the plugin
	plugin, err := svc.Get(ctx, pluginName, "")
	require.NoError(t, err)
	require.NotNil(t, plugin)

	assert.Equal(t, pluginName, plugin.Name)
	assert.NotEmpty(t, plugin.Title)
	assert.NotEmpty(t, plugin.Version)
}

// TestWasmPluginService_Get_NotFound tests getting a non-existent plugin
func TestWasmPluginService_Get_NotFound(t *testing.T) {
	svc, err := NewWasmPluginService(
		nil,
		nil,
		&config.WasmPluginServiceConfig{},
	)
	require.NoError(t, err)

	ctx := context.Background()

	plugin, err := svc.Get(ctx, "non-existent-plugin-xyz", "")
	// Get returns nil, nil for not found
	assert.NoError(t, err)
	assert.Nil(t, plugin)
}

// TestWasmPluginService_GetConfig tests getting plugin config schema
func TestWasmPluginService_GetConfig(t *testing.T) {
	svc, err := NewWasmPluginService(
		nil,
		nil,
		&config.WasmPluginServiceConfig{},
	)
	require.NoError(t, err)

	ctx := context.Background()

	// First list to get a valid plugin name
	list, err := svc.List(ctx, &model.WasmPluginPageQuery{
		CommonPageQuery: model.CommonPageQuery{PageNum: 1, PageSize: 1},
	})
	require.NoError(t, err)
	require.NotEmpty(t, list.Data, "Should have at least one plugin")

	pluginName := list.Data[0].Name

	cfg, err := svc.GetConfig(ctx, pluginName, "")
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.NotEmpty(t, cfg.Schema)
}

// TestWasmPluginService_GetReadme tests getting plugin readme
func TestWasmPluginService_GetReadme(t *testing.T) {
	svc, err := NewWasmPluginService(
		nil,
		nil,
		&config.WasmPluginServiceConfig{},
	)
	require.NoError(t, err)

	ctx := context.Background()

	// First list to get a valid plugin name
	list, err := svc.List(ctx, &model.WasmPluginPageQuery{
		CommonPageQuery: model.CommonPageQuery{PageNum: 1, PageSize: 1},
	})
	require.NoError(t, err)
	require.NotEmpty(t, list.Data, "Should have at least one plugin")

	pluginName := list.Data[0].Name

	readme, err := svc.GetReadme(ctx, pluginName, "zh-CN")
	require.NoError(t, err)
	// README might be empty for some plugins, so we just check no error
	_ = readme
}

// TestWasmPluginService_GetReadme_English tests getting English readme
func TestWasmPluginService_GetReadme_English(t *testing.T) {
	svc, err := NewWasmPluginService(
		nil,
		nil,
		&config.WasmPluginServiceConfig{},
	)
	require.NoError(t, err)

	ctx := context.Background()

	// First list to get a valid plugin name
	list, err := svc.List(ctx, &model.WasmPluginPageQuery{
		CommonPageQuery: model.CommonPageQuery{PageNum: 1, PageSize: 1},
	})
	require.NoError(t, err)
	require.NotEmpty(t, list.Data, "Should have at least one plugin")

	pluginName := list.Data[0].Name

	readme, err := svc.GetReadme(ctx, pluginName, "en-US")
	require.NoError(t, err)
	// README might be empty for some plugins, so we just check no error
	_ = readme
}

// TestWasmPluginService_List_WithLang tests listing plugins with language filter
func TestWasmPluginService_List_WithLang(t *testing.T) {
	svc, err := NewWasmPluginService(
		nil,
		nil,
		&config.WasmPluginServiceConfig{},
	)
	require.NoError(t, err)

	ctx := context.Background()
	query := &model.WasmPluginPageQuery{
		CommonPageQuery: model.CommonPageQuery{
			PageNum:  1,
			PageSize: 10,
		},
		Lang: "zh-CN",
	}

	result, err := svc.List(ctx, query)
	require.NoError(t, err)
	require.NotNil(t, result)
}

// TestWasmPluginInstanceService_Interface tests that the implementation satisfies the interface
func TestWasmPluginInstanceService_Interface(t *testing.T) {
	var _ WasmPluginInstanceService = (*WasmPluginInstanceServiceImpl)(nil)
}

// TestWasmPluginInstanceService_CreateEmptyInstance tests creating an empty instance
func TestWasmPluginInstanceService_CreateEmptyInstance(t *testing.T) {
	// Create WasmPluginService first
	pluginSvc, err := NewWasmPluginService(
		nil,
		nil,
		&config.WasmPluginServiceConfig{},
	)
	require.NoError(t, err)

	// First list to get a valid plugin name
	ctx := context.Background()
	list, err := pluginSvc.List(ctx, &model.WasmPluginPageQuery{
		CommonPageQuery: model.CommonPageQuery{PageNum: 1, PageSize: 1},
	})
	require.NoError(t, err)
	require.NotEmpty(t, list.Data, "Should have at least one plugin")

	pluginName := list.Data[0].Name

	// Create WasmPluginInstanceService
	instanceSvc := NewWasmPluginInstanceService(
		pluginSvc,
		nil, // kubernetesClient
		nil, // modelConverter
	)

	instance, err := instanceSvc.CreateEmptyInstance(ctx, pluginName)
	require.NoError(t, err)
	require.NotNil(t, instance)

	assert.Equal(t, pluginName, instance.PluginName)
}

// TestWasmPluginInstanceService_CreateEmptyInstance_NotFound tests creating an instance for non-existent plugin
func TestWasmPluginInstanceService_CreateEmptyInstance_NotFound(t *testing.T) {
	pluginSvc, err := NewWasmPluginService(
		nil,
		nil,
		&config.WasmPluginServiceConfig{},
	)
	require.NoError(t, err)

	instanceSvc := NewWasmPluginInstanceService(
		pluginSvc,
		nil,
		nil,
	)

	ctx := context.Background()

	instance, err := instanceSvc.CreateEmptyInstance(ctx, "non-existent-plugin-xyz")
	assert.Error(t, err)
	assert.Nil(t, instance)
}

// TestWasmPluginInstanceService_ListByPlugin tests listing instances by plugin name
func TestWasmPluginInstanceService_ListByPlugin(t *testing.T) {
	pluginSvc, err := NewWasmPluginService(
		nil,
		nil,
		&config.WasmPluginServiceConfig{},
	)
	require.NoError(t, err)

	instanceSvc := NewWasmPluginInstanceService(
		pluginSvc,
		nil, // kubernetesClient - nil means no cluster connection
		nil,
	)

	ctx := context.Background()

	// With nil kubernetes client, should return empty list
	instances, err := instanceSvc.ListByPlugin(ctx, "test-plugin", nil)
	require.NoError(t, err)
	assert.Empty(t, instances)
}

// TestWasmPluginInstanceService_ListByScope tests listing instances by scope
func TestWasmPluginInstanceService_ListByScope(t *testing.T) {
	pluginSvc, err := NewWasmPluginService(
		nil,
		nil,
		&config.WasmPluginServiceConfig{},
	)
	require.NoError(t, err)

	instanceSvc := NewWasmPluginInstanceService(
		pluginSvc,
		nil,
		nil,
	)

	ctx := context.Background()

	// With nil kubernetes client, should return empty list
	instances, err := instanceSvc.ListByScope(ctx, model.WasmPluginInstanceScopeGlobal, "")
	require.NoError(t, err)
	assert.Empty(t, instances)
}

// TestWasmPluginInstanceService_Query tests querying a specific instance
func TestWasmPluginInstanceService_Query(t *testing.T) {
	pluginSvc, err := NewWasmPluginService(
		nil,
		nil,
		&config.WasmPluginServiceConfig{},
	)
	require.NoError(t, err)

	instanceSvc := NewWasmPluginInstanceService(
		pluginSvc,
		nil,
		nil,
	)

	ctx := context.Background()

	// With nil kubernetes client, should return nil
	instance, err := instanceSvc.Query(ctx, model.WasmPluginInstanceScopeGlobal, "", "test-plugin", nil)
	require.NoError(t, err)
	assert.Nil(t, instance)
}

// TestWasmPluginInstance_Validate tests validation of WasmPluginInstance
func TestWasmPluginInstance_Validate(t *testing.T) {
	// Test missing plugin name
	instance := &model.WasmPluginInstance{}
	err := instance.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "pluginName")

	// Test missing scope
	instance = &model.WasmPluginInstance{
		PluginName: "test-plugin",
	}
	err = instance.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "scope")

	// Test valid instance
	instance = &model.WasmPluginInstance{
		PluginName: "test-plugin",
		Scope:      model.WasmPluginInstanceScopeGlobal,
		Target:     "",
	}
	err = instance.Validate()
	assert.NoError(t, err)
}

// TestWasmPluginInstanceScope_Priority tests scope priority
func TestWasmPluginInstanceScope_Priority(t *testing.T) {
	assert.Equal(t, 0, model.WasmPluginInstanceScopeGlobal.Priority())
	assert.Equal(t, 10, model.WasmPluginInstanceScopeDomain.Priority())
	assert.Equal(t, 100, model.WasmPluginInstanceScopeRoute.Priority())
	assert.Equal(t, 1000, model.WasmPluginInstanceScopeService.Priority())
}
