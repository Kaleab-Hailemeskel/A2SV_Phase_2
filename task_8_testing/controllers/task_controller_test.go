package controllers_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"task_8_testing/controllers"
	"task_8_testing/infrastructure"
	"task_8_testing/mocks"
	"task_8_testing/models"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type UnitTaskControllerTest struct {
	suite.Suite
	mockDB *mocks.ITaskDataBase
	route  *gin.Engine

	controller *controllers.TaskController
}

var currTestUser *models.User = new(models.User)

func TestUnitTaskControllerTest(t *testing.T) {
	suite.Run(t, &UnitTaskControllerTest{})
}

func (un *UnitTaskControllerTest) SetupTest() {
	un.mockDB = new(mocks.ITaskDataBase)
	un.controller = controllers.NewTaskController(un.mockDB)

	un.route = gin.Default()
	un.route.POST("/task/:id", un.controller.DeleteTaskByID)
	un.route.GET("/task/:id", un.controller.GetTaskByID)
	un.route.PUT("/task/:id", un.controller.PutTaskByID)

	un.route.GET("/task", func(ctx *gin.Context) {
		ctx.Set(infrastructure.CURR_USER, *currTestUser)
	}, un.controller.GetTasks)
	un.route.POST("/task", func(ctx *gin.Context) {
		ctx.Set(infrastructure.CURR_USER, *currTestUser)
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
	un.mockDB.On("DeleteOne", "task1").Return(nil)

	መልስ := በዩርል_በኩል_ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ("POST", un.route, "/task/task1", "")
	un.Equal(http.StatusAccepted, መልስ.Code, "ዝርዝር ሊደመሰስ አልተቻለም")
}
func (un *UnitTaskControllerTest) TestDeleteByID_Negative() {
	un.mockDB.On("DeleteOne", "task1").Return(fmt.Errorf(""))

	መልስ := በዩርል_በኩል_ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ("POST", un.route, "/task/task1", "")
	un.NotContains([]int{http.StatusAccepted, http.StatusOK}, መልስ.Code)
}

func (un *UnitTaskControllerTest) TestGetByID_Postive() {
	un.mockDB.On("FindByID", "task1").Return(&models.Task{
		OwnerEmail:  "oneLove@gmail.com",
		ID:          "task_1",
		Title:       "testing",
		Description: "I am testing today",
		DueDate:     time.Now(),
		Status:      "In Progress",
	}, nil)

	መልስ := በዩርል_በኩል_ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ("GET", un.route, "/task/task1", "")
	un.Equal(http.StatusOK, መልስ.Code, መልስ.Body)
}
func (un *UnitTaskControllerTest) TestGetByID_Negative() {
	un.mockDB.On("FindByID", "task1").Return(nil, nil)

	መልስ := በዩርል_በኩል_ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ("GET", un.route, "/task/task1", "")
	un.NotContains([]int{http.StatusAccepted, http.StatusOK}, መልስ.Code, መልስ.Body)
}

func (un *UnitTaskControllerTest) TestGetTask_Positive() {
	un.mockDB.On("FindAllTasks", "").Return(&[]models.Task{
		{
			OwnerEmail:  "oneLove@gmail.com",
			ID:          "task_1",
			Title:       "testing",
			Description: "I am testing today",
			DueDate:     time.Now(),
			Status:      "In Progress",
		},
	}, nil)

	un.mockDB.On("FindAllTasks", "oneLove@gmail.com").Return(&[]models.Task{
		{
			OwnerEmail:  "oneLove@gmail.com",
			ID:          "task_1",
			Title:       "testing",
			Description: "I am testing today",
			DueDate:     time.Now(),
			Status:      "In Progress",
		},
	}, nil)

	currTestUser.Email = "oneLove@gmail.com" // demo initaialization
	currTestUser.Password = "pass"
	currTestUser.Role = models.USER

	መልስ := በዩርል_በኩል_ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ("GET", un.route, "/task", "")
	un.Equal(http.StatusOK, መልስ.Code, መልስ.Body)
	currTestUser.Role = models.ADMIN

	መልስ = በዩርል_በኩል_ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ("GET", un.route, "/task", "")
	un.Equal(http.StatusOK, መልስ.Code, መልስ.Body)

}
func (un *UnitTaskControllerTest) TestGetTask_Negative() {
	un.mockDB.On("FindAllTasks", "").Return(nil, fmt.Errorf("mock Error"))
	un.mockDB.On("FindAllTasks", "oneLove@gmail.com").Return(nil, fmt.Errorf("mock Error"))

	currTestUser.Email = "oneLove@gmail.com" // demo initaialization
	currTestUser.Password = "pass"
	currTestUser.Role = models.ADMIN

	መልስ := በዩርል_በኩል_ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ("GET", un.route, "/task", "")
	un.NotEqual(http.StatusOK, መልስ.Code, መልስ.Body)
}
func (un *UnitTaskControllerTest) TestPutTaskByID_Positive() {
	t, _ := time.Parse(time.RFC3339, "2025-07-30T14:00:00Z")
	un.mockDB.On("UpdateOne", "task1", models.Task{
		OwnerEmail:  "three@gmail.com",
		ID:          "task_003",
		Title:       "Schedule team sync",
		Description: "Fin a suitable time for weekly stand-up meeting",
		DueDate:     t,
		Status:      "todo",
	}).Return(nil)

	currTestUser.Email = "oneLove@gmail.com" // demo initaialization
	currTestUser.Password = "pass"
	currTestUser.Role = models.ADMIN
	updatedTaskJson := `
    {
        "id": "task_003",
        "ownerEmail": "three@gmail.com",
        "title": "Schedule team sync",
        "description": "Fin a suitable time for weekly stand-up meeting",
        "due_date": "2025-07-30T14:00:00Z",
        "status": "todo"
    }
	`
	መልስ := በዩርል_በኩል_ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ("PUT", un.route, "/task/task1", updatedTaskJson)
	un.Equal(http.StatusOK, መልስ.Code, መልስ.Body)
}
func (un *UnitTaskControllerTest) TestPutTaskByID_Negative() {
	t, _ := time.Parse(time.RFC3339, "2025-07-30T14:00:00Z")
	un.mockDB.On("UpdateOne", "task1", models.Task{
		OwnerEmail:  "three@gmail.com",
		ID:          "task_003",
		Title:       "Schedule team sync",
		Description: "Fin a suitable time for weekly stand-up meeting",
		DueDate:     t,
		Status:      "todo",
	}).Return(fmt.Errorf("mock Error"))

	currTestUser.Email = "oneLove@gmail.com" // demo initaialization
	currTestUser.Password = "pass"
	currTestUser.Role = models.ADMIN
	updatedTaskJson := `
    {
        "id": "task_003",
        "ownerEmail": "three@gmail.com",
        "title": "Schedule team sync",
        "description": "Fin a suitable time for weekly stand-up meeting",
        "due_date": "2025-07-30T14:00:00Z",
        "status": "todo"
    }
	`
	መልስ := በዩርል_በኩል_ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ("PUT", un.route, "/task/task1", updatedTaskJson)
	un.NotEqual(http.StatusOK, መልስ.Code, መልስ.Body)
}
func (un *UnitTaskControllerTest) TestPostTask_Positive() {
	t, _ := time.Parse(time.RFC3339, "2025-07-30T14:00:00Z")
	un.mockDB.On("InsertOne", models.Task{
		OwnerEmail:  "three@gmail.com",
		ID:          "task_003",
		Title:       "Schedule team sync",
		Description: "Fin a suitable time for weekly stand-up meeting",
		DueDate:     t,
		Status:      "todo",
	}).Return(nil)

	currTestUser.Email = "three@gmail.com" // demo initaialization
	currTestUser.Password = "pass"
	currTestUser.Role = models.ADMIN

	newTask := `
    {
        "id": "task_003",
        "title": "Schedule team sync",
        "description": "Fin a suitable time for weekly stand-up meeting",
        "due_date": "2025-07-30T14:00:00Z",
        "status": "todo"
    }
	`
	
	መልስ := በዩርል_በኩል_ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ("POST", un.route, "/task", newTask)
	un.Equal(http.StatusCreated, መልስ.Code, መልስ.Body)
}
func (un *UnitTaskControllerTest) TestPostTask_Negative() {
	t, _ := time.Parse(time.RFC3339, "2025-07-30T14:00:00Z")
	un.mockDB.On("InsertOne", models.Task{
		OwnerEmail:  "three@gmail.com",
		ID:          "task_003",
		Title:       "Schedule team sync",
		Description: "Fin a suitable time for weekly stand-up meeting",
		DueDate:     t,
		Status:      "todo",
	}).Return(nil)

	currTestUser.Email = "three@gmail.com" // demo initaialization
	currTestUser.Password = "pass"
	currTestUser.Role = models.ADMIN

	newTask := `
    {
        "id": "task_003",
        "title": "Schedule team sync",
        "description": "Fin a suitable time for weekly stand-up meeting",
        "due_date": "2025-07-30T14:00:00Z",
    }
	`
	
	መልስ := በዩርል_በኩል_ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ("POST", un.route, "/task", newTask)
	un.NotEqual(http.StatusCreated, መልስ.Code, መልስ.Body)
}