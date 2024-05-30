package bootstrap

import (
	"github.com/gin-gonic/gin"
	"main.go/config"
	"main.go/config/app_config"
	"main.go/database"
	"main.go/routes"
)

func BootstrapApp() {
	app := gin.Default()
	routes.InitRoute(app)
	config.InitConfig()
	database.ConnectDatabase()

	app.Run(app_config.APP_PORT)
}
