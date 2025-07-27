package models

import "time"

// model for task
type Task struct {
	ID          string    `json:"id" binding:"required"`
	OwnerEmail  string    `json:"ownerEmail" binding:"required,email"`   
	Title       string    `json:"title" binding:"required,min=3,max=100"`
	Description string    `json:"description" binding:"required"`        
	DueDate     time.Time `json:"due_date" binding:"required"`           
	Status      string    `json:"status" binding:"required"`             
}
