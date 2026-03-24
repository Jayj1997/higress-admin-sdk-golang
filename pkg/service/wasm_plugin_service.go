// Package service provides business services for the SDK
package service

import (
	"context"

	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
)

// WasmPluginService WASM插件管理服务接口
// 注意：此接口将在里程碑7中完整实现
type WasmPluginService interface {
	// List 列出WASM插件
	List(ctx context.Context, name, version string, builtIn *bool) ([]model.WasmPlugin, error)

	// Get 获取WASM插件详情
	Get(ctx context.Context, name string) (*model.WasmPlugin, error)

	// GetSchema 获取插件配置Schema
	GetSchema(ctx context.Context, name string) (map[string]interface{}, error)
}
