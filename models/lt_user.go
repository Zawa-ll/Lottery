package models

type LtUser struct {
	Id         int    `xorm:"not null pk autoincr INT(10)"`
	Username   string `xorm:"not null default '' comment('User ID') VARCHAR(50)"`
	Blacktime  int    `xorm:"not null default 0 comment('Blacklist restriction expiration time') INT(10)"`
	Realname   string `xorm:"not null default '' comment('Associates') VARCHAR(50)"`
	Mobile     string `xorm:"not null default '' comment('Cell phone number') VARCHAR(50)"`
	Address    string `xorm:"not null default '' comment('Contact address') VARCHAR(255)"`
	SysCreated int    `xorm:"not null default 0 comment('Creation time') INT(10)"`
	SysUpdated int    `xorm:"not null default 0 comment('Modify Time') INT(10)"`
	SysIp      string `xorm:"not null default '' comment('IP address') VARCHAR(50)"`
}
