// Package ai provides AI-related services for Higress Admin SDK.
package ai

import (
	"strconv"
	"strings"

	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/constant"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/errors"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model/route"
)

// DefaultLlmProviderHandler 默认LLM提供商处理器
type DefaultLlmProviderHandler struct {
	providerType string
	domain       string
	port         int
	protocol     string
	contextPath  string
}

// NewDefaultLlmProviderHandler 创建默认处理器
func NewDefaultLlmProviderHandler(providerType, domain string, port int, protocol string) *DefaultLlmProviderHandler {
	return NewDefaultLlmProviderHandlerWithContextPath(providerType, domain, port, protocol, "/")
}

// NewDefaultLlmProviderHandlerWithContextPath 创建带上下文路径的默认处理器
func NewDefaultLlmProviderHandlerWithContextPath(providerType, domain string, port int, protocol, contextPath string) *DefaultLlmProviderHandler {
	return &DefaultLlmProviderHandler{
		providerType: providerType,
		domain:       domain,
		port:         port,
		protocol:     protocol,
		contextPath:  contextPath,
	}
}

// GetType 获取处理器类型
func (h *DefaultLlmProviderHandler) GetType() string {
	return h.providerType
}

// CreateProvider 创建提供商实例
func (h *DefaultLlmProviderHandler) CreateProvider() *model.LlmProvider {
	return &model.LlmProvider{
		Type:     h.providerType,
		Protocol: model.LlmProviderProtocolOpenaiV1,
	}
}

// LoadConfig 从配置加载提供商信息
func (h *DefaultLlmProviderHandler) LoadConfig(provider *model.LlmProvider, configurations map[string]interface{}) bool {
	id := getString(configurations, constant.AiProxyConfigProviderId)
	if id == "" {
		return false
	}

	tokens := getStringSlice(configurations, constant.AiProxyConfigProviderApiTokens)

	failoverConfig := h.buildTokenFailoverConfig(configurations)

	protocol := model.LlmProviderProtocolFromValue(getString(configurations, constant.AiProxyConfigProtocol))
	if protocol == "" {
		protocol = model.LlmProviderProtocolOpenaiV1
	}

	provider.Name = id
	provider.Type = h.providerType
	provider.Protocol = protocol
	provider.Tokens = tokens
	provider.TokenFailoverConfig = failoverConfig
	provider.RawConfigs = copyMap(configurations)

	return true
}

// SaveConfig 保存提供商配置
func (h *DefaultLlmProviderHandler) SaveConfig(provider *model.LlmProvider, configurations map[string]interface{}) {
	configurations[constant.AiProxyConfigProviderId] = provider.Name
	configurations[constant.AiProxyConfigProviderType] = h.providerType

	protocol := model.LlmProviderProtocolFromValue(provider.Protocol)
	if protocol == "" {
		protocol = model.LlmProviderProtocolOpenaiV1
	}
	configurations[constant.AiProxyConfigProtocol] = protocol

	if len(provider.Tokens) > 0 {
		configurations[constant.AiProxyConfigProviderApiTokens] = provider.Tokens
	} else {
		delete(configurations, constant.AiProxyConfigProviderApiTokens)
	}

	if provider.TokenFailoverConfig == nil {
		delete(configurations, constant.AiProxyConfigFailover)
		delete(configurations, constant.AiProxyConfigRetryOnFailure)
	} else {
		failoverMap := make(map[string]interface{})
		h.saveTokenFailoverConfig(provider.TokenFailoverConfig, failoverMap)
		configurations[constant.AiProxyConfigFailover] = failoverMap
		configurations[constant.AiProxyConfigRetryOnFailure] = map[string]interface{}{
			constant.AiProxyConfigRetryEnabled: provider.TokenFailoverConfig.Enabled,
		}
	}
}

// NormalizeConfigs 规范化配置
func (h *DefaultLlmProviderHandler) NormalizeConfigs(configurations map[string]interface{}) {
	// 默认不做任何处理
}

// BuildServiceSource 构建服务来源
func (h *DefaultLlmProviderHandler) BuildServiceSource(providerName string, providerConfig map[string]interface{}) (*model.ServiceSource, error) {
	endpoints := h.GetProviderEndpoints(providerConfig)
	if len(endpoints) == 0 {
		return nil, errors.NewValidationError("No endpoints found for provider: " + providerName)
	}

	serviceSource := &model.ServiceSource{
		Name: h.GetServiceSourceName(providerName),
	}

	var sourceType string
	var domains []string
	var port int

	for _, endpoint := range endpoints {
		if err := endpoint.Validate(); err != nil {
			return nil, err
		}

		// 判断是IP还是域名
		if isIPAddress(endpoint.Address) {
			sourceType = constant.RegistryTypeStatic
			domains = append(domains, endpoint.Address+":"+strconv.Itoa(endpoint.Port))
			port = constant.StaticPort
		} else {
			if len(endpoints) > 1 {
				return nil, errors.NewValidationError("Multiple endpoints only work with static IP addresses, provider: " + providerName)
			}
			port = endpoint.Port
			sourceType = constant.RegistryTypeDNS
			domains = append(domains, endpoint.Address)
		}
	}

	serviceSource.Type = sourceType
	serviceSource.Domain = strings.Join(domains, ",")
	serviceSource.Port = &port

	return serviceSource, nil
}

// BuildUpstreamService 构建上游服务
func (h *DefaultLlmProviderHandler) BuildUpstreamService(providerName string, providerConfig map[string]interface{}) (*route.UpstreamService, error) {
	serviceSource, err := h.BuildServiceSource(providerName, providerConfig)
	if err != nil {
		return nil, err
	}

	service := &route.UpstreamService{
		Name: serviceSource.Name + "." + serviceSource.Type,
	}
	if serviceSource.Port != nil {
		service.Port = *serviceSource.Port
	}
	weight := 100
	service.Weight = &weight

	return service, nil
}

// GetServiceSourceName 获取服务来源名称
func (h *DefaultLlmProviderHandler) GetServiceSourceName(providerName string) string {
	return constant.LlmServiceNamePrefix + providerName + constant.InternalResourceNameSuffix
}

// GetExtraServiceSources 获取额外的服务来源
func (h *DefaultLlmProviderHandler) GetExtraServiceSources(providerName string, providerConfig map[string]interface{}, forDelete bool) []model.ServiceSource {
	return nil
}

// NeedSyncRouteAfterUpdate 更新后是否需要同步路由
func (h *DefaultLlmProviderHandler) NeedSyncRouteAfterUpdate() bool {
	return false
}

// GetProviderEndpoints 获取提供商端点
func (h *DefaultLlmProviderHandler) GetProviderEndpoints(providerConfig map[string]interface{}) []model.LlmProviderEndpoint {
	return []model.LlmProviderEndpoint{
		{
			Protocol:    h.protocol,
			Address:     h.domain,
			Port:        h.port,
			ContextPath: h.contextPath,
		},
	}
}

// buildTokenFailoverConfig 构建故障转移配置
func (h *DefaultLlmProviderHandler) buildTokenFailoverConfig(configurations map[string]interface{}) *model.TokenFailoverConfig {
	failoverObj, ok := configurations[constant.AiProxyConfigFailover]
	if !ok {
		return nil
	}

	failoverMap, ok := failoverObj.(map[string]interface{})
	if !ok {
		return nil
	}

	return &model.TokenFailoverConfig{
		Enabled:             getBool(failoverMap, constant.AiProxyConfigFailoverEnabled, false),
		FailureThreshold:    getInt(failoverMap, constant.AiProxyConfigFailoverFailureThreshold),
		SuccessThreshold:    getInt(failoverMap, constant.AiProxyConfigFailoverSuccessThreshold),
		HealthCheckInterval: getInt(failoverMap, constant.AiProxyConfigFailoverHealthCheckInterval),
		HealthCheckTimeout:  getInt(failoverMap, constant.AiProxyConfigFailoverHealthCheckTimeout),
		HealthCheckModel:    getString(failoverMap, constant.AiProxyConfigFailoverHealthCheckModel),
	}
}

// saveTokenFailoverConfig 保存故障转移配置
func (h *DefaultLlmProviderHandler) saveTokenFailoverConfig(config *model.TokenFailoverConfig, failoverMap map[string]interface{}) {
	failoverMap[constant.AiProxyConfigFailoverEnabled] = config.Enabled
	failoverMap[constant.AiProxyConfigFailoverFailureThreshold] = config.FailureThreshold
	failoverMap[constant.AiProxyConfigFailoverSuccessThreshold] = config.SuccessThreshold
	failoverMap[constant.AiProxyConfigFailoverHealthCheckInterval] = config.HealthCheckInterval
	failoverMap[constant.AiProxyConfigFailoverHealthCheckTimeout] = config.HealthCheckTimeout
	failoverMap[constant.AiProxyConfigFailoverHealthCheckModel] = config.HealthCheckModel
}

// 辅助函数

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getInt(m map[string]interface{}, key string) int {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case int:
			return val
		case int64:
			return int(val)
		case float64:
			return int(val)
		case string:
			if i, err := strconv.Atoi(val); err == nil {
				return i
			}
		}
	}
	return 0
}

func getBool(m map[string]interface{}, key string, defaultValue bool) bool {
	if v, ok := m[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return defaultValue
}

func getStringSlice(m map[string]interface{}, key string) []string {
	if v, ok := m[key]; ok {
		if slice, ok := v.([]interface{}); ok {
			result := make([]string, 0, len(slice))
			for _, item := range slice {
				if s, ok := item.(string); ok {
					result = append(result, s)
				}
			}
			return result
		}
		if slice, ok := v.([]string); ok {
			return slice
		}
	}
	return nil
}

func copyMap(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}

func isIPAddress(address string) bool {
	// 简单检查是否为IP地址（IPv4）
	parts := strings.Split(address, ".")
	if len(parts) != 4 {
		return false
	}
	for _, part := range parts {
		if _, err := strconv.Atoi(part); err != nil {
			return false
		}
	}
	return true
}
