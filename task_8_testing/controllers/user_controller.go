package controllers

import (
	"net/http"
	"task_8_testing/infrastructure"
	"task_8_testing/models"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	useCase models.IUseCase
}

func NewUserController(ucase models.IUseCase) *UserController {
	return &UserController{
		useCase: ucase,
	}
}
func (us *UserController) Register(c *gin.Context) {
	var user models.UserDTO
	if c.ShouldBindBodyWithJSON(&user) != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Error": "Invalid User type"})
		return
	}
	user.Role = models.USER
	// from usecase try to register a user
	if user, err := us.useCase.Register(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Error": err.Error()})
		return
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{
			"message": "User registered successfully",
			"user":    user,
		})
	}

}
func (us *UserController) Login(c *gin.Context) {
	var user models.UserDTO
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}
	accessToken, expTime, err := us.useCase.LoginHandler(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(infrastructure.HEADER, accessToken, int((expTime).Seconds()), "", "", false, true) // the int(time.Now().Add(auth.TokenExpirationTime).Unix()) part could be a field of the jwtAuth structure

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
