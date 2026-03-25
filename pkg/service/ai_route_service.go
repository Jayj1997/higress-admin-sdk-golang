// Package service provides business services for the SDK
package service

import (
	"context"
	"sort"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/constant"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/errors"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model/route"
)

// AiRouteServiceImpl AI路由服务实现
type AiRouteServiceImpl struct {
	wasmPluginInstanceService WasmPluginInstanceService
	llmProviderService        LlmProviderService
}

// NewAiRouteServiceImpl 创建AI路由服务
func NewAiRouteServiceImpl(
	wasmPluginInstanceService WasmPluginInstanceService,
	llmProviderService LlmProviderService,
) *AiRouteServiceImpl {
	return &AiRouteServiceImpl{
		wasmPluginInstanceService: wasmPluginInstanceService,
		llmProviderService:        llmProviderService,
	}
}

// List 列出所有AI路由
func (s *AiRouteServiceImpl) List(ctx context.Context) ([]model.AiRoute, error) {
	routes := s.getRoutes(ctx)

	// 转换为列表
	routeList := make([]model.AiRoute, 0, len(routes))
	for _, r := range routes {
		routeList = append(routeList, *r)
	}

	// 排序
	sort.Slice(routeList, func(i, j int) bool {
		return routeList[i].Name < routeList[j].Name
	})

	return routeList, nil
}

// Get 获取AI路由详情
func (s *AiRouteServiceImpl) Get(ctx context.Context, name string) (*model.AiRoute, error) {
	routes := s.getRoutes(ctx)
	aiRoute := routes[name]
	if aiRoute == nil {
		return nil, errors.NewNotFoundError("AI route", name)
	}
	return aiRoute, nil
}

// Add 添加AI路由
func (s *AiRouteServiceImpl) Add(ctx context.Context, aiRoute *model.AiRoute) (*model.AiRoute, error) {
	if err := aiRoute.Validate(); err != nil {
		return nil, err
	}

	// 验证提供商是否存在
	if err := s.validateProviders(ctx, aiRoute); err != nil {
		return nil, err
	}

	// 获取现有的插件实例列表
	internalTrue := true
	instances, err := s.wasmPluginInstanceService.ListByPlugin(ctx, constant.BuiltInPluginModelRouter, &internalTrue)
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
		instance, err = s.wasmPluginInstanceService.CreateEmptyInstance(ctx, constant.BuiltInPluginModelRouter)
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

	// 获取路由列表
	routesObj, ok := configurations[constant.ModelRouterConfigRoutes].([]interface{})
	if !ok {
		routesObj = []interface{}{}
	}

	// 检查路由是否已存在
	for _, r := range routesObj {
		if rMap, ok := r.(map[string]interface{}); ok {
			if aiRoute.Name == rMap[constant.ModelRouterConfigRouteName] {
				return nil, errors.NewValidationError("AI route already exists: " + aiRoute.Name)
			}
		}
	}

	// 添加路由配置
	routeConfig := s.buildRouteConfig(aiRoute)
	routesObj = append(routesObj, routeConfig)
	configurations[constant.ModelRouterConfigRoutes] = routesObj

	// 保存插件实例
	_, err = s.wasmPluginInstanceService.AddOrUpdate(ctx, instance)
	if err != nil {
		return nil, err
	}

	return s.Get(ctx, aiRoute.Name)
}

// Update 更新AI路由
func (s *AiRouteServiceImpl) Update(ctx context.Context, aiRoute *model.AiRoute) (*model.AiRoute, error) {
	if err := aiRoute.Validate(); err != nil {
		return nil, err
	}

	// 验证提供商是否存在
	if err := s.validateProviders(ctx, aiRoute); err != nil {
		return nil, err
	}

	// 获取现有的插件实例列表
	internalTrue := true
	instances, err := s.wasmPluginInstanceService.ListByPlugin(ctx, constant.BuiltInPluginModelRouter, &internalTrue)
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

	if instance == nil {
		return nil, errors.NewNotFoundError("AI route", aiRoute.Name)
	}

	configurations := instance.Configurations
	if configurations == nil {
		return nil, errors.NewNotFoundError("AI route", aiRoute.Name)
	}

	routesObj, ok := configurations[constant.ModelRouterConfigRoutes].([]interface{})
	if !ok {
		return nil, errors.NewNotFoundError("AI route", aiRoute.Name)
	}

	// 查找并更新路由配置
	found := false
	routeConfig := s.buildRouteConfig(aiRoute)
	for i, r := range routesObj {
		if rMap, ok := r.(map[string]interface{}); ok {
			if aiRoute.Name == rMap[constant.ModelRouterConfigRouteName] {
				routesObj[i] = routeConfig
				found = true
				break
			}
		}
	}

	if !found {
		return nil, errors.NewNotFoundError("AI route", aiRoute.Name)
	}

	configurations[constant.ModelRouterConfigRoutes] = routesObj

	// 保存插件实例
	_, err = s.wasmPluginInstanceService.AddOrUpdate(ctx, instance)
	if err != nil {
		return nil, err
	}

	return s.Get(ctx, aiRoute.Name)
}

// Delete 删除AI路由
func (s *AiRouteServiceImpl) Delete(ctx context.Context, name string) error {
	internalTrue := true
	instances, err := s.wasmPluginInstanceService.ListByPlugin(ctx, constant.BuiltInPluginModelRouter, &internalTrue)
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

	routesObj, ok := configurations[constant.ModelRouterConfigRoutes].([]interface{})
	if !ok {
		return nil
	}

	// 查找并删除路由配置
	found := false
	for i := len(routesObj) - 1; i >= 0; i-- {
		if rMap, ok := routesObj[i].(map[string]interface{}); ok {
			if name == rMap[constant.ModelRouterConfigRouteName] {
				routesObj = append(routesObj[:i], routesObj[i+1:]...)
				found = true
				break
			}
		}
	}

	if !found {
		return nil
	}

	configurations[constant.ModelRouterConfigRoutes] = routesObj

	// 保存更新后的全局实例
	_, err = s.wasmPluginInstanceService.AddOrUpdate(ctx, globalInstance)
	return err
}

// getRoutes 获取所有AI路由
func (s *AiRouteServiceImpl) getRoutes(ctx context.Context) map[string]*model.AiRoute {
	result := make(map[string]*model.AiRoute)

	internalTrue := true
	instance, err := s.wasmPluginInstanceService.Query(ctx, model.WasmPluginInstanceScopeGlobal, "", constant.BuiltInPluginModelRouter, &internalTrue)
	if err != nil || instance == nil {
		return result
	}

	configurations := instance.Configurations
	if configurations == nil {
		return result
	}

	routesObj, ok := configurations[constant.ModelRouterConfigRoutes].([]interface{})
	if !ok {
		return result
	}

	for _, r := range routesObj {
		rMap, ok := r.(map[string]interface{})
		if !ok {
			continue
		}

		aiRoute := s.parseRouteConfig(rMap)
		if aiRoute != nil {
			result[aiRoute.Name] = aiRoute
		}
	}

	return result
}

// buildRouteConfig 构建路由配置
func (s *AiRouteServiceImpl) buildRouteConfig(aiRoute *model.AiRoute) map[string]interface{} {
	config := make(map[string]interface{})
	config[constant.ModelRouterConfigRouteName] = aiRoute.Name

	if len(aiRoute.Domains) > 0 {
		config[constant.ModelRouterConfigRouteDomains] = aiRoute.Domains
	}

	if aiRoute.PathPredicate != nil {
		config[constant.ModelRouterConfigRoutePathPredicate] = map[string]interface{}{
			"matchType": aiRoute.PathPredicate.MatchType,
			"path":      aiRoute.PathPredicate.Path,
		}
	}

	if len(aiRoute.Upstreams) > 0 {
		upstreams := make([]interface{}, 0, len(aiRoute.Upstreams))
		for _, u := range aiRoute.Upstreams {
			upstream := map[string]interface{}{
				"provider": u.Provider,
				"weight":   u.Weight,
			}
			if len(u.ModelMapping) > 0 {
				upstream["modelMapping"] = u.ModelMapping
			}
			upstreams = append(upstreams, upstream)
		}
		config[constant.ModelRouterConfigRouteUpstreams] = upstreams
	}

	return config
}

// parseRouteConfig 解析路由配置
func (s *AiRouteServiceImpl) parseRouteConfig(config map[string]interface{}) *model.AiRoute {
	aiRoute := &model.AiRoute{}

	if name, ok := config[constant.ModelRouterConfigRouteName].(string); ok {
		aiRoute.Name = name
	} else {
		return nil
	}

	if domains, ok := config[constant.ModelRouterConfigRouteDomains].([]interface{}); ok {
		aiRoute.Domains = make([]string, 0, len(domains))
		for _, d := range domains {
			if domain, ok := d.(string); ok {
				aiRoute.Domains = append(aiRoute.Domains, domain)
			}
		}
	}

	if pathPred, ok := config[constant.ModelRouterConfigRoutePathPredicate].(map[string]interface{}); ok {
		aiRoute.PathPredicate = &route.RoutePredicate{}
		if matchType, ok := pathPred["matchType"].(string); ok {
			aiRoute.PathPredicate.MatchType = matchType
		}
		if path, ok := pathPred["path"].(string); ok {
			aiRoute.PathPredicate.Path = path
		}
	}

	if upstreams, ok := config[constant.ModelRouterConfigRouteUpstreams].([]interface{}); ok {
		aiRoute.Upstreams = make([]model.AiUpstream, 0, len(upstreams))
		for _, u := range upstreams {
			if uMap, ok := u.(map[string]interface{}); ok {
				upstream := model.AiUpstream{}
				if provider, ok := uMap["provider"].(string); ok {
					upstream.Provider = provider
				}
				if weight, ok := uMap["weight"].(int); ok {
					upstream.Weight = weight
				} else if weightFloat, ok := uMap["weight"].(float64); ok {
					upstream.Weight = int(weightFloat)
				}
				if modelMapping, ok := uMap["modelMapping"].(map[string]interface{}); ok {
					upstream.ModelMapping = make(map[string]string)
					for k, v := range modelMapping {
						if vs, ok := v.(string); ok {
							upstream.ModelMapping[k] = vs
						}
					}
				}
				aiRoute.Upstreams = append(aiRoute.Upstreams, upstream)
			}
		}
	}

	return aiRoute
}

// validateProviders 验证提供商是否存在
func (s *AiRouteServiceImpl) validateProviders(ctx context.Context, aiRoute *model.AiRoute) error {
	if s.llmProviderService == nil {
		return nil
	}

	for _, upstream := range aiRoute.Upstreams {
		_, err := s.llmProviderService.Get(ctx, upstream.Provider)
		if err != nil {
			return errors.NewValidationError("provider not found: " + upstream.Provider)
		}
	}
	return nil
}
