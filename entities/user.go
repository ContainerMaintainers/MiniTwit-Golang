package entities

type User struct {
	User_ID  uint   `gorm:"autoIncrement; primaryKey"`
	Username string `gorm:"not null"`
	Email    string `gorm:"not null"`
	PW_Hash  string `gorm:"not null"`
}
