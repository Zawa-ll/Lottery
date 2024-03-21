package models

type LtCode struct {
	Id         int    `xorm:"not null pk autoincr INT(10)"`
	GiftId     int    `xorm:"not null default 0 comment('Prize ID, associated with the lt_gift table') INT(10)"`
	Code       string `xorm:"not null default '' comment('Virtual Coupon Code') VARCHAR(255)"`
	SysCreated int    `xorm:"not null default 0 comment('Virtual Coupon Code') INT(10)"`
	SysUpdated int    `xorm:"not null default 0 comment('update time') INT(10)"`
	SysStatus  int    `xorm:"not null default 0 comment('Status, 0 Normal, 1 Voided, 2 Issued') SMALLINT(5)"`
}
