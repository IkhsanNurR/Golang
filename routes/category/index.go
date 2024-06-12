package category_routers

import (
	"github.com/gin-gonic/gin"
	category_controller "main.go/controllers/category"
)

func SetupCategoryRoute(router *gin.Engine) {
	router.GET("/category", category_controller.GetAllCategory)
	router.GET("/category/:id", category_controller.GetOneCategory)
	router.POST("/category", category_controller.CreateOneCategory)
	router.PUT("/category/:id", category_controller.UpdateOneCategory)
	router.DELETE("/category/:id", category_controller.DeleteOneCategory)
}
