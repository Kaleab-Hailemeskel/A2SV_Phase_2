package infrastructure

import (
	"task_7_clean_architecture/models"

	"golang.org/x/crypto/bcrypt"
)

type PasswordService struct {
	passCost int
}

func NewPasswordService() models.IPasswordService {

	return &PasswordService{
		passCost: bcrypt.DefaultCost,
	}
}

func (passService *PasswordService) HashPassword(orginalPass string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(orginalPass), passService.passCost)
	return string(hashedPassword), err
}

func (passService *PasswordService) IsCorrectPass(orginalPass string, hashedPass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(orginalPass)) == nil
}
