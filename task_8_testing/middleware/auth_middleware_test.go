package middleware_test

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"task_8_testing/infrastructure"
	"task_8_testing/middleware"
	"task_8_testing/mocks"
	"task_8_testing/models"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var invalidTestCases_authorization_authorization = []authorizationTestCase{
	{
		name:           "non-existent task",
		taskID:         "invalidID",
		mockTask:       nil,
		mockError:      errors.New("not found"),
		expectedStatus: http.StatusNotFound,
	},
	{
		name:           "internal DB error",
		taskID:         "task2",
		mockTask:       nil,
		mockError:      errors.New("db error"),
		expectedStatus: http.StatusInternalServerError,
	},
}

type UnitTestAuthMiddleware struct {
	suite.Suite
	mockJwt    *mocks.IAuthentication
	mockUserDB *mocks.IUserDataBase
	mockTaskDB *mocks.ITaskDataBase
	mainAuth   middleware.UserAuth

	route *gin.Engine
}

type authorizationTestCase struct {
	name           string
	taskID         string
	mockTask       *models.Task
	mockError      error
	expectedStatus int
}

var validTestCases_authorization = []authorizationTestCase{
	{
		name:   "valid task authorization",
		taskID: "task1",
		mockTask: &models.Task{
			OwnerEmail:  "oneLove@gmail.com",
			ID:          "task1",
			Title:       "testing",
			Description: "I am testing today",
			DueDate:     time.Now(),
			Status:      "In Progress",
		},
		mockError:      nil,
		expectedStatus: http.StatusOK,
	},
}

var validTestCases_authentication = []struct {
	name           string
	token          string
	email          string
	role           string
	expectedStatus int
}{
	{
		name:           "Valid token with ADMIN role",
		token:          "validAdminToken",
		email:          "admin@gmail.com",
		role:           models.ADMIN,
		expectedStatus: http.StatusOK,
	},
	{
		name:           "Valid token with USER role",
		token:          "validUserToken",
		email:          "user@gmail.com",
		role:           models.USER,
		expectedStatus: http.StatusOK,
	},
}
var invalidTestCases_authentication = []struct {
	name           string
	token          string
	email          string
	role           string
	expectedStatus int
}{
	{
		name:           "Token expired",
		token:          "expiredToken",
		email:          "expired@gmail.com",
		role:           models.USER,
		expectedStatus: http.StatusUnauthorized,
	},
	{
		name:           "User not found in DB",
		token:          "ghostToken",
		email:          "ghost@gmail.com",
		role:           models.ADMIN,
		expectedStatus: http.StatusUnauthorized,
	},
	{
		name:           "Malformed token",
		token:          "invalidTokenFormat",
		email:          "",
		role:           models.USER,
		expectedStatus: http.StatusUnauthorized,
	},
	{
		name:           "Missing email claim",
		token:          "tokenMissingEmail",
		email:          "",
		role:           models.USER,
		expectedStatus: http.StatusUnauthorized,
	},
}

func TestUnitTestAuthMiddleware(t *testing.T) {
	suite.Run(t, new(UnitTestAuthMiddleware))
}
func (un *UnitTestAuthMiddleware) SetupSuite() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("cannot initailize env variable")
	}
	infrastructure.InitEnv()
}
func (un *UnitTestAuthMiddleware) SetupTest() {
	un.mockUserDB = new(mocks.IUserDataBase)
	un.mockTaskDB = new(mocks.ITaskDataBase)
	un.mockJwt = new(mocks.IAuthentication)
	un.mainAuth = *middleware.NewUserAuth(un.mockJwt, un.mockUserDB, un.mockTaskDB)

	un.route = gin.Default()

	un.route.GET(ticationRoute, middlewareFiller, un.mainAuth.Authentication)
	un.route.GET(torizationRoute, middlewareFiller, un.mainAuth.Authorization)
}

const (
	ticationRoute   = "/authentication_route"
	torizationRoute = "/authorization_route"
)

func middlewareFiller(ctx *gin.Context) {
	ctx.Set(infrastructure.CURR_USER, models.User{Email: "oneLove@gmail.com", Role: models.USER, Password: "pass"})
	ctx.Next()
}
func (uts *UnitTestAuthMiddleware) sendJSONRequest(method, url string, jsonString string, cookies *http.Cookie) *httptest.ResponseRecorder {
	var reqBody io.Reader

	if jsonString != "" {
		reqBody = bytes.NewBufferString(jsonString)
	}
	req, err := http.NewRequest(method, url, reqBody)
	uts.Require().NoError(err, "Failed to create HTTP request")
	req.Header.Set("Content-Type", "application/json")

	if req.Header.Get("Content-Type") == "" && reqBody != nil {
		// Default to application/json if body exists and no Content-Type set
		req.Header.Set("Content-Type", "application/json")
	}

	// Add cookies to the request
	if cookies != nil {
		req.AddCookie(cookies)
	}

	w := httptest.NewRecorder()
	uts.route.ServeHTTP(w, req)

	return w
}

// stack on sending cookie from the client to server in code. Unable to send a cookie from httptest

func (un *UnitTestAuthMiddleware) TestAuthentication() {
	parsedToken := &jwt.Token{
		Header: map[string]interface{}{"alg": "HS256"},
		Claims: jwt.MapClaims{"email": "userEmail@gmail.com", "exp": float64(time.Now().Add(time.Hour).Unix())},
		Method: jwt.SigningMethodHS256,
		Valid:  true,
	}

	un.mockJwt.On("ParseToken", "validTokenString").Return(parsedToken, nil)
	un.mockJwt.On("TokenExpired", mock.AnythingOfType("*jwt.Token")).Return(false, nil)
	un.mockJwt.On("GetUserEmailFromSecurityToken", mock.AnythingOfType("*jwt.Token")).Return("userEmail@gmail.com", nil)
	un.mockUserDB.On("FindUserByEmail", "userEmail@gmail.com").Return(
		&models.User{
			Email:    "userEmail@gmail.com",
			Password: "pass",
			Role:     models.ADMIN,
		}, nil)

	w := un.sendJSONRequest("GET", ticationRoute, `{"Hi":"hallo"}`, &http.Cookie{
		Name:     infrastructure.HEADER,
		Value:    "validTokenString",
		Path:     "/", // Important for cookie to be sent with request
		HttpOnly: true,
		Secure:   false,
	})

	un.Equal(http.StatusOK, w.Code, w.Body)
}
func (un *UnitTestAuthMiddleware) TestAuthentication_Positive() {
	for _, tc := range validTestCases_authentication {
		un.T().Run(tc.name, func(t *testing.T) {
			parsedToken := &jwt.Token{
				Header: map[string]interface{}{"alg": "HS256"},
				Claims: jwt.MapClaims{"email": tc.email, "exp": float64(time.Now().Add(time.Hour).Unix())},
				Method: jwt.SigningMethodHS256,
				Valid:  true,
			}

			un.mockJwt.On("ParseToken", tc.token).Return(parsedToken, nil)
			un.mockJwt.On("TokenExpired", parsedToken).Return(false, nil)
			un.mockJwt.On("GetUserEmailFromSecurityToken", parsedToken).Return(tc.email, nil)
			un.mockUserDB.On("FindUserByEmail", tc.email).Return(&models.User{
				Email:    tc.email,
				Password: "securePass",
				Role:     tc.role,
			}, nil)

			w := un.sendJSONRequest("GET", ticationRoute, `{"test":"data"}`, &http.Cookie{
				Name:     infrastructure.HEADER,
				Value:    tc.token,
				Path:     "/",
				HttpOnly: true,
				Secure:   false,
			})

			un.Equal(tc.expectedStatus, w.Code)

		})
	}
}

func (un *UnitTestAuthMiddleware) TestAuthentication_Negative() {
	for _, tc := range invalidTestCases_authentication {
		un.T().Run(tc.name, func(t *testing.T) {
			parsedToken := &jwt.Token{
				Header: map[string]interface{}{"alg": "HS256"},
				Claims: jwt.MapClaims{"email": tc.email, "exp": float64(time.Now().Add(-time.Hour).Unix())},
				Method: jwt.SigningMethodHS256,
				Valid:  true,
			}

			switch tc.name {
			case "Token expired":
				un.mockJwt.On("ParseToken", tc.token).Return(parsedToken, nil)
				un.mockJwt.On("TokenExpired", parsedToken).Return(true, nil)
			case "User not found in DB":
				un.mockJwt.On("ParseToken", tc.token).Return(parsedToken, nil)
				un.mockJwt.On("TokenExpired", parsedToken).Return(false, nil)
				un.mockJwt.On("GetUserEmailFromSecurityToken", parsedToken).Return(tc.email, nil)
				un.mockUserDB.On("FindUserByEmail", tc.email).Return(nil, errors.New("not found"))
			case "Malformed token":
				un.mockJwt.On("ParseToken", tc.token).Return(nil, errors.New("token malformed"))
			case "Missing email claim":
				parsedToken.Claims = jwt.MapClaims{"exp": float64(time.Now().Add(time.Hour).Unix())}
				un.mockJwt.On("ParseToken", tc.token).Return(parsedToken, nil)
				un.mockJwt.On("TokenExpired", parsedToken).Return(false, nil)
				un.mockJwt.On("GetUserEmailFromSecurityToken", parsedToken).Return("", errors.New("email missing"))
			}

			w := un.sendJSONRequest("GET", ticationRoute, `{"test":"data"}`, &http.Cookie{
				Name:     infrastructure.HEADER,
				Value:    tc.token,
				Path:     "/",
				HttpOnly: true,
				Secure:   false,
			})

			un.NotEqual(http.StatusOK, w.Code)
		})
	}
}

func (un *UnitTestAuthMiddleware) TestAuthorization_Positive() {
	for _, tc := range validTestCases_authorization {
		un.Run(tc.name, func() {
			un.mockTaskDB.On("FindByID", tc.taskID).Return(tc.mockTask, tc.mockError)
			w := un.sendJSONRequest("GET", torizationRoute, "", nil)

			un.Equal(http.StatusOK, w.Code)

		})
	}
}
func (un *UnitTestAuthMiddleware) TestAuthorization_Negative() {
	for _, tc := range invalidTestCases_authorization_authorization {
		un.Run(tc.name, func() {
			un.mockTaskDB.On("FindByID", tc.taskID).Return(tc.mockTask, tc.mockError)

			w := un.sendJSONRequest("GET", torizationRoute+"/"+tc.taskID, "", nil)

			un.NotEqual(http.StatusOK, w.Code)

		})
	}
}
