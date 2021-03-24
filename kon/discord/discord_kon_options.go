package discord

import (
	"github.com/accelforce/lemokon/db"
)

type (
	DiscordKonOptions struct {
		Uniq      bool `gorm:"primaryKey"`
		ChannelID string
	}
)

var (
	Options = &DiscordKonOptions{}
)

func (DiscordKonOptions) TableName() string {
	return "discord_kon_options"
}

func (o *DiscordKonOptions) Load() error {
	return db.DB.First(o).Error
}

func (o *DiscordKonOptions) SetChannelID(channelID string) error {
	o.ChannelID = channelID
	return db.DB.Save(o).Error
}
