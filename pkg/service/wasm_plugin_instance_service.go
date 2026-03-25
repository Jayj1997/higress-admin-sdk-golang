// Package service provides business services for the SDK
package service

import (
	"context"
	"strings"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/internal/kubernetes"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/internal/kubernetes/crd/wasm"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/errors"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model"
)

// WasmPluginInstanceScope 定义插件实例的作用域类型
type WasmPluginInstanceScope = model.WasmPluginInstanceScope

// WasmPluginInstanceService WASM插件实例服务接口
type WasmPluginInstanceService interface {
	// CreateEmptyInstance 创建空插件实例
	CreateEmptyInstance(ctx context.Context, pluginName string) (*model.WasmPluginInstance, error)

	// ListByPlugin 按插件名列出实例
	ListByPlugin(ctx context.Context, pluginName string, internal *bool) ([]model.WasmPluginInstance, error)

	// ListByScope 按作用域列出实例
	ListByScope(ctx context.Context, scope model.WasmPluginInstanceScope, target string) ([]model.WasmPluginInstance, error)

	// Query 查询特定插件实例
	Query(ctx context.Context, scope model.WasmPluginInstanceScope, target, pluginName string, internal *bool) (*model.WasmPluginInstance, error)

	// AddOrUpdate 添加或更新实例
	AddOrUpdate(ctx context.Context, instance *model.WasmPluginInstance) (*model.WasmPluginInstance, error)

	// Delete 删除实例
	Delete(ctx context.Context, scope model.WasmPluginInstanceScope, target, pluginName string, internal *bool) error

	// DeleteAll 删除指定作用域和目标的所有插件实例
	DeleteAll(ctx context.Context, scope model.WasmPluginInstanceScope, target string) error
}

// WasmPluginInstanceServiceImpl implements WasmPluginInstanceService
type WasmPluginInstanceServiceImpl struct {
	wasmPluginService WasmPluginService
	kubernetesClient  *kubernetes.KubernetesClientService
	modelConverter    *kubernetes.KubernetesModelConverter
}

// NewWasmPluginInstanceService creates a new WasmPluginInstanceService
func NewWasmPluginInstanceService(
	wasmPluginService WasmPluginService,
	kubernetesClient *kubernetes.KubernetesClientService,
	modelConverter *kubernetes.KubernetesModelConverter,
) WasmPluginInstanceService {
	return &WasmPluginInstanceServiceImpl{
		wasmPluginService: wasmPluginService,
		kubernetesClient:  kubernetesClient,
		modelConverter:    modelConverter,
	}
}

// CreateEmptyInstance 创建空插件实例
func (s *WasmPluginInstanceServiceImpl) CreateEmptyInstance(ctx context.Context, pluginName string) (*model.WasmPluginInstance, error) {
	plugin, err := s.wasmPluginService.Get(ctx, pluginName, "")
	if err != nil {
		return nil, err
	}
	if plugin == nil {
		return nil, errors.NewBusinessError("plugin " + pluginName + " not found")
	}

	instance := &model.WasmPluginInstance{
		PluginName:    plugin.Name,
		PluginVersion: plugin.Version,
	}
	return instance, nil
}

// ListByPlugin 按插件名列出实例
func (s *WasmPluginInstanceServiceImpl) ListByPlugin(ctx context.Context, pluginName string, internal *bool) ([]model.WasmPluginInstance, error) {
	if s.kubernetesClient == nil {
		return []model.WasmPluginInstance{}, nil
	}

	plugins, err := s.kubernetesClient.ListWasmPlugins(ctx, pluginName, "", nil)
	if err != nil {
		return nil, err
	}

	var result []model.WasmPluginInstance
	for _, plugin := range plugins {
		// Filter by internal if specified
		if internal != nil {
			isInternal := isInternalWasmPlugin(plugin)
			if *internal != isInternal {
				continue
			}
		}

		instances, err := s.modelConverter.GetWasmPluginInstancesFromCR(plugin)
		if err != nil {
			continue
		}
		result = append(result, instances...)
	}

	return result, nil
}

// ListByScope 按作用域列出实例
func (s *WasmPluginInstanceServiceImpl) ListByScope(ctx context.Context, scope model.WasmPluginInstanceScope, target string) ([]model.WasmPluginInstance, error) {
	if s.kubernetesClient == nil {
		return []model.WasmPluginInstance{}, nil
	}

	plugins, err := s.kubernetesClient.ListWasmPlugins(ctx, "", "", nil)
	if err != nil {
		return nil, err
	}

	var result []model.WasmPluginInstance
	for _, plugin := range plugins {
		instances, err := s.modelConverter.GetWasmPluginInstancesFromCR(plugin)
		if err != nil {
			continue
		}
		// Filter by scope and target
		for _, instance := range instances {
			if instance.Scope == scope && instance.Target == target {
				result = append(result, instance)
			}
			if instance.Targets != nil {
				if t, ok := instance.Targets[scope]; ok && t == target {
					result = append(result, instance)
				}
			}
		}
	}

	return result, nil
}

// Query 查询特定插件实例
func (s *WasmPluginInstanceServiceImpl) Query(ctx context.Context, scope model.WasmPluginInstanceScope, target, pluginName string, internal *bool) (*model.WasmPluginInstance, error) {
	if s.kubernetesClient == nil {
		return nil, nil
	}

	plugins, err := s.kubernetesClient.ListWasmPlugins(ctx, pluginName, "", nil)
	if err != nil {
		return nil, err
	}

	for _, plugin := range plugins {
		// Filter by internal if specified
		if internal != nil {
			isInternal := isInternalWasmPlugin(plugin)
			if *internal != isInternal {
				continue
			}
		}

		instances, err := s.modelConverter.GetWasmPluginInstancesFromCR(plugin)
		if err != nil {
			continue
		}
		// Find matching instance
		for _, instance := range instances {
			if instance.Scope == scope && instance.Target == target {
				return &instance, nil
			}
			if instance.Targets != nil {
				if t, ok := instance.Targets[scope]; ok && t == target {
					return &instance, nil
				}
			}
		}
	}

	return nil, nil
}

// AddOrUpdate 添加或更新实例
func (s *WasmPluginInstanceServiceImpl) AddOrUpdate(ctx context.Context, instance *model.WasmPluginInstance) (*model.WasmPluginInstance, error) {
	if s.kubernetesClient == nil {
		return nil, errors.NewBusinessError("kubernetes client is not available")
	}

	// Validate instance
	if err := instance.Validate(); err != nil {
		return nil, err
	}

	// Sync deprecated fields
	instance.SyncDeprecatedFields()

	// Get plugin info
	plugin, err := s.wasmPluginService.Get(ctx, instance.PluginName, "")
	if err != nil {
		return nil, err
	}
	if plugin == nil {
		return nil, errors.NewBusinessError("unknown plugin: " + instance.PluginName)
	}

	// Determine version
	version := instance.PluginVersion
	if version == "" {
		version = plugin.Version
	}

	// Find existing CR
	existingCRs, err := s.kubernetesClient.ListWasmPlugins(ctx, instance.PluginName, version, nil)
	if err != nil {
		return nil, err
	}

	var existingCR *wasm.V1alpha1WasmPlugin
	internal := false
	if instance.Internal != nil {
		internal = *instance.Internal
	}

	for _, cr := range existingCRs {
		isInternal := isInternalWasmPlugin(cr)
		if internal == isInternal {
			existingCR = cr
			break
		}
	}

	if existingCR != nil {
		// Update existing CR
		err = s.modelConverter.SetWasmPluginInstanceToCR(existingCR, instance)
		if err != nil {
			return nil, err
		}
		updated, err := s.kubernetesClient.UpdateWasmPlugin(ctx, existingCR)
		if err != nil {
			return nil, err
		}
		instances, err := s.modelConverter.GetWasmPluginInstancesFromCR(updated)
		if err != nil {
			return nil, err
		}
		if len(instances) > 0 {
			return &instances[0], nil
		}
		return instance, nil
	}

	// Create new CR
	newCR := wasm.NewV1alpha1WasmPlugin()
	newCR.Metadata.Name = instance.PluginName
	if version != "" {
		newCR.Metadata.Name = instance.PluginName + "-" + version
	}

	// Add .internal suffix for internal resources
	if internal {
		newCR.Metadata.Name += ".internal"
	}

	// Set spec from plugin
	if plugin.ImageRepository != "" {
		newCR.Spec.Url = plugin.ImageRepository
		if plugin.ImageVersion != "" {
			newCR.Spec.Url += ":" + plugin.ImageVersion
		}
	}

	err = s.modelConverter.SetWasmPluginInstanceToCR(newCR, instance)
	if err != nil {
		return nil, err
	}

	created, err := s.kubernetesClient.CreateWasmPlugin(ctx, newCR)
	if err != nil {
		return nil, err
	}
	instances, err := s.modelConverter.GetWasmPluginInstancesFromCR(created)
	if err != nil {
		return nil, err
	}
	if len(instances) > 0 {
		return &instances[0], nil
	}
	return instance, nil
}

// Delete 删除实例
func (s *WasmPluginInstanceServiceImpl) Delete(ctx context.Context, scope model.WasmPluginInstanceScope, target, pluginName string, internal *bool) error {
	if s.kubernetesClient == nil {
		return errors.NewBusinessError("kubernetes client is not available")
	}

	plugins, err := s.kubernetesClient.ListWasmPlugins(ctx, pluginName, "", nil)
	if err != nil {
		return err
	}

	for _, plugin := range plugins {
		// Filter by internal if specified
		if internal != nil {
			isInternal := isInternalWasmPlugin(plugin)
			if *internal != isInternal {
				continue
			}
		}

		// Check if this CR has the instance
		instances, err := s.modelConverter.GetWasmPluginInstancesFromCR(plugin)
		if err != nil {
			continue
		}
		for _, inst := range instances {
			if inst.Scope == scope && inst.Target == target {
				// Found matching instance, remove it
				// For simplicity, delete the entire CR if it only has this instance
				if len(instances) == 1 {
					if err := s.kubernetesClient.DeleteWasmPlugin(ctx, plugin.Metadata.Name); err != nil {
						return err
					}
				}
				break
			}
		}
	}

	return nil
}

// DeleteAll 删除指定作用域和目标的所有插件实例
func (s *WasmPluginInstanceServiceImpl) DeleteAll(ctx context.Context, scope model.WasmPluginInstanceScope, target string) error {
	if s.kubernetesClient == nil {
		return errors.NewBusinessError("kubernetes client is not available")
	}

	plugins, err := s.kubernetesClient.ListWasmPlugins(ctx, "", "", nil)
	if err != nil {
		return err
	}

	for _, plugin := range plugins {
		// Check if this CR has instances for the scope/target
		instances, err := s.modelConverter.GetWasmPluginInstancesFromCR(plugin)
		if err != nil {
			continue
		}
		for _, inst := range instances {
			if inst.Scope == scope && inst.Target == target {
				// Found matching instance, remove it
				// For simplicity, delete the entire CR if it only has this instance
				if len(instances) == 1 {
					if err := s.kubernetesClient.DeleteWasmPlugin(ctx, plugin.Metadata.Name); err != nil {
						return err
					}
				}
				break
			}
		}
	}

	return nil
}

// isInternalWasmPlugin checks if a WasmPlugin is an internal resource
func isInternalWasmPlugin(plugin *wasm.V1alpha1WasmPlugin) bool {
	if plugin == nil || plugin.Metadata == nil {
		return false
	}
	// Check if name ends with ".internal" suffix (Java SDK behavior)
	return strings.HasSuffix(plugin.Metadata.Name, ".internal")
}

// MockWasmPluginInstanceService 用于测试的Mock实现
type MockWasmPluginInstanceService struct{}

// NewMockWasmPluginInstanceService 创建Mock服务实例
func NewMockWasmPluginInstanceService() WasmPluginInstanceService {
	return &MockWasmPluginInstanceService{}
}

func (s *MockWasmPluginInstanceService) CreateEmptyInstance(ctx context.Context, pluginName string) (*model.WasmPluginInstance, error) {
	return &model.WasmPluginInstance{PluginName: pluginName}, nil
}

func (s *MockWasmPluginInstanceService) ListByPlugin(ctx context.Context, pluginName string, internal *bool) ([]model.WasmPluginInstance, error) {
	return []model.WasmPluginInstance{}, nil
}

func (s *MockWasmPluginInstanceService) ListByScope(ctx context.Context, scope model.WasmPluginInstanceScope, target string) ([]model.WasmPluginInstance, error) {
	return []model.WasmPluginInstance{}, nil
}

func (s *MockWasmPluginInstanceService) Query(ctx context.Context, scope model.WasmPluginInstanceScope, target, pluginName string, internal *bool) (*model.WasmPluginInstance, error) {
	return nil, nil
}

func (s *MockWasmPluginInstanceService) AddOrUpdate(ctx context.Context, instance *model.WasmPluginInstance) (*model.WasmPluginInstance, error) {
	return instance, nil
}

func (s *MockWasmPluginInstanceService) Delete(ctx context.Context, scope model.WasmPluginInstanceScope, target, pluginName string, internal *bool) error {
	return nil
}

func (s *MockWasmPluginInstanceService) DeleteAll(ctx context.Context, scope model.WasmPluginInstanceScope, target string) error {
	return nil
}
