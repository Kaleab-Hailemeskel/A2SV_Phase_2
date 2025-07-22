## Task Manager API Documentation

This document describes the RESTful API for managing tasks, built with Go and leveraging MongoDB for data persistence.

### 1. API Overview

The Task Manager API provides a set of endpoints to perform standard Create, Read, Update, and Delete (CRUD) operations on task resources. Each task is identified by a unique ID and includes properties such as title, description, due date, and status.

### 2. API Endpoints

The API base URL is assumed to be `http://localhost:8081`.

---

#### `POST /tasks/` - Create a New Task

- **Description:** Creates a new task record in the database. The task's ID must be unique.
    
- **HTTP Method:** `POST`
    
- **Request Path:** `/tasks/`
    
- **Request Body (JSON):**
    
    JSON
    
    ```
    {
        "id": "string",
        "title": "string",
        "description": "string",
        "due_date": "YYYY-MM-DDTHH:MM:SSZ",
        "status": "string"
    }
    ```
    
    - `id`: (Required, string) A unique identifier for the task.
        
    - `title`: (Required, string) The title of the task.
        
    - `description`: (Optional, string) A detailed description of the task.
        
    - `due_date`: (Required, string, ISO 8601 format) The deadline for the task.
        
    - `status`: (Required, string) The current status of the task (e.g., "pending", "completed").
        
- **Responses:**
    
    - **`201 Created`**: Task successfully created.
        
        JSON
        
        ```
        {
            "id": "task_001",
            "title": "Complete Project Report",
            "description": "Write and finalize the quarterly project report for Q3.",
            "due_date": "2025-07-30T17:00:00Z",
            "status": "pending"
        }
        ```
        
    - **`400 Bad Request`**: If the request body is malformed or invalid JSON.
        
        JSON
        
        ```
        {"message": "Invalid request body: json: cannot unmarshal string into Go struct field Task.due_date of type time.Time"}
        ```
        
    - **`409 Conflict`**: If a task with the provided `id` already exists.
        
        JSON
        
        ```
        {"error": "task with ID < task_001 > Already exists."}
        ```
        
    - **`500 Internal Server Error`**: If a database operation fails.
        
        JSON
        
        ```
        {"error": "Failed to save new task: <database_error_details>"}
        ```
        
    - **`503 Service Unavailable`**: If the server cannot connect to the database.
        
        JSON
        
        ```
        {"Error": "Database Connection Unavailable"}
        ```
        

---

#### `GET /tasks` - Retrieve All Tasks

- **Description:** Fetches all task records currently stored in the database.
    
- **HTTP Method:** `GET`
    
- **Request Path:** `/tasks`
    
- **Request Body:** None
    
- **Responses:**
    
    - **`200 OK`**: Successfully retrieved tasks.
        
        JSON
        
        ```
        [
            {
                "id": "task_001",
                "title": "Complete Project Report",
                "description": "Write and finalize the quarterly project report for Q3.",
                "due_date": "2025-07-30T17:00:00Z",
                "status": "pending"
            },
            {
                "id": "task_002",
                "title": "Review Code",
                "description": "Review pull request #123 for feature X.",
                "due_date": "2025-07-25T10:00:00Z",
                "status": "in_progress"
            }
        ]
        ```
        
    - **`200 OK`**: If no tasks are found.
        
        JSON
        
        ```
        {"message": "No task Found"}
        ```
        
    - **`500 Internal Server Error`**: If a database operation fails.
        
        JSON
        
        ```
        {"error": "Error while fetching All Data from DB: <database_error_details>"}
        ```
        
    - **`503 Service Unavailable`**: If the server cannot connect to the database.
        
        JSON
        
        ```
        {"Error": "Database Connection Unavailable"}
        ```
        

---

#### `GET /tasks/:id` - Retrieve Task by ID

- **Description:** Fetches a single task record using its unique ID.
    
- **HTTP Method:** `GET`
    
- **Request Path:** `/tasks/{id}` (e.g., `/tasks/task_001`)
    
- **Request Body:** None
    
- **Responses:**
    
    - **`200 OK`**: Task successfully retrieved.
        
        JSON
        
        ```
        {
            "id": "task_001",
            "title": "Complete Project Report",
            "description": "Write and finalize the quarterly project report for Q3.",
            "due_date": "2025-07-30T17:00:00Z",
            "status": "pending"
        }
        ```
        
    - **`200 OK`**: If no task is found with the given ID.
        
        JSON
        
        ```
        {"message": "No task Found with ID task_005"}
        ```
        
    - **`500 Internal Server Error`**: If a database operation fails.
        
        JSON
        
        ```
        {"error": "Failed to find task: <database_error_details>"}
        ```
        
    - **`503 Service Unavailable`**: If the server cannot connect to the database.
        
        JSON
        
        ```
        {"Error": "Database Connection Unavailable"}
        ```
        

---

#### `PUT /tasks/:id` - Update Task by ID

- **Description:** Updates an existing task identified by its ID with the provided data. The entire task object should be sent.
    
- **HTTP Method:** `PUT`
    
- **Request Path:** `/tasks/{id}` (e.g., `/tasks/task_001`)
    
- **Request Body (JSON):** Same structure as `POST /tasks/`.
    
    - Note: The `id` in the URL path (`:id`) specifies the task to be updated. The `id` in the request body (`updatedTask.ID`) is the _new_ ID for the task if you intend to change it.
        
- **Responses:**
    
    - **`200 OK`**: Task successfully updated.
        
        JSON
        
        ```
        {"message": "Task updated successfully"}
        ```
        
    - **`400 Bad Request`**: If the request body is malformed or invalid JSON.
        
        JSON
        
        ```
        {"error": "Invalid request body: json: cannot unmarshal number into Go struct field Task.id of type string"}
        ```
        
    - **`404 Not Found`**: If no task with the original ID is found to update.
        
        JSON
        
        ```
        {"error": "no task with ID task_999 found to update"}
        ```
        
    - **`409 Conflict`**: If the `id` in the request body is changed and already exists for another task.
        
        JSON
        
        ```
        {"error": "task With The new ID already exists, Use Unique ID"}
        ```
        
    - **`500 Internal Server Error`**: If a database operation fails, or no data was actually modified (e.g., you sent the exact same data).
        
        JSON
        
        ```
        {"error": "Failed to update task: <database_error_details>"}
        {"error": "no Data modified for task with ID task_001"}
        ```
        
    - **`503 Service Unavailable`**: If the server cannot connect to the database.
        
        JSON
        
        ```
        {"Error": "Database Connection Unavailable"}
        ```
        

---

#### `DELETE /tasks/:id` - Delete Task by ID

- **Description:** Deletes a task record from the database using its unique ID.
    
- **HTTP Method:** `DELETE`
    
- **Request Path:** `/tasks/{id}` (e.g., `/tasks/task_001`)
    
- **Request Body:** None
    
- **Responses:**
    
    - **`202 Accepted`**: Task successfully deleted.
        
        JSON
        
        ```
        {"message": "Task with ID task_001 got Deleted"}
        ```
        
    - **`404 Not Found`**: If no task with the provided ID is found.
        
        JSON
        
        ```
        {"Error": "task with ID task_999 isn't found"}
        ```
        
    - **`500 Internal Server Error`**: If a database operation fails.
        
        JSON
        
        ```
        {"Error": "Failed to delete task: <database_error_details>"}
        ```
        
    - **`503 Service Unavailable`**: If the server cannot connect to the database.
        
        JSON
        
        ```
        {"Error": "Database Connection Unavailable"}
        ```
        

---

### 3. Database Integration

The Task Manager API utilizes **MongoDB** as its persistent data store.

- **Connection Details:**
    
    - **Connection String:** `mongodb://localhost:27017`
        
    - **Database Name:** `TaskBase`
        
    - **Collection Name:** `Tasks`
        
    - These details are configured as constants within the `data` package (`task_service.go`).
        
- **Connection Management:**
    
    - The `data.ConnectToMongo()` function establishes the initial connection to the MongoDB server during application startup. It includes a ping to ensure the connection is active.
        
    - The `data.IsClientConnected()` function provides a basic check of the client's connection status before any database operation is attempted by the controllers.
        
    - Crucially, `data.CloseMongoDB()` is called during the application's graceful shutdown process to ensure that all open database connections are properly closed and resources are released.
            
- **Data Operations (CRUD):** The `data` package (`task_service.go`) encapsulates all direct interactions with the MongoDB collection:
    
    - `InsertOne(task models.Task)`: Inserts a single new task document. Includes a check to prevent duplicate `id` values.
        
    - `FindALL() ([]models.Task, error)`: Retrieves all documents from the `Tasks` collection.
        
    - `FindByID(taskID string) (*models.Task, error)`: Retrieves a single document by its `id`. Returns `mongo.ErrNoDocuments` if not found.
        
    - `UpdateOne(taskID string, updatedTask models.Task)`: Updates a single document identified by `taskID`. Includes checks for existing new `ID` during update and whether any data was actually modified.
        
    - `DeleteOne(taskID string)`: Deletes a single document by its `id`.


<div style="text-align: center;"> 
	![alt text](image.png)
</div>

👍