package models

import "github.com/gin-gonic/gin"

type IUserDataBase interface {
	FindUserByEmail(userEmail string) (*User, error)
	StoreUser(currUser *User) error
	CheckUserExistance(userEmail string) bool
	CloseDataBase() error
}

type ITaskDataBase interface {
	FindAllTasks(userEmail string) (*[]Task, error)
	FindByID(taskID string) (*Task, error)
	DeleteOne(taskID string) error
	UpdateOne(taskID string, updatedTask Task) error
	InsertOne(t Task) error
	CheckTaskExistance(taskID string) bool
	CloseDataBase() error
}

type IAuthentication interface {
	GetSecurityTokenFromClinet(ctx *gin.Context) (string, error)
	SendSecurityTokenToClinet(ctx *gin.Context, JWTBody map[string]interface{}) error
}

type IPasswordService interface {
	HashPassword(orginalPass string) (string, error)
	IsCorrectPass(orginalPass string, hashedPass string) bool
}
