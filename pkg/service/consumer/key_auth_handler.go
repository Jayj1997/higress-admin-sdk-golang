package consumer

import (
	"strings"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model"
)

// Key Auth配置常量
const (
	KeyAuthConfigConsumers           = "consumers"
	KeyAuthConfigConsumerName        = "name"
	KeyAuthConfigConsumerCredentials = "credentials"
	KeyAuthConfigConsumerCredential  = "credential" // 已废弃，保留用于兼容
	KeyAuthConfigKeys                = "keys"
	KeyAuthConfigInHeader            = "in_header"
	KeyAuthConfigInQuery             = "in_query"
	KeyAuthConfigAllow               = "allow"
	KeyAuthConfigGlobalAuth          = "global_auth"
)

// 插件名称常量
const (
	PluginNameKeyAuth = "key-auth"
)

// Bearer Token前缀
const bearerTokenPrefix = "Bearer "

// KeyAuthCredentialHandler Key Auth凭证处理器
type KeyAuthCredentialHandler struct{}

// NewKeyAuthCredentialHandler 创建Key Auth凭证处理器
func NewKeyAuthCredentialHandler() *KeyAuthCredentialHandler {
	return &KeyAuthCredentialHandler{}
}

// GetType 返回凭证类型
func (h *KeyAuthCredentialHandler) GetType() string {
	return model.CredentialTypeKeyAuth
}

// GetPluginName 返回关联的WASM插件名称
func (h *KeyAuthCredentialHandler) GetPluginName() string {
	return PluginNameKeyAuth
}

// IsConsumerInUse 检查消费者是否正在使用
func (h *KeyAuthCredentialHandler) IsConsumerInUse(consumerName string, instances []*model.WasmPluginInstance) bool {
	if len(instances) == 0 {
		return false
	}
	for _, instance := range instances {
		if instance.Configurations == nil {
			continue
		}
		allowObj, ok := instance.Configurations[KeyAuthConfigAllow]
		if !ok {
			continue
		}
		allowList, ok := allowObj.([]interface{})
		if !ok {
			continue
		}
		for _, item := range allowList {
			if str, ok := item.(string); ok && str == consumerName {
				return true
			}
		}
	}
	return false
}

// ExtractConsumers 从插件实例中提取消费者列表
func (h *KeyAuthCredentialHandler) ExtractConsumers(instance *model.WasmPluginInstance) []*model.Consumer {
	if instance == nil || instance.Configurations == nil {
		return []*model.Consumer{}
	}

	consumersObj, ok := instance.Configurations[KeyAuthConfigConsumers]
	if !ok {
		return []*model.Consumer{}
	}

	consumerList, ok := consumersObj.([]interface{})
	if !ok {
		return []*model.Consumer{}
	}

	consumers := make([]*model.Consumer, 0, len(consumerList))
	for _, consumerObj := range consumerList {
		consumerMap, ok := consumerObj.(map[string]interface{})
		if !ok {
			continue
		}
		consumer := h.extractConsumer(consumerMap)
		if consumer != nil {
			consumers = append(consumers, consumer)
		}
	}
	return consumers
}

// InitDefaultGlobalConfigs 初始化默认全局配置
func (h *KeyAuthCredentialHandler) InitDefaultGlobalConfigs(instance *model.WasmPluginInstance) {
	if instance.Configurations == nil {
		instance.Configurations = make(map[string]interface{})
	}

	if _, ok := instance.Configurations[KeyAuthConfigGlobalAuth]; !ok {
		instance.Configurations[KeyAuthConfigGlobalAuth] = false
	}
	if _, ok := instance.Configurations[KeyAuthConfigAllow]; !ok {
		instance.Configurations[KeyAuthConfigAllow] = []string{}
	}
	// 添加一个虚拟key，因为插件要求至少有一个全局key
	if _, ok := instance.Configurations[KeyAuthConfigKeys]; !ok {
		instance.Configurations[KeyAuthConfigKeys] = []string{"x-higress-dummy-key"}
	}
	if _, ok := instance.Configurations[KeyAuthConfigConsumers]; !ok {
		instance.Configurations[KeyAuthConfigConsumers] = []interface{}{}
	}
}

// SaveConsumer 保存消费者到插件实例
func (h *KeyAuthCredentialHandler) SaveConsumer(instance *model.WasmPluginInstance, consumer *model.Consumer) bool {
	if len(consumer.Credentials) == 0 {
		return false
	}

	// 查找KeyAuth凭证
	var keyAuthCredential *model.KeyAuthCredential
	for _, cred := range consumer.Credentials {
		if kc, ok := cred.(*model.KeyAuthCredential); ok {
			keyAuthCredential = kc
			break
		}
	}

	if keyAuthCredential == nil {
		return h.DeleteConsumer(instance, consumer.Name)
	}

	if instance.Configurations == nil {
		h.InitDefaultGlobalConfigs(instance)
	}

	// 获取或创建consumers列表
	consumersObj, ok := instance.Configurations[KeyAuthConfigConsumers]
	if !ok {
		consumersObj = []interface{}{}
	}

	consumers, ok := consumersObj.([]interface{})
	if !ok {
		consumers = []interface{}{}
	}

	// 查找现有消费者配置
	var consumerConfig map[string]interface{}
	for _, c := range consumers {
		cMap, ok := c.(map[string]interface{})
		if !ok {
			continue
		}
		existingConsumer := h.extractConsumer(cMap)
		if existingConsumer == nil {
			continue
		}
		if existingConsumer.Name == consumer.Name {
			consumerConfig = cMap
		} else if h.hasSameCredential(existingConsumer, keyAuthCredential) {
			panic("Key auth credential already in use by consumer: " + existingConsumer.Name)
		}
	}

	if consumerConfig == nil {
		consumerConfig = map[string]interface{}{
			KeyAuthConfigConsumerName: consumer.Name,
		}
		consumers = append(consumers, consumerConfig)
	} else {
		keyAuthCredential = h.mergeExistedConfig(keyAuthCredential, consumerConfig)
	}

	// 验证凭证
	if err := keyAuthCredential.Validate(false); err != nil {
		panic(err.Error())
	}

	source := model.ParseKeyAuthCredentialSource(keyAuthCredential.Source)
	if source == "" {
		panic("Invalid key auth credential source: " + keyAuthCredential.Source)
	}

	key := keyAuthCredential.Key
	credentials := keyAuthCredential.Values

	switch source {
	case model.KeyAuthCredentialSourceBearer, model.KeyAuthCredentialSourceHeader:
		consumerConfig[KeyAuthConfigInHeader] = true
		consumerConfig[KeyAuthConfigInQuery] = false
		if source == model.KeyAuthCredentialSourceBearer {
			key = "Authorization"
			bearerCredentials := make([]string, len(credentials))
			for i, c := range credentials {
				bearerCredentials[i] = bearerTokenPrefix + c
			}
			credentials = bearerCredentials
		}
	case model.KeyAuthCredentialSourceQuery:
		consumerConfig[KeyAuthConfigInHeader] = false
		consumerConfig[KeyAuthConfigInQuery] = true
	}

	consumerConfig[KeyAuthConfigKeys] = []string{key}
	consumerConfig[KeyAuthConfigConsumerCredentials] = credentials
	delete(consumerConfig, KeyAuthConfigConsumerCredential)

	instance.Configurations[KeyAuthConfigConsumers] = consumers
	instance.Configurations[KeyAuthConfigGlobalAuth] = false

	return true
}

// DeleteConsumer 从插件实例删除消费者
func (h *KeyAuthCredentialHandler) DeleteConsumer(globalInstance *model.WasmPluginInstance, consumerName string) bool {
	if globalInstance == nil || globalInstance.Configurations == nil {
		return false
	}

	consumersObj, ok := globalInstance.Configurations[KeyAuthConfigConsumers]
	if !ok {
		return false
	}

	consumers, ok := consumersObj.([]interface{})
	if !ok {
		return false
	}

	deleted := false
	newConsumers := make([]interface{}, 0, len(consumers))
	for _, c := range consumers {
		cMap, ok := c.(map[string]interface{})
		if !ok {
			newConsumers = append(newConsumers, c)
			continue
		}
		name, _ := cMap[KeyAuthConfigConsumerName].(string)
		if name == consumerName {
			deleted = true
		} else {
			newConsumers = append(newConsumers, c)
		}
	}

	if deleted {
		globalInstance.Configurations[KeyAuthConfigConsumers] = newConsumers
	}
	return deleted
}

// GetAllowedConsumers 获取允许的消费者列表
func (h *KeyAuthCredentialHandler) GetAllowedConsumers(instance *model.WasmPluginInstance) []string {
	if instance == nil || instance.Configurations == nil {
		return []string{}
	}

	allowObj, ok := instance.Configurations[KeyAuthConfigAllow]
	if !ok {
		return []string{}
	}

	allowList, ok := allowObj.([]interface{})
	if !ok {
		return []string{}
	}

	result := make([]string, 0, len(allowList))
	for _, item := range allowList {
		if str, ok := item.(string); ok {
			result = append(result, str)
		}
	}
	return result
}

// UpdateAllowList 更新允许列表
func (h *KeyAuthCredentialHandler) UpdateAllowList(operation model.AllowListOperation, instance *model.WasmPluginInstance, consumerNames []string) {
	if instance.Configurations == nil {
		instance.Configurations = make(map[string]interface{})
	}

	newAllowList := h.GetAllowedConsumers(instance)

	switch operation {
	case model.AllowListOperationAdd:
		for _, name := range consumerNames {
			found := false
			for _, existing := range newAllowList {
				if existing == name {
					found = true
					break
				}
			}
			if !found {
				newAllowList = append(newAllowList, name)
			}
		}
	case model.AllowListOperationRemove:
		for _, name := range consumerNames {
			for i, existing := range newAllowList {
				if existing == name {
					newAllowList = append(newAllowList[:i], newAllowList[i+1:]...)
					break
				}
			}
		}
	case model.AllowListOperationReplace:
		newAllowList = consumerNames
	case model.AllowListOperationToggleOnly:
		if len(newAllowList) == 0 {
			newAllowList = []string{}
		}
	}

	instance.Configurations[KeyAuthConfigAllow] = newAllowList
}

// extractConsumer 从配置map中提取消费者
func (h *KeyAuthCredentialHandler) extractConsumer(consumerMap map[string]interface{}) *model.Consumer {
	if consumerMap == nil {
		return nil
	}

	name, _ := consumerMap[KeyAuthConfigConsumerName].(string)
	if name == "" {
		return nil
	}

	credential := h.parseCredential(consumerMap)
	if credential == nil {
		return nil
	}

	return &model.Consumer{
		Name:        name,
		Credentials: []model.Credential{credential},
	}
}

// parseCredential 从配置map中解析凭证
func (h *KeyAuthCredentialHandler) parseCredential(consumerMap map[string]interface{}) *model.KeyAuthCredential {
	if consumerMap == nil {
		return nil
	}

	keyObj, ok := consumerMap[KeyAuthConfigKeys]
	if !ok {
		return nil
	}

	keyList, ok := keyObj.([]interface{})
	if !ok || len(keyList) == 0 {
		return nil
	}

	var key string
	for _, keyItem := range keyList {
		if str, ok := keyItem.(string); ok && str != "" {
			key = str
			break
		}
	}

	if key == "" {
		return nil
	}

	inHeader, _ := consumerMap[KeyAuthConfigInHeader].(bool)
	inQuery, _ := consumerMap[KeyAuthConfigInQuery].(bool)

	credentials := []string{}
	if credObj, ok := consumerMap[KeyAuthConfigConsumerCredentials]; ok {
		if credList, ok := credObj.([]interface{}); ok {
			for _, c := range credList {
				if str, ok := c.(string); ok {
					credentials = append(credentials, str)
				}
			}
		}
	}

	// 兼容旧的单凭证字段
	if cred, ok := consumerMap[KeyAuthConfigConsumerCredential].(string); ok && cred != "" {
		found := false
		for _, c := range credentials {
			if c == cred {
				found = true
				break
			}
		}
		if !found {
			credentials = append(credentials, cred)
		}
	}

	var source model.KeyAuthCredentialSource
	if inHeader {
		if key == "Authorization" && len(credentials) > 0 {
			allBearer := true
			for _, c := range credentials {
				if !strings.HasPrefix(c, bearerTokenPrefix) {
					allBearer = false
					break
				}
			}
			if allBearer {
				source = model.KeyAuthCredentialSourceBearer
				key = ""
				for i, c := range credentials {
					credentials[i] = strings.TrimSpace(strings.TrimPrefix(c, bearerTokenPrefix))
				}
			} else {
				source = model.KeyAuthCredentialSourceHeader
			}
		} else {
			source = model.KeyAuthCredentialSourceHeader
		}
	} else if inQuery {
		source = model.KeyAuthCredentialSourceQuery
	} else {
		return nil
	}

	return model.NewKeyAuthCredential(string(source), key, credentials)
}

// hasSameCredential 检查是否有相同的凭证
func (h *KeyAuthCredentialHandler) hasSameCredential(existingConsumer *model.Consumer, credential *model.KeyAuthCredential) bool {
	if credential == nil || existingConsumer == nil {
		return false
	}

	var existingCredential *model.KeyAuthCredential
	for _, cred := range existingConsumer.Credentials {
		if kc, ok := cred.(*model.KeyAuthCredential); ok {
			existingCredential = kc
			break
		}
	}

	if existingCredential == nil {
		return false
	}

	if !strings.EqualFold(credential.Source, existingCredential.Source) {
		return false
	}
	if credential.Key != existingCredential.Key {
		return false
	}
	if len(credential.Values) == 0 || len(existingCredential.Values) == 0 {
		return false
	}

	for _, v := range credential.Values {
		for _, ev := range existingCredential.Values {
			if v == ev {
				return true
			}
		}
	}
	return false
}

// mergeExistedConfig 合并现有配置
func (h *KeyAuthCredentialHandler) mergeExistedConfig(credential *model.KeyAuthCredential, consumerConfig map[string]interface{}) *model.KeyAuthCredential {
	existedCredential := h.parseCredential(consumerConfig)
	if existedCredential == nil {
		return credential
	}

	source := credential.Source
	if source == "" {
		source = existedCredential.Source
	}

	key := credential.Key
	if key == "" {
		key = existedCredential.Key
	}

	values := credential.Values
	if len(values) == 0 {
		values = existedCredential.Values
	}

	return model.NewKeyAuthCredential(source, key, values)
}
