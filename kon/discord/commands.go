package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	SlashCommands = []*discordgo.ApplicationCommand{
		{
			Name:        "lemokon",
			Description: "Configure LemoKon (Discord Kon)",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "summon",
					Description: "Register this channel as the destination of bot messages",
				},
			},
		},
	}

	CommandsHandler = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Data.Name {
		case "lemokon":
			switch i.Data.Options[0].Name {
			case "summon":
				if err := Options.SetChannelID(i.ChannelID); err != nil {
					respond(s, i, "Failed to set this channel as the destination of bot messages.")
					fmt.Printf("Failed to update options.ChannelID: %s\n", err)
					return
				}
				respond(s, i, "Successfully configured this channel as the destination of bot messages.")
				return
			}
		}

		respond(s, i, "Unknown command received.")
	}
)

func respond(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: content,
			Flags:   1 << 6,
		},
	}); err != nil {
		fmt.Printf("Error sending followup message: %s\n", err)
	}
}
