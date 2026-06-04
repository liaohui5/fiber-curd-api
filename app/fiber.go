package app

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type structValidator struct {
	validate *validator.Validate
}

func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

// InitFiberApp 初始化 Fiber App 实例
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

// ShutdownFiberServer 关闭 fiber 服务
func ShutdownFiberServer(server *fiber.App) {
	if err := server.Shutdown(); err != nil {
		log.Fatal("[Shutdown] error: ", err)
		os.Exit(1)
	}
}
