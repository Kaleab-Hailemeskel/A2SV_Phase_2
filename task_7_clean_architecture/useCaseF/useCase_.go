package useCaseF

import (
	"fmt"
	"task_7_clean_architecture/models"
)

type UseCase struct {
	userDataBase models.IUserDataBase
	taskDataBase models.ITaskDataBase
	JwtAuth      models.IAuthentication
	passService  models.IPasswordService
}

func NewUseCase(userDB models.IUserDataBase, taskDB models.ITaskDataBase, jwtAU models.IAuthentication, passSer models.IPasswordService) *UseCase {
	return &UseCase{
		userDataBase: userDB,
		taskDataBase: taskDB,
		JwtAuth:      jwtAU,
		passService:  passSer,
	}
}

func (uc *UseCase) Register(user *models.User) error {
	exists := uc.userDataBase.CheckUserExistance(user.Email)
	if exists {
		return fmt.Errorf("user already exists")
	}
	hashedPassword, err := uc.passService.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("can not hash the password")
	}
	user.Password = hashedPassword

	return uc.userDataBase.StoreUser(user)
}
func (uc *UseCase) LoginHandler(user *models.User) error {
	userFromDB, err := uc.userDataBase.FindUserByEmail(user.Email)
	if err != nil {
		return fmt.Errorf("user didn't exist, better register")
	}
	correctPassword := uc.passService.IsCorrectPass(user.Password, userFromDB.Password)
	if !correctPassword {
		return fmt.Errorf("incorrect email or Password")
	}

	return nil
}

func (uc *UseCase) GetUserWithEmail(userEmail string) (*models.User, error) {
	return uc.userDataBase.FindUserByEmail(userEmail)
}
func (uc *UseCase) DelteTask(requestID string) error {
	return uc.taskDataBase.DeleteOne(requestID)
}
func (uc *UseCase) CreatNewTask(newTask models.Task) error {
	return uc.taskDataBase.InsertOne(newTask)
}
func (uc *UseCase) EditTaskByID(taskID string, updatedTask models.Task) error {
	return uc.taskDataBase.UpdateOne(taskID, updatedTask)
}
func (uc *UseCase) GetAllTask(userEmail string) (*[]models.Task, error) {
	return uc.taskDataBase.FindAllTasks(userEmail)
}
func (uc *UseCase) GetTaskByID(taskID string) (*models.Task, error) {
	return uc.taskDataBase.FindByID(taskID)
}
func (uc *UseCase) CloseDBConnection() error {
	err := uc.userDataBase.CloseDataBase()
	if err != nil {
		return err
	}
	return uc.taskDataBase.CloseDataBase()
}
