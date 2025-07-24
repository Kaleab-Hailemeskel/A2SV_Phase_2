# 🛡️ Task Management API with Authentication & Authorization

This project is a secure task management API built using **Go**, **Gin**, **MongoDB**, and **JWT-based authentication**. It provides user registration, login, task creation, editing, deletion, and retrieval — all gated by access control mechanisms.

---
## 🚀 Technologies Used

- **Go** — main programming language
- **Gin** — web framework for routing and middleware
- **MongoDB** — persistent storage for users and tasks
- **bcrypt** — password hashing
- **JWT (JSON Web Tokens)** — authentication mechanism
---
## 🗺️ Project Structure
```
task-6_authentication_and_authorization/ 
	│ 
	├── controllers/ → API endpoint handlers (task & user operations) 
	├── data/ → MongoDB connectivity and database actions 
	├── middleware/ → JWT authentication and role-based access control 
	├── models/ → Struct definitions (User, Task) 
	├── router/ → Gin engine setup and graceful shutdown 
	└── main.go → App entry point

```
---
## 👥 User Model

```go
type User struct {
  Email    string `json:"email" binding:"required,email"`
  Password string `json:"password" binding:"required"`
  Role     string `json:"role" binding:"required"` // ADMIN or USER
}
````

---
## 📋 Task Model

```go
type Task struct {
  ID          string    `json:"id" binding:"required"`
  OwnerEmail  string    `json:"ownerEmail" binding:"required,email"`
  Title       string    `json:"title" binding:"required,min=3,max=100"`
  Description string    `json:"description" binding:"required"`
  DueDate     time.Time `json:"due_date" binding:"required"`
  Status      string    `json:"status" binding:"required"`
}
```

---
## ✅ API Endpoints

### 🔐 Auth Routes

|Method|Route|Description|
|---|---|---|
|POST|`/register`|Register a new user|
|POST|`/login`|Login & get JWT|

### 👤 User Info

|Method|Route|Middleware|Description|
|---|---|---|---|
|GET|`/whoAmI`|`Authentication`|Return user info|

### 📋 Task Routes

All under `/tasks` and gated by `Authentication` middleware.

|Method|Route|Extra Middleware|Description|
|---|---|---|---|
|GET|`/tasks`|-|Get tasks (admin sees all)|
|POST|`/tasks`|-|Create new task|
|GET|`/tasks/:id`|`IsAuthorizedUserForTaskManipulation`|Get task by ID|
|PUT|`/tasks/:id`|`IsAuthorizedUserForTaskManipulation`|Update task|
|DELETE|`/tasks/:id`|`IsAuthorizedUserForTaskManipulation`|Delete task|

---
## 🔐 Middleware Logic

### Authentication

- Verifies JWT from cookie
- Validates expiration time
- Fetches user details from DB
- Injects `currUser` into context

### Authorization

- **IsAdmin**: ensures `ADMIN` role
- **IsAuthorizedUserForTaskManipulation**: user must be owner or admin

---

## 📦 Database Logic

### Task Operations (`data/task_service.go`)

- Connects to `TaskBase.Tasks`
- Functions:
    - `InsertOne(task)`
    - `FindALL(email)`
    - `FindByID(id)`
    - `UpdateOne(id, updatedTask)`
    - `DeleteOne(id)`

### User Operations (`data/user_service.go`)

- Connects to `UserBase.users`
- Functions:
    - `InsertOneUser(user)`
    - `FindOneUser(email)`

---

## 🔄 Server Lifecycle

Defined in `router/StartEngine(port)`:

- Initializes Gin routes
- Connects to MongoDB
- Awaits user ENTER key to gracefully shut down
- Closes MongoDB connection

---

## 🏁 Getting Started

1. Make sure MongoDB is running locally at `mongodb://localhost:27017`
2. Run the app:
    
    ```bash
    go run main.go
    ```
    
3. Register a user via `/register`
4. Login via `/login` to receive JWT cookie
5. Interact with `/tasks` endpoints

---
## 🙌 Author

This project was built by **Kaleab** — a fourth-year Software Engineering student at AASTU, deeply passionate about embedded systems, backend architecture, and building secure APIs with Golang.

---

