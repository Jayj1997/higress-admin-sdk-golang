// Package service provides business services for the SDK
package service

import (
	"context"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model"
)

// ConsumerService 消费者管理服务接口
type ConsumerService interface {
	// List 列出所有消费者
	List(ctx context.Context) ([]model.Consumer, error)

	// Get 获取消费者详情
	Get(ctx context.Context, name string) (*model.Consumer, error)

	// AddOrUpdate 添加或更新消费者
	AddOrUpdate(ctx context.Context, consumer *model.Consumer) (*model.Consumer, error)

	// Delete 删除消费者
	Delete(ctx context.Context, name string) error

	// ListAllowLists 列出所有允许列表
	ListAllowLists(ctx context.Context) ([]model.AllowList, error)

	// GetAllowList 获取指定目标的允许列表
	GetAllowList(ctx context.Context, targets map[model.WasmPluginInstanceScope]string) (*model.AllowList, error)

	// UpdateAllowList 更新允许列表
	UpdateAllowList(ctx context.Context, operation model.AllowListOperation, allowList *model.AllowList) error
}
