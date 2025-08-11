package models

func ChangeUserDTO(userDto *UserDTO) *User {
	return &User{
		Email:    userDto.Email,
		Password: userDto.Password,
		Role:     userDto.Role,
	}
}

func ChangeUserModel(userDto *User) *UserDTO {
	return &UserDTO{
		Email:    userDto.Email,
		Password: userDto.Password,
		Role:     userDto.Role,
	}
}

func ChangeTaskDTO(taskDto *TaskDTO) *Task {
	return &Task{
		Status:      taskDto.Status,
		Description: taskDto.Description,
		Title:       taskDto.Title,
		DueDate:     taskDto.DueDate,
	}
}
func ChangeTaskModel(taskDto *Task) *TaskDTO {
	return &TaskDTO{
		Status:      taskDto.Status,
		Description: taskDto.Description,
		Title:       taskDto.Title,
		DueDate:     taskDto.DueDate,
	}
}
