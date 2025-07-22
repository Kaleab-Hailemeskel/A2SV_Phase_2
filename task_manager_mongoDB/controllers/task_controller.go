package controllers

import (
	"net/http"
	"task_manager/data"
	"task_manager/models"

	"github.com/gin-gonic/gin"
)

func unableToConnectMassage(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Server Not Found"})
}
func DeleteTaskByID(ctx *gin.Context) {
	if !data.IsClientConnected() {
		unableToConnectMassage(ctx)
		return
	}
	requestID := ctx.Param("id") // get the id from the link parameter
	err := data.DeleteOne(requestID)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
	} else {
		ctx.IndentedJSON(http.StatusAccepted, gin.H{"Error": "Task with ID " + requestID + " got Deleted"})
	}

}
func PostTask(ctx *gin.Context) {
	var newtask models.Task
	if !data.IsClientConnected() {
		unableToConnectMassage(ctx)
		return
	}
	var err error
	if err = ctx.BindJSON(&newtask); err == nil { // check if there were no error while binding ctx BODY to newTask AND after that check if insertion went successful
		err = data.InsertOne(newtask)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.IndentedJSON(http.StatusCreated, newtask)
		return
	}
	ctx.IndentedJSON(http.StatusConflict, gin.H{"message": "Can't save a new task"}) // only excuted when there is a Binding Problem in the ctx.BindJSON()

}
func PutTaskByID(ctx *gin.Context) {
	if !data.IsClientConnected() {
		unableToConnectMassage(ctx)
		return
	}
	id := ctx.Param("id")

	var updatedTask models.Task

	if err := ctx.ShouldBindJSON(&updatedTask); err != nil { // type mismatch got handled
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := data.UpdateOne(id, updatedTask)
	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "Task updated"})
		return
	}
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // if no task with the same ID found send this message
}
func GetTasks(ctx *gin.Context) {
	if !data.IsClientConnected() {
		unableToConnectMassage(ctx)
		return
	}
	listOfTasks, err := data.FindALL()
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else if len(listOfTasks) == 0 {
		ctx.IndentedJSON(http.StatusOK, gin.H{"message": "No task Found"})
	} else {
		ctx.IndentedJSON(http.StatusOK, listOfTasks)
	}

}
func GetTaskByID(ctx *gin.Context) {
	if !data.IsClientConnected() {
		unableToConnectMassage(ctx)
		return
	}
	urlID := ctx.Param("id")
	task, err := data.FindByID(urlID)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else if task == nil {
		ctx.IndentedJSON(http.StatusOK, gin.H{"message": "No task Found with ID" + urlID})
	} else {
		ctx.IndentedJSON(http.StatusOK, *task)
	}
}
