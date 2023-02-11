package database

import (
	"fmt"
	"os"
	
	"github.com/ContainerMaintainers/MiniTwit-Golang/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() {
	var err error
	dsnString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	DB, err = gorm.Open(postgres.Open(dsnString), &gorm.Config{})

	if err != nil {
		fmt.Println("error connecting to database %s", err)
	}
}

func MigrateEntities() {
	DB.AutoMigrate(&entities.User{}, &entities.Message{}, &entities.Follower{})
}