package database

import (
	"database/sql"
	"fmt"
)

type LeaderboardEntry struct {
	Username string
	UserID   string
	Points   int
}

func Leaderboard(guildID string, db *sql.DB) ([]LeaderboardEntry, error) {
	rows, err := db.Query(`
		SELECT username, user_id, points
		FROM points
		WHERE guild_id = ?
		ORDER BY points DESC
		LIMIT 10
	`, guildID)
	if err != nil {
		return nil, fmt.Errorf("DB error getting leaderboard: %w", err)
	}
	defer rows.Close()

	var leaderboard []LeaderboardEntry
	for rows.Next() {
		var entry LeaderboardEntry
		if err := rows.Scan(&entry.Username, &entry.UserID, &entry.Points); err != nil {
			return nil, fmt.Errorf("DB error scanning leaderboard: %w", err)
		}
		leaderboard = append(leaderboard, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("DB error iterating leaderboard rows: %w", err)
	}

	return leaderboard, nil
}
