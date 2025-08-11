package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// model for task
type Task struct {
	OwnerEmail  string
	Title       string
	Description string
	DueDate     time.Time
	Status      string
}

type TaskDTO struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id" binding:"required"`
	Title       string             `json:"title" bson:"title" binding:"required,min=3,max=100"`
	Description string             `json:"description" bson:"description" binding:"required"`
	OwnerEmail  string             `json:"owner_email" bson:"owner_email" binding:"required"`
	DueDate     time.Time          `json:"due_date" bson:"due_date" binding:"required"`
	Status      string             `json:"status" bson:"status" binding:"required"`
}
