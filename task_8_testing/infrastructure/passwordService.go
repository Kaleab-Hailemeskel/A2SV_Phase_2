package infrastructure

import (
	"task_8_testing/models"

	"golang.org/x/crypto/bcrypt"
)

type passwordService struct {
	passCost int
}

func NewPasswordService() models.IPasswordService {

	return &passwordService{
		passCost: bcrypt.DefaultCost,
	}
}

func (passService *passwordService) HashPassword(orginalPass string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(orginalPass), passService.passCost)
	return string(hashedPassword), err
}

func (passService *passwordService) IsCorrectPass(orginalPass string, hashedPass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(orginalPass)) == nil
}
