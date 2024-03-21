package models

type LtGift struct {
	Id           int    `xorm:"not null pk autoincr INT(10)" json:"id"`
	Title        string `xorm:"not null default '' comment('Prize name') VARCHAR(255)" json:"title"`
	PrizeNum     int    `xorm:"not null default -1 comment('Prize quantity, 0 unlimited, >0 limited, <0 no prize') INT(11)" json:"-"`
	LeftNum      int    `xorm:"not null default 0 comment('Remaining quantity') INT(11)" json:"-"`
	PrizeCode    string `xorm:"not null default '' comment('0-9999 means 100% chance, 0-0 means one in ten thousand chance of winning') VARCHAR(50)" json:"-"`
	PrizeTime    int    `xorm:"not null default 0 comment('Award cycle, D days') INT(10)" json:"-"`
	Img          string `xorm:"not null default '' comment('Prize image') VARCHAR(255)" json:"img"`
	Displayorder int    `xorm:"not null default 0 comment('Position number, smaller comes first') INT(10)" json:"displayorder"`
	Gtype        int    `xorm:"not null default 0 comment('Prize type, 0 virtual coin, 1 virtual voucher, 2 physical-small prize, 3 physical-big prize') INT(10)" json:"gtype"`
	Gdata        string `xorm:"not null default '' comment('Extended data, e.g., virtual coin amount') VARCHAR(255)" json:"-"`
	TimeBegin    int    `xorm:"not null default 0 comment('Start time') INT(11)" json:"-"`
	TimeEnd      int    `xorm:"not null default 0 comment('End time') INT(11)" json:"-"`
	PrizeData    string `xorm:"comment('Award plan, [[Time1, Quantity1],[Time2, Quantity2]]') MEDIUMTEXT" json:"-"`
	PrizeBegin   int    `xorm:"not null default 0 comment('Start of the award plan cycle') INT(11)" json:"-"`
	PrizeEnd     int    `xorm:"not null default 0 comment('End of the award plan cycle') INT(11)" json:"-"`
	SysStatus    int    `xorm:"not null default 0 comment('Status, 0 normal, 1 deleted') SMALLINT(5)" json:"-"`
	SysCreated   int    `xorm:"not null default 0 comment('Creation time') INT(10)" json:"-"`
	SysUpdated   int    `xorm:"not null default 0 comment('Modification time') INT(10)" json:"-"`
	SysIp        string `xorm:"not null default '' comment('Operator IP') VARCHAR(50)" json:"-"`
}
