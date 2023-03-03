package entities

type Follower struct {
	Who_ID  uint // Should be a reference to User rather than uint
	Whom_ID uint
}
