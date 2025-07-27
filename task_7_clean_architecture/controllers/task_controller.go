package controllers

import (
	"log"
	"net/http"

	"task_7_clean_architecture/infrastructure"
	"task_7_clean_architecture/models"
	"task_7_clean_architecture/useCaseF"
	"time"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUseCase useCaseF.UseCase
}

func NewUserController(useCase *useCaseF.UseCase) *UserController {
	return &UserController{
		userUseCase: *useCase,
	}
}

type UserDTO struct {
	ID          string    `json:"id" binding:"required"`
	Title       string    `json:"title" binding:"required,min=3,max=100"`
	Description string    `json:"description" binding:"required"`
	DueDate     time.Time `json:"due_date" binding:"required"`
	Status      string    `json:"status" binding:"required"`
}

func (uc *UserController) DeleteTaskByID(ctx *gin.Context) {
	requestID := ctx.Param("id") // get the id from the link parameter
	err := uc.userUseCase.DelteTask(requestID)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
	} else {
		ctx.IndentedJSON(http.StatusAccepted, gin.H{"message": "Task with ID " + requestID + " got Deleted"})
	}

}
func (uc *UserController) PostTask(ctx *gin.Context) {
	currUser, exists := ctx.Get(infrastructure.CURR_USER)

	if !exists {
		ctx.IndentedJSON(http.StatusConflict, gin.H{"Error": "This should never happen in MILLION YEARS !!!"})
		return
	}

	var tempTask UserDTO
	if err := ctx.BindJSON(&tempTask); err == nil { // check if there were no error while binding ctx BODY to tempTask AND after that check if insertion went successful
		newTask := models.Task{
			ID:          tempTask.ID,
			OwnerEmail:  currUser.(models.User).Email,
			DueDate:     tempTask.DueDate,
			Status:      tempTask.Status,
			Description: tempTask.Description,
			Title:       tempTask.Title}

		err = uc.userUseCase.CreatNewTask(newTask)

		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.IndentedJSON(http.StatusCreated, newTask)
	}
	ctx.IndentedJSON(http.StatusConflict, gin.H{"message": "Can't save a new task"}) // only excuted when there is a Binding Problem in the ctx.BindJSON()

}
func (uc *UserController) PutTaskByID(ctx *gin.Context) {
	id := ctx.Param("id")

	var updatedTask models.Task

	if err := ctx.ShouldBindJSON(&updatedTask); err != nil { // type mismatch got handled
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := uc.userUseCase.EditTaskByID(id, updatedTask)
	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "Task updated"})
		return
	}
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // if no task with the same ID found send this message
}
func (uc *UserController) GetTasks(ctx *gin.Context) {
	userResult, exist := ctx.Get(infrastructure.CURR_USER)
	if !exist {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "unblivable"})
		log.Fatal("user Wasn't found while Getting Tasks in GetTasks")
	}

	var listOfTasks *[]models.Task
	var err error

	if userResult.(models.User).Role == models.ADMIN {
		listOfTasks, err = uc.userUseCase.GetAllTask("")
	} else {
		listOfTasks, err = uc.userUseCase.GetAllTask(userResult.(models.User).Email)

	}

	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else if len(*listOfTasks) == 0 {
		ctx.IndentedJSON(http.StatusOK, gin.H{"message": "No task Found"})
	} else {
		ctx.IndentedJSON(http.StatusOK, *listOfTasks)
	}

}
func (uc *UserController) GetTaskByID(ctx *gin.Context) {

	urlID := ctx.Param("id")
	task, err := uc.userUseCase.GetTaskByID(urlID)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else if task == nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"message": "No task Found with ID" + urlID})
	} else {
		ctx.IndentedJSON(http.StatusOK, *task)
	}

}
func (us *UserController) Register(c *gin.Context) {
	var user models.User
	if c.ShouldBindBodyWithJSON(&user) != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Error": "Invalid User type"})
		return
	}
	err := us.userUseCase.Register(&user)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "User registered successfully"})

}
func (us *UserController) Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}
	err := us.userUseCase.LoginHandler(&user)

	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	jwtBody := map[string]interface{}{
		"email": user.Email,
		"role":  user.Role,
	}
	securityToken, timeDuration := us.userUseCase.JwtAuth.GenerateSecurityToken(jwtBody)
	// sending jwt to client as cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(infrastructure.CURR_USER, securityToken, int(time.Now().Add(*timeDuration).Unix()), "", "", false, true) // the int(time.Now().Add(auth.TokenExpirationTime).Unix()) part could be a field of the jwtAuth structure

}
func (us *UserController) GiveMeMyInfo(c *gin.Context) {
	userResult, exists := c.Get(infrastructure.CURR_USER)
	if !exists {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "IMPOSSIBLEEEEE"})
		return
	}
	var realUser = userResult.(models.User)
	realUser.Password = "**HIDDEN**"
	c.IndentedJSON(http.StatusAccepted, gin.H{"Current User": realUser})
}
