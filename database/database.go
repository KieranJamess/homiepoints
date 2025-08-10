package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"local/common"
)

var DB *sql.DB

func InitDB(filepath string) error {
	common.Logger.Infof("Opening database file: %s", filepath)

	var err error
	DB, err = sql.Open("sqlite3", filepath)
	if err != nil {
		common.Logger.Errorf("Failed to open DB: %v", err)
		return err
	}

	if err = DB.Ping(); err != nil {
		common.Logger.Errorf("Failed to ping DB: %v", err)
		return err
	}

	common.Logger.Info("Successfully connected to the database")

	if err := createTables(); err != nil {
		return err
	}

	return nil
}

func createTables() error {
	common.Logger.Info("Creating tables if they don't exist...")

	createUserTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        discord_id TEXT NOT NULL UNIQUE,
        points INTEGER DEFAULT 0
    );`

	_, err := DB.Exec(createUserTable)
	if err != nil {
		common.Logger.Errorf("Failed to create users table: %v", err)
		return err
	}

	common.Logger.Info("Users table created or already exists")
	return nil
}
