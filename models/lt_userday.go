package models

type LtUserday struct {
	Id         int `xorm:"not null pk autoincr INT(10)"`
	Uid        int `xorm:"not null default 0 comment('User ID') INT(10)"`
	Day        int `xorm:"not null default 0 comment('Date, e.g.20180725') INT(10)"`
	Num        int `xorm:"not null default 0 comment('Ordinal number') INT(10)"`
	SysCreated int `xorm:"not null default 0 comment('Creation time') INT(10)"`
	SysUpdated int `xorm:"not null default 0 comment('Modify Time') INT(10)"`
}
