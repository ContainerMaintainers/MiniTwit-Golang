package entities

type Message struct {
	Message_id uint   `gorm:"autoIncrement; primaryKey"`
	Author_id  uint   `gorm:"not null"`
	Text       string `gorm:"not null"`
	Pub_Date   uint   `gorm:"not null"`
	Flagged    bool   `gorm:"not null"`
}
