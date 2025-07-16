package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/task", getTasks)
	router.GET("/task/:id", getTaskByID)
	router.POST("/task/", postTask)
	router.PUT("/task/:id", postTaskByID)
	router.DELETE("/task/:id", deleteTaskByID)
	

	router.Run() // Listen and serve on 0.0.0.0:8080
	fmt.Println("live server")
}

func deleteTaskByID(ctx *gin.Context) {
	requestID := ctx.Param("id")
	for index, task := range tasks {
		if task.ID == requestID {
			tasks = append(tasks[:index], tasks[index+1:]...)
			ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Task removed"})
			return
		}
	}
	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Task Not Found"})
}
func postTask(ctx *gin.Context) {
	var newtask Task
	if err := ctx.BindJSON(&newtask); err == nil {
		tasks = append(tasks, newtask)
		ctx.IndentedJSON(http.StatusCreated, newtask)
		return
	}
	ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Can't save a new task"})

}
func postTaskByID(ctx *gin.Context) {
	taskID := ctx.Param("id")
	var requestTask Task
	if bindErr := ctx.BindJSON(requestTask); bindErr != nil {
		ctx.IndentedJSON(http.StatusConflict, gin.H{"message": "Invalid task type"})
		return
	}
	for index, task := range tasks {
		if task.ID == taskID {
			// only couple of fields need to be edited
			if requestTask.Title != "" {
				tasks[index].Title = requestTask.Title
			}
			if requestTask.Description != "" {
				tasks[index].Description = requestTask.Description
			}
			if requestTask.Status != "" {
				tasks[index].Status = requestTask.Status
			}
			ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Task Updated!"})
			return
		}
	}
	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Task with ID %s not found", taskID)})
}
func getTasks(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, tasks)
}
func getTaskByID(ctx *gin.Context) {
	urlID := ctx.Param("id")
	for _, task := range tasks {
		if task.ID == urlID {
			ctx.IndentedJSON(http.StatusOK, task)
			return
		}
	}
	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Task with task id <%s> is not found", urlID)})
}
