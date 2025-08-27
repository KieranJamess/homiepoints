package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func handleHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	commands := []struct {
		Name        string
		Description string
	}{
		{"/help", "Show this help message"},
		{"/give [user] [amount] [reason]", "Give homie points to another user. Reason is only required if giving more than 1 point"},
		{"/get [user]", "Get points for a specific user"},
		{"/leaderboard", "Show the server's current homie points leaderboard"},
	}

	msg := "**📖 Available Commands:**\n\n"
	for _, cmd := range commands {
		msg += fmt.Sprintf("**%s**\n%s\n\n", cmd.Name, cmd.Description)
	}

	// Send response
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
