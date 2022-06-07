package gottp_server

import "log"

func LogInfo(msg string) {
	log.Printf("INFO - %s\n", msg)
}
func LogWarning(msg string) {
	log.Printf("WARN - %s\n", msg)
}
func LogError(msg string) {
	log.Printf("ERROR - %s\n", msg)
}
