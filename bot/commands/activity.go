package commands

import (
	"fmt"
	"github.com/KieranJamess/homiepoints/bot/database"
	"github.com/KieranJamess/homiepoints/common"
	"github.com/bwmarrin/discordgo"
)

func handleActivity(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var userID *string

	for _, opt := range i.ApplicationCommandData().Options {
		if opt.Name == "user" {
			if u, ok := i.ApplicationCommandData().Resolved.Users[opt.Value.(string)]; ok {
				userID = &u.ID
			}
		}
	}

	activities, err := database.GetRecentActivities(database.DB, i.GuildID, userID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "⚠️ Couldn't fetch activities!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	if len(activities) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No recent activity found.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	msg := "**📝 Recent Activity:**\n\n"
	for _, a := range activities {
		reason := ""
		if a.Reason.Valid && a.Reason.String != "" {
			reason = fmt.Sprintf(" — *%s*", a.Reason.String)
		}

		msg += fmt.Sprintf(
			"• **%s** gave **%d** homie %s to **%s**%s\n",
			common.CapitalizeFirst(a.GivingUsername),
			a.Points,
			map[bool]string{true: "point", false: "points"}[a.Points == 1],
			common.CapitalizeFirst(a.ReceivingUsername),
			reason,
		)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
}
