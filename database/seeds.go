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
			SeedName: "CreateUser1",
			Run: func(db *gorm.DB) error {
				db.Create(&entities.User{User_ID: 1, Username: "user1", Email: "user1@gmail.com", PW_Hash "user1iscool"})
			}
		},
		Seed{
			SeedName: "CreateUser2",
			Run: func(db *gorm.DB) error {
				db.Create(&entities.User{User_ID: 2, Username: "user2", Email: "user2@gmail.com", PW_Hash "user2iscool"})
			}
		},
		Seed{
			SeedName: "CreateUser3",
			Run: func(db *gorm.DB) error {
				db.Create(&entities.User{User_ID: 3, Username: "user3", Email: "user3@gmail.com", PW_Hash "user3iscool"})
			}
		},
		Seed{
			SeedName: "CreateMessage1",
			Run: func(db *gorm.DB) error {
				db.Create(&entities.Message{Message_ID: 1, Author_ID: 1, Text: "Hello World! From user1", Pub_Date: 123456, Flagged: false})
			}
		},
		Seed{
			SeedName: "CreateMessage2",
			Run: func(db *gorm.DB) error {
				db.Create(&entities.Message{Message_ID: 2, Author_ID: 2, Text: "Hello World! From user2", Pub_Date: 123456, Flagged: false})
			}
		},
		Seed{
			SeedName: "CreateMessage3",
			Run: func(db *gorm.DB) error {
				db.Create(&entities.Message{Message_ID: 3, Author_ID: 3, Text: "Hello World! From user3", Pub_Date: 123456, Flagged: false})
			}
		},
		Seed{
			SeedName: "CreateFollower1",
			Run: func(db *gorm.DB) error {
				db.Create(&entities.Follower{Who_ID: 1, Whom_ID: 2})
			}
		},
		Seed{
			SeedName: "CreateFollower2",
			Run: func(db *gorm.DB) error {
				db.Create(&entities.Follower{Who_ID: 1, Whom_ID: 3})
			}
		},
		Seed{
			SeedName: "CreateFollower3",
			Run: func(db *gorm.DB) error {
				db.Create(&entities.Follower{Who_ID: 3, Whom_ID: 2})
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