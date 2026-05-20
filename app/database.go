package app

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 数据库链接对象
var Connection *gorm.DB

// ConnectDB 链接数据库
func ConnectDB() *gorm.DB {
	if Connection != nil {
		return Connection
	}

	dsn := Config.Get("database.dsn").(string)
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 禁止使用物理外键(仅支持逻辑外键即可)
	})
	if err != nil {
		panic("连接数据库失败")
	}

	Connection = db
	return db
}

// InitDatabase 初始化数据库链接
func InitDatabase() {
	ConnectDB()
}
