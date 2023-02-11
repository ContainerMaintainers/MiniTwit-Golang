package initializers

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDatabase() {
	dsn := "host=localhost user=admin password=admin dbname=postgres port=9920 sslmode=disable TimeZone=Europe/Copenhagen"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
}