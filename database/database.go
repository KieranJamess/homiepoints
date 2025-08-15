package database

import (
	"database/sql"
	"github.com/KieranJamess/homiepoints/common"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB(filepath string) error {
	common.Log.Infof("Opening database file: %s", filepath)

	var err error
	DB, err = sql.Open("sqlite", filepath) // <-- driver name is "sqlite"
	if err != nil {
		common.Log.Errorf("Failed to open DB: %v", err)
		return err
	}

	if err = DB.Ping(); err != nil {
		common.Log.Errorf("Failed to ping DB: %v", err)
		return err
	}

	common.Log.Info("Successfully connected to the database")

	if err := createTables(); err != nil {
		return err
	}

	return nil
}

func createTables() error {
	common.Log.Info("Creating tables if they don't exist...")

	createUserTable := `
	CREATE TABLE IF NOT EXISTS points (
		user_id TEXT PRIMARY KEY,
		username TEXT,
		points INTEGER DEFAULT 0
	);`

	_, err := DB.Exec(createUserTable)
	if err != nil {
		common.Log.Errorf("Failed to create users table: %v", err)
		return err
	}

	common.Log.Info("Users table created or already exists")
	return nil
}
