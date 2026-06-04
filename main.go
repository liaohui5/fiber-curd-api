package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"fiber_curd_api/app"
	"fiber_curd_api/cli"
	"fiber_curd_api/routes"
)

// @title fiber-app
// @version 1.0
// @description This is an API for practice fiber
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api
func main() {
	if cli.Run() {
		return
	}

	// go app.StartAsynqServer()

	fiberApp := app.InitFiberApp()
	routes.InitMiddlewares(fiberApp)
	routes.InitRoutes(fiberApp)
	routes.InitSwaggerRoutes(fiberApp)
	go app.StartFiberServer(fiberApp)

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, syscall.SIGTERM)
	<-exitChan

	app.CloseAsynqClient()
	app.CloseAsynqServer()
	app.ShutdownFiberServer(fiberApp)

	fmt.Println(">>> Server stopped <<<")
}
