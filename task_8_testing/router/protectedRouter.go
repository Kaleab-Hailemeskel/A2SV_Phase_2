package router

import (
	"task_8_testing/controllers"
	"task_8_testing/middleware"

	"github.com/gin-gonic/gin"
)

func StartProtectedRouter(router *gin.Engine, userAuth *middleware.UserAuth, controllers *controllers.TaskController, control *controllers.UserController) {

	router.GET("/whoAmI", userAuth.Authentication, control.GiveMeMyInfo)
	taskNeedAuthentication := router.Group("/tasks")
	{
		taskNeedAuthentication.Use(userAuth.Authentication)

		taskNeedAuthentication.POST("", controllers.PostTask)
		taskNeedAuthentication.GET("", controllers.GetTasks)

		tasksGroup := taskNeedAuthentication.Group("/:id")
		{
			tasksGroup.Use(userAuth.Authorization)

			tasksGroup.GET("", controllers.GetTaskByID)
			tasksGroup.PUT("", controllers.PutTaskByID)
			tasksGroup.DELETE("", controllers.DeleteTaskByID)
		}
	}
}
