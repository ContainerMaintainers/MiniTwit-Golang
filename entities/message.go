package entities

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	Author_id uint   `gorm:"not null"`
	Text      string `gorm:"not null"`
	Pub_Date  uint   `gorm:"not null"`
	Flagged   bool   `gorm:"not null"`
}
