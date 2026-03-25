// Package mcp provides MCP server related services
package mcp

import (
	"context"
	"fmt"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/constant"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/errors"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
)

// McpServerConfigMapHelper MCP服务器ConfigMap辅助工具
type McpServerConfigMapHelper struct {
	kubernetesClient KubernetesClientInterface
}

// KubernetesClientInterface Kubernetes客户端接口
type KubernetesClientInterface interface {
	GetConfigMap(ctx context.Context, name string) (*corev1.ConfigMap, error)
	UpdateConfigMap(ctx context.Context, cm *corev1.ConfigMap) (*corev1.ConfigMap, error)
	CreateConfigMap(ctx context.Context, cm *corev1.ConfigMap) (*corev1.ConfigMap, error)
}

// NewMcpServerConfigMapHelper 创建MCP服务器ConfigMap辅助工具
func NewMcpServerConfigMapHelper(client KubernetesClientInterface) *McpServerConfigMapHelper {
	return &McpServerConfigMapHelper{
		kubernetesClient: client,
	}
}

// GenerateMcpServerPath 生成MCP服务器路径
func (h *McpServerConfigMapHelper) GenerateMcpServerPath(mcpServerName string) string {
	return constant.McpServerPathPre + mcpServerName
}

// GenerateMatchList 生成匹配规则列表
func (h *McpServerConfigMapHelper) GenerateMatchList(mcpServer *model.McpServer) *model.McpServerConfigMapMatchList {
	result := &model.McpServerConfigMapMatchList{
		MatchRulePath:   h.GenerateMcpServerPath(mcpServer.Name),
		MatchRuleDomain: "*",
		MatchRuleType:   "prefix",
	}
	return result
}

// GetMcpConfig 从ConfigMap获取MCP配置
func (h *McpServerConfigMapHelper) GetMcpConfig(ctx context.Context) (*model.McpServerConfigMap, error) {
	cm, err := h.kubernetesClient.GetConfigMap(ctx, constant.HigressConfig)
	if err != nil {
		return nil, err
	}
	return h.ParseMcpConfigFromConfigMap(cm)
}

// ParseMcpConfigFromConfigMap 从ConfigMap解析MCP配置
func (h *McpServerConfigMapHelper) ParseMcpConfigFromConfigMap(cm *corev1.ConfigMap) (*model.McpServerConfigMap, error) {
	if cm == nil || cm.Data == nil {
		return &model.McpServerConfigMap{}, nil
	}

	higressData, ok := cm.Data[constant.McpConfigKey]
	if !ok {
		return &model.McpServerConfigMap{}, nil
	}

	var higressConfig map[string]interface{}
	if err := yaml.Unmarshal([]byte(higressData), &higressConfig); err != nil {
		return nil, fmt.Errorf("failed to parse higress config: %w", err)
	}

	mcpConfig := &model.McpServerConfigMap{}

	// 解析mcpServer配置
	if mcpServerData, ok := higressConfig[constant.McpServerKey]; ok {
		if mcpServerMap, ok := mcpServerData.(map[string]interface{}); ok {
			// 解析servers
			if serversData, ok := mcpServerMap[constant.ServersKey]; ok {
				if serversList, ok := serversData.([]interface{}); ok {
					servers := make([]model.McpServerConfigMapServer, 0, len(serversList))
					for _, s := range serversList {
						if serverMap, ok := s.(map[string]interface{}); ok {
							server := model.McpServerConfigMapServer{}
							if name, ok := serverMap[constant.ServerNameKey].(string); ok {
								server.Name = name
							}
							if config, ok := serverMap["config"].(map[string]interface{}); ok {
								server.Config = config
							}
							servers = append(servers, server)
						}
					}
					mcpConfig.Servers = servers
				}
			}

			// 解析match_list
			if matchListData, ok := mcpServerMap[constant.MatchListKey]; ok {
				if matchListSlice, ok := matchListData.([]interface{}); ok {
					matchList := make([]model.McpServerConfigMapMatchList, 0, len(matchListSlice))
					for _, m := range matchListSlice {
						if matchMap, ok := m.(map[string]interface{}); ok {
							match := model.McpServerConfigMapMatchList{}
							if path, ok := matchMap[constant.MatchRulePathKey].(string); ok {
								match.MatchRulePath = path
							}
							if domain, ok := matchMap[constant.MatchRuleDomainKey].(string); ok {
								match.MatchRuleDomain = domain
							}
							if matchType, ok := matchMap[constant.MatchRuleTypeKey].(string); ok {
								match.MatchRuleType = matchType
							}
							matchList = append(matchList, match)
						}
					}
					mcpConfig.MatchList = matchList
				}
			}
		}
	}

	return mcpConfig, nil
}

// InitMcpServerConfig 初始化MCP服务器配置
func (h *McpServerConfigMapHelper) InitMcpServerConfig(ctx context.Context) error {
	cm, err := h.kubernetesClient.GetConfigMap(ctx, constant.HigressConfig)
	if err != nil {
		return err
	}

	// 检查是否已存在配置
	if cm.Data == nil {
		cm.Data = make(map[string]string)
	}

	higressData, ok := cm.Data[constant.McpConfigKey]
	if !ok || higressData == "" {
		// 初始化空配置
		initialConfig := map[string]interface{}{
			constant.McpServerKey: map[string]interface{}{
				constant.ServersKey:   []interface{}{},
				constant.MatchListKey: []interface{}{},
			},
		}
		data, err := yaml.Marshal(initialConfig)
		if err != nil {
			return fmt.Errorf("failed to marshal initial config: %w", err)
		}
		cm.Data[constant.McpConfigKey] = string(data)
		_, err = h.kubernetesClient.UpdateConfigMap(ctx, cm)
		return err
	}

	return nil
}

// UpdateServerConfig 更新服务器配置
func (h *McpServerConfigMapHelper) UpdateServerConfig(ctx context.Context, updateFunc func([]model.McpServerConfigMapServer) []model.McpServerConfigMapServer) error {
	cm, err := h.kubernetesClient.GetConfigMap(ctx, constant.HigressConfig)
	if err != nil {
		return err
	}

	mcpConfig, err := h.ParseMcpConfigFromConfigMap(cm)
	if err != nil {
		return err
	}

	// 应用更新函数
	mcpConfig.Servers = updateFunc(mcpConfig.Servers)

	// 保存回ConfigMap
	err = h.saveMcpConfigToConfigMap(cm, mcpConfig)
	if err != nil {
		return err
	}

	_, err = h.kubernetesClient.UpdateConfigMap(ctx, cm)
	return err
}

// UpdateMatchList 更新匹配规则列表
func (h *McpServerConfigMapHelper) UpdateMatchList(ctx context.Context, updateFunc func([]model.McpServerConfigMapMatchList) []model.McpServerConfigMapMatchList) error {
	cm, err := h.kubernetesClient.GetConfigMap(ctx, constant.HigressConfig)
	if err != nil {
		return err
	}

	mcpConfig, err := h.ParseMcpConfigFromConfigMap(cm)
	if err != nil {
		return err
	}

	// 应用更新函数
	mcpConfig.MatchList = updateFunc(mcpConfig.MatchList)

	// 保存回ConfigMap
	err = h.saveMcpConfigToConfigMap(cm, mcpConfig)
	if err != nil {
		return err
	}

	_, err = h.kubernetesClient.UpdateConfigMap(ctx, cm)
	return err
}

// saveMcpConfigToConfigMap 保存MCP配置到ConfigMap
func (h *McpServerConfigMapHelper) saveMcpConfigToConfigMap(cm *corev1.ConfigMap, mcpConfig *model.McpServerConfigMap) error {
	if cm.Data == nil {
		cm.Data = make(map[string]string)
	}

	// 解析现有的higress配置
	higressData := cm.Data[constant.McpConfigKey]
	var higressConfig map[string]interface{}
	if higressData != "" {
		if err := yaml.Unmarshal([]byte(higressData), &higressConfig); err != nil {
			return fmt.Errorf("failed to parse existing higress config: %w", err)
		}
	} else {
		higressConfig = make(map[string]interface{})
	}

	// 构建mcpServer配置
	servers := make([]interface{}, len(mcpConfig.Servers))
	for i, s := range mcpConfig.Servers {
		server := map[string]interface{}{
			constant.ServerNameKey: s.Name,
		}
		if s.Config != nil {
			server["config"] = s.Config
		}
		servers[i] = server
	}

	matchList := make([]interface{}, len(mcpConfig.MatchList))
	for i, m := range mcpConfig.MatchList {
		matchList[i] = map[string]interface{}{
			constant.MatchRulePathKey:   m.MatchRulePath,
			constant.MatchRuleDomainKey: m.MatchRuleDomain,
			constant.MatchRuleTypeKey:   m.MatchRuleType,
		}
	}

	higressConfig[constant.McpServerKey] = map[string]interface{}{
		constant.ServersKey:   servers,
		constant.MatchListKey: matchList,
	}

	data, err := yaml.Marshal(higressConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal higress config: %w", err)
	}

	cm.Data[constant.McpConfigKey] = string(data)
	return nil
}

// AddServer 添加服务器配置
func (h *McpServerConfigMapHelper) AddServer(ctx context.Context, server *model.McpServerConfigMapServer) error {
	return h.UpdateServerConfig(ctx, func(servers []model.McpServerConfigMapServer) []model.McpServerConfigMapServer {
		// 检查是否已存在
		for i, s := range servers {
			if s.Name == server.Name {
				servers[i] = *server
				return servers
			}
		}
		return append(servers, *server)
	})
}

// RemoveServer 删除服务器配置
func (h *McpServerConfigMapHelper) RemoveServer(ctx context.Context, name string) error {
	return h.UpdateServerConfig(ctx, func(servers []model.McpServerConfigMapServer) []model.McpServerConfigMapServer {
		result := make([]model.McpServerConfigMapServer, 0, len(servers))
		for _, s := range servers {
			if s.Name != name {
				result = append(result, s)
			}
		}
		return result
	})
}

// AddMatchList 添加匹配规则
func (h *McpServerConfigMapHelper) AddMatchList(ctx context.Context, match *model.McpServerConfigMapMatchList) error {
	return h.UpdateMatchList(ctx, func(matchList []model.McpServerConfigMapMatchList) []model.McpServerConfigMapMatchList {
		// 检查是否已存在相同路径的规则
		for i, m := range matchList {
			if m.MatchRulePath == match.MatchRulePath {
				matchList[i] = *match
				return matchList
			}
		}
		return append(matchList, *match)
	})
}

// RemoveMatchList 删除匹配规则
func (h *McpServerConfigMapHelper) RemoveMatchList(ctx context.Context, path string) error {
	return h.UpdateMatchList(ctx, func(matchList []model.McpServerConfigMapMatchList) []model.McpServerConfigMapMatchList {
		result := make([]model.McpServerConfigMapMatchList, 0, len(matchList))
		for _, m := range matchList {
			if m.MatchRulePath != path {
				result = append(result, m)
			}
		}
		return result
	})
}

// GetServer 获取服务器配置
func (h *McpServerConfigMapHelper) GetServer(ctx context.Context, name string) (*model.McpServerConfigMapServer, error) {
	mcpConfig, err := h.GetMcpConfig(ctx)
	if err != nil {
		return nil, err
	}

	for _, s := range mcpConfig.Servers {
		if s.Name == name {
			return &s, nil
		}
	}

	return nil, errors.NewNotFoundError("MCP server", name)
}
