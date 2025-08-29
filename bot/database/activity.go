package database

import (
	"database/sql"
	"fmt"
)

type Activity struct {
	GivingUsername    string
	GivingUserID      string
	ReceivingUsername string
	ReceivingUserID   string
	Reason            sql.NullString
	Points            int
	OccurredAt        string
}

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

func GetRecentActivities(db *sql.DB, guildID string, userID *string) ([]Activity, error) {
	var rows *sql.Rows
	var err error

	if userID != nil {
		rows, err = db.Query(`
            SELECT giving_username, giving_user_id, receiving_username, receiving_user_id, reason, points
            FROM activity_points
            WHERE guild_id = ?
              AND (giving_user_id = ?)
            ORDER BY id DESC
            LIMIT 10
        `, guildID, *userID, *userID)
	} else {
		rows, err = db.Query(`
            SELECT giving_username, giving_user_id, receiving_username, receiving_user_id, reason, points
            FROM activity_points
            WHERE guild_id = ?
            ORDER BY id DESC
            LIMIT 10
        `, guildID)
	}

	if err != nil {
		return nil, fmt.Errorf("DB error getting activity: %w", err)
	}
	defer rows.Close()

	activities := make([]Activity, 0)
	for rows.Next() {
		var a Activity
		if err := rows.Scan(&a.GivingUsername, &a.GivingUserID, &a.ReceivingUsername, &a.ReceivingUserID, &a.Reason, &a.Points); err != nil {
			return nil, err
		}
		activities = append(activities, a)
	}

	return activities, nil
}
