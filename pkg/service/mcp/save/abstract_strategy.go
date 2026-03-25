// Package save provides MCP server save strategies
package save

import (
	"context"
	"errors"
	"fmt"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/constant"
	sdkerrors "github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/errors"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model/route"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/service/mcp"
)

// RouteServiceInterface 路由服务接口
type RouteServiceInterface interface {
	Get(ctx context.Context, name string) (*model.Route, error)
	Add(ctx context.Context, r *model.Route) (*model.Route, error)
	Update(ctx context.Context, r *model.Route) (*model.Route, error)
	Delete(ctx context.Context, name string) error
}

// ConsumerServiceInterface 消费者服务接口
type ConsumerServiceInterface interface {
	UpdateAllowList(ctx context.Context, operation model.AllowListOperation, allowList *model.AllowList) error
}

// AbstractMcpServerSaveStrategy 抽象保存策略
type AbstractMcpServerSaveStrategy struct {
	helper          *mcp.McpServerHelper
	configMapHelper *mcp.McpServerConfigMapHelper
	routeService    RouteServiceInterface
	consumerService ConsumerServiceInterface
}

// NewAbstractMcpServerSaveStrategy 创建抽象保存策略
func NewAbstractMcpServerSaveStrategy(
	configMapHelper *mcp.McpServerConfigMapHelper,
	routeService RouteServiceInterface,
	consumerService ConsumerServiceInterface,
) *AbstractMcpServerSaveStrategy {
	return &AbstractMcpServerSaveStrategy{
		helper:          mcp.NewMcpServerHelper(),
		configMapHelper: configMapHelper,
		routeService:    routeService,
		consumerService: consumerService,
	}
}

// Save 保存MCP服务器（模板方法）
func (s *AbstractMcpServerSaveStrategy) Save(ctx context.Context, mcpServer *model.McpServer) (*model.McpServer, error) {
	// 保存路由
	if err := s.saveRoute(ctx, mcpServer); err != nil {
		return nil, err
	}

	// 保存MCP服务器配置（由子类实现）
	if err := s.DoSaveMcpServerConfig(ctx, mcpServer); err != nil {
		return nil, err
	}

	return mcpServer, nil
}

// SaveWithAuthorization 带授权保存
func (s *AbstractMcpServerSaveStrategy) SaveWithAuthorization(ctx context.Context, mcpServer *model.McpServer) (*model.McpServer, error) {
	// 保存路由
	if err := s.saveRoute(ctx, mcpServer); err != nil {
		return nil, err
	}

	// 保存MCP服务器配置（由子类实现）
	if err := s.DoSaveMcpServerConfig(ctx, mcpServer); err != nil {
		return nil, err
	}

	// 保存授权信息
	if err := s.saveAuthInfo(ctx, mcpServer); err != nil {
		return nil, err
	}

	return mcpServer, nil
}

// DoSaveMcpServerConfig 保存MCP服务器配置（抽象方法，由子类实现）
func (s *AbstractMcpServerSaveStrategy) DoSaveMcpServerConfig(ctx context.Context, mcpServer *model.McpServer) error {
	return nil
}

// saveRoute 保存路由
func (s *AbstractMcpServerSaveStrategy) saveRoute(ctx context.Context, mcpServer *model.McpServer) error {
	routeRequest := s.buildRouteRequest(mcpServer)
	if err := routeRequest.Validate(); err != nil {
		return err
	}

	routeName := routeRequest.Name

	// 检查路由是否已存在
	existingRoute, err := s.routeService.Get(ctx, routeName)
	if err != nil {
		var notFoundErr *sdkerrors.NotFoundError
		if !errors.As(err, &notFoundErr) {
			return fmt.Errorf("failed to check existing route: %w", err)
		}
	}

	if existingRoute == nil {
		// 创建新路由
		_, err = s.routeService.Add(ctx, routeRequest)
		if err != nil {
			return fmt.Errorf("failed to create route: %w", err)
		}
	} else {
		// 更新现有路由
		routeRequest.Version = existingRoute.Version
		_, err = s.routeService.Update(ctx, routeRequest)
		if err != nil {
			return fmt.Errorf("failed to update route: %w", err)
		}
	}

	return nil
}

// saveAuthInfo 保存授权信息
func (s *AbstractMcpServerSaveStrategy) saveAuthInfo(ctx context.Context, mcpServer *model.McpServer) error {
	if mcpServer.ConsumerAuthInfo == nil {
		return nil
	}

	routeName := s.helper.McpServerName2RouteName(mcpServer.Name)
	authEnabled := mcpServer.ConsumerAuthInfo.Enable
	allowList := &model.AllowList{
		AuthEnabled:     &authEnabled,
		CredentialTypes: []string{mcpServer.ConsumerAuthInfo.Type},
		ConsumerNames:   mcpServer.ConsumerAuthInfo.AllowedConsumers,
	}
	allowList.Targets = map[model.WasmPluginInstanceScope]string{
		model.WasmPluginInstanceScopeRoute: routeName,
	}

	return s.consumerService.UpdateAllowList(ctx, model.AllowListOperationReplace, allowList)
}

// buildRouteRequest 构建路由请求
func (s *AbstractMcpServerSaveStrategy) buildRouteRequest(mcpServer *model.McpServer) *model.Route {
	routeName := s.helper.McpServerName2RouteName(mcpServer.Name)

	// 转换Services类型 []route.UpstreamService -> []*route.UpstreamService
	services := make([]*route.UpstreamService, len(mcpServer.Services))
	for i, svc := range mcpServer.Services {
		svcCopy := svc
		services[i] = &svcCopy
	}

	r := &model.Route{
		Name:     routeName,
		Services: services,
		Domains:  mcpServer.Domains,
		Path: &route.RoutePredicate{
			MatchType: route.MatchTypePrefix,
			Path:      s.configMapHelper.GenerateMcpServerPath(mcpServer.Name),
		},
	}

	s.setDefaultConfigs(r, mcpServer)
	s.setDefaultLabels(r, mcpServer)

	return r
}

// setDefaultConfigs 设置默认配置
func (s *AbstractMcpServerSaveStrategy) setDefaultConfigs(r *model.Route, mcpServer *model.McpServer) {
	if r.CustomConfigs == nil {
		r.CustomConfigs = make(map[string]string)
	}

	r.CustomConfigs[constant.AnnotationResourceDescriptionKey] = mcpServer.Description
	r.CustomConfigs[constant.AnnotationResourceMcpServerKey] = "true"

	// 设置匹配规则域名
	matchRuleDomains := "*"
	if len(r.Domains) > 0 {
		// 过滤空域名
		domains := make([]string, 0, len(r.Domains))
		for _, d := range r.Domains {
			if d != "" {
				domains = append(domains, d)
			}
		}
		if len(domains) > 0 {
			matchRuleDomains = domains[0]
			if len(domains) > 1 {
				for i := 1; i < len(domains); i++ {
					matchRuleDomains += "," + domains[i]
				}
			}
		}
	}
	r.CustomConfigs[constant.AnnotationResourceMcpServerMatchRuleDomainsKey] = matchRuleDomains
	r.CustomConfigs[constant.AnnotationResourceMcpServerMatchRuleTypeKey] = "prefix"
	r.CustomConfigs[constant.AnnotationResourceMcpServerMatchRuleValueKey] = s.configMapHelper.GenerateMcpServerPath(mcpServer.Name)
}

// setDefaultLabels 设置默认标签
func (s *AbstractMcpServerSaveStrategy) setDefaultLabels(r *model.Route, mcpServer *model.McpServer) {
	if r.CustomLabels == nil {
		r.CustomLabels = make(map[string]string)
	}

	r.CustomLabels[constant.LabelResourceDefinerKey] = constant.LabelResourceDefinerValue
	r.CustomLabels[constant.LabelInternalKey] = "true"
	r.CustomLabels[constant.LabelResourceBizTypeKey] = constant.LabelMcpServerBizTypeValue
	r.CustomLabels[constant.LabelResourceMcpServerTypeKey] = string(mcpServer.Type)
}

// GetHelper 获取辅助工具
func (s *AbstractMcpServerSaveStrategy) GetHelper() *mcp.McpServerHelper {
	return s.helper
}

// GetConfigMapHelper 获取ConfigMap辅助工具
func (s *AbstractMcpServerSaveStrategy) GetConfigMapHelper() *mcp.McpServerConfigMapHelper {
	return s.configMapHelper
}

// GetRouteService 获取路由服务
func (s *AbstractMcpServerSaveStrategy) GetRouteService() RouteServiceInterface {
	return s.routeService
}

// GetConsumerService 获取消费者服务
func (s *AbstractMcpServerSaveStrategy) GetConsumerService() ConsumerServiceInterface {
	return s.consumerService
}
