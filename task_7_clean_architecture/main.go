package main

import (
	"log"
	"task_7_clean_architecture/infrastructure"
	"task_7_clean_architecture/router"

	"github.com/joho/godotenv"
)

func main() { //
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	infrastructure.InitEnv()

	port_number := "8081"
	router.StartEngine(port_number)

}
