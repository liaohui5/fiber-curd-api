package middlewares

import (
	"fiber_curd_api/app"
	"fiber_curd_api/tools/fmtres"
	"strings"

	"github.com/gofiber/fiber/v3"
)

// 登录认证中间件
func Auth(c fiber.Ctx) error {
	// 获取 Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fmtres.ErrorStr("please login first"))
	}

	// 验证 jwt 是否正确 authorization => Bearer <token>
	accessToken := strings.TrimPrefix(authHeader, "Bearer ")
	app.Logger.Debug("Auth middleware got accessToken: " + accessToken)

	uid, err := app.JWTService.VerifyAccessToken(accessToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fmtres.ErrorFmt("invalid accessToken", err))
	}

	c.Set("uid", string(uid))

	return c.Next()
}
