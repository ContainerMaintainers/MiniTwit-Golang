package database

import (
	"github.com/jinzhu/gorm"

	"github.com/ContainerMaintainers/MiniTwit-Golang/entities"
)

type Seed struct {
	SeedName string
	Run func(*gorm.DB) error
}

func AllSeeds() []Seed {
	return []Seed{
		Seed{
			SeedName: "Users",
			Run: func(db *gorm.DB) error {
				db.Create(&entities.User{User_ID: 1, Username: "user1", Email: ""})
			}

func CreateUser(db *gorm.DB, userID uint, username string, email string, pw string) error {
	return db.Create(&entities.User{
					User_ID: userID, 
					Username: username, 
					Email: email, 
					PW_Hash: pw}).Error
}

func CreateFollower(db *gorm.DB, who uint, whom uint) error {
	return db.Create(&entities.Follower{
					Who_ID: who,
					Whom_ID: whom}).Error
}

func CreateMessage(db *gorm.DB, author uint, text string, date uint, flagged bool) error {
	return db.Create(&entities.User{
					Author_ID: author, 
					Text: text, 
					Pub_Date: date, 
					Flagged: flagged}).Error
}