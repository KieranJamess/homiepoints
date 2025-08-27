package commands

import (
	"github.com/bwmarrin/discordgo"
)

const DEV_GUILD_ID = "1015678358190817442"

var commandHandlers = map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
	"help":        handleHelp,
	"give":        handleGive,
	"get":         handleGet,
	"leaderboard": handleLeaderboard,
}

// Register registers all application commands
func Register(s *discordgo.Session) error {
	var err error
	// /help command
	_, err = s.ApplicationCommandCreate(s.State.User.ID, DEV_GUILD_ID, &discordgo.ApplicationCommand{
		Name:        "help",
		Description: "Shows the help message",
	})
	if err != nil {
		return err
	}

	// /give command
	_, err = s.ApplicationCommandCreate(s.State.User.ID, DEV_GUILD_ID, &discordgo.ApplicationCommand{
		Name:        "give",
		Description: "Give homie points to another user",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user to give points to",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "amount",
				Description: "Number of points to give",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "reason",
				Description: "Reason for the points but not required if giving 1 point",
				Required:    false,
			},
		},
	})
	if err != nil {
		return err
	}

	// /get command
	_, err = s.ApplicationCommandCreate(s.State.User.ID, DEV_GUILD_ID, &discordgo.ApplicationCommand{
		Name:        "get",
		Description: "Get points for a user",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user to get points for",
				Required:    true,
			},
		},
	})
	if err != nil {
		return err
	}

	// /leaderboard command
	_, err = s.ApplicationCommandCreate(s.State.User.ID, DEV_GUILD_ID, &discordgo.ApplicationCommand{
		Name:        "leaderboard",
		Description: "Shows a leaderboard",
	})
	if err != nil {
		return err
	}

	return nil
}

func InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
		handler(s, i)
	}
}
