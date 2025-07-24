package controllers

import (
	"net/http"
	"task-6_authentication_and_authorization/data"
	"task-6_authentication_and_authorization/middleware"
	"task-6_authentication_and_authorization/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
		ctx.IndentedJSON(http.StatusAccepted, gin.H{"message": "Task with ID " + requestID + " got Deleted"})
	}

}
func PostTask(ctx *gin.Context) {
	// I neeed to use an inplace temporary struct because,
	// unless it will cause a trouble with the real Task struct
	// having a required email while marshalling and unmarshalling
	var tempTask struct {
		ID          string    `json:"id" binding:"required"`
		Title       string    `json:"title" binding:"required,min=3,max=100"`
		Description string    `json:"description" binding:"required"`
		DueDate     time.Time `json:"due_date" binding:"required"`
		Status      string    `json:"status" binding:"required"`
	}
	if !data.IsClientConnected() {
		unableToConnectMassage(ctx)
		return
	}
	currUser, exists := ctx.Get("currUser")
	if !exists {
		ctx.IndentedJSON(http.StatusConflict, gin.H{"Error": "This should never happen in MILLION YEARS !!!"})
		return
	}

	var err error

	if err = ctx.BindJSON(&tempTask); err == nil { // check if there were no error while binding ctx BODY to tempTask AND after that check if insertion went successful
		newTask := models.Task{
			ID:          tempTask.ID,
			OwnerEmail:  currUser.(models.User).Email,
			DueDate:     tempTask.DueDate,
			Status:      tempTask.Status,
			Description: tempTask.Description,
			Title:       tempTask.Title}
		err = data.InsertOne(newTask)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.IndentedJSON(http.StatusCreated, newTask)
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

	userResult, _ := ctx.Get("currUser")

	var listOfTasks *[]models.Task
	var err error

	if userResult.(models.User).Role == models.ADMIN {
		listOfTasks, err = data.FindALL("")
	} else {
		listOfTasks, err = data.FindALL(userResult.(models.User).Email)
	}

	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else if len(*listOfTasks) == 0 {
		ctx.IndentedJSON(http.StatusOK, gin.H{"message": "No task Found"})
	} else {
		ctx.IndentedJSON(http.StatusOK, *listOfTasks)
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
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else if task == nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"message": "No task Found with ID" + urlID})
	} else {
		ctx.IndentedJSON(http.StatusOK, *task)
	}

}


func Register(c *gin.Context) {
	var user models.User
	if c.ShouldBindBodyWithJSON(&user) != nil || (user.Email == "" || user.Password == "" || user.Role == "") {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Error": "Invalid User type"})
		return
	}

	//? checking for Existing User should be writen on data/user_service.go BUT for the time being
	_, err := data.FindOneUser(user.Email)
	if err == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "User Already Exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": "Error Happened while hashing"})
		return
	}
	user.Password = string(hashedPassword)
	err = data.InsertOneUser(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": "Unable to Register a user"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "User registered successfully"})

}
func Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	// searching user existance in database
	existingUser, err := data.FindOneUser(user.Email)

	//? This Side wasn't working because of the `json : "-"` marshaling takes place without storing the hashed password
	if (err != nil) || (bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)) != nil) {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": existingUser.Email,
		"exp":   time.Now().Add(time.Hour).Unix(),
	})

	jwtToken, err := token.SignedString(middleware.JWTSECRET)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", jwtToken, 3600*3, "", "", false, true)
	c.JSON(200, gin.H{"message": "Cookies were sent"})

}
func GiveMeMyInfo(c *gin.Context) {
	userResult, exists := c.Get("currUser")
	if !exists {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "IMPOSSIBLEEEEE"})
		return
	}
	var realUser = userResult.(models.User)
	realUser.Password = "**HIDDEN**"
	c.IndentedJSON(http.StatusAccepted, gin.H{"Current User": realUser})
}
