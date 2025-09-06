package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/MananLed/evalProjectRailway/constants"
	"github.com/MananLed/evalProjectRailway/internal/utils"
	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	connectionStr := utils.GetDBConnString()

	db, err := sql.Open("postgres", connectionStr)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %v", err)
	}

	fmt.Println("Connected to DB successfully!")
	return db, nil
}

func RunInitialSetup(db *sql.DB) error {
	script, err := os.ReadFile(string(constants.DBScriptFilePath))

	if err != nil {
		return fmt.Errorf("failed to read SQL file: %v", err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		return fmt.Errorf("failed to execute SQL script: %v", err)
	}

	fmt.Println("All tables created initialized.")
	return nil
}
