package router

import (
	"task_7_clean_architecture/controllers"

	"github.com/gin-gonic/gin"
)

func StartPublicRouter(router *gin.Engine, controller *controllers.UserController) {
	router.POST("/register", controller.Register)
	router.POST("/login", controller.Login)
}
