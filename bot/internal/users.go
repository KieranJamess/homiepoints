package internal

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func GetDisplayName(s *discordgo.Session, guildID string, user *discordgo.User) string {
	displayName := user.Username

	member, err := s.State.Member(guildID, user.ID)
	if err != nil || member == nil {
		member, err = s.GuildMember(guildID, user.ID)
		if err != nil || member == nil {
			return displayName
		}
	}

	if member.Nick != "" {
		displayName = member.Nick
	}

	return displayName
}

func GetUserContext(s *discordgo.Session, userID string) (*discordgo.User, error) {
	user, err := s.User(userID)
	if err != nil || user == nil {
		return nil, fmt.Errorf("Can't find user with ID %s", userID)
	}
	return user, nil
}
