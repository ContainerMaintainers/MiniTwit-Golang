package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"not null"`
	Email    string `gorm:"not null"`
	PW_Hash  string `gorm:"not null"`
}
