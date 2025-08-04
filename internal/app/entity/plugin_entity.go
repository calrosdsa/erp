package entity

import "time"

type Plugin struct {
	ID          uint      `gorm:"primarykey"`
	CreatedAt   time.Time `gorm:"default:now()"`
	Name        string    `gorm:"not null;default:null"`
	Code        string    `gorm:"not null;default:null"`
}
		