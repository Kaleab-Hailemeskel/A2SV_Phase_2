

# Task Management API Documentation

📄 **Get started here** This document provides a comprehensive guide to using the Task Management API, including how to authenticate, manage tasks, and handle common scenarios.




## Overview of the Task Management API

The **Task Management API** provides a robust set of tools and resources that enable you to **programmatically manage tasks and user accounts**. This includes creating, reading, updating, and deleting tasks, as well as handling user registration and authentication.

## Getting Started Guide

To start using the Task Management API, you need to follow these essential steps:

1. **Understand Authentication**: All protected API endpoints require proper authentication using a JSON Web Token (JWT).
    
2. **Obtain a JWT**: After registering, you'll need to log in to receive a JWT as an `HttpOnly` cookie. This cookie will automatically be sent with subsequent requests to authenticated endpoints by most HTTP clients (like Postman or web browsers).
    
3. **Use HTTPS**: The API is designed to respond only to HTTPS-secured communications. Any requests sent via HTTP may be redirected. (Note: Your provided `router.go` runs on HTTP. For production, ensure HTTPS is configured.)
    
4. **JSON Format**: The API expects request bodies and returns responses primarily in JSON format. Ensure your `Content-Type` header is set to `application/json` for requests with a body.
    
5. **Error Handling**: When an API request returns an error, it is sent in the JSON response, typically with an `"error"` or `"message"` key, along with an appropriate HTTP status code.
    

## Authentication

The **Task Management API** uses **JWT (JSON Web Token) via HttpOnly cookies** for authentication.

### How to Authenticate:

1. **Register**: First, register a new user by sending a `POST` request to `/register` with your `email` and `password`.
    
    - **Endpoint**: `POST /register`
        
    - **Request Body**:
        
        JSON
        
        ```
        {
            "email": "your_email@example.com",
            "password": "your_strong_password"
        }
        ```
        
    - **Password Requirements**: Minimum 8 characters, maximum 20 characters.
        
    - **Email Requirements**: Must be a valid email format.
        
2. **Log In**: After successful registration (or if you already have an account), log in by sending a `POST` request to `/login` with your `email` and `password`.
    
    - **Endpoint**: `POST /login`
        
    - **Request Body**:
        
        JSON
        
        ```
        {
            "email": "your_email@example.com",
            "password": "your_password"
        }
        ```
        
    - **Response**: A successful login will not return a JSON body, but it will set an `HttpOnly` cookie in your browser/client containing the JWT. This cookie is automatically managed by clients and sent with subsequent authenticated requests.
        

### Authentication Error Response

If a JWT cookie is missing, expired, or invalid, you will typically receive an **HTTP 401 Unauthorized** response code, and requests to protected endpoints will fail.

### Authorization (Roles)

The API supports two user roles: `USER` and `ADMIN`.

- **USER**: Can only manage (create, read, update, delete) tasks that they own. They can view their own user information.
    
- **ADMIN**: Can manage all tasks in the system (create, read, update, delete any task) and view any user's information.
    

The `Authorization` middleware (`middleware/user_auth.go`) handles checking if the authenticated user has the necessary permissions to perform an action on a specific task (e.g., ensuring a `USER` only modifies their own tasks).

## Rate and Usage Limits

This API does not currently implement explicit rate limiting. However, it's good practice to design your client applications to be respectful of server resources. Unexpected spikes in traffic may lead to server instability.

### 503 Response

An HTTP **503 Service Unavailable** response from the server indicates an unexpected issue, possibly due to an unexpected spike in API access traffic or a temporary server outage. If you encounter a 503 error, please wait a few minutes and retry your request. If the outage persists or you receive any other form of an HTTP 5XX error, contact support.

## API Endpoints

### User Management

#### `POST /register`

- **Description**: Registers a new user.
    
- **Authentication**: None required.
    
- **Request Body (UserDTO)**:
    
    JSON
    
    ```
    {
        "email": "string",  // required, valid email format
        "password": "string" // required, min=8, max=20
    }
    ```
    
- **Responses**:
    
    - `200 OK`: `{"message": "User registered successfully"}`
        
    - `406 Not Acceptable`: `{"Error": "Invalid User type"}` or `{"Error": "User already exists"}`
        
    - `400 Bad Request`: `{"Error": "validation error message"}`
        

#### `POST /login`

- **Description**: Logs in a user and sets an HttpOnly JWT cookie.
    
- **Authentication**: None required.
    
- **Request Body (UserDTO)**:
    
    JSON
    
    ```
    {
        "email": "string",  // required
        "password": "string" // required
    }
    ```
    
- **Responses**:
    
    - `200 OK`: Sets `jwt_token` cookie. No JSON response body.
        
    - `400 Bad Request`: `{"error": "Invalid request payload"}` or `{"error": "incorrect email or Password"}`
        
    - `404 Not Found`: `{"error": "user didn't exist, better register"}`
        

#### `GET /whoAmI`

- **Description**: Retrieves information about the currently authenticated user.
    
- **Authentication**: Required (valid JWT cookie).
    
- **Responses**:
    
    - `202 Accepted`: `{"Current User": {"Email": "string", "Password": "**HIDDEN**", "Role": "string"}}`
        
    - `401 Unauthorized`: (If JWT is missing/invalid)
        
    - `400 Bad Request`: `{"message": "IMPOSSIBLEEEEE"}` (Internal error, should not occur normally)
        

### Task Management

#### `POST /tasks`

- **Description**: Creates a new task for the authenticated user.
    
- **Authentication**: Required (valid JWT cookie).
    
- **Request Body (TaskDTO)**:
    
    JSON
    
    ```
    {
        "id": "string",        // required, unique task ID
        "title": "string",     // required, min=3, max=100
        "description": "string", // required
        "due_date": "2006-01-02T15:04:05Z07:00", // required, ISO 8601 format
        "status": "string"     // required
    }
    ```
    
- **Responses**:
    
    - `201 Created`: Returns the created `Task` object.
        
    - `400 Bad Request`: `{"error": "validation error message"}`
        
    - `409 Conflict`: `{"message": "Can't save a new task"}` (Binding error)
        
    - `500 Internal Server Error`: `{"error": "task with ID <id> Already exists."}` or other database errors.
        

#### `GET /tasks`

- **Description**: Retrieves a list of tasks. For regular users, this returns only tasks owned by them. For ADMIN users, it returns all tasks in the database.
    
- **Authentication**: Required (valid JWT cookie).
    
- **Responses**:
    
    - `200 OK`: `[Task, Task, ...]` (Array of Task objects).
        
    - `404 Not Found`: `{"message": "No task Found"}`
        
    - `500 Internal Server Error`: `{"error": "database error message"}`
        

#### `GET /tasks/:id`

- **Description**: Retrieves a single task by its ID. Users can only access their own tasks unless they are an ADMIN.
    
- **Authentication**: Required (valid JWT cookie).
    
- **Path Parameters**:
    
    - `id`: The unique ID of the task.
        
- **Responses**:
    
    - `200 OK`: Returns the `Task` object.
        
    - `404 Not Found`: `{"message": "No task Found with ID<id>"}`
        
    - `403 Forbidden`: `{"error": "You are not authorized to access this task"}` (If authenticated but not owner/admin).
        
    - `500 Internal Server Error`: `{"error": "database error message"}`
        

#### `PUT /tasks/:id`

- **Description**: Updates an existing task by its ID. Users can only update their own tasks unless they are an ADMIN.
    
- **Authentication**: Required (valid JWT cookie).
    
- **Path Parameters**:
    
    - `id`: The unique ID of the task to update.
        
- **Request Body (TaskDTO)**: (Same structure as `POST /tasks`)
    
    JSON
    
    ```
    {
        "id": "string",        // required, can be the same as path ID or a new unique ID
        "title": "string",     // required, min=3, max=100
        "description": "string", // required
        "due_date": "2006-01-02T15:04:05Z07:00", // required, ISO 8601 format
        "status": "string"     // required
    }
    ```
    
- **Responses**:
    
    - `200 OK`: `{"message": "Task updated"}`
        
    - `400 Bad Request`: `{"error": "validation error message"}`
        
    - `403 Forbidden`: `{"error": "You are not authorized to access this task"}` (If authenticated but not owner/admin).
        
    - `500 Internal Server Error`: `{"error": "no Data Modified"}` or `{"error": "task With The new ID already exists, Use Unique ID"}` or other database errors.
        

#### `DELETE /tasks/:id`

- **Description**: Deletes a task by its ID. Users can only delete their own tasks unless they are an ADMIN.
    
- **Authentication**: Required (valid JWT cookie).
    
- **Path Parameters**:
    
    - `id`: The unique ID of the task to delete.
        
- **Responses**:
    
    - `202 Accepted`: `{"message": "Task with ID <id> got Deleted"}`
        
    - `400 Bad Request`: `{"Error": "task with ID <id> isn't found"}`
        
    - `403 Forbidden`: `{"error": "You are not authorized to access this task"}` (If authenticated but not owner/admin).
        

## Need Some Help?

In case you have questions or encounter issues:

- **Review the Code**: Refer to the Go source code you provided for detailed implementation logic in `controllers/`, `models/`, `data/`, `infrastructure/`, and `middleware/`.
    
- **Check Server Logs**: The server outputs useful information to the console, including environment variable loading and database connection status. Errors are also logged.
    
- **Database Inspection**: If you have direct access to your MongoDB instance, you can inspect the `users` and `tasks` collections for data integrity.