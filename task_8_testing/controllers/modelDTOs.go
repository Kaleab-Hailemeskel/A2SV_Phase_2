package controllers

import (
	"task_8_testing/models"
	"time"
)

type UserDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=20"`
}

type TaskDTO struct {
	ID          string    `json:"id" binding:"required"`
	Title       string    `json:"title" binding:"required,min=3,max=100"`
	Description string    `json:"description" binding:"required"`
	DueDate     time.Time `json:"due_date" binding:"required"`
	Status      string    `json:"status" binding:"required"`
}

func changeUserDTO(userDto *UserDTO) *models.User {
	return &models.User{
		Email: userDto.Email,
		Password: userDto.Password,
		Role:  models.USER, // this always assign the current registree to a User type
	}
}

func changeTaskDTO(taskDto *TaskDTO) *models.Task {
	return &models.Task{
		ID:          taskDto.ID,
		Status:      taskDto.Status,
		Description: taskDto.Description,
		Title:       taskDto.Title,
		DueDate:     taskDto.DueDate,
	}
}
