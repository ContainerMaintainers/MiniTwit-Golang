package database

import (
	"fmt"
	"github.com/ContainerMaintainers/MiniTwit-Golang/infrastructure/entities"
	"github.com/ContainerMaintainers/MiniTwit-Golang/initializers"
	"os"

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
		fmt.Println("error connecting to database ", err)
	}
}

func MigrateEntities() {
	DB.AutoMigrate(&entities.User{}, &entities.Message{}, &entities.Follower{})
}

func SeedDatabase() {
	for _, seed := range initializers.AllSeeds() {
		seed.Run(DB)
	}
}
