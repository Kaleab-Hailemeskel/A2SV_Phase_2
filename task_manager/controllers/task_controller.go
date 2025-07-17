package controllers

import (
	"fmt"
	"net/http"
	"task_manager/data"
	"github.com/gin-gonic/gin"
	"task_manager/models"
)

func DeleteTaskByID(ctx *gin.Context) {
	requestID := ctx.Param("id") // get the id from the link parameter
	for index, task := range data.Tasks {
		if task.ID == requestID { // if a match is found delete it and send a comfirmation json file with 200 OK
			data.Tasks = append(data.Tasks[:index], data.Tasks[index+1:]...)
			ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Task removed"})
			return
		}
	}
	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Task Not Found"}) // not found, send not found json message
}
func PostTask(ctx *gin.Context) {
	var newtask models.Task
	_foundDuplicate := func(t models.Task) bool { //: a simple func that check if there is already a task with the same ID.
		for _, each_task := range data.Tasks {
			if each_task.ID == t.ID {
				return true
			}
		}
		return false
	}
	if err := ctx.BindJSON(&newtask); err == nil && !_foundDuplicate(newtask) { // check if the json request's type match with Task type and if that is true, then check if there isn't a task with the same ID number
		data.Tasks = append(data.Tasks, newtask)
		ctx.IndentedJSON(http.StatusCreated, newtask)
		return
	}
	ctx.IndentedJSON(http.StatusConflict, gin.H{"message": "Can't save a new task"}) // if there is a type mismatch or duplicate found this message will be sent to client

}
func PutTaskByID(ctx *gin.Context) {
	id := ctx.Param("id")

	var updatedTask models.Task

	if err := ctx.ShouldBindJSON(&updatedTask); err != nil { // type mismatch got handled
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, task := range data.Tasks {
		if task.ID == id { // if there is a match update the Title, Description then send a comfimation massage with 200 Ok status code
			if updatedTask.Title != "" {
				data.Tasks[i].Title = updatedTask.Title
			}
			if updatedTask.Description != "" {
				data.Tasks[i].Description = updatedTask.Description
			}
			ctx.JSON(http.StatusOK, gin.H{"message": "Task updated"})
			return
		}
	}

	ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"}) // if no task with the same ID found send this message
}
func GetTasks(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, data.Tasks) // just a simple get method
}
func GetTaskByID(ctx *gin.Context) {
	urlID := ctx.Param("id")
	for _, task := range data.Tasks {
		if task.ID == urlID {// search for a task specified by the taskID and if there is a match send the task to client
			ctx.IndentedJSON(http.StatusOK, task)
			return
		}
	}
	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Task with task id <%s> is not found", string(urlID))}) // if the task didn't exist sent this one
}
