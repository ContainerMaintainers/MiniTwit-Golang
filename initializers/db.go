package initializers

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)



func ConnectToDatabase() {
	var dsnString = "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Stockholm", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT")
	dsn := dsnString
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("error connecting to database")
	}
}