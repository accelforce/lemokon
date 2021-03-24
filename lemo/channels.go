package lemo

import (
	"github.com/accelforce/lemokon/db"
)

var (
	newRecording                      = make(chan *db.Recording)
	NewRecording <-chan *db.Recording = newRecording

	endedRecording                      = make(chan *db.Recording)
	EndedRecording <-chan *db.Recording = endedRecording
)
