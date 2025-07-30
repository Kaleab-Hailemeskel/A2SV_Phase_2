package middleware

import (
	"net/http"
	"task_8_testing/infrastructure"
	"task_8_testing/models"

	"github.com/gin-gonic/gin"
)

type UserAuth struct {
	jwtHandler models.IAuthentication
	userDB     models.IUserDataBase
	taskDB     models.ITaskDataBase
}

func NewUserAuth(jwt models.IAuthentication, userDataBase models.IUserDataBase, taskDataBase models.ITaskDataBase) *UserAuth {
	return &UserAuth{
		jwtHandler: jwt,
		userDB:     userDataBase,
		taskDB:     taskDataBase,
	}
}

func (uc *UserAuth) Authentication(ctx *gin.Context) {

	tokenString, err := ctx.Cookie(infrastructure.HEADER)
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	jwtToken, err := uc.jwtHandler.ParseToken(tokenString)
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	tokenExpired, err := uc.jwtHandler.TokenExpired(jwtToken)
	if tokenExpired {
		ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	userEmail, err := uc.jwtHandler.GetUserEmailFromSecurityToken(jwtToken)

	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	userResult, err := uc.userDB.FindUserByEmail(userEmail)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "user wasn't found in DataBase"})
		return
	}

	ctx.Set(infrastructure.CURR_USER, *userResult)
	ctx.Status(http.StatusOK)
	ctx.Next()
}
func (uc *UserAuth) Authorization(ctx *gin.Context) {
	userResult, exists := ctx.Get(infrastructure.CURR_USER)

	if !exists {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Error": "Panic if this appears"})
		return
	}

	var user models.User = userResult.(models.User)
	requestID := ctx.Param("id") // get the id from the link parameter
	if user.Role == models.ADMIN || requestID == "" {
		ctx.Next()
		return
	}

	taskResult, err := uc.taskDB.FindByID(requestID)

	if err != nil || (user.Email != taskResult.OwnerEmail) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": "Unautorized access"})
		return
	}

	ctx.Status(http.StatusOK)
	ctx.Next()
}
