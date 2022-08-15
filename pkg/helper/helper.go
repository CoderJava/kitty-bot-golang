package helper

import (
	"fmt"
	"time"
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
