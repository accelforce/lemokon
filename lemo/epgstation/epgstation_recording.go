package epgstation

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/accelforce/lemokon/db"
)

type (
	EPGStationRecording struct {
		EPGStationID int `gorm:"column:epgstation_id;primaryKey"`
		RecordingID  int
		Recording    db.Recording
	}
)

func (EPGStationRecording) TableName() string {
	return "epgstation_recordings"
}

func RecordingExists(epgstationRecordingID int) (bool, error) {
	if err := db.DB.First(&EPGStationRecording{}, epgstationRecordingID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func NewRecording(record Record, channel *db.Channel) (*db.Recording, error) {
	recording := db.Recording{
		ProgramName: record.Name,
		StartedAt:   time.Time(record.StartedAt),
		EndsAt:      time.Time(record.EndsAt),
		ChannelID:   &channel.ID,
		Channel:     *channel,
	}
	if err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&recording).Error; err != nil {
			return err
		}

		epgstationRecording := EPGStationRecording{
			EPGStationID: record.ID,
			RecordingID:  recording.ID,
		}
		return tx.Create(&epgstationRecording).Error
	}); err != nil {
		return nil, err
	}

	return &recording, nil
}

func FinishRunning(without []int) ([]db.Recording, error) {
	var epgstationRecordings []EPGStationRecording
	if err := db.DB.
		Joins("Recording").
		Where("Recording.ended", false).
		Not(without).
		Preload("Recording.Channel").
		Find(&epgstationRecordings).Error; err != nil {
		return nil, err
	}

	if len(epgstationRecordings) < 1 {
		return make([]db.Recording, 0), nil
	}

	ids := make([]int, len(epgstationRecordings))
	recordings := make([]db.Recording, len(epgstationRecordings))
	for i, r := range epgstationRecordings {
		r.Recording.Ended = true
		ids[i] = r.RecordingID
		recordings[i] = r.Recording
	}
	if err := db.DB.
		Model(&db.Recording{}).
		Where("id IN ?", ids).
		Update("ended", true).Error; err != nil {
		return nil, err
	}

	return recordings, nil
}
