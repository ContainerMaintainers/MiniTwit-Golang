package initializers

import (
	"github.com/ContainerMaintainers/MiniTwit-Golang/infrastructure/entities"
	"gorm.io/gorm"
)

type Seed struct {
	SeedName string
	Run      func(*gorm.DB)
}

func AllSeeds() []Seed {
	return []Seed{
		{
			SeedName: "CreateUser1",
			Run: func(db *gorm.DB) {
				CreateUser(db, "user1", "user1@gmail.com", "user1iscool")
			},
		},
		{
			SeedName: "CreateUser2",
			Run: func(db *gorm.DB) {
				CreateUser(db, "user2", "user2@gmail.com", "user2iscool")
			},
		},
		{
			SeedName: "CreateUser3",
			Run: func(db *gorm.DB) {
				CreateUser(db, "user3", "user3@gmail.com", "user3iscool")
			},
		},
		{
			SeedName: "CreateMessage1",
			Run: func(db *gorm.DB) {
				CreateMessage(db, 1, "Hello World! From user1", 123456, false)
			},
		},
		{
			SeedName: "CreateMessage2",
			Run: func(db *gorm.DB) {
				CreateMessage(db, 2, "Hello World! From user2", 123456, false)
			},
		},
		{
			SeedName: "CreateMessage3",
			Run: func(db *gorm.DB) {
				CreateMessage(db, 3, "Hello World! From user3", 123456, false)
			},
		},
		{
			SeedName: "CreateFollower1",
			Run: func(db *gorm.DB) {
				CreateFollower(db, 1, 2)
			},
		},
		{
			SeedName: "CreateFollower2",
			Run: func(db *gorm.DB) {
				CreateFollower(db, 1, 3)
			},
		},
		{
			SeedName: "CreateFollower3",
			Run: func(db *gorm.DB) {
				CreateFollower(db, 3, 2)
			},
		},
	}
}

func CreateUser(db *gorm.DB, username string, email string, pwd string) {
	db.Create(&entities.User{
		Username: username,
		Email:    email,
		Password: pwd})
}

func CreateMessage(db *gorm.DB, author uint, text string, date uint, flagged bool) {
	db.Create(&entities.Message{
		Author_id: author,
		Text:      text,
		Pub_Date:  date,
		Flagged:   flagged})
}

func CreateFollower(db *gorm.DB, who uint, whom uint) {
	db.Create(&entities.Follower{
		Who_ID:  who,
		Whom_ID: whom})
}
