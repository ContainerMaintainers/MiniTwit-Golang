package database

import (
	"fmt"
	"os"
	"log"

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
		fmt.Println("error connecting to database %s", err)
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
	//closes the database connection after the seeding is done
	defer DB.Close()

	for _, seed := range AllSeeds() {
		err := seed.Run(DB)
		if err != nil {
			log.Fatalf("Error when seeding database with the following seed: %s\nError message: %s", seed.SeedName, err)
		}
	}
}
