package main

import (
	"database/sql"
	"fmt"
	"github.com/KieranJamess/homiepoints/bot/commands"
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
		user := i.ApplicationCommandData().Options[0].UserValue(s)
		amount := i.ApplicationCommandData().Options[1].IntValue()

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

		err := commands.AddPoints(user.ID, user.Username, int(amount), database.DB)
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

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("%s gave %d homie %s to %s!",
					common.CapitalizeFirst(i.Member.User.Username),
					amount,
					func() string {
						if amount == 1 {
							return "point"
						}
						return "points"
					}(),
					user.Username),
			},
		})
	}

	if i.ApplicationCommandData().Name == "get" {
		user := i.ApplicationCommandData().Options[0].UserValue(s)

		points, err := commands.GetPoints(user.ID, database.DB)

		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("⚠️ Can't get points for %s!", user.Username),
					Flags:   discordgo.MessageFlagsEphemeral, // only visible to the user
				},
			})
		} else {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Points for %s is %v!", user.Username, points),
					Flags:   discordgo.MessageFlagsEphemeral, // only visible to the user
				},
			})
		}
	}
}
