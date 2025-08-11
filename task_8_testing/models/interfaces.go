package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IUserDataBase interface {
	CloseDataBase() error
	FindUserByID(userID primitive.ObjectID) (*UserDTO, error)
	FindUserByEmail(userEmail string) (*UserDTO, error)
	StoreUser(user *User) (*UserDTO, error)
	CheckUserExistance(userEmail string) bool
	CheckUserExistanceByID(userID primitive.ObjectID) bool
}

type ITaskDataBase interface {
	CloseDataBase() error
	CheckTaskExistance(taskID primitive.ObjectID) bool
	FindAllTasks(userEmail string) ([]*TaskDTO, error)
	FindByID(taskID primitive.ObjectID) (*TaskDTO, error)
	DeleteOne(taskID primitive.ObjectID) error
	UpdateOne(taskID primitive.ObjectID, updatedTask *TaskDTO) (*TaskDTO, error)
	InsertOne(t *TaskDTO) (*TaskDTO, error)
}

type IAuthentication interface {
	TokenExpired(token *jwt.Token) (bool, error)
	ParseToken(tokenString string) (*jwt.Token, error)
	GetUserEmailFromSecurityToken(token *jwt.Token) (string, error)
	GenerateSecurityToken(JWTBody map[string]interface{}) (string, time.Duration)
	GetUserID(token *jwt.Token) (*primitive.ObjectID, error)
}
type IUserAuth interface {
	Authentication(ctx *gin.Context)
}
type IPasswordService interface {
	HashPassword(orginalPass string) (string, error)
	IsCorrectPass(orginalPass string, hashedPass string) bool
}

type IUseCase interface {
	Register(user *UserDTO) (*UserDTO, error)
	LoginHandler(user *UserDTO) (string, *time.Duration, error)
	GetUserWithID(userID string) (*UserDTO, error)
	GetUserWithEmail(userEmail string) (*UserDTO, error)
	DeleteTask(requestID, userEmail string) error
	CreatNewTask(newTask *TaskDTO) (*TaskDTO, error)
	EditTaskByID(taskID, userEmail string, updatedTask *TaskDTO) (*TaskDTO, error)
	GetAllTask(userEmail string) ([]*TaskDTO, error)
	GetTaskByID(taskID, userEmail string) (*TaskDTO, error)
	CheckOwnership(taskID primitive.ObjectID, userEmail string) error
	UpdateTask(taskID, userEmail string, task *TaskDTO) (*TaskDTO, error)
	CloseALLDBConnection() error
}
