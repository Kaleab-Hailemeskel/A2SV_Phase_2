package controllers_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"task_8_testing/controllers"
	"task_8_testing/mocks"
	"task_8_testing/models"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UnitUserControllerTest struct {
	suite.Suite
	mockDB         *mocks.IUserDataBase
	mockPassServ   *mocks.IPasswordService
	mockJwtHandler *mocks.IAuthentication

	controller *controllers.UserController
	router     *gin.Engine
}

var validTestCases = []struct {
	userFromDB   models.User
	originalPass string
	jsonBody     string
}{
	{
		models.User{
			Email:    "admin@example.com",
			Role:     models.ADMIN,
			Password: "hashed_password_admin",
		},
		"orginal_Pass",
		`{
    	"email": "admin@example.com",
    	"password": "orginal_Pass",
    	"role": "ADMIN"
	}`,
	}, {models.User{
		Email:    "user_one@domain.com",
		Role:     models.USER,
		Password: "hashed_password_user1",
	},
		"original_Pass_user1", // Assuming this is the unhashed password
		`{
		"email": "user_one@domain.com",
		"password": "original_Pass_user1",
		"role": "USER"
	}`}, {
		models.User{
			Email:    "support@company.org",
			Role:     models.ADMIN,
			Password: "secure_pass_support",
		},
		"original_Pass_support", // Assuming this is the unhashed password
		`{
		"email": "support@company.org",
		"password": "original_Pass_support",
		"role": "ADMIN"
	}`,
	}, {models.User{
		Email:    "jane.doe@service.net",
		Role:     models.USER,
		Password: "hashed_password_jane",
	},
		"original_Pass_jane", // Assuming this is the unhashed password
		`{
		"email": "jane.doe@service.net",
		"password": "original_Pass_jane",
		"role": "USER"
	}`},
}
var invalidTestCases = []struct {
	userFromDB   models.User
	originalPass string
	jsonBody     string
}{
	// --- Variant 1: Password mismatch in jsonBody vs. originalPass ---
	{
		models.User{
			Email:    "test@example.com",
			Role:     models.USER,
			Password: "hashed_test_pass",
		},
		"correct_password", // This is the 'originalPass' that should be used
		`{
            "email": "test@example.com",
            "password": "incorrect_password_in_json", // Mismatch here
            "role": "USER"
        }`,
	},
	{
		models.User{
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

	// --- Variant 2: Email in jsonBody is not a valid format ---
	{
		models.User{
			Email:    "user_one@domain.com",
			Role:     models.USER,
			Password: "hashed_password_user1",
		},
		"original_Pass_user1",
		`{
            "email": "invalid-email-format", // Invalid email
            "password": "original_Pass_user1",
            "role": "USER"
        }`,
	},
	{
		models.User{
			Email:    "test@example.com",
			Role:     models.USER,
			Password: "hashed_test_pass",
		},
		"correct_password",
		`{
            "email": "missing.at.sign.com", // Invalid email
            "password": "correct_password",
            "role": "USER"
        }`,
	},
	{
		models.User{
			Email:    "test@example.com",
			Role:     models.USER,
			Password: "hashed_test_pass",
		},
		"correct_password",
		`{
            "email": "no-dot@com", // Invalid email (missing TLD)
            "password": "correct_password",
            "role": "USER"
        }`,
	},
	{
		models.User{
			Email:    "test@example.com",
			Role:     models.USER,
			Password: "hashed_test_pass",
		},
		"correct_password",
		`{
            "email": "@example.com", // Invalid email (missing local part)
            "password": "correct_password",
            "role": "USER"
        }`,
	},

	// --- Variant 3: Empty instances of fields (jsonBody and/or originalPass and/or userFromDB fields) ---

	// Empty originalPass (but valid userFromDB and jsonBody matching)
	{
		models.User{
			Email:    "user_empty_pass@domain.com",
			Role:     models.USER,
			Password: "hashed_empty_pass",
		},
		"", // originalPass is empty
		`{
            "email": "user_empty_pass@domain.com",
            "password": "",
            "role": "USER"
        }`,
	},

	// jsonBody with empty email
	{
		models.User{
			Email:    "some@user.com",
			Role:     models.USER,
			Password: "some_hashed_pass",
		},
		"some_pass",
		`{
            "email": "", // Empty email in JSON
            "password": "some_pass",
            "role": "USER"
        }`,
	},
	// jsonBody with empty password
	{
		models.User{
			Email:    "another@user.com",
			Role:     models.USER,
			Password: "another_hashed_pass",
		},
		"another_pass",
		`{
            "email": "another@user.com",
            "password": "", // Empty password in JSON
            "role": "USER"
        }`,
	},
	// jsonBody with empty role (assuming string type for role, will be empty string)
	{
		models.User{
			Email:    "admin@test.com",
			Role:     models.ADMIN,
			Password: "admin_hashed_pass",
		},
		"admin_pass",
		`{
            "email": "admin@test.com",
            "password": "admin_pass",
            "role": "" // Empty role in JSON
        }`,
	},

	// userFromDB with empty email
	{
		models.User{
			Email:    "", // Empty email in userFromDB
			Role:     models.USER,
			Password: "hashed_empty_db_email",
		},
		"some_pass",
		`{
            "email": "test@user.com",
            "password": "some_pass",
            "role": "USER"
        }`,
	},
	// userFromDB with empty password
	{
		models.User{
			Email:    "test@user.com",
			Role:     models.USER,
			Password: "", // Empty password in userFromDB
		},
		"some_pass",
		`{
            "email": "test@user.com",
            "password": "some_pass",
            "role": "USER"
        }`,
	},
	// userFromDB with empty role (if Role is string-based, it could be empty)
	{
		models.User{
			Email:    "test@user.com",
			Role:     "", // Empty role in userFromDB (if models.UserRole can be empty)
			Password: "some_hashed_pass",
		},
		"some_pass",
		`{
            "email": "test@user.com",
            "password": "some_pass",
            "role": "USER"
        }`,
	},

	// --- Variant 4: jsonBody loses some fields or has additional fields ---

	// jsonBody missing 'email' field
	{
		models.User{
			Email:    "test@example.com",
			Role:     models.USER,
			Password: "hashed_pass",
		},
		"test_pass",
		`{
            "password": "test_pass",
            "role": "USER"
        }`,
	},
	// jsonBody missing 'password' field
	{
		models.User{
			Email:    "test@example.com",
			Role:     models.USER,
			Password: "hashed_pass",
		},
		"test_pass",
		`{
            "email": "test@example.com",
            "role": "USER"
        }`,
	},
	// jsonBody missing 'role' field
	{
		models.User{
			Email:    "test@example.com",
			Role:     models.USER,
			Password: "hashed_pass",
		},
		"test_pass",
		`{
            "email": "test@example.com",
            "password": "test_pass"
        }`,
	},
	// jsonBody with an additional field
	{
		models.User{
			Email:    "test@example.com",
			Role:     models.USER,
			Password: "hashed_pass",
		},
		"test_pass",
		`{
            "email": "test@example.com",
            "password": "test_pass",
            "role": "USER",
            "extra_field": "some_value" // Additional field
        }`,
	},
	// jsonBody with multiple additional fields
	{
		models.User{
			Email:    "test@example.com",
			Role:     models.USER,
			Password: "hashed_pass",
		},
		"test_pass",
		`{
            "email": "test@example.com",
            "password": "test_pass",
            "role": "USER",
            "field1": "value1",
            "field2": 123
        }`,
	},
	// jsonBody with malformed JSON (not a missing field, but structural issue)
	{
		models.User{
			Email:    "test@example.com",
			Role:     models.USER,
			Password: "hashed_pass",
		},
		"test_pass",
		`{
            "email": "test@example.com",
            "password": "test_pass",
            "role": "USER" // Missing closing brace, or other syntax error
        `, // Malformed JSON
	},
	{
		models.User{
			Email:    "test@example.com",
			Role:     models.USER,
			Password: "hashed_pass",
		},
		"test_pass",
		`{"email": "test@example.com", "password": "test_pass", "role": "USER",}`, // Trailing comma (invalid in some parsers, though Go's JSON unmarshaler is forgiving)
	},
}

func TestUnitUserControllerTest(t *testing.T) {
	suite.Run(t, &UnitUserControllerTest{})
}

func (un *UnitUserControllerTest) SetupTest(){
	un.mockDB = new(mocks.IUserDataBase)
	un.mockPassServ = new(mocks.IPasswordService)
	un.mockJwtHandler = new(mocks.IAuthentication)
	un.controller = controllers.NewUserController(un.mockDB, un.mockPassServ, un.mockJwtHandler)
	
	un.router = gin.Default()
	un.router.GET("/login", un.controller.Login)
	un.router.GET("/register", un.controller.Register)

}
func ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ(jsonLiteral string, router *gin.Engine, lastUrl string) *httptest.ResponseRecorder {
	jsonBuffer := bytes.NewBufferString(jsonLiteral)

	req, _ := http.NewRequest("GET", lastUrl, jsonBuffer)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	return w
}

func (un *UnitUserControllerTest) TestLogin_Positive() {
	for _, each_test_case := range validTestCases {
		un.mockDB.On("FindUserByEmail", each_test_case.userFromDB.Email).Return(&each_test_case.userFromDB, nil)
		un.mockPassServ.On("IsCorrectPass", each_test_case.originalPass, (each_test_case.userFromDB.Password)).Return(true)
		un.mockJwtHandler.On("GenerateSecurityToken", mock.Anything).Return("simple_hash", time.Duration(24*time.Hour))
	}
	for _, each_test_case := range validTestCases {
		የሙከራ_ደንበኛ_መቀበያ := ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ(each_test_case.jsonBody, un.router, "/login")
		un.Require().Equal(የሙከራ_ደንበኛ_መቀበያ.Code, 200, "failure in the login", የሙከራ_ደንበኛ_መቀበያ.Body)

	}
}
func (un *UnitUserControllerTest) TestLogin_Negative() {
	for _, each_test_case := range invalidTestCases {
		un.mockDB.On("FindUserByEmail", each_test_case.userFromDB.Email).Return(&each_test_case.userFromDB, fmt.Errorf("some error"))
		un.mockPassServ.On("IsCorrectPass", each_test_case.originalPass, (each_test_case.userFromDB.Password)).Return(false)
		un.mockJwtHandler.On("GenerateSecurityToken", mock.Anything).Return("", time.Duration(0))
	}
	for index, each_test_case := range invalidTestCases {
		የሙከራ_ደንበኛ_መቀበያ := ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ(each_test_case.jsonBody, un.router, "/login")
		un.Require().NotEqual(የሙከራ_ደንበኛ_መቀበያ.Code, 200, "failure in the login in testcase ", index + 1)
	}
}
func (un *UnitUserControllerTest) TestRegister_Positive() {
	un.mockDB.On("StoreUser", mock.Anything).Return(nil) // since the address of two variables are different it won't work hear, that's why I used mock.Anything

	for _, each := range validTestCases {
		የአማርኛ_ቫርያብል_ይቻላል_እንዴ := ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ(each.jsonBody, un.router, "/register")
		un.EqualValues(http.StatusOK, የአማርኛ_ቫርያብል_ይቻላል_እንዴ.Code, የአማርኛ_ቫርያብል_ይቻላል_እንዴ.Body)
	}

}
func (un *UnitUserControllerTest) TestRegister_Negative() {
	un.mockDB.On("StoreUser", mock.Anything).Return(fmt.Errorf("")) // since the address of two variables are different it won't work hear, that's why I used mock.Anything

	for _, each := range invalidTestCases{
		የአማርኛ_ቫርያብል_ይቻላል_እንዴ := ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ(each.jsonBody, un.router, "/register")
		un.NotEqual(http.StatusOK, የአማርኛ_ቫርያብል_ይቻላል_እንዴ.Code, የአማርኛ_ቫርያብል_ይቻላል_እንዴ.Body)
	}

}