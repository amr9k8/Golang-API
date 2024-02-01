package pkg

import (
	"github.com/gin-gonic/gin"
	"test/pkg/api/routes"
)

type App struct {
	*gin.Engine
}

func NewApp() *App {
	app := &App{Engine: gin.Default()}
	return app
}

func (app *App) Run(address string) {
	routes.SetupRoutes(app.Engine)
	app.Engine.Run(address)
}
