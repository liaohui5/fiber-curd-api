package db

import (
	"fiber_curd_api/app"
	"fiber_curd_api/models"
	"fmt"
)

// 数据库迁移
func Migrate() {
	app.ConnectDB().AutoMigrate(&models.User{}, &models.Article{}, &models.Task{})
	fmt.Println("===数据库迁移完成===")
}
