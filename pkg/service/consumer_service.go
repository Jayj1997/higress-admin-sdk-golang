// Package service provides business services for the SDK
package service

import (
	"context"

	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
)

// ConsumerService 消费者管理服务接口
// 注意：此接口将在里程碑9中完整实现
type ConsumerService interface {
	// List 列出所有消费者
	List(ctx context.Context) ([]model.Consumer, error)

	// Get 获取消费者详情
	Get(ctx context.Context, name string) (*model.Consumer, error)

	// Add 添加消费者
	Add(ctx context.Context, consumer *model.Consumer) (*model.Consumer, error)

	// Update 更新消费者
	Update(ctx context.Context, consumer *model.Consumer) (*model.Consumer, error)

	// Delete 删除消费者
	Delete(ctx context.Context, name string) error
}
