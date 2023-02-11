package entities

type User struct {
	user_id  uint   `gorm:"autoIncrement; primaryKey"`
	username string `gorm:"not null"`
	email    string `gorm:"not null"`
	pw_hash  string `gorm:"not null"`
}
