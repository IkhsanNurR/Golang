package users_privilege_routers

import (
	"github.com/gin-gonic/gin"
	users_privilege_controller "main.go/controllers/users_privilege"
)

func SetupUserPrivilegeRoutes(router *gin.Engine) {
	router.GET("/users_privilege", users_privilege_controller.GetAllUsersPrivilege)
	router.GET("/users_privilege/:id", users_privilege_controller.GetOneUsersPrivilege)
	router.GET("/users_privilege/users/:id", users_privilege_controller.GetUsersPrivilegeByUser)
	router.POST("/users_privilege", users_privilege_controller.CreateOneUsers)
	router.PUT("/users_privilege/:id", users_privilege_controller.UpdateOneUsers)
	router.DELETE("/users_privilege/:id", users_privilege_controller.DeleteOneUsers)
}
