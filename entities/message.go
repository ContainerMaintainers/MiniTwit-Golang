package entities

type Message struct {
	Message_ID uint   `gorm:"autoIncrement; primaryKey"`
	Author_ID  uint   `gorm:"not null"`
	Text       string `gorm:"not null"`
	Pub_Date   uint   `gorm:"not null"`
	Flagged    bool   `gorm:"not null"`
}
