package entities

type Follower struct {
	Who_ID  uint `gorm:"primaryKey"` // Should be a reference to User rather than uint
	Whom_ID uint `gorm:"primaryKey"`
}
