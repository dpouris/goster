package gottp_server

import (
	"log"
)

func LogInfo(msg string, logger *log.Logger) {
	logger.Printf("INFO - %s\n", msg)
}
func LogWarning(msg string, logger *log.Logger) {
	logger.Printf("WARN - %s\n", msg)
}
func LogError(msg string, logger *log.Logger) {
	logger.Printf("ERROR - %s\n", msg)
}
