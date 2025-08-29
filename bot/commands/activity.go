package commands

import (
	"fmt"
	"github.com/KieranJamess/homiepoints/bot/database"
	"github.com/KieranJamess/homiepoints/bot/internal"
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

	// Due to activity database load times and also having to do a lookup on each user (giving and receiving) to get their nicknames
	// We have to defer the message so the all the data can be loaded!
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		common.Log.Errorf("Failed to defer interaction: %v", err)
		return
	}

	activities, err := database.GetRecentActivities(database.DB, i.GuildID, userID)
	if err != nil {
		content := "⚠️ Couldn't fetch activities!"
		_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &content,
		})
		return
	}

	if len(activities) == 0 {
		content := "No recent activity found."
		_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &content,
		})
		return
	}

	msg := "**📝 Recent Activity:**\n\n"
	for _, a := range activities {
		givingUser, _ := internal.GetUserContext(s, a.GivingUserID)
		receivingUser, _ := internal.GetUserContext(s, a.ReceivingUserID)

		reason := ""
		if a.Reason.Valid && a.Reason.String != "" {
			reason = fmt.Sprintf(" — *%s*", a.Reason.String)
		}

		msg += fmt.Sprintf(
			"• **%s** gave **%d** homie %s to **%s**%s\n",
			common.CapitalizeFirst(internal.GetDisplayName(s, i.GuildID, givingUser)),
			a.Points,
			map[bool]string{true: "point", false: "points"}[a.Points == 1],
			common.CapitalizeFirst(internal.GetDisplayName(s, i.GuildID, receivingUser)),
			reason,
		)
	}

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
	if err != nil {
		common.Log.Errorf("Failed to edit interaction response: %v", err)
	}
}
