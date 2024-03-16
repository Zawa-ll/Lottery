package models

// The user model that interacts with the browser on the site
type ObjLoginuser struct {
	Uid      int
	Username string
	Now      int
	Ip       string
	Sign     string
}
