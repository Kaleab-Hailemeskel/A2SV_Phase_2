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
}

func NewUserAuth(jwt models.IAuthentication, userDataBase models.IUserDataBase, taskDataBase models.ITaskDataBase) *UserAuth {
	return &UserAuth{
		jwtHandler: jwt,
		userDB:     userDataBase,
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
	UserID, err := uc.jwtHandler.GetUserID(jwtToken)

	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	userResult, err := uc.userDB.FindUserByID(*UserID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "user wasn't found in DataBase"})
		return
	}

	ctx.Set(infrastructure.CURR_USER, *userResult)
	ctx.Status(http.StatusOK)
	ctx.Next()
}

