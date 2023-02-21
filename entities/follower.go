package entities

type Follower struct {
	Who_ID  uint `gorm:"primaryKey"`
	Whom_ID uint `gorm:"primaryKey"`
}
