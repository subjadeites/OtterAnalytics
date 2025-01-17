package models

import (
	"time"
)

type Event struct {
	ID                   int64     `gorm:"primaryKey;autoIncrement"`
	ClientID             string    `gorm:"type:varchar(64);not null;index"`
	EventType            string    `gorm:"type:text;not null"`
	Timestamp            time.Time `gorm:"type:timestampz;not null;default:now();primaryKey"`
	UserID               string    `gorm:"type:varchar(64);not null"`
	HomeWorld            int       `gorm:"not null"`
	CheatBannedHashValid bool      `gorm:"not null"`
	Client               string    `gorm:"type:varchar(64);not null;index"`
	OS                   string    `gorm:"type:text;not null"`
	DalamudVersion       string    `gorm:"type:text;not null"`
	IsTesting            bool      `gorm:"default:false"`
	Plugin3rdCount       int       `gorm:"default:0;not null;index"`
	MachineID            string    `gorm:"type:varchar(255);not null;index"`
}

type MachineIDPlugin struct {
	MachineID      string    `gorm:"type:varchar(255);primaryKey;index"`
	Plugin3rdList  string    `gorm:"type:jsonb;not null;default:'[]'"`
	Plugin3rdCount int       `gorm:"default:0;not null;index"`
	LastSeen       time.Time `gorm:"type:timestamptz;not null;default:now()"`
}

var AllModels = []interface{}{
	&Event{},
	&MachineIDPlugin{},
}
