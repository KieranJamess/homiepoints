package database

import (
	"database/sql"
	"fmt"
)

func AddPointActivity(givingUserID, givingUsername, receivingUserID, receivingUsername string, reason *string, amount int, guildID string, db *sql.DB) error {
	_, err := db.Exec(`
        INSERT INTO activity_points (guild_id, giving_user_id, giving_username, receiving_user_id, receiving_username, reason, points)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `, guildID, givingUserID, givingUsername, receivingUserID, receivingUsername, reason, amount)
	if err != nil {
		return fmt.Errorf("DB Error: %v", err)
	}
	return nil
}
