package database

import (
	"fmt"
	"os"

	"github.com/ContainerMaintainers/MiniTwit-Golang/entities"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() {
	var err error
	dsnString := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	DB, err = gorm.Open(postgres.Open(dsnString), &gorm.Config{})

	if err != nil {
		fmt.Println("error connecting to database ", err)
	}
}

func ConnectToTestDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})

	if err != nil {
		fmt.Println("error connecting to database %s", err)
	}
}

func MigrateEntities() {
	DB.AutoMigrate(&entities.User{}, &entities.Message{}, &entities.Follower{})
}

func SeedDatabase() {
	//Users
	DB.Create(&entities.User{User_ID: 1, Username: "user1", Email: "user1@gmail.com", PW_Hash: "user1iscool"})
	DB.Create(&entities.User{User_ID: 1, Username: "user2", Email: "user1@gmail.com", PW_Hash: "user2iscool"})
	DB.Create(&entities.User{User_ID: 1, Username: "user3", Email: "user1@gmail.com", PW_Hash: "user3iscool"})

	//Messages
	DB.Create(&entities.Message{Message_ID: 1, Author_ID: 1, Text: "Hello World! From user1", Pub_Date: 123456, Flagged: false})
	DB.Create(&entities.Message{Message_ID: 2, Author_ID: 2, Text: "Hello World! From user2", Pub_Date: 123456, Flagged: false})
	DB.Create(&entities.Message{Message_ID: 3, Author_ID: 3, Text: "Hello World! From user3", Pub_Date: 123456, Flagged: false})

	//Followers
	DB.Create(&entities.Follower{Who_ID: 1, Whom_ID: 2})
	DB.Create(&entities.Follower{Who_ID: 1, Whom_ID: 3})
	DB.Create(&entities.Follower{Who_ID: 3, Whom_ID: 2})
}
