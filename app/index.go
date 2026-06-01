package app

import "sync"

var once sync.Once

func init() {
	once.Do(func() {
		InitConfig()     // 初始化配置模块
		InitLogger()     // 初始化日志模块
		InitDatabase()   // 初始化数据库模块
		InitJWTService() // 初始化 JWT
	})
}

// IsDevMode 是否为开发模式
func IsDevMode() bool {
	return Config.Get("app.debug").(bool)
}
