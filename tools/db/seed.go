package db

import (
	"fmt"
	"math/rand/v2"

	"fiber_curd_api/app"
	"fiber_curd_api/models"
)

// MaxSeedLimit 最大填充个数
const MaxSeedLimit = 30

// Seed 数据库填充
func Seed() {
	CreateUsers()
	CreateArticles()
	fmt.Println("===数据库填充完成===")
}

// CreateUsers 给 users 表填充一些假数据
func CreateUsers() {
	adminPasswd := "$2a$10$0Nw/9fpYupjK9pYyXzellO.ZRwDR1scBZdcGrQj/harVIbrsdjgZe"
	var users []models.User
	for i := 1; i <= MaxSeedLimit; i++ {
		users = append(users, models.User{
			Email:    fmt.Sprintf("test%d@test.com", i),
			Password: adminPasswd, // 123456 -> md5 -> bcrypt(10)
		})
	}

	// 固定一个数据的参数方便调试
	users[0].Email = "admin@example.com"

	app.ConnectDB().CreateInBatches(&users, len(users))
}

// CreateArticles 给 articles 表填充一些数据
func CreateArticles() {
	var articles []models.Article
	for i := 1; i <= MaxSeedLimit; i++ {
		articles = append(articles, models.Article{
			Title:    fmt.Sprintf("title-%d", i),
			Contents: fmt.Sprintf("contents-%d", i),
			UserID:   uint(rand.Uint32N(MaxSeedLimit)),
		})
	}
	app.ConnectDB().CreateInBatches(&articles, len(articles))
}
