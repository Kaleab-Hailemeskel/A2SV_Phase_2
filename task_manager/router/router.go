package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"task_manager/controllers"
)
// GET, GET_ID, POST, PUT, DELETE methods are mapped with their counter func
func StartEngine(port_number string) {
	router := gin.Default()
	router.GET("/tasks", controllers.GetTasks)
	router.GET("/tasks/:id", controllers.GetTaskByID)
	router.POST("/tasks/", controllers.PostTask)
	router.PUT("/tasks/:id", controllers.PutTaskByID)
	router.DELETE("/tasks/:id", controllers.DeleteTaskByID)

	router.Run(fmt.Sprintf("localhost:%s", port_number)) // Listen and serve on port_number
	
}
