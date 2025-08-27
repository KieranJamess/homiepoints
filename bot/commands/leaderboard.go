package commands

import (
	"fmt"
	"github.com/KieranJamess/homiepoints/bot/database"
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
				Content: "⚠️ No leaderboard data found for this server!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	msg := "**🏆 Leaderboard 🏆**\n\n"

	for i, entry := range leaderboardData {
		var rankEmoji string
		switch i {
		case 0:
			rankEmoji = ":first_place:"
		case 1:
			rankEmoji = ":second_place:"
		case 2:
			rankEmoji = ":third_place:"
		default:
			rankEmoji = ":bust_in_silhouette:"
		}

		msg += fmt.Sprintf("%s **%s** — %d points\n", rankEmoji, common.CapitalizeFirst(entry.Username), entry.Points)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
}
