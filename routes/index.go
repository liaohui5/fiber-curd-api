package routes

import (
	"fiber_curd_api/app"
	"fiber_curd_api/handlers"
	"fiber_curd_api/routes/middlewares"

	swaggerUI "github.com/gofiber/contrib/v3/swaggerui"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"
)

// InitMiddlewares 初始化全局中间件
func InitMiddlewares(server *fiber.App) {
	server.Use(cors.New())
	server.Use("/*", static.New("./static"))

	if app.IsDevMode() {
		server.Use(logger.New(logger.Config{
			Format:     "${cyan}[${time}] ${white}${pid} ${red}${status} ${blue}[${method}] ${white}${path}\n",
			TimeFormat: app.DATETIME_FMT,
			TimeZone:   "UTC",
		}))
	}
}

// InitSwaggerRoutes 初始化 swagger 插件路由
func InitSwaggerRoutes(server *fiber.App) {
	// rebuild swagger docs: swag init
	// swagger ui base path: /swagger/docs
	config := swaggerUI.Config{
		BasePath: "/swagger",
		FilePath: "./docs/swagger.json",
	}
	server.Use(swaggerUI.New(config))
}

// InitRoutes 初始化路由
func InitRoutes(server *fiber.App) {
	api := server.Group("/api")

	// 检查接口运行状态
	api.Get("/health", handlers.CheckHealth)

	// 认证模块
	api.Post("/login", handlers.Login)
	api.Post("/register", handlers.Register)
	api.Get("/refresh_access_token", handlers.RenewAccessToken)

	// 文章模块 CURD
	api.Get("/articles", middlewares.Auth, handlers.SearchArticles)
	api.Get("/articles/:id", middlewares.Auth, handlers.ArticleDetail)
	api.Post("/articles", middlewares.Auth, handlers.CreateArtilce)
	api.Patch("/articles/:id", middlewares.Auth, handlers.UpdateArticle)
	api.Delete("/articles/:id", middlewares.Auth, handlers.DeleteArticle)

	// 异步任务 CURD + asynq
	api.Get("/tasks", middlewares.Auth, handlers.SearchTasks)
	api.Get("/tasks/:id", middlewares.Auth, handlers.TaskDetail)
	api.Post("/tasks", middlewares.Auth, handlers.CreateTasks)
	api.Delete("/tasks/:id", middlewares.Auth, handlers.DeleteTasks)
	api.Get("/tasks/:id/retry", middlewares.Auth, handlers.RetryTasks)
}
