package privilege_routers

import (
	"github.com/gin-gonic/gin"
	privilege_controller "main.go/controllers/privilege"
)

func SetupPrivilegeRoutes(router *gin.Engine) {
	router.GET("/privilege", privilege_controller.GetAllPrivilege)
	router.GET("/privilege/:id", privilege_controller.GetOnePrivilege)
	router.POST("/privilege", privilege_controller.CreateOnePrivilege)
	router.PUT("/privilege/:id", privilege_controller.UpdateOnePrivilege)
	router.DELETE("/privilege/:id", privilege_controller.DeleteOnePrivilege)
}
