package goster

import (
	"log"
)

// Supply msg with a string value to be logged to the console and a logger
func LogInfo(msg string, logger *log.Logger) {
	logger.Printf("INFO - %s\n", msg)
}

// Supply msg with a string value to be logged to the console and a logger
func LogWarning(msg string, logger *log.Logger) {
	logger.Printf("WARN - %s\n", msg)
}

// Supply msg with a string value to be logged to the console and a logger
func LogError(msg string, logger *log.Logger) {
	logger.Printf("ERROR - %s\n", msg)
}
