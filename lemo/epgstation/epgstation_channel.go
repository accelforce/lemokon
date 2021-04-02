package epgstation

import (
	"gorm.io/gorm"

	"github.com/accelforce/lemokon/db"
)

type (
	EPGStationChannel struct {
		EPGStationID int `gorm:"column:epgstation_id;primaryKey"`
		ChannelID    int
		HasLogoData  bool
	}
)

func (EPGStationChannel) TableName() string {
	return "epgstation_channels"
}

func UpdateChannel(channel *Channel) (*db.Channel, error) {
	var c db.Channel
	if err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where(db.Channel{Name: channel.Name}).
			Assign(db.Channel{RemoteControlKeyID: &channel.RemoteControlKeyID}).
			FirstOrCreate(&c).Error; err != nil {
			return err
		}

		return tx.Where(EPGStationChannel{EPGStationID: channel.ID}).
			Assign(EPGStationChannel{ChannelID: c.ID, HasLogoData: channel.HasLogoData}).
			FirstOrCreate(&EPGStationChannel{}).Error
	}); err != nil {
		return nil, err
	}
	return &c, nil
}
