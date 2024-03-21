package models

type LtBlackip struct {
	Id         int    `xorm:"not null pk autoincr INT(10)"`
	Ip         string `xorm:"not null default '' comment('IP address') VARCHAR(50)"`
	Blacktime  int    `xorm:"not null default 0 comment('Blacklist restriction expiration time') INT(10)"`
	SysCreated int    `xorm:"not null default 0 comment('Creation time') INT(10)"`
	SysUpdated int    `xorm:"not null default 0 comment('Modify Time') INT(10)"`
}
