package router

import (
	"task_8_testing/controllers"

	"github.com/gin-gonic/gin"
)

func StartPublicRouter(router *gin.Engine, controller *controllers.UserController) {
	router.POST("/register", controller.Register)
	router.POST("/login", controller.Login)
	
}
