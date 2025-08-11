package infrastructure

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	CONNECTION_STRING    string
	USER_DB              string
	USER_COLLECTION_NAME string

	TASK_DB              string
	TASK_COLLECTION_NAME string

	JWTSECRET string
	HEADER    string
	CURR_USER string
)

func InitEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	CONNECTION_STRING = getEnv("CONNECTION_STRING")
	USER_DB = getEnv("USER_DB")
	USER_COLLECTION_NAME = getEnv("USER_COLLECTION_NAME")

	TASK_DB = getEnv("TASK_DB")
	TASK_COLLECTION_NAME = getEnv("TASK_COLLECTION_NAME")

	JWTSECRET = getEnv("JWTSECRET")
	HEADER = getEnv("HEADER")
	CURR_USER = getEnv("CURR_USER")

	log.Println("Loaded Environment Variables:")
	log.Println("CONNECTION_STRING:", CONNECTION_STRING)
	log.Println("USER_DB:", USER_DB)
	log.Println("USER_COLLECTION_NAME:", USER_COLLECTION_NAME)
	log.Println("TASK_DB:", TASK_DB)
	log.Println("TASK_COLLECTION_NAME:", TASK_COLLECTION_NAME)
	log.Println("JWTSECRET:", JWTSECRET)
	log.Println("HEADER:", HEADER)
	log.Println("CURR_USER:", CURR_USER)
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return val
}
