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
	log.Println("🥈 CurrUSER", infrastructure.CURR_USER)
	currUser, exists := ctx.Get(infrastructure.CURR_USER)
	log.Println("✅ ", currUser)
	if !exists {
		ctx.IndentedJSON(http.StatusConflict, gin.H{"Error": "This should never happen in MILLION YEARS !!! while deleteing by task ID"})
		return
	}
	err := tc.useCase.DeleteTask(requestID, currUser.(models.UserDTO).Email)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
	} else {
		ctx.IndentedJSON(http.StatusAccepted, gin.H{"message": "Task with ID " + requestID + " got Deleted"})
	}

}
func (tc *TaskController) PostTask(ctx *gin.Context) {
	currUser, exists := ctx.Get(infrastructure.CURR_USER)

	if !exists {
		ctx.IndentedJSON(http.StatusConflict, gin.H{"Error": "This should never happen in MILLION YEARS !!! post task in controller" })
		return
	}

	var task models.TaskDTO
	if err := ctx.ShouldBindJSON(&task); err == nil { // check if there were no error while binding ctx BODY to task AND after that check if insertion went stc cessful
		task.OwnerEmail = currUser.(models.UserDTO).Email
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
		ctx.IndentedJSON(http.StatusConflict, gin.H{"Error": "This should never happen in MILLION YEARS !!! while => put in controller"})
		return
	}
	var updatedTask models.TaskDTO

	if err := ctx.ShouldBindJSON(&updatedTask); err != nil { // type mismatch got handled
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedRes, err := tc.useCase.UpdateTask(id, currUser.(models.UserDTO).Email, &updatedTask)
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
	log.Println("✅ BEFORE____----")
	log.Println("📩 ", userResult.(models.UserDTO).Email)
	if userResult.(models.UserDTO).Role == models.ADMIN {
		listOfTasks, err = tc.useCase.GetAllTask("")
	} else {
		listOfTasks, err = tc.useCase.GetAllTask(userResult.(models.UserDTO).Email)
	}
	log.Println("✅ GET ALL")

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
		ctx.IndentedJSON(http.StatusConflict, gin.H{"Error": "This should never happen in MILLION YEARS !!! while Get Task By ID controller "})
		return
	}

	task, err := tc.useCase.GetTaskByID(urlID, currUser.(models.UserDTO).Email)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else if task == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "No task Found with ID" + urlID})
	} else {
		ctx.IndentedJSON(http.StatusOK, *task)
	}

}
