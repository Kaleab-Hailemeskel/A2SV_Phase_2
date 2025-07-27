package router

import (
	"task_7_clean_architecture/controllers"
	"task_7_clean_architecture/middleware"

	"github.com/gin-gonic/gin"
)

func StartProtectedRouter(router *gin.Engine, userAuth *middleware.UserAuth, controllers *controllers.UserController) {

	router.GET("/whoAmI", userAuth.Authentication, controllers.GiveMeMyInfo)
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
