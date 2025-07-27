package middleware

import (
	"net/http"
	"task_7_clean_architecture/infrastructure"
	"task_7_clean_architecture/models"
	useCase "task_7_clean_architecture/useCaseF"

	"github.com/gin-gonic/gin"
)

type UserAuth struct {
	userUseCase *useCase.UseCase
}

func NewUserAuth(userUseCase_ *useCase.UseCase) *UserAuth {
	return &UserAuth{
		userUseCase: userUseCase_,
	}
}
func (uc *UserAuth) Authentication(ctx *gin.Context) {

	userEmail, err := uc.userUseCase.JwtAuth.GetSecurityTokenFromClinet(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userResult, err := uc.userUseCase.GetUserWithEmail(userEmail)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "user wasn't found in DataBase"})
		return
	}

	ctx.Set(infrastructure.CURR_USER, *userResult)

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

	taskResult, err := uc.userUseCase.GetTaskByID(requestID)

	if err != nil || (user.Email != taskResult.OwnerEmail) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": "Unautorized access"})
		return
	}

	ctx.Next()
}
