package service

import (
	"context"
	"sort"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/errors"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/service/consumer"
)

// 默认凭证类型列表
var defaultCredentialTypes = []string{model.CredentialTypeKeyAuth}

// ConsumerServiceImpl 消费者服务实现
type ConsumerServiceImpl struct {
	wasmPluginInstanceService WasmPluginInstanceService
	credentialHandlers        map[string]consumer.CredentialHandler
}

// NewConsumerServiceImpl 创建消费者服务实例
func NewConsumerServiceImpl(wasmPluginInstanceService WasmPluginInstanceService) *ConsumerServiceImpl {
	handlers := map[string]consumer.CredentialHandler{
		model.CredentialTypeKeyAuth: consumer.NewKeyAuthCredentialHandler(),
	}
	return &ConsumerServiceImpl{
		wasmPluginInstanceService: wasmPluginInstanceService,
		credentialHandlers:        handlers,
	}
}

// List 列出所有消费者
func (s *ConsumerServiceImpl) List(ctx context.Context) ([]model.Consumer, error) {
	consumers := s.getConsumers(ctx)
	result := make([]model.Consumer, 0, len(consumers))
	for _, c := range consumers {
		result = append(result, *c)
	}
	return result, nil
}

// Get 获取消费者详情
func (s *ConsumerServiceImpl) Get(ctx context.Context, name string) (*model.Consumer, error) {
	if name == "" {
		return nil, errors.NewValidationError("consumerName cannot be empty")
	}
	consumers := s.getConsumers(ctx)
	if c, ok := consumers[name]; ok {
		return c, nil
	}
	return nil, nil
}

// AddOrUpdate 添加或更新消费者
func (s *ConsumerServiceImpl) AddOrUpdate(ctx context.Context, consumer *model.Consumer) (*model.Consumer, error) {
	instancesToUpdate := []*model.WasmPluginInstance{}

	for _, handler := range s.credentialHandlers {
		instance := s.getGlobalPluginInstance(ctx, handler)
		if instance == nil {
			instance = s.createGlobalPluginInstance(ctx, handler)
		}
		instance.Enabled = boolPtr(true)
		if handler.SaveConsumer(instance, consumer) {
			instancesToUpdate = append(instancesToUpdate, instance)
		}
	}

	for _, instance := range instancesToUpdate {
		if _, err := s.wasmPluginInstanceService.AddOrUpdate(ctx, instance); err != nil {
			return nil, err
		}
	}

	return s.Get(ctx, consumer.Name)
}

// Delete 删除消费者
func (s *ConsumerServiceImpl) Delete(ctx context.Context, name string) error {
	if name == "" {
		return errors.NewValidationError("consumerName cannot be empty")
	}

	instancesCache := make(map[string][]*model.WasmPluginInstance)

	// 检查消费者是否正在使用
	for _, handler := range s.credentialHandlers {
		instances := s.getAllPluginInstances(ctx, handler)
		if handler.IsConsumerInUse(name, instances) {
			return errors.NewBusinessError("Consumer " + name + " is still in use")
		}
		instancesCache[handler.GetType()] = instances
	}

	// 删除消费者
	for _, handler := range s.credentialHandlers {
		instances := instancesCache[handler.GetType()]
		var globalInstance *model.WasmPluginInstance
		for _, inst := range instances {
			if inst.HasScopedTarget(model.WasmPluginInstanceScopeGlobal, "") {
				globalInstance = inst
				break
			}
		}
		if globalInstance == nil {
			continue
		}
		if handler.DeleteConsumer(globalInstance, name) {
			if _, err := s.wasmPluginInstanceService.AddOrUpdate(ctx, globalInstance); err != nil {
				return err
			}
		}
	}

	return nil
}

// ListAllowLists 列出所有允许列表
func (s *ConsumerServiceImpl) ListAllowLists(ctx context.Context) ([]model.AllowList, error) {
	allowLists := []*model.AllowList{}

	for _, handler := range s.credentialHandlers {
		instances := s.getAllPluginInstances(ctx, handler)
		if len(instances) == 0 {
			continue
		}

		for _, instance := range instances {
			if instance.HasScopedTarget(model.WasmPluginInstanceScopeGlobal, "") {
				continue
			}

			// 查找或创建AllowList
			var allowList *model.AllowList
			for _, al := range allowLists {
				if mapsEqual(al.Targets, instance.Targets) {
					allowList = al
					break
				}
			}
			if allowList == nil {
				allowList = &model.AllowList{
					Targets:         instance.Targets,
					AuthEnabled:     instance.Enabled,
					CredentialTypes: []string{},
					ConsumerNames:   []string{},
				}
				allowLists = append(allowLists, allowList)
			}

			consumerNames := handler.GetAllowedConsumers(instance)
			if !contains(allowList.CredentialTypes, handler.GetType()) {
				allowList.CredentialTypes = append(allowList.CredentialTypes, handler.GetType())
			}
			for _, name := range consumerNames {
				if !contains(allowList.ConsumerNames, name) {
					allowList.ConsumerNames = append(allowList.ConsumerNames, name)
				}
			}
		}
	}

	result := make([]model.AllowList, len(allowLists))
	for i, al := range allowLists {
		result[i] = *al
	}
	return result, nil
}

// GetAllowList 获取指定目标的允许列表
func (s *ConsumerServiceImpl) GetAllowList(ctx context.Context, targets map[model.WasmPluginInstanceScope]string) (*model.AllowList, error) {
	if len(targets) == 0 {
		return nil, errors.NewValidationError("targets cannot be null or empty")
	}
	if _, ok := targets[model.WasmPluginInstanceScopeGlobal]; ok {
		return nil, errors.NewValidationError("targets cannot contain GLOBAL scope")
	}

	var credentialTypes []string
	var allConsumerNames []string
	var authEnabled bool

	for _, handler := range s.credentialHandlers {
		instance := s.getPluginInstance(ctx, handler, targets)
		if instance == nil {
			continue
		}
		consumerNames := handler.GetAllowedConsumers(instance)
		if instance.Enabled != nil && *instance.Enabled {
			authEnabled = true
		}
		credentialTypes = append(credentialTypes, handler.GetType())
		for _, name := range consumerNames {
			if !contains(allConsumerNames, name) {
				allConsumerNames = append(allConsumerNames, name)
			}
		}
	}

	if len(allConsumerNames) == 0 {
		return nil, nil
	}

	return &model.AllowList{
		Targets:         targets,
		AuthEnabled:     &authEnabled,
		CredentialTypes: credentialTypes,
		ConsumerNames:   allConsumerNames,
	}, nil
}

// UpdateAllowList 更新允许列表
func (s *ConsumerServiceImpl) UpdateAllowList(ctx context.Context, operation model.AllowListOperation, allowList *model.AllowList) error {
	if operation == "" {
		return errors.NewValidationError("operation cannot be null")
	}
	if allowList == nil {
		return errors.NewValidationError("allowList cannot be null")
	}

	targets := allowList.Targets
	consumerNames := allowList.ConsumerNames

	if len(targets) == 0 {
		return errors.NewValidationError("targets cannot be null or empty")
	}
	if _, ok := targets[model.WasmPluginInstanceScopeGlobal]; ok {
		return errors.NewValidationError("targets cannot contain GLOBAL scope")
	}

	credentialTypes := allowList.CredentialTypes
	if len(credentialTypes) == 0 {
		credentialTypes = defaultCredentialTypes
	} else {
		credentialTypes = unique(credentialTypes)
	}

	switch operation {
	case model.AllowListOperationAdd, model.AllowListOperationRemove:
		if len(consumerNames) == 0 && allowList.AuthEnabled == nil {
			return nil
		}
	case model.AllowListOperationToggleOnly:
		if allowList.AuthEnabled == nil {
			return nil
		}
	case model.AllowListOperationReplace:
		// 继续执行
	default:
		return errors.NewValidationError("Unsupported operation: " + string(operation))
	}

	for _, credType := range credentialTypes {
		handler, ok := s.credentialHandlers[credType]
		if !ok {
			return errors.NewValidationError("Unsupported credential type: " + credType)
		}

		instancesToSave := []*model.WasmPluginInstance{}
		instances := s.getAllPluginInstances(ctx, handler)

		// 查找目标实例
		var targetInstance *model.WasmPluginInstance
		for _, inst := range instances {
			if mapsEqual(inst.Targets, targets) {
				targetInstance = inst
				break
			}
		}

		if targetInstance == nil {
			var err error
			targetInstance, err = s.wasmPluginInstanceService.CreateEmptyInstance(ctx, handler.GetPluginName())
			if err != nil {
				return err
			}
			targetInstance.Internal = boolPtr(true)
			targetInstance.Enabled = boolPtr(false)
			targetInstance.Targets = targets
		}

		if allowList.AuthEnabled != nil {
			targetInstance.Enabled = allowList.AuthEnabled
		}

		handler.UpdateAllowList(operation, targetInstance, consumerNames)
		instancesToSave = append(instancesToSave, targetInstance)

		// 查找全局实例
		var globalInstance *model.WasmPluginInstance
		for _, inst := range instances {
			if inst.HasScopedTarget(model.WasmPluginInstanceScopeGlobal, "") {
				globalInstance = inst
				break
			}
		}

		if globalInstance == nil && targetInstance.Enabled != nil && *targetInstance.Enabled {
			globalInstance = s.createGlobalPluginInstance(ctx, handler)
			instancesToSave = append(instancesToSave, globalInstance)
		}

		for _, inst := range instancesToSave {
			if _, err := s.wasmPluginInstanceService.AddOrUpdate(ctx, inst); err != nil {
				return err
			}
		}
	}

	return nil
}

// getConsumers 获取所有消费者
func (s *ConsumerServiceImpl) getConsumers(ctx context.Context) map[string]*model.Consumer {
	consumers := make(map[string]*model.Consumer)

	for _, handler := range s.credentialHandlers {
		instance := s.getGlobalPluginInstance(ctx, handler)
		if instance == nil {
			continue
		}

		extractedConsumers := handler.ExtractConsumers(instance)
		for _, c := range extractedConsumers {
			if existing, ok := consumers[c.Name]; ok {
				// 合并凭证
				existing.Credentials = append(existing.Credentials, c.Credentials...)
			} else {
				consumers[c.Name] = c
			}
		}
	}

	return consumers
}

// getGlobalPluginInstance 获取全局插件实例
func (s *ConsumerServiceImpl) getGlobalPluginInstance(ctx context.Context, handler consumer.CredentialHandler) *model.WasmPluginInstance {
	instance, err := s.wasmPluginInstanceService.Query(ctx, model.WasmPluginInstanceScopeGlobal, "", handler.GetPluginName(), boolPtr(true))
	if err != nil {
		return nil
	}
	return instance
}

// getAllPluginInstances 获取所有插件实例
func (s *ConsumerServiceImpl) getAllPluginInstances(ctx context.Context, handler consumer.CredentialHandler) []*model.WasmPluginInstance {
	instances, err := s.wasmPluginInstanceService.ListByPlugin(ctx, handler.GetPluginName(), boolPtr(true))
	if err != nil {
		return []*model.WasmPluginInstance{}
	}
	// 转换为指针切片
	result := make([]*model.WasmPluginInstance, len(instances))
	for i := range instances {
		result[i] = &instances[i]
	}
	return result
}

// getPluginInstance 获取指定目标的插件实例
func (s *ConsumerServiceImpl) getPluginInstance(ctx context.Context, handler consumer.CredentialHandler, targets map[model.WasmPluginInstanceScope]string) *model.WasmPluginInstance {
	// 从targets获取scope和target
	var scope model.WasmPluginInstanceScope
	var target string
	for s, t := range targets {
		scope = s
		target = t
		break
	}

	instance, err := s.wasmPluginInstanceService.Query(ctx, scope, target, handler.GetPluginName(), boolPtr(true))
	if err != nil {
		return nil
	}
	return instance
}

// createGlobalPluginInstance 创建全局插件实例
func (s *ConsumerServiceImpl) createGlobalPluginInstance(ctx context.Context, handler consumer.CredentialHandler) *model.WasmPluginInstance {
	instance, err := s.wasmPluginInstanceService.CreateEmptyInstance(ctx, handler.GetPluginName())
	if err != nil {
		return nil
	}
	instance.Internal = boolPtr(true)
	instance.Targets = map[model.WasmPluginInstanceScope]string{
		model.WasmPluginInstanceScopeGlobal: "",
	}
	handler.InitDefaultGlobalConfigs(instance)
	return instance
}

// 辅助函数

func boolPtr(b bool) *bool {
	return &b
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func unique(slice []string) []string {
	keys := make(map[string]bool)
	result := []string{}
	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}
	sort.Strings(result)
	return result
}

func mapsEqual(m1, m2 map[model.WasmPluginInstanceScope]string) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v := range m1 {
		if v2, ok := m2[k]; !ok || v != v2 {
			return false
		}
	}
	return true
}
