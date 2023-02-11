package database

type Message struct {
	message_id uint   `gorm:"autoIncrement; primaryKey"`
	author_id  uint   `gorm:"not null"`
	text       string `gorm:"not null"`
	pub_date   uint   `gorm:"not null"`
	flagged    bool   `gorm:"not null"`
}
