package app

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type structValidator struct {
	validate *validator.Validate
}

func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

// 初始化 Fiber App 实例
func InitFiberApp() *fiber.App {
	return fiber.New(fiber.Config{
		StructValidator: &structValidator{validate: validator.New()},
	})
}

// StartFiberServer 启动 fiber 服务
func StartFiberServer(server *fiber.App) {
	server.Hooks().OnPreShutdown(func() error {
		// fiber 服务器退出的时候结束 asynq 服务
		CloseAsynqServer()
		return nil
	})

	if err := server.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
