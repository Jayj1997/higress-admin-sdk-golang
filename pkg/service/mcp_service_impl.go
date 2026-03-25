// Package service provides business services for the SDK
package service

import (
	"context"
	"sort"
	"strings"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/errors"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/service/mcp"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/service/mcp/detail"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/service/mcp/save"
)

// McpKubernetesClientInterface MCP服务Kubernetes客户端接口
type McpKubernetesClientInterface interface {
	ListIngresses(ctx context.Context, labelSelectors map[string]string) (interface{}, error)
	GetConfigMap(ctx context.Context, name string) (interface{}, error)
	UpdateConfigMap(ctx context.Context, cm interface{}) (interface{}, error)
	CreateConfigMap(ctx context.Context, cm interface{}) (interface{}, error)
}

// McpServiceContextImpl MCP服务器服务实现
type McpServiceContextImpl struct {
	routeService          RouteService
	consumerService       ConsumerService
	helper                *mcp.McpServerHelper
	configMapHelper       *mcp.McpServerConfigMapHelper
	saveStrategyFactory   *save.McpServerSaveStrategyFactory
	detailStrategyFactory *detail.McpServerDetailStrategyFactory
}

// NewMcpServiceContextImpl 创建MCP服务器服务实现
func NewMcpServiceContextImpl(
	kubernetesClient mcp.KubernetesClientInterface,
	routeService RouteService,
	consumerService ConsumerService,
) *McpServiceContextImpl {
	helper := mcp.NewMcpServerHelper()
	configMapHelper := mcp.NewMcpServerConfigMapHelper(kubernetesClient)

	return &McpServiceContextImpl{
		routeService:          routeService,
		consumerService:       consumerService,
		helper:                helper,
		configMapHelper:       configMapHelper,
		saveStrategyFactory:   save.NewMcpServerSaveStrategyFactory(configMapHelper, routeService, consumerService),
		detailStrategyFactory: detail.NewMcpServerDetailStrategyFactory(configMapHelper),
	}
}

// List 列出MCP服务器
func (s *McpServiceContextImpl) List(ctx context.Context, query *model.McpServerPageQuery) (*model.PaginatedResult[model.McpServer], error) {
	// 通过路由服务获取MCP服务器路由
	routes, err := s.routeService.List(ctx, &model.RoutePageQuery{})
	if err != nil {
		return nil, err
	}

	// 过滤MCP服务器路由
	resultList := make([]model.McpServer, 0)
	if routes != nil && routes.Data != nil {
		for _, route := range routes.Data {
			if s.helper.IsMcpServerRoute(route.CustomLabels) {
				mcpServer := s.helper.RouteToMcpServer(&route)
				if mcpServer != nil {
					resultList = append(resultList, *mcpServer)
				}
			}
		}
	}

	// 过滤
	resultList = s.filterMcpServers(resultList, query)

	// 排序
	sort.Slice(resultList, func(i, j int) bool {
		return resultList[i].Name < resultList[j].Name
	})

	return model.PaginatedResultFromFullList(resultList, &query.CommonPageQuery), nil
}

// Get 获取MCP服务器详情
func (s *McpServiceContextImpl) Get(ctx context.Context, name string) (*model.McpServer, error) {
	routeName := s.helper.McpServerName2RouteName(name)
	route, err := s.routeService.Get(ctx, routeName)
	if err != nil {
		return nil, err
	}

	if route == nil || !s.helper.IsMcpServerRoute(route.CustomLabels) {
		return nil, errors.NewNotFoundError("MCP server", name)
	}

	mcpServer := s.helper.RouteToMcpServer(route)
	if mcpServer == nil {
		return nil, errors.NewNotFoundError("MCP server", name)
	}

	// 获取详情
	strategy := s.detailStrategyFactory.GetService(mcpServer)
	if strategy != nil {
		detail, err := strategy.Query(ctx, name)
		if err == nil && detail != nil {
			// 合并详情信息
			if detail.RawConfigurations != "" {
				mcpServer.RawConfigurations = detail.RawConfigurations
			}
			if detail.DBConfig != nil {
				mcpServer.DBConfig = detail.DBConfig
			}
		}
	}

	return mcpServer, nil
}

// Add 添加MCP服务器
func (s *McpServiceContextImpl) Add(ctx context.Context, server *model.McpServer) (*model.McpServer, error) {
	strategy := s.saveStrategyFactory.GetService(server)
	if strategy == nil {
		return nil, errors.NewValidationError("unsupported MCP server type")
	}
	return strategy.Save(ctx, server)
}

// Update 更新MCP服务器
func (s *McpServiceContextImpl) Update(ctx context.Context, server *model.McpServer) (*model.McpServer, error) {
	strategy := s.saveStrategyFactory.GetService(server)
	if strategy == nil {
		return nil, errors.NewValidationError("unsupported MCP server type")
	}
	return strategy.Save(ctx, server)
}

// AddOrUpdate 添加或更新MCP服务器
func (s *McpServiceContextImpl) AddOrUpdate(ctx context.Context, server *model.McpServer) (*model.McpServer, error) {
	strategy := s.saveStrategyFactory.GetService(server)
	if strategy == nil {
		return nil, errors.NewValidationError("unsupported MCP server type")
	}
	return strategy.Save(ctx, server)
}

// AddOrUpdateWithAuthorization 带授权添加或更新MCP服务器
func (s *McpServiceContextImpl) AddOrUpdateWithAuthorization(ctx context.Context, server *model.McpServer) (*model.McpServer, error) {
	strategy := s.saveStrategyFactory.GetService(server)
	if strategy == nil {
		return nil, errors.NewValidationError("unsupported MCP server type")
	}
	return strategy.SaveWithAuthorization(ctx, server)
}

// Delete 删除MCP服务器
func (s *McpServiceContextImpl) Delete(ctx context.Context, name string) error {
	// 删除路由
	routeName := s.helper.McpServerName2RouteName(name)
	if err := s.routeService.Delete(ctx, routeName); err != nil {
		return err
	}

	// 删除匹配规则
	path := s.configMapHelper.GenerateMcpServerPath(name)
	_ = s.configMapHelper.RemoveMatchList(ctx, path)

	// 删除服务器配置
	_ = s.configMapHelper.RemoveServer(ctx, name)

	return nil
}

// AddConsumer 添加消费者到MCP服务器
func (s *McpServiceContextImpl) AddConsumer(ctx context.Context, serverName string, consumer *model.McpServerConsumer) error {
	route, err := s.getMcpServerBoundRoute(ctx, serverName)
	if err != nil {
		return err
	}

	allowList := model.ForTarget(model.WasmPluginInstanceScopeRoute, route.Name)
	allowList.ConsumerNames = []string{consumer.ConsumerName}

	return s.consumerService.UpdateAllowList(ctx, model.AllowListOperationAdd, allowList)
}

// RemoveConsumer 从MCP服务器移除消费者
func (s *McpServiceContextImpl) RemoveConsumer(ctx context.Context, serverName, consumerName string) error {
	route, err := s.getMcpServerBoundRoute(ctx, serverName)
	if err != nil {
		return err
	}

	allowList := model.ForTarget(model.WasmPluginInstanceScopeRoute, route.Name)
	allowList.ConsumerNames = []string{consumerName}

	return s.consumerService.UpdateAllowList(ctx, model.AllowListOperationRemove, allowList)
}

// ListConsumers 列出MCP服务器的消费者
func (s *McpServiceContextImpl) ListConsumers(ctx context.Context, query *model.McpServerPageQuery) (*model.PaginatedResult[model.McpServerConsumerDetail], error) {
	route, err := s.getMcpServerBoundRoute(ctx, query.McpServerName)
	if err != nil {
		return nil, err
	}

	allowedConsumers := route.AuthConfig.AllowedConsumers
	if allowedConsumers == nil {
		allowedConsumers = []string{}
	}

	// 过滤
	resultList := make([]model.McpServerConsumerDetail, 0)
	for _, consumer := range allowedConsumers {
		if query.McpServerName == "" || strings.Contains(consumer, query.McpServerName) {
			resultList = append(resultList, model.McpServerConsumerDetail{
				McpServerConsumer: model.McpServerConsumer{ConsumerName: consumer},
				McpServerName:     query.McpServerName,
			})
		}
	}

	return model.PaginatedResultFromFullList(resultList, &query.CommonPageQuery), nil
}

// AddAllowConsumers 添加允许的消费者
func (s *McpServiceContextImpl) AddAllowConsumers(ctx context.Context, consumers *model.McpServerConsumers) error {
	route, err := s.getMcpServerBoundRoute(ctx, consumers.McpServerName)
	if err != nil {
		return err
	}

	allowList := model.ForTarget(model.WasmPluginInstanceScopeRoute, route.Name)
	allowList.ConsumerNames = consumers.Consumers

	return s.consumerService.UpdateAllowList(ctx, model.AllowListOperationAdd, allowList)
}

// RemoveAllowConsumers 移除允许的消费者
func (s *McpServiceContextImpl) RemoveAllowConsumers(ctx context.Context, consumers *model.McpServerConsumers) error {
	route, err := s.getMcpServerBoundRoute(ctx, consumers.McpServerName)
	if err != nil {
		return err
	}

	allowList := model.ForTarget(model.WasmPluginInstanceScopeRoute, route.Name)
	allowList.ConsumerNames = consumers.Consumers

	return s.consumerService.UpdateAllowList(ctx, model.AllowListOperationRemove, allowList)
}

// ListAllowConsumers 列出允许的消费者
func (s *McpServiceContextImpl) ListAllowConsumers(ctx context.Context, query *model.McpServerConsumersPageQuery) (*model.PaginatedResult[model.McpServerConsumerDetail], error) {
	route, err := s.getMcpServerBoundRoute(ctx, query.McpServerName)
	if err != nil {
		return nil, err
	}

	allowedConsumers := route.AuthConfig.AllowedConsumers
	if allowedConsumers == nil {
		allowedConsumers = []string{}
	}

	// 过滤
	resultList := make([]model.McpServerConsumerDetail, 0)
	for _, consumer := range allowedConsumers {
		if query.ConsumerName == "" || strings.Contains(consumer, query.ConsumerName) {
			resultList = append(resultList, model.McpServerConsumerDetail{
				McpServerConsumer: model.McpServerConsumer{ConsumerName: consumer},
				McpServerName:     query.McpServerName,
			})
		}
	}

	return model.PaginatedResultFromFullList(resultList, &query.CommonPageQuery), nil
}

// getMcpServerBoundRoute 获取MCP服务器绑定的路由
func (s *McpServiceContextImpl) getMcpServerBoundRoute(ctx context.Context, serverName string) (*model.Route, error) {
	routeName := s.helper.McpServerName2RouteName(serverName)
	route, err := s.routeService.Get(ctx, routeName)
	if err != nil {
		return nil, err
	}
	if route == nil {
		return nil, errors.NewBusinessError("No MCP-bound route found for server: " + serverName)
	}
	return route, nil
}

// filterMcpServers 过滤MCP服务器列表
func (s *McpServiceContextImpl) filterMcpServers(list []model.McpServer, query *model.McpServerPageQuery) []model.McpServer {
	if query == nil {
		return list
	}

	result := make([]model.McpServer, 0)
	for _, server := range list {
		// 名称过滤
		if query.McpServerName != "" && !strings.Contains(server.Name, query.McpServerName) {
			continue
		}
		// 类型过滤
		if query.Type != "" && string(server.Type) != query.Type {
			continue
		}
		result = append(result, server)
	}
	return result
}
