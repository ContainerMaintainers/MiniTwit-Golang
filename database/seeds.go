package database

import (
	"gorm.io/gorm"

	"github.com/ContainerMaintainers/MiniTwit-Golang/entities"
)

type Seed struct {
	SeedName string
	Run func(*gorm.DB) error
}

func AllSeeds() []Seed {
	return []Seed{
		Seed{
			SeedName: "CreateUser1",
			Run: func(db *gorm.DB) error {
				CreateUser(db, 1, "user1", "user1@gmail.com", "user1iscool")
			}
		},
		Seed{
			SeedName: "CreateUser2",
			Run: func(db *gorm.DB) error {
				CreateUser(db, 2, "user2", "user2@gmail.com", "user2iscool")
			}
		},
		Seed{
			SeedName: "CreateUser3",
			Run: func(db *gorm.DB) error {
				CreateUser(db, 3, "user3", "user3@gmail.com", "user3iscool")
			}
		},
		Seed{
			SeedName: "CreateMessage1",
			Run: func(db *gorm.DB) error {
				CreateMessage(db, 1, 1, "Hello World! From user1", 123456, false)
			}
		},
		Seed{
			SeedName: "CreateMessage2",
			Run: func(db *gorm.DB) error {
				CreateMessage(db, 2, 2, "Hello World! From user2", 123456, false)
			}
		},
		Seed{
			SeedName: "CreateMessage3",
			Run: func(db *gorm.DB) error {
				CreateMessage(db, 3, 3, "Hello World! From user3", 123456, false)
			}
		},
		Seed{
			SeedName: "CreateFollower1",
			Run: func(db *gorm.DB) error {
				CreateFollower(db, 1, 2)
			}
		},
		Seed{
			SeedName: "CreateFollower2",
			Run: func(db *gorm.DB) error {
				CreateFollower(db, 1, 3)
			}
		},
		Seed{
			SeedName: "CreateFollower3",
			Run: func(db *gorm.DB) error {
				CreateFollower(db, 3, 2)
			}
		},
	}
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