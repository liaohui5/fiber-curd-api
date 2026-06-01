package models

import (
	"time"

	"gorm.io/gorm"
)

// Article 文章表模型
type Article struct {
	ID        uint           `gorm:"primarykey"         json:"id"`                              // 主键ID
	CreatedAt time.Time      `gorm:""                   json:"created_at"`                      // 创建时间
	UpdatedAt time.Time      `gorm:""                   json:"updated_at"`                      // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index"              json:"deleted_at" swaggerignore:"true"` // 删除时间
	Title     string         `gorm:"size:255"           json:"title"      validate:"required"`  // 文章标题
	UserID    uint           `gorm:""                   json:"user_id"    validate:"required"`  // 作者ID
	Contents  string         `gorm:"type:text"          json:"contents"   validate:"required"`  // 文章内容
	Likes     uint           `gorm:"default:0"          json:"likes"      validate:"number"`    // 点赞数
	Stars     uint           `gorm:"default:0"          json:"stars"      validate:"number"`    // 收藏数
	Shares    uint           `gorm:"default:0"          json:"shares"     validate:"number"`    // 转发(分享)数
	Comments  uint           `gorm:"default:0"          json:"comments"   validate:"number"`    // 评论数
	Author    *User          `gorm:"foreignKey:UserID;" json:"author"`                          // 文章作者信息
}
