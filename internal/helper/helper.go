package helper

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func ContainString(values []string, search string) bool {
	for _, value := range values {
		if value == search {
			return true
		}
	}
	return false
}

func PrintLog(message string) {
	now := time.Now()
	formattedNow := now.Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] %s\n", formattedNow, message)
}

func SprintLog(message string) string {
	now := time.Now()
	formattedNow := now.Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s] %s\n", formattedNow, message)
}

func LoadEnvVariable(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv(key)
}
