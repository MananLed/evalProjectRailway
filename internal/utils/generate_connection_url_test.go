package utils

import (
	"fmt"
	"testing"

	"github.com/MananLed/evalProjectRailway/internal/config"
)

func TestGetDBConnString(t *testing.T) {

	got := GetDBConnString()

	// Expected format
	want := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
		config.User, config.Password, config.Dbname, config.Host, config.Port, "disable",
	)

	// Assert
	if got != want {
		t.Errorf("unexpected connection string:\n got  %q\n want %q", got, want)
	}
}
