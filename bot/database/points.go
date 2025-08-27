package database

import (
	"database/sql"
	"fmt"
)

func AddPoints(givingUserID, givingUsername, receivingUserID, receivingUsername, guildID string, amount int, reason *string, db *sql.DB) error {
	_, err := db.Exec(`
    INSERT INTO points (guild_id, user_id, username, points)
    VALUES (?, ?, ?, ?)
    ON CONFLICT(guild_id, user_id) DO UPDATE SET
        points = points + excluded.points,
        username = excluded.username
`, guildID, receivingUserID, receivingUsername, amount)
	if err != nil {
		return fmt.Errorf("DB Error on adding points: %v", err)
	}
	err = AddPointActivity(givingUserID, givingUsername, receivingUserID, receivingUsername, reason, amount, guildID, db)
	if err != nil {
		return fmt.Errorf("DB Error on adding to Activity log: %v", err)
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
