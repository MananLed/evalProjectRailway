package utils

import (
	"fmt"

	"github.com/MananLed/evalProjectRailway/internal/config"
)

func GetDBConnString() string {
	return fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
		config.User, config.Password, config.Dbname, config.Host, config.Port, "disable",
	)
}
