# Task Management API Test Documentation

This document outlines the unit tests implemented for the Task Management API. These tests ensure the correct functionality of individual components (controllers and middleware) by isolating them and mocking their dependencies.

---

## Testing Frameworks and Tools

- **Go's `testing` package**: The standard library for writing unit tests in Go.
    
- **`github.com/stretchr/testify/suite`**: Provides a testing suite framework, allowing for better test organization and setup/teardown methods.
    
- **`github.com/stretchr/testify/mock`**: Used for creating mock objects, which are essential for isolating units under test and controlling their dependencies.
    
- **`github.com/gin-gonic/gin`**: The Gin web framework's testing utilities are used to simulate HTTP requests and responses.
    
- **`github.com/joho/godotenv`**: Used to load environment variables for testing, similar to how the main application loads them.
    

---

## Unit Test Coverage

The following sections detail the unit tests for key components of the Task Management API:

### 1. Task Controller Tests (`task_controller_test.go`)

These tests focus on the `TaskController`'s ability to handle task-related requests (Create, Read, Update, Delete) and interact with the database layer. A mock database (`mocks.ITaskDataBase`) is used to control the database responses, allowing for isolated testing of the controller logic.

#### `UnitTaskControllerTest` Suite

This suite sets up a Gin router and injects a mocked `ITaskDataBase` into the `TaskController`. A helper function `በዩርል_በኩል_ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ` (which translates to "send_test_client_to_json_via_url") is used to simplify sending HTTP requests with JSON payloads.

- **`SetupTest()`**:
    
    - Initializes a new mock `ITaskDataBase`.
        
    - Creates a `TaskController` with the mock.
        
    - Sets up a `gin.Engine` with specific routes (`/task/:id` for DELETE, GET, PUT and `/task` for GET, POST) that point to the controller's methods.
        
    - Crucially, it sets a `currTestUser` in the Gin context for `/task` routes, simulating an authenticated user for testing authorization logic.
        

#### Test Cases:

- **`TestDeleteByID_Positive`**:
    
    - **Scenario**: Tests successful deletion of a task by ID.
        
    - **Mocks**: `mockDB.On("DeleteOne", "task1").Return(nil)`
        
    - **Assertion**: Expects an `HTTP 202 Accepted` status code.
        
- **`TestDeleteByID_Negative`**:
    
    - **Scenario**: Tests deletion when the database operation fails (e.g., task not found, internal DB error).
        
    - **Mocks**: `mockDB.On("DeleteOne", "task1").Return(fmt.Errorf(""))`
        
    - **Assertion**: Expects a status code other than `HTTP 202 Accepted` or `HTTP 200 OK`.
        
- **`TestGetByID_Positive`**:
    
    - **Scenario**: Tests successful retrieval of a single task by ID.
        
    - **Mocks**: `mockDB.On("FindByID", "task1").Return(&models.Task{}, nil)` (returns a sample task)
        
    - **Assertion**: Expects an `HTTP 200 OK` status code.
        
- **`TestGetByID_Negative`**:
    
    - **Scenario**: Tests retrieval when the task is not found in the database.
        
    - **Mocks**: `mockDB.On("FindByID", "task1").Return(nil, nil)`
        
    - **Assertion**: Expects a status code other than `HTTP 202 Accepted` or `HTTP 200 OK`.
        
- **`TestGetTask_Positive`**:
    
    - **Scenario**: Tests successful retrieval of all tasks. Covers both `USER` and `ADMIN` roles, verifying that `FindAllTasks` is called appropriately (with `OwnerEmail` for `USER`, or empty for `ADMIN`).
        
    - **Mocks**: `mockDB.On("FindAllTasks", "").Return(&[]models.Task{}, nil)` and `mockDB.On("FindAllTasks", "oneLove@gmail.com").Return(&[]models.Task{}, nil)`
        
    - **Assertion**: Expects `HTTP 200 OK` for both user roles.
        
- **`TestGetTask_Negative`**:
    
    - **Scenario**: Tests task retrieval when a database error occurs.
        
    - **Mocks**: `mockDB.On("FindAllTasks", "").Return(nil, fmt.Errorf("mock Error"))` and `mockDB.On("FindAllTasks", "oneLove@gmail.com").Return(nil, fmt.Errorf("mock Error"))`
        
    - **Assertion**: Expects a status code other than `HTTP 200 OK`.
        
- **`TestPutTaskByID_Positive`**:
    
    - **Scenario**: Tests successful update of an existing task.
        
    - **Mocks**: `mockDB.On("UpdateOne", "task1", mock.AnythingOfType("models.Task")).Return(nil)`
        
    - **Assertion**: Expects an `HTTP 200 OK` status code.
        
- **`TestPutTaskByID_Negative`**:
    
    - **Scenario**: Tests task update when a database error occurs.
        
    - **Mocks**: `mockDB.On("UpdateOne", "task1", mock.AnythingOfType("models.Task")).Return(fmt.Errorf("mock Error"))`
        
    - **Assertion**: Expects a status code other than `HTTP 200 OK`.
        
- **`TestPostTask_Positive`**:
    
    - **Scenario**: Tests successful creation of a new task.
        
    - **Mocks**: `mockDB.On("InsertOne", mock.AnythingOfType("models.Task")).Return(nil)`
        
    - **Assertion**: Expects an `HTTP 201 Created` status code.
        
- **`TestPostTask_Negative`**:
    
    - **Scenario**: Tests task creation when the request body is invalid (e.g., missing required fields).
        
    - **Mocks**: None (the test focuses on Gin's binding/validation errors before DB interaction).
        
    - **Assertion**: Expects a status code other than `HTTP 201 Created`.
        

---

### 2. User Controller Tests (`user_controller_test.go`)

These tests validate the `UserController`'s functionality for user registration and login, ensuring proper interaction with the user database, password hashing service, and JWT generation. Mock versions of `IUserDataBase`, `IPasswordService`, and `IAuthentication` are used.

#### `UnitUserControllerTest` Suite

This suite sets up a Gin router and injects mocked services into the `UserController`. A helper function `ወደ_ጄሰን_የሙከራ_ደንበኛ_ላክ` (which translates to "send_to_json_test_client") is used for sending requests.

- **`SetupTest()`**:
    
    - Initializes mocks for `IUserDataBase`, `IPasswordService`, and `IAuthentication`.
        
    - Creates a `UserController` with these mocks.
        
    - Sets up a `gin.Engine` with `/login` and `/register` routes.
        

#### Test Data:

- **`validTestCases`**: An array of structs containing valid user data for positive test scenarios, including `User` models, original passwords, and corresponding JSON bodies.
    
- **`invalidTestCases`**: An array of structs containing various invalid scenarios, such as password mismatches, invalid email formats, empty fields, and malformed JSON.
    

#### Test Cases:

- **`TestLogin_Positive`**:
    
    - **Scenario**: Iterates through `validTestCases` to ensure successful user login.
        
    - **Mocks**:
        
        - `mockDB.On("FindUserByEmail", ...).Return(&userFromDB, nil)`
            
        - `mockPassServ.On("IsCorrectPass", originalPass, hashedPass).Return(true)`
            
        - `mockJwtHandler.On("GenerateSecurityToken", ...).Return("simple_hash", ...)`
            
    - **Assertion**: Expects an `HTTP 200 OK` status code for each valid test case.
        
- **`TestLogin_Negative`**:
    
    - **Scenario**: Iterates through `invalidTestCases` to verify that login fails under various erroneous conditions.
        
    - **Mocks**: Configured to simulate database errors (`FindUserByEmail` returns error), incorrect password validation (`IsCorrectPass` returns `false`), and failed token generation.
        
    - **Assertion**: Expects a status code other than `HTTP 200 OK` for each invalid test case.
        
- **`TestRegister_Positive`**:
    
    - **Scenario**: Tests successful user registration.
        
    - **Mocks**: `mockDB.On("StoreUser", mock.Anything).Return(nil)` (using `mock.Anything` because the exact `models.User` object will be different due to hashing).
        
    - **Assertion**: Expects an `HTTP 200 OK` status code for each valid test case.
        
- **`TestRegister_Negative`**:
    
    - **Scenario**: Tests user registration failure due to invalid input.
        
    - **Mocks**: `mockDB.On("StoreUser", mock.Anything).Return(fmt.Errorf(""))` (simulates a database error, or a user already exists).
        
    - **Assertion**: Expects a status code other than `HTTP 200 OK` for each invalid test case.
        

---

### 3. Middleware Tests (`auth_middleware_test.go`)

These tests verify the behavior of the `UserAuth` middleware, which handles authentication (verifying JWTs) and authorization (checking user permissions for tasks). Mock services for JWT handling, user database, and task database are used.

#### `UnitTestAuthMiddleware` Suite

This suite sets up the Gin router and injects mocked authentication and database services into the `UserAuth` middleware. A `sendJSONRequest` helper is used to simulate requests, including the crucial ability to set cookies.

- **`SetupSuite()`**: Loads environment variables using `godotenv.Load()` and initializes application-specific environment variables via `infrastructure.InitEnv()`. This ensures that JWT secret keys and other configurations are available during tests.
    
- **`SetupTest()`**:
    
    - Initializes mocks for `IAuthentication`, `IUserDataBase`, and `ITaskDataBase`.
        
    - Creates a `UserAuth` middleware instance.
        
    - Sets up a `gin.Engine` with two test routes: `/authentication_route` (for testing the `Authentication` middleware) and `/authorization_route` (for testing the `Authorization` middleware). A `middlewareFiller` is used to pre-set a `CURR_USER` in the context before the `Authorization` middleware runs.
        

#### Test Data:

- **`validTestCases_authentication`**: Defines valid JWT scenarios for the authentication middleware.
    
- **`invalidTestCases_authentication`**: Defines invalid JWT scenarios (expired, user not found, malformed token, missing claims) for the authentication middleware.
    
- **`validTestCases_authorization`**: Defines scenarios where a user is authorized to access a task.
    
- **`invalidTestCases_authorization_authorization`**: Defines scenarios where authorization should fail (non-existent task, internal DB error).
    

#### Test Cases:

- **`TestAuthentication` (General Positive Case)**:
    
    - **Scenario**: A basic test to confirm that a valid JWT token in a cookie successfully authenticates a user.
        
    - **Mocks**: Configures `mockJwt` to parse a valid token, verify it's not expired, and extract the user's email. Configures `mockUserDB` to find the user.
        
    - **Assertion**: Expects an `HTTP 200 OK` status code.
        
- **`TestAuthentication_Positive`**:
    
    - **Scenario**: Iterates through `validTestCases_authentication` to test various valid authentication scenarios (e.g., ADMIN, USER roles).
        
    - **Mocks**: Dynamically sets up mocks for `ParseToken`, `TokenExpired`, `GetUserEmailFromSecurityToken`, and `FindUserByEmail` based on the test case data.
        
    - **Assertion**: Expects the `expectedStatus` (which is `HTTP 200 OK`) for each valid case.
        
- **`TestAuthentication_Negative`**:
    
    - **Scenario**: Iterates through `invalidTestCases_authentication` to test various invalid authentication scenarios (expired token, user not found, malformed token, missing email claim).
        
    - **Mocks**: Dynamically sets up mocks to simulate the specific error conditions for each invalid test case (e.g., `TokenExpired` returning `true`, `FindUserByEmail` returning an error, `ParseToken` returning an error).
        
    - **Assertion**: Expects a status code other than `HTTP 200 OK` for each invalid case.
        
- **`TestAuthorization_Positive`**:
    
    - **Scenario**: Tests successful authorization where the user is either the task owner or an ADMIN.
        
    - **Mocks**: `mockTaskDB.On("FindByID", ...).Return(task, nil)` (returning a valid task).
        
    - **Assertion**: Expects an `HTTP 200 OK` status code.
        
- **`TestAuthorization_Negative`**:
    
    - **Scenario**: Tests authorization failures, such as when the task does not exist or a database error occurs during task retrieval.
        
    - **Mocks**: `mockTaskDB.On("FindByID", ...).Return(nil, errors.New(...))` (simulating not found or internal errors).
        
    - **Assertion**: Expects a status code other than `HTTP 200 OK`. (Note: The specific forbidden cases, where a non-owner/non-admin tries to access another user's task, would depend on the `middlewareFiller` setting the `CURR_USER` to a non-admin and the `FindByID` returning a task not owned by that user.)