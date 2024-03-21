package models

type LtResult struct {
	Id         int    `xorm:"not null pk autoincr INT(10)" json:"-"`
	GiftId     int    `xorm:"not null default 0 comment('Prize ID, associated with the lt_gift table') INT(10)" json:"gift_id"`
	GiftName   string `xorm:"not null default '' comment('Prize Name') VARCHAR(255)" json:"gift_name"`
	GiftType   int    `xorm:"not null default 0 comment('Prize type, same as lt_gift.gtype') INT(10)" json:"gift_type"`
	Uid        int    `xorm:"not null default 0 comment('user ID') INT(10)" json:"uid"`
	Username   string `xorm:"not null default '' comment('user Name') VARCHAR(50)" json:"username"`
	PrizeCode  int    `xorm:"not null default 0 comment('Raffle number (4-digit random number)') INT(10)" json:"-"`
	GiftData   string `xorm:"not null default '' comment('Award Information') VARCHAR(255)" json:"-"`
	SysCreated int    `xorm:"not null default 0 comment('Creation time') INT(10)" json:"-"`
	SysIp      string `xorm:"not null default '' comment('User Lucky Draw IP') VARCHAR(50)" json:"-"`
	SysStatus  int    `xorm:"not null default 0 comment('Status, 0 Normal, 1 Delete, 2 Cheat') SMALLINT(5)" json:"-"`
}
