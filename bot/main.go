package main

import (
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

const DEV_GUILD_ID = "1015678358190817442"

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

	dg.AddHandler(commands.InteractionHandler)

	if err = dg.Open(); err != nil {
		common.Log.Errorf("Error opening connection: %v", err)
		os.Exit(1)
	}

	/*if os.Getenv("CLEAR_COMMANDS") == "true" {
		// Clear all commands (use guildID for instant updates)
		cmds, _ := dg.ApplicationCommands(dg.State.User.ID, DEV_GUILD_ID)
		for _, c := range cmds {
			common.Log.Infof("Clearning command: %s", c.Name)
			dg.ApplicationCommandDelete(dg.State.User.ID, DEV_GUILD_ID, c.ID)
		}
	}*/

	common.Log.Info("Homie bot is online!")

	// Register slash commands
	if err := commands.Register(dg); err != nil {
		common.Log.Errorf("Failed to register commands: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	common.Log.Info("Shutting down bot...")
	dg.Close()
}
