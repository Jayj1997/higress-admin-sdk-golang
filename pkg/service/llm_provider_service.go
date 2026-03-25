// Package service provides business services for the SDK
package service

import (
	"context"
	"sort"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/constant"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/errors"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model/route"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/service/ai"
)

// LlmProviderServiceImpl LLM提供商服务实现
type LlmProviderServiceImpl struct {
	serviceSourceService      ServiceSourceService
	wasmPluginInstanceService WasmPluginInstanceService
	aiRouteService            AiRouteService
}

// NewLlmProviderService 创建LLM提供商服务
func NewLlmProviderService(
	serviceSourceService ServiceSourceService,
	wasmPluginInstanceService WasmPluginInstanceService,
) *LlmProviderServiceImpl {
	return &LlmProviderServiceImpl{
		serviceSourceService:      serviceSourceService,
		wasmPluginInstanceService: wasmPluginInstanceService,
	}
}

// SetAiRouteService 设置AI路由服务
func (s *LlmProviderServiceImpl) SetAiRouteService(service AiRouteService) {
	s.aiRouteService = service
}

// AddOrUpdate 添加或更新LLM提供商
func (s *LlmProviderServiceImpl) AddOrUpdate(ctx context.Context, provider *model.LlmProvider) (*model.LlmProvider, error) {
	// 验证提供商配置
	if err := provider.Validate(false); err != nil {
		return nil, err
	}

	// 获取处理器
	handler := ai.GetHandler(provider.Type)
	if handler == nil {
		return nil, errors.NewValidationError("Provider type " + provider.Type + " is not supported")
	}

	// 规范化配置
	if provider.RawConfigs != nil {
		handler.NormalizeConfigs(provider.RawConfigs)
	}

	// 获取现有的插件实例列表
	internalTrue := true
	instances, err := s.wasmPluginInstanceService.ListByPlugin(ctx, constant.BuiltInPluginAiProxy, &internalTrue)
	if err != nil {
		return nil, err
	}

	// 查找全局实例
	var instance *model.WasmPluginInstance
	for i := range instances {
		if instances[i].HasScopedTarget(model.WasmPluginInstanceScopeGlobal, "") {
			instance = &instances[i]
			break
		}
	}

	// 如果不存在全局实例，创建一个新的
	if instance == nil {
		instance, err = s.wasmPluginInstanceService.CreateEmptyInstance(ctx, constant.BuiltInPluginAiProxy)
		if err != nil {
			return nil, err
		}
		instance.Internal = &internalTrue
		instance.Scope = model.WasmPluginInstanceScopeGlobal
	}
	instance.Enabled = &internalTrue

	// 获取或初始化配置
	configurations := instance.Configurations
	if configurations == nil {
		configurations = make(map[string]interface{})
		instance.Configurations = configurations
	}

	// 获取providers列表
	providersObj, ok := configurations[constant.AiProxyConfigProviders].([]interface{})
	if !ok {
		providersObj = []interface{}{}
	}

	// 构建提供商配置
	providerConfig := make(map[string]interface{})
	if provider.RawConfigs != nil {
		for k, v := range provider.RawConfigs {
			providerConfig[k] = v
		}
	}
	handler.SaveConfig(provider, providerConfig)

	// 更新或添加提供商配置
	found := false
	for i, p := range providersObj {
		if pMap, ok := p.(map[string]interface{}); ok {
			if provider.Name == pMap[constant.AiProxyConfigProviderId] {
				providersObj[i] = providerConfig
				found = true
				break
			}
		}
	}
	if !found {
		providersObj = append(providersObj, providerConfig)
	}
	configurations[constant.AiProxyConfigProviders] = providersObj

	// 构建服务来源
	serviceSource, err := handler.BuildServiceSource(provider.Name, providerConfig)
	if err != nil {
		return nil, err
	}

	// 构建上游服务
	upstreamService, err := handler.BuildUpstreamService(provider.Name, providerConfig)
	if err != nil {
		return nil, err
	}

	// 检查服务实例是否已绑定其他提供商
	for i := range instances {
		if instances[i].HasScopedTarget(model.WasmPluginInstanceScopeService, upstreamService.Name) {
			if boundProvider, ok := instances[i].Configurations[constant.AiProxyConfigActiveProviderId].(string); ok {
				if boundProvider != provider.Name {
					return nil, errors.NewValidationError("The service instance for provider " + boundProvider + " is already existed. Cannot bind it to provider " + provider.Name)
				}
			}
		}
	}

	// 创建服务实例配置
	serviceInstance := &model.WasmPluginInstance{
		PluginName:    instance.PluginName,
		PluginVersion: instance.PluginVersion,
		Enabled:       &internalTrue,
		Internal:      &internalTrue,
		Configurations: map[string]interface{}{
			constant.AiProxyConfigActiveProviderId: provider.Name,
		},
		Scope:  model.WasmPluginInstanceScopeService,
		Target: upstreamService.Name,
	}

	// 保存服务来源
	if serviceSource != nil {
		_, err = s.serviceSourceService.Add(ctx, serviceSource)
		if err != nil {
			// 如果已存在，尝试更新
			_, err = s.serviceSourceService.Update(ctx, serviceSource)
			if err != nil {
				return nil, err
			}
		}
	}

	// 保存插件实例
	_, err = s.wasmPluginInstanceService.AddOrUpdate(ctx, instance)
	if err != nil {
		return nil, err
	}

	_, err = s.wasmPluginInstanceService.AddOrUpdate(ctx, serviceInstance)
	if err != nil {
		return nil, err
	}

	// 同步相关AI路由
	if handler.NeedSyncRouteAfterUpdate() {
		s.syncRelatedAiRoutes(ctx, provider)
	}

	return s.Get(ctx, provider.Name)
}

// List 列出所有LLM提供商 - 接口适配器方法
func (s *LlmProviderServiceImpl) List(ctx context.Context) ([]model.LlmProvider, error) {
	providers := s.getProviders(ctx)

	// 转换为列表
	providerList := make([]*model.LlmProvider, 0, len(providers))
	for _, p := range providers {
		providerList = append(providerList, p)
	}

	// 排序
	sort.Slice(providerList, func(i, j int) bool {
		return providerList[i].Name < providerList[j].Name
	})

	// 转换为无指针列表
	resultList := make([]model.LlmProvider, 0, len(providerList))
	for _, p := range providerList {
		resultList = append(resultList, *p)
	}

	return resultList, nil
}

// ListWithQuery 列出所有LLM提供商（带分页查询）
func (s *LlmProviderServiceImpl) ListWithQuery(ctx context.Context, query *model.CommonPageQuery) (*model.PaginatedResult[model.LlmProvider], error) {
	resultList, err := s.List(ctx)
	if err != nil {
		return nil, err
	}

	result := model.PaginatedResult[model.LlmProvider]{
		Data:  resultList,
		Total: len(resultList),
	}
	return &result, nil
}

// Get 获取LLM提供商详情
func (s *LlmProviderServiceImpl) Get(ctx context.Context, name string) (*model.LlmProvider, error) {
	providers := s.getProviders(ctx)
	provider := providers[name]
	if provider == nil {
		return nil, errors.NewNotFoundError("LLM provider", name)
	}
	return provider, nil
}

// Delete 删除LLM提供商
func (s *LlmProviderServiceImpl) Delete(ctx context.Context, name string) error {
	instances, err := s.wasmPluginInstanceService.ListByPlugin(ctx, constant.BuiltInPluginAiProxy, nil)
	if err != nil {
		return err
	}

	if len(instances) == 0 {
		return nil
	}

	// 查找全局实例
	var globalInstance *model.WasmPluginInstance
	for i := range instances {
		if instances[i].HasScopedTarget(model.WasmPluginInstanceScopeGlobal, "") {
			globalInstance = &instances[i]
			break
		}
	}

	if globalInstance == nil {
		return nil
	}

	configurations := globalInstance.Configurations
	if configurations == nil {
		return nil
	}

	providersObj, ok := configurations[constant.AiProxyConfigProviders].([]interface{})
	if !ok {
		return nil
	}

	// 查找并删除提供商配置
	var deletedProvider map[string]interface{}
	for i := len(providersObj) - 1; i >= 0; i-- {
		if pMap, ok := providersObj[i].(map[string]interface{}); ok {
			if name == pMap[constant.AiProxyConfigProviderId] {
				deletedProvider = pMap
				providersObj = append(providersObj[:i], providersObj[i+1:]...)
				break
			}
		}
	}

	if deletedProvider == nil {
		return nil
	}

	configurations[constant.AiProxyConfigProviders] = providersObj

	// 删除相关资源
	if providerType, ok := deletedProvider[constant.AiProxyConfigProviderType].(string); ok {
		handler := ai.GetHandler(providerType)
		if handler != nil {
			upstreamService, err := handler.BuildUpstreamService(name, deletedProvider)
			if err == nil {
				internalTrue := true
				_ = s.wasmPluginInstanceService.Delete(ctx, model.WasmPluginInstanceScopeService, upstreamService.Name, constant.BuiltInPluginAiProxy, &internalTrue)
			}

			serviceSource, err := handler.BuildServiceSource(name, deletedProvider)
			if err == nil && serviceSource != nil {
				_ = s.serviceSourceService.Delete(ctx, serviceSource.Name)
			}
		}
	}

	// 保存更新后的全局实例
	_, err = s.wasmPluginInstanceService.AddOrUpdate(ctx, globalInstance)
	return err
}

// BuildUpstreamService 构建上游服务
func (s *LlmProviderServiceImpl) BuildUpstreamService(ctx context.Context, providerName string) (*route.UpstreamService, error) {
	provider, err := s.Get(ctx, providerName)
	if err != nil {
		return nil, err
	}
	if provider == nil {
		return nil, errors.NewValidationError("Unknown provider: " + providerName)
	}

	handler := ai.GetHandler(provider.Type)
	if handler == nil {
		return nil, errors.NewValidationError("Provider type " + provider.Type + " of provider " + providerName + " is not supported")
	}

	return handler.BuildUpstreamService(provider.Name, provider.RawConfigs)
}

// getProviders 获取所有提供商
func (s *LlmProviderServiceImpl) getProviders(ctx context.Context) map[string]*model.LlmProvider {
	result := make(map[string]*model.LlmProvider)

	internalTrue := true
	instance, err := s.wasmPluginInstanceService.Query(ctx, model.WasmPluginInstanceScopeGlobal, "", constant.BuiltInPluginAiProxy, &internalTrue)
	if err != nil || instance == nil {
		return result
	}

	configurations := instance.Configurations
	if configurations == nil {
		return result
	}

	providersObj, ok := configurations[constant.AiProxyConfigProviders].([]interface{})
	if !ok {
		return result
	}

	for _, p := range providersObj {
		pMap, ok := p.(map[string]interface{})
		if !ok {
			continue
		}

		providerType, ok := pMap[constant.AiProxyConfigProviderType].(string)
		if !ok || providerType == "" {
			continue
		}

		handler := ai.GetHandler(providerType)
		if handler == nil {
			continue
		}

		provider := handler.CreateProvider()
		if !handler.LoadConfig(provider, pMap) {
			continue
		}

		result[provider.Name] = provider
	}

	// 填充代理信息
	s.fillProxyInfo(ctx, result)

	return result
}

// fillProxyInfo 填充代理信息
func (s *LlmProviderServiceImpl) fillProxyInfo(ctx context.Context, providers map[string]*model.LlmProvider) {
	if len(providers) == 0 {
		return
	}

	serviceSources, err := s.serviceSourceService.List(ctx, nil)
	if err != nil || serviceSources == nil || len(serviceSources.Data) == 0 {
		return
	}

	serviceSourceMap := make(map[string]model.ServiceSource)
	for _, ss := range serviceSources.Data {
		serviceSourceMap[ss.Name] = ss
	}

	for _, provider := range providers {
		handler := ai.GetHandler(provider.Type)
		if handler == nil {
			continue
		}

		serviceSourceName := handler.GetServiceSourceName(provider.Name)
		if serviceSource, ok := serviceSourceMap[serviceSourceName]; ok {
			// 从服务来源获取代理名称（如果有的话）
			_ = serviceSource
		}
	}
}

// syncRelatedAiRoutes 同步相关AI路由
func (s *LlmProviderServiceImpl) syncRelatedAiRoutes(ctx context.Context, provider *model.LlmProvider) {
	if s.aiRouteService == nil {
		return
	}

	routes, err := s.aiRouteService.List(ctx)
	if err != nil || len(routes) == 0 {
		return
	}

	for i := range routes {
		aiRoute := &routes[i]
		if len(aiRoute.Upstreams) == 0 {
			continue
		}

		hasProvider := false
		for _, upstream := range aiRoute.Upstreams {
			if upstream.Provider == provider.Name {
				hasProvider = true
				break
			}
		}

		if hasProvider {
			_, _ = s.aiRouteService.Update(ctx, aiRoute)
		}
	}
}

// Add 添加LLM提供商 - 接口适配器方法
func (s *LlmProviderServiceImpl) Add(ctx context.Context, provider *model.LlmProvider) (*model.LlmProvider, error) {
	return s.AddOrUpdate(ctx, provider)
}

// Update 更新LLM提供商 - 接口适配器方法
func (s *LlmProviderServiceImpl) Update(ctx context.Context, provider *model.LlmProvider) (*model.LlmProvider, error) {
	return s.AddOrUpdate(ctx, provider)
}
