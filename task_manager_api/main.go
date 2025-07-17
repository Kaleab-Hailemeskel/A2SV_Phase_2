package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/tasks", getTasks)
	router.GET("/tasks/:id", getTaskByID)
	router.POST("/tasks/", postTask)
	router.PUT("/tasks/:id", putTaskByID)
	router.DELETE("/tasks/:id", deleteTaskByID)

	router.Run("localhost:8080") // Listen and serve on 0.0.0.0:8080
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
	_foundDuplicate := func(t Task) bool { // todo: a simple func that check if there is already a task with the same ID.
		for _, each_task := range tasks {
			if each_task.ID == t.ID {
				return true
			}
		}
		return false
	}
	if err := ctx.BindJSON(&newtask); err == nil && !_foundDuplicate(newtask) {
		tasks = append(tasks, newtask)
		ctx.IndentedJSON(http.StatusCreated, newtask)
		return
	}
	ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Can't save a new task"})

}
func putTaskByID(ctx *gin.Context) {
	id := ctx.Param("id")

	var updatedTask Task

	if err := ctx.ShouldBindJSON(&updatedTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			if updatedTask.Title != "" {
				tasks[i].Title = updatedTask.Title
			}
			if updatedTask.Description != "" {
				tasks[i].Description = updatedTask.Description
			}
			ctx.JSON(http.StatusOK, gin.H{"message": "Task updated"})
			return
		}
	}

	ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
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
	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Task with task id <%s> is not found", string(urlID))})
}
