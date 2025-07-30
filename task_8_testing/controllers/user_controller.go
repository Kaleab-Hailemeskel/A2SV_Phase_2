package controllers

import (
	"fmt"
	"net/http"
	"task_8_testing/infrastructure"
	"task_8_testing/models"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userDB      models.IUserDataBase
	passService models.IPasswordService
	jwtHandler  models.IAuthentication
}

func NewUserController(userDataBase models.IUserDataBase, passwordService models.IPasswordService, jwt models.IAuthentication) *UserController {
	return &UserController{
		userDB:      userDataBase,
		passService: passwordService,
		jwtHandler:  jwt,
	}
}
func (us *UserController) Register(c *gin.Context) {
	var user UserDTO
	if c.ShouldBindBodyWithJSON(&user) != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Error": "Invalid User type"})
		return
	}
	userExists := us.userDB.CheckUserExistance(user.Email)
	if userExists {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Error": "User already exists"})
		return
	}
	err := us.userDB.StoreUser(changeUserDTO(&user))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "User registered successfully"})

}
func (us *UserController) Login(c *gin.Context) {
	var user UserDTO
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}
	userFromDB, err := us.userDB.FindUserByEmail(user.Email)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("user didn't exist, better register"))
		return
	}
	correctPassword := us.passService.IsCorrectPass(user.Password, userFromDB.Password)
	if !correctPassword {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("incorrect email or Password"))
		return
	}

	jwtBody := map[string]interface{}{
		"email": userFromDB.Email,
		"role":  userFromDB.Role,
	}

	securityToken, timeDuration := us.jwtHandler.GenerateSecurityToken(jwtBody)
	if securityToken == "" {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("can't generate jwt"))
		return
	}
	// sending jwt to client as cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(infrastructure.HEADER, securityToken, int((timeDuration).Seconds()), "", "", false, true) // the int(time.Now().Add(auth.TokenExpirationTime).Unix()) part could be a field of the jwtAuth structure

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
