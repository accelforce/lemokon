package discord

import (
	"github.com/accelforce/lemokon/db"
)

type (
	DiscordRecording struct {
		RecordingID int `gorm:"primaryKey"`
		ChannelID   string
		MessageID   string
	}
)

func SaveRecordingMessageID(recordingID int, channelID string, messageID string) error {
	r := &DiscordRecording{
		RecordingID: recordingID,
		ChannelID:   channelID,
		MessageID:   messageID,
	}
	return db.DB.Save(r).Error
}

func GetRecordingMessageID(recordingID int) (string, string, error) {
	var r DiscordRecording
	if err := db.DB.First(&r, recordingID).Error; err != nil {
		return "", "", err
	}
	return r.ChannelID, r.MessageID, nil
}
