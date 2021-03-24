// +build discord

package kon

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/accelforce/lemokon/db"
	"github.com/accelforce/lemokon/kon/discord"
	"github.com/accelforce/lemokon/lemo"
	"github.com/accelforce/lemokon/util"
)

const (
	colorOrange = 15105570
	colorGreen  = 3066993
)

var (
	s *discordgo.Session
)

func Kon(ctx context.Context) error {
	if err := discord.Options.Load(); err != nil {
		return fmt.Errorf("error loading Discord Kon options: %s\n", err)
	}

	token, err := util.CheckEnv("DISCORD_TOKEN")
	if err != nil {
		return err
	}

	go observe()

	return streaming(ctx, token)
}

func observe() {
	for {
		select {
		case r := <-lemo.NewRecording:
			m, err := s.ChannelMessageSendEmbed(discord.Options.ChannelID, getRecordingEmbed(r))
			if err != nil {
				fmt.Printf("Failed to notify new recording: %s\n", err)
				break
			}
			if err := discord.SaveRecordingMessageID(r.ID, m.ChannelID, m.ID); err != nil {
				fmt.Printf("Failed to save Discord message id: %s\n", err)
			}
			break
		case r := <-lemo.EndedRecording:
			cID, mID, err := discord.GetRecordingMessageID(r.ID)
			if err != nil {
				fmt.Printf("Failed to find Discord message id: %s\n", err)
				break
			}
			if err := s.ChannelMessageDelete(cID, mID); err != nil {
				fmt.Printf("Failed to delete Discord message: %s\n", err)
				break
			}
			if _, err := s.ChannelMessageSendEmbed(discord.Options.ChannelID, getRecordingEmbed(r)); err != nil {
				fmt.Printf("Failed to notify ended recording: %s\n", err)
			}
			break
		}
	}
}

func getRecordingEmbed(r *db.Recording) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       r.ProgramName,
		Description: "録画を開始します",
		Color:       colorOrange,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "放送日",
				Value:  r.StartedAt.Format("2006/01/02"),
				Inline: true,
			},
			{
				Name:   "放送開始",
				Value:  r.StartedAt.Format("15:04"),
				Inline: true,
			},
			{
				Name:   "放送終了",
				Value:  r.EndsAt.Format("15:04"),
				Inline: true,
			},
		},
	}
	if r.Ended {
		embed.Description = "録画を終了しました"
		embed.Color = colorGreen
	}
	return embed
}

func streaming(ctx context.Context, token string) error {
	retry := make(chan struct{})
	disconnect := make(chan struct{})

	for {
		var err error
		s, err = discordgo.New(fmt.Sprintf("Bot %s", token))
		if err != nil {
			return fmt.Errorf("error creating Discord client: %w", err)
		}

		s.AddHandler(func(s *discordgo.Session, _ *discordgo.Connect) {
			fmt.Printf("Connected to Discord as %s#%s\n", s.State.User.Username, s.State.User.Discriminator)

			for _, command := range discord.SlashCommands {
				if _, err := s.ApplicationCommandCreate(s.State.User.ID, "", command); err != nil {
					fmt.Printf("Failed to register command: %s\n", err)
				}
			}
		})
		s.AddHandler(func(_ *discordgo.Session, _ *discordgo.Disconnect) {
			fmt.Println("Disconnected from Discord")
			disconnect <- struct{}{}
		})
		s.AddHandler(discord.CommandsHandler)

		if err := s.Open(); err != nil {
			fmt.Printf("Cannot connect to Discord: %s\n", err)
			fmt.Println("Retry in 15 seconds...")
			go func() {
				<-time.After(15 * time.Second)
				retry <- struct{}{}
			}()
		}

		select {
		case <-retry:
		case <-disconnect:
			continue
		case <-ctx.Done():
			fmt.Println("Gracefully shutting down Discord connection...")
			return s.Close()
		}
	}
}
