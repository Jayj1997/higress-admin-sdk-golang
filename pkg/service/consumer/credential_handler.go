// Package consumer provides consumer management services
package consumer

import (
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
)

// CredentialHandler 凭证处理器接口
// 负责处理特定类型凭证的验证、存储和提取
type CredentialHandler interface {
	// GetType 返回凭证类型
	GetType() string

	// GetPluginName 返回关联的WASM插件名称
	GetPluginName() string

	// IsConsumerInUse 检查消费者是否正在使用
	// instances: 所有插件实例列表
	IsConsumerInUse(consumerName string, instances []*model.WasmPluginInstance) bool

	// ExtractConsumers 从插件实例中提取消费者列表
	ExtractConsumers(instance *model.WasmPluginInstance) []*model.Consumer

	// InitDefaultGlobalConfigs 初始化默认全局配置
	InitDefaultGlobalConfigs(instance *model.WasmPluginInstance)

	// SaveConsumer 保存消费者到插件实例
	// 返回true表示实例已修改需要保存
	SaveConsumer(instance *model.WasmPluginInstance, consumer *model.Consumer) bool

	// DeleteConsumer 从插件实例删除消费者
	// 返回true表示实例已修改需要保存
	DeleteConsumer(globalInstance *model.WasmPluginInstance, consumerName string) bool

	// GetAllowedConsumers 获取允许的消费者列表
	GetAllowedConsumers(instance *model.WasmPluginInstance) []string

	// UpdateAllowList 更新允许列表
	UpdateAllowList(operation model.AllowListOperation, instance *model.WasmPluginInstance, consumerNames []string)
}
