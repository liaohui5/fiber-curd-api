package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	TaskStatusPending    = 0 // 等待中
	TaskStatusProcessing = 1 // 处理中
	TaskStatusFinished   = 2 // 已完成
	TaskStatusFailed     = 3 // 已失败
)

type Task struct {
	ID           uint           `gorm:"primarykey"         json:"id"`                              // 主键ID
	Status       uint           `gorm:"default:0"          json:"status"`                          // 任务状态
	Name         string         `gorm:""                   json:"name"       validate:"required"`  // 任务名称
	MetaData     map[string]any `gorm:"serializer:json"    json:"meta_data"  validate:"required"`  // 执行任务需要的元数据
	FailedReason string         `gorm:"default:''"         json:"failed_reason"`                   // 任务失败原因
	CreatedAt    time.Time      `gorm:""                   json:"created_at"`                      // 创建时间
	UpdatedAt    time.Time      `gorm:""                   json:"updated_at"`                      // 更新时间
	DeletedAt    gorm.DeletedAt `gorm:"index"              json:"deleted_at" swaggerignore:"true"` // 删除时间
}
