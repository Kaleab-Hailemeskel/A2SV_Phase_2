package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"task_8_testing/controllers"
	"task_8_testing/infrastructure"
	"task_8_testing/mocks"
	"task_8_testing/models"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UnitTaskControllerTestAI struct {
	suite.Suite
	mockDB     *mocks.ITaskDataBase
	router     *gin.Engine
	controller *controllers.TaskController

	currTestUser *models.User
}

func TestUnitTaskControllerTestAI(t *testing.T) {
	suite.Run(t, new(UnitTaskControllerTestAI))
}

func (uts *UnitTaskControllerTestAI) SetupTest() {
	uts.mockDB = new(mocks.ITaskDataBase)
	uts.controller = controllers.NewTaskController(uts.mockDB)

	gin.SetMode(gin.TestMode)
	uts.router = gin.New()
	uts.router.Use(gin.Recovery())

	uts.currTestUser = &models.User{
		Email:    "test@example.com",
		Password: "test_password",
		Role:     models.USER,
	}

	uts.router.Use(uts.setCurrentUserMiddleware())

	uts.router.POST("/task/:id", uts.controller.DeleteTaskByID)
	uts.router.GET("/task/:id", uts.controller.GetTaskByID)
	uts.router.PUT("/task/:id", uts.controller.PutTaskByID)
	uts.router.GET("/task", uts.controller.GetTasks)
	uts.router.POST("/task", uts.controller.PostTask)
}

func (uts *UnitTaskControllerTestAI) setCurrentUserMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(infrastructure.CURR_USER, *uts.currTestUser)
		ctx.Next()
	}
}

func (uts *UnitTaskControllerTestAI) sendJSONRequest(method, url string, jsonString string) *httptest.ResponseRecorder {
	var reqBody io.Reader

	if jsonString != "" {
		reqBody = bytes.NewBufferString(jsonString)
	}

	req, err := http.NewRequest(method, url, reqBody)
	uts.Require().NoError(err, "Failed to create HTTP request")

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	uts.router.ServeHTTP(w, req)
	return w
}

func (uts *UnitTaskControllerTestAI) TestDeleteByID_Positive() {
	taskID := "task1"
	uts.mockDB.On("DeleteOne", taskID).Return(nil).Once()

	response := uts.sendJSONRequest(http.MethodPost, fmt.Sprintf("/task/%s", taskID), "")
	uts.Equal(http.StatusAccepted, response.Code)
}

func (uts *UnitTaskControllerTestAI) TestDeleteByID_Negative() {
	taskID := "task1"
	uts.mockDB.On("DeleteOne", taskID).Return(fmt.Errorf("database error")).Once()

	response := uts.sendJSONRequest(http.MethodPost, fmt.Sprintf("/task/%s", taskID), "")
	uts.NotEqual(http.StatusAccepted, response.Code)
	uts.NotEqual(http.StatusOK, response.Code)
}

func (uts *UnitTaskControllerTestAI) TestGetByID_Positive() {
	taskID := "task1"
	fixedTime := time.Date(2025, time.July, 30, 14, 0, 0, 0, time.UTC)
	expectedTask := &models.Task{
		OwnerEmail:  "oneLove@gmail.com",
		ID:          "task_1",
		Title:       "testing",
		Description: "I am testing today",
		DueDate:     fixedTime,
		Status:      "In Progress",
	}
	uts.mockDB.On("FindByID", taskID).Return(expectedTask, nil).Once()

	response := uts.sendJSONRequest(http.MethodGet, fmt.Sprintf("/task/%s", taskID), "")
	uts.Equal(http.StatusOK, response.Code)

	var actualTask models.Task
	err := json.Unmarshal(response.Body.Bytes(), &actualTask)
	uts.Require().NoError(err, "Failed to unmarshal response body")
	uts.Equal(expectedTask.ID, actualTask.ID)
	uts.Equal(expectedTask.OwnerEmail, actualTask.OwnerEmail)
	uts.Equal(expectedTask.Title, actualTask.Title)
	uts.Equal(expectedTask.Description, actualTask.Description)
	uts.True(actualTask.DueDate.Equal(expectedTask.DueDate), "DueDate mismatch")
	uts.Equal(expectedTask.Status, actualTask.Status)
}

func (uts *UnitTaskControllerTestAI) TestGetByID_Negative_NotFound() {
	taskID := "nonexistent_task"
	uts.mockDB.On("FindByID", taskID).Return(nil, nil).Once()

	response := uts.sendJSONRequest(http.MethodGet, fmt.Sprintf("/task/%s", taskID), "")
	uts.Equal(http.StatusNotFound, response.Code)
}

func (uts *UnitTaskControllerTestAI) TestGetByID_Negative_DBError() {
	taskID := "task1"
	uts.mockDB.On("FindByID", taskID).Return(nil, fmt.Errorf("database connection lost")).Once()

	response := uts.sendJSONRequest(http.MethodGet, fmt.Sprintf("/task/%s", taskID), "")
	uts.Equal(http.StatusInternalServerError, response.Code)
}

func (uts *UnitTaskControllerTestAI) TestGetTasks_Positive_UserRole() {
	uts.currTestUser.Email = "oneLove@gmail.com"
	uts.currTestUser.Role = models.USER

	expectedTasks := &[]models.Task{
		{
			OwnerEmail:  "oneLove@gmail.com",
			ID:          "task_1",
			Title:       "User Task",
			Description: "This task belongs to the user",
			DueDate:     time.Now(),
			Status:      "In Progress",
		},
	}
	uts.mockDB.On("FindAllTasks", uts.currTestUser.Email).Return(expectedTasks, nil).Once()

	response := uts.sendJSONRequest(http.MethodGet, "/task", "")
	uts.Equal(http.StatusOK, response.Code, response.Body.String())

	var actualTasks []models.Task
	err := json.Unmarshal(response.Body.Bytes(), &actualTasks)
	uts.Require().NoError(err, "Failed to unmarshal response body")
	uts.Len(actualTasks, len(*expectedTasks))
	uts.Equal((*expectedTasks)[0].ID, actualTasks[0].ID)
}

func (uts *UnitTaskControllerTestAI) TestGetTasks_Positive_AdminRole() {
	uts.currTestUser.Email = "admin@example.com"
	uts.currTestUser.Role = models.ADMIN

	expectedTasks := &[]models.Task{
		{OwnerEmail: "user1@example.com", ID: "task_a", Title: "Admin View Task 1"},
		{OwnerEmail: "user2@example.com", ID: "task_b", Title: "Admin View Task 2"},
	}
	uts.mockDB.On("FindAllTasks", "").Return(expectedTasks, nil).Once()

	response := uts.sendJSONRequest(http.MethodGet, "/task", "")
	uts.Equal(http.StatusOK, response.Code, response.Body.String())

	var actualTasks []models.Task
	err := json.Unmarshal(response.Body.Bytes(), &actualTasks)
	uts.Require().NoError(err, "Failed to unmarshal response body")
	uts.Len(actualTasks, len(*expectedTasks))
}

func (uts *UnitTaskControllerTestAI) TestGetTasks_Negative_DBError() {
	uts.currTestUser.Email = "test@example.com"
	uts.currTestUser.Role = models.ADMIN

	uts.mockDB.On("FindAllTasks", mock.AnythingOfType("string")).Return(nil, fmt.Errorf("DB connection failed")).Once()

	response := uts.sendJSONRequest(http.MethodGet, "/task", "")
	uts.Equal(http.StatusInternalServerError, response.Code, response.Body.String())
}

func (uts *UnitTaskControllerTestAI) TestPutTaskByID_Positive() {
	taskID := "task1"
	fixedTime := time.Date(2025, time.July, 30, 14, 0, 0, 0, time.UTC)

	uts.currTestUser.Email = "three@gmail.com"
	uts.currTestUser.Role = models.USER

	expectedTaskUpdate := models.Task{
		OwnerEmail:  "three@gmail.com",
		ID:          "task_003",
		Title:       "Schedule team sync",
		Description: "Fin a suitable time for weekly stand-up meeting",
		DueDate:     fixedTime,
		Status:      "todo",
	}

	uts.mockDB.On("UpdateOne", taskID, expectedTaskUpdate).Return(nil).Once()

	updatedTaskJson := fmt.Sprintf(`
	{
		"id": "task_003",
		"ownerEmail": "three@gmail.com",
		"title": "Schedule team sync",
		"description": "Fin a suitable time for weekly stand-up meeting",
		"due_date": "%s",
		"status": "todo"
	}`, fixedTime.Format(time.RFC3339))

	response := uts.sendJSONRequest(http.MethodPut, fmt.Sprintf("/task/%s", taskID), updatedTaskJson)
	uts.Equal(http.StatusOK, response.Code, response.Body.String())
}

func (uts *UnitTaskControllerTestAI) TestPutTaskByID_Negative_DBError() {
	taskID := "task1"
	fixedTime := time.Date(2025, time.July, 30, 14, 0, 0, 0, time.UTC)

	uts.currTestUser.Email = "three@gmail.com"
	uts.currTestUser.Role = models.USER

	expectedTaskUpdate := models.Task{
		OwnerEmail:  "three@gmail.com",
		ID:          "task_003",
		Title:       "Schedule team sync",
		Description: "Fin a suitable time for weekly stand-up meeting",
		DueDate:     fixedTime,
		Status:      "todo",
	}

	uts.mockDB.On("UpdateOne", taskID, expectedTaskUpdate).Return(fmt.Errorf("DB update failed")).Once()

	updatedTaskJson := fmt.Sprintf(`
	{
		"id": "task_003",
		"ownerEmail": "three@gmail.com",
		"title": "Schedule team sync",
		"description": "Fin a suitable time for weekly stand-up meeting",
		"due_date": "%s",
		"status": "todo"
	}`, fixedTime.Format(time.RFC3339))

	response := uts.sendJSONRequest(http.MethodPut, fmt.Sprintf("/task/%s", taskID), updatedTaskJson)
	uts.Equal(http.StatusInternalServerError, response.Code, response.Body.String())
}

func (uts *UnitTaskControllerTestAI) TestPutTaskByID_Negative_InvalidJSON() {
	taskID := "task1"
	uts.currTestUser.Email = "three@gmail.com"
	uts.currTestUser.Role = models.USER

	invalidJson := `{"id": "task_003", "title": "Missing due_date and status"}`

	response := uts.sendJSONRequest(http.MethodPut, fmt.Sprintf("/task/%s", taskID), invalidJson)
	uts.Equal(http.StatusBadRequest, response.Code, response.Body.String())
	uts.mockDB.AssertNotCalled(uts.T(), "UpdateOne", mock.Anything, mock.Anything)
}

func (uts *UnitTaskControllerTestAI) TestPostTask_Positive() {
	fixedTime := time.Date(2025, time.July, 30, 14, 0, 0, 0, time.UTC)

	uts.currTestUser.Email = "three@gmail.com"
	uts.currTestUser.Role = models.USER

	expectedNewTask := models.Task{
		OwnerEmail:  uts.currTestUser.Email,
		ID:          "task_003",
		Title:       "Schedule team sync",
		Description: "Fin a suitable time for weekly stand-up meeting",
		DueDate:     fixedTime,
		Status:      "todo",
	}

	uts.mockDB.On("InsertOne", expectedNewTask).Return(nil).Once()

	newTaskJson := fmt.Sprintf(`
	{
		"id": "task_003",
		"title": "Schedule team sync",
		"description": "Fin a suitable time for weekly stand-up meeting",
		"due_date": "%s",
		"status": "todo"
	}`, fixedTime.Format(time.RFC3339))

	response := uts.sendJSONRequest(http.MethodPost, "/task", newTaskJson)
	uts.Equal(http.StatusCreated, response.Code, response.Body.String())
}

func (uts *UnitTaskControllerTestAI) TestPostTask_Negative_DBError() {
	fixedTime := time.Date(2025, time.July, 30, 14, 0, 0, 0, time.UTC)

	uts.currTestUser.Email = "three@gmail.com"
	uts.currTestUser.Role = models.USER

	expectedNewTask := models.Task{
		OwnerEmail:  uts.currTestUser.Email,
		ID:          "task_003",
		Title:       "Schedule team sync",
		Description: "Fin a suitable time for weekly stand-up meeting",
		DueDate:     fixedTime,
		Status:      "todo",
	}
	uts.mockDB.On("InsertOne", expectedNewTask).Return(fmt.Errorf("DB insert failed")).Once()

	newTaskJson := fmt.Sprintf(`
	{
		"id": "task_003",
		"title": "Schedule team sync",
		"description": "Fin a suitable time for weekly stand-up meeting",
		"due_date": "%s",
		"status": "todo"
	}`, fixedTime.Format(time.RFC3339))

	response := uts.sendJSONRequest(http.MethodPost, "/task", newTaskJson)
	uts.NotEqual(http.StatusCreated, response.Code, response.Body.String())
}

func (uts *UnitTaskControllerTestAI) TestPostTask_Negative_InvalidJSON() {
	uts.currTestUser.Email = "three@gmail.com"
	uts.currTestUser.Role = models.USER

	invalidJson := `{"id": "task_004", "description": "missing title", "due_date": "invalid-date-format", "status": "todo"}`

	response := uts.sendJSONRequest(http.MethodPost, "/task", invalidJson)
	uts.NotEqual(http.StatusCreated, response.Code, response.Body.String())
	uts.mockDB.AssertNotCalled(uts.T(), "InsertOne", mock.Anything)
}

func (uts *UnitTaskControllerTestAI) TearDownTest() {
	uts.mockDB.AssertExpectations(uts.T())
}
