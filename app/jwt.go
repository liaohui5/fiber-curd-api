package app

import (
	"fiber_curd_api/tools/jwt"
)

// AppConfigAdapter 适配项目的 app.Config，实现 ConfigProvider 接口
type AppConfigAdapter struct {
	getFunc func(key string) any
}

func NewAppConfigAdapter(getFunc func(key string) any) *AppConfigAdapter {
	return &AppConfigAdapter{getFunc: getFunc}
}

func (a *AppConfigAdapter) AccessTokenExpiredTime() int64 {
	return a.getFunc("jwt.access_token_expire_time").(int64)
}

func (a *AppConfigAdapter) AccessTokenSecret() string {
	return a.getFunc("jwt.access_token_secret").(string)
}

func (a *AppConfigAdapter) RefreshTokenExpiredTime() int64 {
	return a.getFunc("jwt.refresh_token_expire_time").(int64)
}

func (a *AppConfigAdapter) RefreshTokenSecret() string {
	return a.getFunc("jwt.refresh_token_secret").(string)
}

var JWTService *jwt.JWTService

// InitJWTService 从全局 app.Config 创建 JWTService
func InitJWTService() {
	adapter := NewAppConfigAdapter(Config.Get)
	JWTService = jwt.NewJWTService(adapter)
}
