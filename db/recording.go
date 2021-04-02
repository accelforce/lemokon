package db

import (
	"time"
)

type (
	Recording struct {
		ID          int `gorm:"primaryKey"`
		ProgramName string
		StartedAt   time.Time
		EndsAt      time.Time
		Ended       bool
		ChannelID   *int
		Channel     Channel
	}
)
