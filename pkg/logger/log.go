package logger

import (
	"log"
	"os"

	"github.com/MananLed/evalProjectRailway/constants"
)

func LogToFile(logMessage string) {
	file, err := os.OpenFile(string(constants.LogFilePath), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

	log.SetOutput(file)

	log.Println(logMessage)
}
