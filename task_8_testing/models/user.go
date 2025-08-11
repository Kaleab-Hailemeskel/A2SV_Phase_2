package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	ADMIN = "ADMIN"
	USER  = "USER"
)

type User struct { // other fields make in this struct will only make it just unnessarly board2
	Email    string
	Password string
	Role     string
}

type UserDTO struct {
	ID       primitive.ObjectID `json:"id" bson:"_id" binding:"required"`
	Email    string             `json:"email" bson:"email" binding:"required,email"`
	Password string             `json:"password" bson:"password" binding:"required,min=8,max=20"`
	Role     string             `json:"role" bson:"role"`
}
