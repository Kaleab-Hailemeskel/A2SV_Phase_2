package useCaseF

import (
	"fmt"
	"task_8_testing/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UseCase struct {
	userDataBase models.IUserDataBase
	taskDataBase models.ITaskDataBase
	JwtAuth      models.IAuthentication
	passService  models.IPasswordService
}

func NewUseCase(userDB models.IUserDataBase, taskDB models.ITaskDataBase, jwtAU models.IAuthentication, passSer models.IPasswordService) models.IUseCase {
	return &UseCase{
		userDataBase: userDB,
		taskDataBase: taskDB,
		JwtAuth:      jwtAU,
		passService:  passSer,
	}
}

func (uc *UseCase) Register(user *models.UserDTO) (*models.UserDTO, error) {
	exists := uc.userDataBase.CheckUserExistance(user.Email)
	if exists {
		return nil, fmt.Errorf("user already exists")
	}
	hashedPassword, err := uc.passService.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("can not hash the password")
	}
	user.Password = hashedPassword

	return uc.userDataBase.StoreUser(models.ChangeUserDTO(user))
}
func (uc *UseCase) LoginHandler(user *models.UserDTO) (string, *time.Duration, error) {
	userFromDB, err := uc.userDataBase.FindUserByEmail(user.Email)
	if err != nil {
		return "", nil, fmt.Errorf("user didn't exist, better register")
	}
	correctPassword := uc.passService.IsCorrectPass(user.Password, userFromDB.Password)
	if !correctPassword {
		return "", nil, fmt.Errorf("incorrect email or Password")
	}

	jwtBody := map[string]interface{}{
		"email": userFromDB.Email,
		"role":  userFromDB.Role,
		"id":    userFromDB.ID,
	}

	securityToken, expTime := uc.JwtAuth.GenerateSecurityToken(jwtBody)
	if securityToken == "" {
		return "", nil, fmt.Errorf("can't generate jwt")
	}

	return securityToken, &expTime, nil

}
func (uc *UseCase) GetUserWithID(userID string) (*models.UserDTO, error) {
	if objID, err := primitive.ObjectIDFromHex(userID); err != nil {
		return nil, fmt.Errorf("invalid user id")
	} else {
		return uc.userDataBase.FindUserByID(objID)
	}
}
func (uc *UseCase) GetUserWithEmail(userEmail string) (*models.UserDTO, error) {
	return uc.userDataBase.FindUserByEmail(userEmail)
}
func (uc *UseCase) DeleteTask(requestID, userEmail string) error {
	taskID, err := primitive.ObjectIDFromHex(requestID)
	if err != nil {
		return fmt.Errorf("invalid taskID")
	}

	err = uc.CheckOwnership(taskID, userEmail)
	if err != nil {
		return err
	}
	return uc.taskDataBase.DeleteOne(taskID)
}
func (uc *UseCase) CreatNewTask(newTask *models.TaskDTO) (*models.TaskDTO, error) {
	return uc.taskDataBase.InsertOne(newTask)
}
func (uc *UseCase) EditTaskByID(taskID, userEmail string, updatedTask *models.TaskDTO) (*models.TaskDTO, error) {

	taskIDObj, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return nil, fmt.Errorf("invalid taskID")
	}

	err = uc.CheckOwnership(taskIDObj, userEmail)
	if err != nil {
		return nil, err
	}
	return uc.taskDataBase.UpdateOne(taskIDObj, updatedTask)
}
func (uc *UseCase) GetAllTask(userEmail string) ([]*models.TaskDTO, error) {
	return uc.taskDataBase.FindAllTasks(userEmail)
}
func (uc *UseCase) GetTaskByID(taskID, userEmail string) (*models.TaskDTO, error) {
	taskIDObj, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return nil, fmt.Errorf("invalid taskID")
	}
	err = uc.CheckOwnership(taskIDObj, userEmail)
	if err != nil {
		return nil, err
	}
	return uc.taskDataBase.FindByID(taskIDObj)
}
func (uc *UseCase) CheckOwnership(taskID primitive.ObjectID, userEmail string) error {
	task, err := uc.taskDataBase.FindByID(taskID)
	if err != nil {
		return err
	}
	if task.OwnerEmail != userEmail {
		return fmt.Errorf("access Denied")
	}
	return nil
}
func (uc *UseCase) UpdateTask(taskID, userEmail string, task *models.TaskDTO) (*models.TaskDTO, error) {
	// first check user ownership
	taskIDObj, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return nil, fmt.Errorf("invalid taskID")
	}
	err = uc.CheckOwnership(taskIDObj, userEmail)
	if err != nil {
		return nil, err
	}
	return uc.taskDataBase.UpdateOne(taskIDObj, task)
}
func (uc *UseCase) CloseALLDBConnection() error {
	err := uc.userDataBase.CloseDataBase()
	if err != nil {
		return err
	}
	return uc.taskDataBase.CloseDataBase()
}
