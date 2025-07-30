package models

import "time"

// model for task
type Task struct {
	ID          string    
	OwnerEmail  string    
	Title       string    
	Description string    
	DueDate     time.Time 
	Status      string    
}
