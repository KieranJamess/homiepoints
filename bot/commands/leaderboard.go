package commands

import (
	"fmt"
	"github.com/KieranJamess/homiepoints/bot/database"
	"github.com/KieranJamess/homiepoints/bot/internal"
	"github.com/KieranJamess/homiepoints/common"
	"github.com/bwmarrin/discordgo"
)

func handleLeaderboard(s *discordgo.Session, i *discordgo.InteractionCreate) {
	leaderboardData, err := database.Leaderboard(i.GuildID, database.DB)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "⚠️ Can't get leaderboard!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	if len(leaderboardData) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "⚠️ No leaderboard data found for this server! Give some points...",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	msg := "**🏆 Leaderboard 🏆**\n\n"

	for idx, entry := range leaderboardData {
		var rankEmoji string
		switch idx {
		case 0:
			rankEmoji = ":first_place:"
		case 1:
			rankEmoji = ":second_place:"
		case 2:
			rankEmoji = ":third_place:"
		default:
			rankEmoji = ":bust_in_silhouette:"
		}

		user, err := internal.GetUserContext(s, entry.UserID)
		if err != nil {
			common.Log.Errorf("Error getting user from context: %v", err)
		}

		msg += fmt.Sprintf("%s **%s** — %d points\n", rankEmoji, common.CapitalizeFirst(internal.GetDisplayName(s, i.GuildID, user)), entry.Points)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
}
