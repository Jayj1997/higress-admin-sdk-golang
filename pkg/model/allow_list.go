package model

// AllowListOperation 允许列表操作类型
type AllowListOperation string

const (
	// AllowListOperationAdd 添加消费者到允许列表
	// 认证开关如果非空则更新
	AllowListOperationAdd AllowListOperation = "ADD"
	// AllowListOperationRemove 从允许列表移除消费者
	// 认证开关如果非空则更新
	AllowListOperationRemove AllowListOperation = "REMOVE"
	// AllowListOperationReplace 替换整个允许列表
	// 认证开关如果非空则更新
	AllowListOperationReplace AllowListOperation = "REPLACE"
	// AllowListOperationToggleOnly 仅切换认证开关，保持允许列表不变
	AllowListOperationToggleOnly AllowListOperation = "TOGGLE_ONLY"
)

// AllowList 消费者允许列表
type AllowList struct {
	// Targets 目标范围 (作用域 -> 目标名称)
	Targets map[WasmPluginInstanceScope]string `json:"targets,omitempty"`
	// AuthEnabled 是否启用认证
	AuthEnabled *bool `json:"authEnabled,omitempty"`
	// CredentialTypes 凭证类型列表
	CredentialTypes []string `json:"credentialTypes,omitempty"`
	// ConsumerNames 消费者名称列表
	ConsumerNames []string `json:"consumerNames,omitempty"`
}

// NewAllowList 创建允许列表
func NewAllowList() *AllowList {
	return &AllowList{
		Targets:         make(map[WasmPluginInstanceScope]string),
		CredentialTypes: []string{},
		ConsumerNames:   []string{},
	}
}

// ForTarget 创建指定目标的允许列表构建器
func ForTarget(scope WasmPluginInstanceScope, target string) *AllowList {
	allowList := NewAllowList()
	allowList.Targets[scope] = target
	return allowList
}
