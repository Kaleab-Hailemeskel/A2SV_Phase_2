package router

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"task_manager/controllers"
	"task_manager/data"
	"time"

	"github.com/gin-gonic/gin"
)

// GET, GET_ID, POST, PUT, DELETE methods are mapped with their counter func
func StartEngine(port_number string) {
	router := gin.Default()

	router.GET("/tasks", controllers.GetTasks)
	router.GET("/tasks/:id", controllers.GetTaskByID)
	router.POST("/tasks/", controllers.PostTask)
	router.PUT("/tasks/:id", controllers.PutTaskByID)
	router.DELETE("/tasks/:id", controllers.DeleteTaskByID)

	data.ConnectToMongo()
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

	data.CloseMongoDB()
	fmt.Println()

	fmt.Println("********************************")
	log.Println("*      DataBase Shutdown       *")
	fmt.Println("********************************")

	fmt.Println()

}
