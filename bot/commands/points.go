package commands

import (
	"fmt"
	"github.com/KieranJamess/homiepoints/bot/database"
	"github.com/KieranJamess/homiepoints/bot/internal"
	"github.com/KieranJamess/homiepoints/common"
	"github.com/bwmarrin/discordgo"
)

func handleGive(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	if (reason == nil || *reason == "") && amount > 1 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "⚠️ Adding more than 1 point requires a reason",
				Flags:   discordgo.MessageFlagsEphemeral,
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
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "⚠️ Something went wrong while giving points. Please try again.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	msg := fmt.Sprintf("%s gave %d homie %s to %s!",
		common.CapitalizeFirst(i.Member.Nick),
		amount,
		map[bool]string{true: "point", false: "points"}[amount == 1],
		common.CapitalizeFirst(internal.GetDisplayName(s, guildId, user)),
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

func handleGet(s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := i.ApplicationCommandData().Options[0].UserValue(s)

	points, err := database.GetPoints(user.ID, database.DB)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("⚠️ Can't get points for %s!", common.CapitalizeFirst(internal.GetDisplayName(s, i.GuildID, user))),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Points for %s is %v!", common.CapitalizeFirst(internal.GetDisplayName(s, i.GuildID, user)), points),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
