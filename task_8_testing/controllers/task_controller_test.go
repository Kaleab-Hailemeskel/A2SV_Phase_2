package controllers_test

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"task_8_testing/controllers"
	"task_8_testing/infrastructure"
	"task_8_testing/models"
	models_mocks "task_8_testing/models/mocks"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UnitTaskControllerTest struct {
	suite.Suite
	useCase    *models_mocks.MockIUseCase
	route      *gin.Engine
	controller *controllers.TaskController
}

var currTestUser *models.UserDTO = new(models.UserDTO)
var (
	ErrConst       = fmt.Errorf("tEST ERROR")
	timeTT, _        = time.Parse(time.RFC3339, "2025-07-30T14:00:00Z")
	listOfValidObj   = []*primitive.ObjectID{}
	listOfObjStrings = []string{
		"689a4d408ef3d58c84aa1b92",
		" 689a4d408ef3d58c84aa1b93",
		" 689a4d408ef3d58c84aa1b94",
		"689a4d408ef3d58c84aa1b95",
		"689a4d408ef3d58c84aa1b96",
		"689a4d408ef3d58c84aa1b97",
		"689a4d408ef3d58c84aa1b98",
		"689a4d408ef3d58c84aa1b99",
		"689a4d408ef3d58c84aa1b9a",
		"689a4d408ef3d58c84aa1b9b",
	}
	listOfValidTasks = []*models.TaskDTO{
		{
			OwnerEmail:  "oneLove@gmail.com",
			ID:          primitive.NewObjectID(),
			Title:       "testing",
			Description: "I am testing today",
			DueDate:     time.Now(),
			Status:      "In Progress",
		},
		{
			OwnerEmail:  "oneLove@gmail.com",
			ID:          primitive.NewObjectID(),
			Title:       "testing",
			Description: "I am testing today",
			DueDate:     time.Now(),
			Status:      "In Progress",
		},
		{
			OwnerEmail:  "three@gmail.com",
			ID:          primitive.NewObjectID(),
			Title:       "Schedule team sync",
			Description: "Fin a suitable time for weekly stand-up meeting",
			DueDate:     timeTT,
			Status:      "todo",
		},
	}
)

func TestUnitTaskControllerTest(t *testing.T) {
	suite.Run(t, &UnitTaskControllerTest{})
}

func (un *UnitTaskControllerTest) SetupSuite() {
	infrastructure.InitEnv()
	for _, each_string := range listOfObjStrings {
		resObj, _ := primitive.ObjectIDFromHex(each_string)
		listOfValidObj = append(listOfValidObj, &resObj)
	}
	for index := range min(len(listOfObjStrings), len(listOfValidTasks)) {
		resObj, _ := primitive.ObjectIDFromHex(listOfObjStrings[index])
		listOfValidTasks[index].ID = resObj
	}
}
func (un *UnitTaskControllerTest) SetupTest() {
	un.useCase = new(models_mocks.MockIUseCase)
	un.controller = controllers.NewTaskController(un.useCase)
	un.route = gin.Default()

	currTestUser.Email = "kaleabExample@gmail.com"
	currTestUser.ID = primitive.NewObjectID()
	currTestUser.Password = "hallo world"
	currTestUser.Role = models.USER
	log.Println("📝: ", currTestUser)

	un.route.POST("/task/:id", func(ctx *gin.Context) {
		ctx.Set(infrastructure.CURR_USER, *currTestUser)
		ctx.Next()
	}, un.controller.DeleteTaskByID)
	un.route.GET("/task/:id", func(ctx *gin.Context) {
		ctx.Set(infrastructure.CURR_USER, *currTestUser)
		ctx.Next()
	}, un.controller.GetTaskByID)
	un.route.PUT("/task/:id", func(ctx *gin.Context) {
		ctx.Set(infrastructure.CURR_USER, *currTestUser)
		ctx.Next()
	}, un.controller.PutTaskByID)
	un.route.GET("/task", func(ctx *gin.Context) {
		ctx.Set(infrastructure.CURR_USER, *currTestUser)
		ctx.Next()
	}, un.controller.GetTasks)
	un.route.POST("/task", func(ctx *gin.Context) {
		ctx.Set(infrastructure.CURR_USER, *currTestUser)
		ctx.Next()
	}, un.controller.PostTask)
}
func በዩርል_በኩል_ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ(method string, router *gin.Engine, lastUrl string, jsonString string) *httptest.ResponseRecorder {
	jsonBuffer := bytes.NewBufferString(jsonString)

	req, _ := http.NewRequest(method, lastUrl, jsonBuffer)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	return w
}

func (un *UnitTaskControllerTest) TestDeleteByID_Positive() {
	un.useCase.On("DeleteTask", listOfObjStrings[0], currTestUser.Email).Return(nil)

	log.Println("🥇 CurrUSER", infrastructure.CURR_USER)
	መልስ := በዩርል_በኩል_ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ("POST", un.route, "/task/"+listOfObjStrings[0], "")
	un.Equal(http.StatusAccepted, መልስ.Code, "ዝርዝር ሊደመሰስ አልተቻለም", መልስ.Body)
}
func (un *UnitTaskControllerTest) TestDeleteByID_Negative() {

	un.useCase.On("DeleteTask", listOfObjStrings[0], currTestUser.Email).Return(ErrConst)

	መልስ := በዩርል_በኩል_ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ("POST", un.route, "/task/"+listOfObjStrings[0], "")
	un.NotContains([]int{http.StatusAccepted, http.StatusOK}, መልስ.Code)
}

func (un *UnitTaskControllerTest) TestGetByID_Postive() {
	un.useCase.On("GetTaskByID", listOfObjStrings[1], currTestUser.Email).Return(listOfValidTasks[0], nil)

	መልስ := በዩርል_በኩል_ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ("GET", un.route, "/task/"+listOfObjStrings[1], "")
	un.Equal(http.StatusOK, መልስ.Code, መልስ.Body)
}
func (un *UnitTaskControllerTest) TestGetByID_Negative() {
	un.useCase.On("GetTaskByID", listOfObjStrings[1], currTestUser.Email).Return(nil, ErrConst)

	መልስ := በዩርል_በኩል_ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ("GET", un.route, "/task/"+listOfObjStrings[1], "")
	un.NotContains([]int{http.StatusAccepted, http.StatusOK}, መልስ.Code, መልስ.Body)
}

func (un *UnitTaskControllerTest) TestGetTask_Positive() {
	un.useCase.On("GetAllTask", currTestUser.Email).Return(listOfValidTasks, nil)

	መልስ := በዩርል_በኩል_ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ("GET", un.route, "/task", "")
	un.Equal(http.StatusOK, መልስ.Code, መልስ.Body)

	currTestUser.Role = models.ADMIN
	un.useCase.On("GetAllTask", "").Return(listOfValidTasks, nil)

	መልስ = በዩርል_በኩል_ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ("GET", un.route, "/task", "")
	un.Equal(http.StatusOK, መልስ.Code, መልስ.Body)
}
func (un *UnitTaskControllerTest) TestGetTask_Negative() {
	un.useCase.On("GetAllTask", currTestUser.Email).Return(nil, ErrConst)
	un.useCase.On("GetAllTask", "").Return(nil, ErrConst)
	መልስ := በዩርል_በኩል_ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ("GET", un.route, "/task", "")
	un.NotEqual(http.StatusOK, መልስ.Code, መልስ.Body)

	currTestUser.Role = models.ADMIN
}


// func (un *UnitTaskControllerTest) TestPostTask_Positive() {
// 	un.useCase.On("InsertOne", listOfValidTasks[2]).Return(nil)

// 	currTestUser.Email = "three@gmail.com" // demo initaialization
// 	currTestUser.Password = "pass"
// 	currTestUser.Role = models.ADMIN

// 	newTask := `
//     {
//         "id": "task_003",
//         "title": "Schedule team sync",
//         "description": "Fin a suitable time for weekly stand-up meeting",
//         "due_date": "2025-07-30T14:00:00Z",
//         "status": "todo"
//     }
// 	`

// 	መልስ := በዩርል_በኩል_ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ("POST", un.route, "/task", newTask)
// 	un.Equal(http.StatusCreated, መልስ.Code, መልስ.Body)
// }
// func (un *UnitTaskControllerTest) TestPostTask_Negative() {

// 	un.useCase.On("InsertOne", listOfValidTasks[2]).Return(nil)

// 	currTestUser.Email = "three@gmail.com" // demo initaialization
// 	currTestUser.Password = "pass"
// 	currTestUser.Role = models.ADMIN

// 	newTask := `
//     {
//         "id": "task_003",
//         "title": "Schedule team sync",
//         "description": "Fin a suitable time for weekly stand-up meeting",
//         "due_date": "2025-07-30T14:00:00Z",
//     }
// 	`

// 	መልስ := በዩርል_በኩል_ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ("POST", un.route, "/task", newTask)
// 	un.NotEqual(http.StatusCreated, መልስ.Code, መልስ.Body)
// }
