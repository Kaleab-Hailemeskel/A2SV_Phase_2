package middleware

import (
	"fmt"
	"net/http"
	"task-6_authentication_and_authorization/data"
	"task-6_authentication_and_authorization/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var JWTSECRET = []byte("what123321")

func Authentication(c *gin.Context) {

	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signging method %v", token.Header["alg"])
		}
		return JWTSECRET, nil
	})
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) >= claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusRequestTimeout)
		}

		userEmail := claims["email"].(string)
		userResult, err := data.FindOneUser(userEmail)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "user wasn't found in DataBase"})
			return
		}
		c.Set("currUser", *userResult)

		c.Next()
		return
	}
	c.AbortWithStatus(http.StatusUnauthorized)
}
func IsAuthorizedUserForTaskManipulation(ctx *gin.Context) {
	// check if the user is the rightful owner OR he is an Admin
	userResult, exists := ctx.Get("currUser")
	// since a user that doesn't exist in the database even won't pass the
	// authorization, here I don't have to
	// check for the existance. But for formality

	if !exists {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"Error": "This part wasn't supposed to be sent in any curcumstance"})
		return
	}
	var user models.User = userResult.(models.User)

	requestID := ctx.Param("id") // get the id from the link parameter
	if user.Role == models.ADMIN || requestID == "" {

		ctx.Next()
		return
	}

	taskResult, err := data.FindByID(requestID)

	if !exists || err != nil || (user.Email != taskResult.OwnerEmail) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": "Unautorized access"})
		return
	}

	ctx.Next()
}
