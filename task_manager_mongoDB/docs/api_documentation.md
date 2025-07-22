# Task Manager API Documentation

This documentation describes the API endpoints for a simple task manager application, enabling standard CRUD (Create, Read, Update, Delete) operations on tasks.

---

## Data Model

### Task

The core data structure is the **Task** model. Each task is uniquely identified by an **ID** (string) and includes a **Title** (string), a **Description** (string), a **DueDate** (time.Time type from Go's `time` package), and a **Status** (string, e.g., "Pending", "In Progress", "Completed").

---

## API Endpoints

The API is accessible at the base URL, typically `http://localhost:8081`.

### Get All Tasks

- **Endpoint:** `/tasks`
    
- **Method:** `GET`
    
- **Description:** Retrieves a complete list of all tasks currently stored in the system.
    
- **Response:**
    
    - **200 OK:** Returns a JSON array, where each element is a Task object.
        

### Get Task by ID

- **Endpoint:** `/tasks/:id`
    
- **Method:** `GET`
    
- **Description:** Fetches a single task based on its unique identifier. The `:id` in the URL should be replaced with the actual ID of the task you want to retrieve.
    
- **Response:**
    
    - **200 OK:** Returns a single JSON Task object if a matching task is found.
        
    - **404 Not Found:** Indicates that no task exists with the specified ID.
        

### Create New Task

- **Endpoint:** `/tasks/`
    
- **Method:** `POST`
    
- **Description:** Adds a new task to the system. The details of the new task are provided in the request body as a JSON object. The API checks for duplicate IDs to prevent conflicts.
    
- **Request Body:** Expects a JSON object conforming to the Task structure. The `ID` must be unique.
    
- **Response:**
    
    - **201 Created:** The new Task object is returned in JSON format, confirming successful creation.
        
    - **409 Conflict:** This status is returned if a task with the provided ID already exists, or if the JSON in the request body does not correctly match the Task structure.
        

### Update Task by ID

- **Endpoint:** `/tasks/:id`
    
- **Method:** `PUT`
    
- **Description:** Modifies an existing task identified by its ID. Currently, this endpoint allows for updates to the task's **Title** and **Description**. Other fields like `DueDate` or `Status` are not updated through this method in the current implementation.
    
- **Request Body:** A JSON object containing the fields to be updated (e.g., `title`, `description`).
    
- **Response:**
    
    - **200 OK:** A confirmation message indicating the task was successfully updated.
        
    - **400 Bad Request:** Occurs if the JSON format of the request body is incorrect or cannot be processed.
        
    - **404 Not Found:** If no task with the given ID is found to update.
        

### Delete Task by ID

- **Endpoint:** `/tasks/:id`
    
- **Method:** `DELETE`
    
- **Description:** Removes a task from the system using its unique identifier.
    
- **Response:**
    
    - **200 OK:** A confirmation message indicating that the task was successfully removed.
        
    - **404 Not Found:** If no task with the specified ID is found for deletion.
        

---

## Running the API

To start the API server, simply run the `main.go` file. The server will begin listening for requests on port `8081`. You can then interact with the API endpoints using tools like cURL or Postman.
##### PS: I didn't too much from the last task, I simply refactred it with the required folder structure 👍![alt text](image.png)