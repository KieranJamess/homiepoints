package database

import (
	"database/sql"
	"github.com/KieranJamess/homiepoints/common"
	"os"
)

var DB *sql.DB

// Init opens the database and verifies the connection
func Init(filepath string) error {
	common.Log.Infof("Opening database file: %s", filepath)

	var err error
	DB, err = sql.Open("sqlite", filepath)
	if err != nil {
		common.Log.Errorf("Failed to open DB: %v", err)
		return err
	}

	if err = DB.Ping(); err != nil {
		common.Log.Errorf("Failed to connect to DB: %v", err)
		return err
	}

	common.Log.Info("Database connection established")
	return nil
}

// Close shuts down the database connection
func Close() {
	if DB != nil {
		common.Log.Info("Closing database connection")
		DB.Close()
	}
}

// Exists checks if the DB file exists
func Exists(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil
}
