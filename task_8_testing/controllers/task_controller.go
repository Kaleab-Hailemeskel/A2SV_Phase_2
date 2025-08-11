package controllers

import (
	"log"
	"net/http"

	"task_8_testing/infrastructure"
	"task_8_testing/models"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	useCase models.IUseCase
}

func NewTaskController(newuc models.IUseCase) *TaskController {
	return &TaskController{
		useCase: newuc,
	}
}

func (tc *TaskController) DeleteTaskByID(ctx *gin.Context) {
	requestID := ctx.Param("id") // get the id from the link parameter
	currUser, exists := ctx.Get(infrastructure.CURR_USER)
	if !exists {
		ctx.IndentedJSON(http.StatusConflict, gin.H{"Error": "This should never happen in MILLION YEARS !!!"})
		return
	}
	err := tc.useCase.DeleteTask(requestID, currUser.(models.User).Email)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
	} else {
		ctx.IndentedJSON(http.StatusAccepted, gin.H{"message": "Task with ID " + requestID + " got Deleted"})
	}

}
func (tc *TaskController) PostTask(ctx *gin.Context) {
	currUser, exists := ctx.Get(infrastructure.CURR_USER)

	if !exists {
		ctx.IndentedJSON(http.StatusConflict, gin.H{"Error": "This should never happen in MILLION YEARS !!!"})
		return
	}

	var task models.TaskDTO
	if err := ctx.ShouldBindJSON(&task); err == nil { // check if there were no error while binding ctx BODY to task AND after that check if insertion went stc cessful
		task.OwnerEmail = currUser.(models.User).Email
		res, err := tc.useCase.CreatNewTask(&task)

		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.IndentedJSON(http.StatusCreated, *models.ChangeTaskDTO(res))
	}
	ctx.IndentedJSON(http.StatusConflict, gin.H{"message": "Can't save a new task"}) // only excuted when there is a Binding Problem in the ctx.BindJSON()

}
func (tc *TaskController) PutTaskByID(ctx *gin.Context) {
	id := ctx.Param("id")
	currUser, exists := ctx.Get(infrastructure.CURR_USER)
	if !exists {
		ctx.IndentedJSON(http.StatusConflict, gin.H{"Error": "This should never happen in MILLION YEARS !!!"})
		return
	}
	var updatedTask models.TaskDTO

	if err := ctx.ShouldBindJSON(&updatedTask); err != nil { // type mismatch got handled
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedRes, err := tc.useCase.UpdateTask(id, currUser.(models.User).Email, &updatedTask)
	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message":      "Task updated",
			"updated_task": updatedRes,
		})
		return
	}
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // if no task with the same ID found send this message
}
func (tc *TaskController) GetTasks(ctx *gin.Context) {
	userResult, exist := ctx.Get(infrastructure.CURR_USER)

	if !exist {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "unblivable"})
		log.Fatal("user Wasn't found while Getting Tasks in GetTasks")
	}

	var listOfTasks []*models.TaskDTO
	var err error

	if userResult.(models.User).Role == models.ADMIN {
		listOfTasks, err = tc.useCase.GetAllTask("")
	} else {
		listOfTasks, err = tc.useCase.GetAllTask(userResult.(models.User).Email)

	}

	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else if len(listOfTasks) == 0 {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "No task Found"})
	} else {
		ctx.IndentedJSON(http.StatusOK, listOfTasks)
	}

}
func (tc *TaskController) GetTaskByID(ctx *gin.Context) {

	urlID := ctx.Param("id")
	if urlID == "" {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "no id found"})
	}
	currUser, exists := ctx.Get(infrastructure.CURR_USER)
	if !exists {
		ctx.IndentedJSON(http.StatusConflict, gin.H{"Error": "This should never happen in MILLION YEARS !!!"})
		return
	}

	task, err := tc.useCase.GetTaskByID(urlID, currUser.(models.User).Email)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else if task == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "No task Found with ID" + urlID})
	} else {
		ctx.IndentedJSON(http.StatusOK, *task)
	}

}
