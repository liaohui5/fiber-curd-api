package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户表模型
type User struct {
	ID        uint           `gorm:"primarykey"          json:"id"`                              // 主键ID
	CreatedAt time.Time      `gorm:""                    json:"created_at"`                      // 创建时间
	UpdatedAt time.Time      `gorm:""                    json:"updated_at"`                      // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index"               json:"deleted_at" swaggerignore:"true"` // 删除时间
	Username  string         `gorm:"size:64;"            json:"username"`                        // 用户名
	Telepone  string         `gorm:"size:16;"            json:"telephone"`                       // 用户手机号
	Email     string         `gorm:"size:64;unique;"     json:"email"`                           // 用户邮箱
	Password  string         `gorm:"size:128;"           json:"-"`                               // 用户密码(-json序列化时忽略字段)
	AvatarURL string         `gorm:"size:255"            json:"avatar_url"`                      // 用户头像URL
	Articles  []Article      `gorm:"foreignKey:UserID;"  json:"articles"`                        // 该用户发表的文章
}

// gorm.Model 带有默认约定字段(id/created_at/updated_at/delete_at)
// > 为什么不使用 gorm.Model 默认约定字段?
// 1. json 字段控制问题
// 如果使用默认的约定字段,则无法控制字段在返回 json 中的字段大小写
//
// 2. swagger 文档生成问题
// swag init 会忽略外部依赖包定义的 struct 导致报错
// swag init --parseDependency --parseInternal 则会导致文档生成变慢
//
// > swaggerignore 标签是什么?
// 在使用 swag init 时候忽略这个字段(因为这是做软删除功能,不应该展示到接口文档中)
