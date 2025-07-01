package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)





func GetEnv(key string) string {

  // load .env file
  err := godotenv.Load(".env")

  if err != nil {
    fmt.Println("Error loading .env file")
  }

  return os.Getenv(key)
}