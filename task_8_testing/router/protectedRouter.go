package router

import (
	"task_8_testing/controllers"
	"task_8_testing/models"

	"github.com/gin-gonic/gin"
)

func StartProtectedRouter(router *gin.Engine, userAuth models.IUserAuth, controllers *controllers.TaskController, control *controllers.UserController) {

	needAuthentication := router.Group("")
	{
		needAuthentication.Use(userAuth.Authentication)

		needAuthentication.GET("/whoAmI", control.GiveMeMyInfo)

		tasksGroup := needAuthentication.Group("/tasks")
		{
			tasksGroup.POST("", controllers.PostTask)
			tasksGroup.GET("", controllers.GetTasks)
			taskIDGroup := tasksGroup.Group("/:id")
			{
				taskIDGroup.GET("", controllers.GetTaskByID)
				taskIDGroup.PUT("", controllers.PutTaskByID)
				taskIDGroup.DELETE("", controllers.DeleteTaskByID)
			}
		}

	}
}
