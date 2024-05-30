package routes

import (
	"github.com/gin-gonic/gin"
	privilege_routers "main.go/routes/privilege"
	users_routers "main.go/routes/users"
	users_privilege_routers "main.go/routes/users_privilege"
)

func InitRoute(app *gin.Engine) {
	router := app
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	users_routers.SetupUserRoutes(router)
	privilege_routers.SetupPrivilegeRoutes(router)
	users_privilege_routers.SetupUserPrivilegeRoutes(router)
}
