package commands

import (
	"database/sql"
	"fmt"
)

func AddPoints(userID, username string, amount int, db *sql.DB) error {
	_, err := db.Exec(`
        INSERT INTO points (user_id, username, points)
        VALUES (?, ?, ?)
        ON CONFLICT(user_id) DO UPDATE SET
            points = points + excluded.points,
            username = excluded.username
    `, userID, username, amount)
	if err != nil {
		return fmt.Errorf("DB Error: %v", err)
	}
	return nil
}

func GetPoints(userID string, db *sql.DB) (int, error) {
	var points int
	err := db.QueryRow("SELECT points FROM points WHERE user_id = ?", userID).Scan(&points)
	if err != nil {
		return 0, fmt.Errorf("DB error: %w", err)
	}
	return points, nil
}
