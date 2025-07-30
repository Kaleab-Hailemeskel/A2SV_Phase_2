package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)


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
	TokenExpired(token *jwt.Token) (bool, error)
	ParseToken(tokenString string) (*jwt.Token, error)
	GetUserEmailFromSecurityToken(token *jwt.Token) (string, error)
	GenerateSecurityToken(JWTBody map[string]interface{}) (string, time.Duration)
}

type IPasswordService interface {
	HashPassword(orginalPass string) (string, error)
	IsCorrectPass(orginalPass string, hashedPass string) bool
}
