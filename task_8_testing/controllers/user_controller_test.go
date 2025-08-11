package controllers_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"task_8_testing/controllers"
	"task_8_testing/models"
	models_mocks "task_8_testing/models/mocks"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UnitUserControllerTest struct {
	suite.Suite
	mockUseCase *models_mocks.MockIUseCase
	controller  *controllers.UserController
	router      *gin.Engine
}

var validTestCases = []struct {
	userFromDB   *models.UserDTO
	originalPass string
	jsonBody     string
}{
	
	{&models.UserDTO{
		Email:    "user_one@domain.com",
		Role:     models.USER,
		Password: "original_Pass_user1",
	},
		"original_Pass_user1", // Assuming this is the unhashed password
		`{
		"email": "user_one@domain.com",
		"password": "original_Pass_user1",
		"role": "USER"
	}`},
}
var invalidTestCases = []struct {
	userFromDB   *models.UserDTO
	originalPass string
	jsonBody     string
}{
	{
		&models.UserDTO{
			Email:    "test@example.com",
			Role:     models.USER,
			Password: "hashed_test_pass",
		},
		"correct_password", // This is the 'originalPass' that should be used
		`{
            "email": "test@example.com",
            "password": "hashed_test_pass", // Mismatch here
            "role": "USER"
        }`,
	},
	{
		&models.UserDTO{
			Email:    "admin@company.com",
			Role:     models.ADMIN,
			Password: "hashed_admin_pass",
		},
		"admin_password",
		`{
            "email": "admin@company.com",
            "password": "wrong_password_in_json", // Mismatch here
            "role": "ADMIN"
        }`,
	},
}
var listOFTime = []time.Duration{}

func TestUnitUserControllerTest(t *testing.T) {
	suite.Run(t, &UnitUserControllerTest{})
}

func (un *UnitUserControllerTest) SetupTest() {

	un.mockUseCase = new(models_mocks.MockIUseCase)
	un.controller = controllers.NewUserController(un.mockUseCase)
	un.router = gin.Default()
	un.router.POST("/login", un.controller.Login)
	un.router.POST("/register", un.controller.Register)
	for range len(validTestCases) {
		listOFTime = append(listOFTime, time.Duration(24*time.Microsecond))
	}

}
func ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ(jsonLiteral string, router *gin.Engine, lastUrl string) *httptest.ResponseRecorder {
	jsonBuffer := bytes.NewBufferString(jsonLiteral)

	req, _ := http.NewRequest("POST", lastUrl, jsonBuffer)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	return w
}

func (un *UnitUserControllerTest) TestLogin_Positive() {
	for index, each_test_case := range validTestCases {
		un.mockUseCase.On("LoginHandler", each_test_case.userFromDB).Return(each_test_case.originalPass, &listOFTime[index], nil)
	}
	for _, each_test_case := range validTestCases {
		የሙከራ_ደንበኛ_መቀበያ := ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ(each_test_case.jsonBody, un.router, "/login")
		un.Require().Equal(200, የሙከራ_ደንበኛ_መቀበያ.Code, "failure in the login", የሙከራ_ደንበኛ_መቀበያ.Body)

	}
}
func (un *UnitUserControllerTest) TestLogin_Negative() {

	for _, each_test_case := range invalidTestCases {
		un.mockUseCase.On("LoginHandler", each_test_case.userFromDB).Return("hasedPass", nil, nil)
	}
	for index, each_test_case := range invalidTestCases {
		የሙከራ_ደንበኛ_መቀበያ := ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ(each_test_case.jsonBody, un.router, "/login")
		un.Require().NotEqual(200, የሙከራ_ደንበኛ_መቀበያ.Code, "failure in the login in testcase ", index+1)
	}
}
func (un *UnitUserControllerTest) TestRegister_Positive() {
	// since the address of two variables are different it won't work hear, that's why I used mock.Anything

	for _, each := range validTestCases {
		un.mockUseCase.On("Register", each.userFromDB).Return(each.userFromDB, nil)
		የአማርኛ_ቫርያብል_ይቻላል_እንዴ := ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ(each.jsonBody, un.router, "/register")
		un.EqualValues(http.StatusOK, የአማርኛ_ቫርያብል_ይቻላል_እንዴ.Code, የአማርኛ_ቫርያብል_ይቻላል_እንዴ.Body)
	}

}
func (un *UnitUserControllerTest) TestRegister_Negative() {
	un.mockUseCase.On("Regiser", mock.Anything).Return(nil, fmt.Errorf("")) // since the address of two variables are different it won't work hear, that's why I used mock.Anything

	for _, each := range invalidTestCases {
		የአማርኛ_ቫርያብል_ይቻላል_እንዴ := ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ(each.jsonBody, un.router, "/register")
		un.NotEqual(http.StatusOK, የአማርኛ_ቫርያብል_ይቻላል_እንዴ.Code, የአማርኛ_ቫርያብል_ይቻላል_እንዴ.Body)
	}

}
