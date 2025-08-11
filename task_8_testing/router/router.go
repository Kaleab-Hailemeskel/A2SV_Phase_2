package router

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"task_8_testing/useCaseF"
	"task_8_testing/controllers"
	"task_8_testing/data"
	"task_8_testing/infrastructure"
	"task_8_testing/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

// GET, GET_ID, POST, PUT, DELETE methods are mapped with their counter func
func StartEngine(port_number string) {
	router := gin.Default()
	log.Println("before Init")

	newUserDB := data.NewUserDataBase()
	newTaskDB := data.NewTaskDataBaseService()
	newJwtAuth := infrastructure.NewJWTAuth()
	newPassServ := infrastructure.NewPasswordService()
	usecase := useCaseF.NewUseCase(newUserDB, newTaskDB, newJwtAuth, newPassServ)
	userAuth := middleware.NewUserAuth(newJwtAuth, newUserDB, newTaskDB)
	userC := controllers.NewUserController(usecase)
	taskC := controllers.NewTaskController(usecase)
	StartPublicRouter(router, userC)
	StartProtectedRouter(router, userAuth, taskC, userC)

	// ? CODES AFTER THIS ARE HELPERS

	fmt.Println()

	fmt.Println("********************************")
	fmt.Println("* DataBase Connection Initiated*")
	fmt.Println("********************************")

	fmt.Println()

	addr := fmt.Sprintf(":%s", port_number)

	// Create an http.Server instance. This allows us to control its lifecycle.
	srv := &http.Server{
		Addr:    addr,
		Handler: router, // Gin router acts as the HTTP handler
	}

	// Channel to signal errors from the server's goroutine
	serverErrors := make(chan error, 1)

	// --- Start the Gin server in a goroutine ---
	go func() {
		log.Printf("Gin server listening on %s\n", addr)
		// ListenAndServe blocks until server is shut down or an error occurs.
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			// If the error is not http.ErrServerClosed (which is expected during shutdown),
			// then it's a real startup or runtime error we should capture.
			serverErrors <- err
		}
	}()

	// --- Interaction to trigger shutdown ---
	fmt.Println("Press ENTER in THIS console to gracefully shut down the server...")

	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n') // Read until newline (user presses Enter)

	log.Println("User initiated shutdown. Attempting graceful shutdown...")
	// This gives active connections/requests a chance to complete.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 10-second timeout
	defer cancel()                                                                  // Important: Release resources associated with the context

	// Attempt to shut down the server gracefully
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown (timeout or error during shutdown): %v\n", err)
	}

	newTaskDB.CloseDataBase()
	newUserDB.CloseDataBase()
	
	fmt.Println()

	fmt.Println("********************************")
	fmt.Println("*      DataBase Shutdown       *")
	fmt.Println("********************************")

	fmt.Println()

}
