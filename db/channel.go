package db

type (
	Channel struct {
		ID                 int `gorm:"primaryKey"`
		Name               string
		RemoteControlKeyID *int
	}
)
