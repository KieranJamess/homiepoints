package main

import (
	"database/sql"
	"fmt"
	"github.com/KieranJamess/homiepoints/bot/database"
	"github.com/KieranJamess/homiepoints/common"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
	"os"
	"os/signal"
	"syscall"
)

var db *sql.DB

func main() {
	common.Log.Info("Loading .env file")
	if err := godotenv.Load(".env"); err != nil {
		common.Log.Fatalf(".env file not found: %v", err)
	}

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		common.Log.Fatal("DISCORD_TOKEN not set")
	}

	var err error

	// Open DB
	if err := database.Init("./homiepoints.db"); err != nil {
		common.Log.Fatalf("Error initializing database: %v", err)
	}
	defer database.Close()

	common.Log.Info("Successfully connected to the database")

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		common.Log.Errorf("Error creating bot: %v", err)
		os.Exit(1)
	}

	dg.AddHandler(interactionHandler)

	if err = dg.Open(); err != nil {
		common.Log.Errorf("Error opening connection: %v", err)
		os.Exit(1)
	}

	common.Log.Info("Homie bot is online!")

	_, err = dg.ApplicationCommandCreate(dg.State.User.ID, "", &discordgo.ApplicationCommand{
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
		common.Log.Errorf("Failed to register /give command: %v", err)
	}

	_, err = dg.ApplicationCommandCreate(dg.State.User.ID, "", &discordgo.ApplicationCommand{
		Name:        "get",
		Description: "Get points for a user",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user to give points to",
				Required:    true,
			},
		},
	})
	if err != nil {
		common.Log.Errorf("Failed to register /get command: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	common.Log.Info("Shutting down bot...")
	dg.Close()
}

func interactionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name == "give" {

		var reason *string
		for _, opt := range i.ApplicationCommandData().Options {
			if opt.Name == "reason" {
				val := opt.StringValue()
				reason = &val
				break
			}
		}

		user := i.ApplicationCommandData().Options[0].UserValue(s)
		amount := i.ApplicationCommandData().Options[1].IntValue()
		guildId := i.GuildID

		if i.Member.User.ID == user.ID {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⚠️ You can't give points to yourself!",
					Flags:   discordgo.MessageFlagsEphemeral, // only visible to the user
				},
			})
			return
		}

		if (reason == nil || *reason == "") && amount > 1 {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⚠️ Adding more than 1 point requires a reason",
					Flags:   discordgo.MessageFlagsEphemeral, // only visible to the user
				},
			})
			return
		}

		err := database.AddPoints(
			i.Member.User.ID,       // Giving User
			i.Member.User.Username, // Giving User
			user.ID,                // Receiving User
			user.Username,          // Receiving User
			guildId,
			int(amount),
			reason,
			database.DB,
		)
		if err != nil {
			// Send ephemeral error response
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⚠️ Something went wrong while giving points. Please try again.",
					Flags:   discordgo.MessageFlagsEphemeral, // only visible to the user
				},
			})
			return
		}

		msg := fmt.Sprintf("%s gave %d homie %s to %s!",
			common.CapitalizeFirst(i.Member.User.DisplayName()),
			amount,
			map[bool]string{true: "point", false: "points"}[amount == 1],
			common.CapitalizeFirst(user.DisplayName()),
		)

		if reason != nil && *reason != "" {
			msg = fmt.Sprintf("%s Reason: %s", msg, common.CapitalizeFirst(*reason))
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
	}

	if i.ApplicationCommandData().Name == "get" {
		user := i.ApplicationCommandData().Options[0].UserValue(s)

		points, err := database.GetPoints(user.ID, database.DB)

		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("⚠️ Can't get points for %s!", common.CapitalizeFirst(user.DisplayName())),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		} else {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Points for %s is %v!", common.CapitalizeFirst(user.DisplayName()), points),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
	}
}
