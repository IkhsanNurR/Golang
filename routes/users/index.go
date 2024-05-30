package users_routers

import (
	"github.com/gin-gonic/gin"
	users_controller "main.go/controllers/users"
)

func SetupUserRoutes(router *gin.Engine) {
	router.GET("/users", users_controller.GetAllUsers)
	router.GET("/users/:id", users_controller.GetOneUsers)
	router.POST("/users", users_controller.CreateOneUsers)
	router.PUT("/users/:id", users_controller.UpdateOneUsers)
	router.DELETE("/users/:id", users_controller.DeleteOneUsers)
}
