package commands

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

func handleHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	commands := []string{
		"/help - Show this help message",
		"/give [user] [amount] [reason] - Give homie points to another user. Reason is only required if giving more than 1 point",
		"/get [user] - Get points for a specific user",
	}

	msg := "**Available Commands:**\n" + strings.Join(commands, "\n")

	// Send response
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
