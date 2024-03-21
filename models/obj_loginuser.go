package models

// The user model used on browser
type ObjLoginuser struct {
	Uid      int
	Username string
	Now      int // -- timestamp
	Ip       string
	Sign     string // -- signature
}
