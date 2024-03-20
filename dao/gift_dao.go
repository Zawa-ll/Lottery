package dao

import (
	"github.com/Zawa-ll/raffle/comm"
	"github.com/Zawa-ll/raffle/models"
	"github.com/go-xorm/xorm"
)

type GiftDao struct {
	engine *xorm.Engine
}

func NewGiftDao(engine *xorm.Engine) *GiftDao {
	return &GiftDao{
		engine: engine,
	}
}

func (d *GiftDao) Get(id int) *models.LtGift {
	// Instance of models.LtGift with the provided id in pointer
	data := &models.LtGift{Id: id}
	ok, err := d.engine.Get(data) // Noted: pass pointer to .Get() for in-place modification
	if ok && err == nil {
		return data
	} else {
		data.Id = 0
		return data
	}
}

func (d *GiftDao) GetAll() []models.LtGift {
	datalist := make([]models.LtGift, 0) // initialize as a slice of models.LtGift with a length of 0
	// make()'s return type is the specified type (first argument)
	err := d.engine.
		Asc("sys_status").
		Asc("displayorder").
		Find(&datalist) // Fills datalist with the results, in-place modification
	if err != nil {
		return datalist
	} else {
		return datalist
	}
}

func (d *GiftDao) CountAll() int64 {
	num, err := d.engine.
		Count(&models.LtGift{})
	if err != nil {
		return 0
	} else {
		return num
	}
}

//func (d *GiftDao) Search(country string) []models.LtGift {
//	datalist := make([]models.LtGift, 0)
//	err := d.engine.
//		Where("country=?", country).
//		Desc("id").
//		Find(&datalist)
//	if err != nil {
//		return datalist
//	} else {
//		return datalist
//	}
//}

func (d *GiftDao) Delete(id int) error {
	data := &models.LtGift{Id: id, SysStatus: 1} // SysStatus to 1 as deleted: Soft Deletion
	_, err := d.engine.ID(data.Id).Update(data)
	return err
}

func (d *GiftDao) Update(data *models.LtGift, columns []string) error {
	_, err := d.engine.ID(data.Id).MustCols(columns...).Update(data) // MustCols: columns to be forcefully included in update
	return err
}

func (d *GiftDao) Create(data *models.LtGift) error {
	_, err := d.engine.Insert(data)
	return err
}

// Get a list of currently available prizes.
// Prize-qualified, status-normal, time-duration
func (d *GiftDao) GetAllUse() []models.LtGift {
	now := comm.NowUnix()
	datalist := make([]models.LtGift, 0)
	err := d.engine.
		Cols("id", "title", "prize_num", "left_num", "prize_code",
			"prize_time", "img", "displayorder", "gtype", "gdata").
		Desc("gtype").
		Asc("displayorder").
		Where("prize_num>=?", 0).    // Qualified prizes
		Where("sys_status=?", 0).    // Valid prizes
		Where("time_begin<=?", now). // Time period
		Where("time_end>=?", now).   // Within the time period
		Find(&datalist)              // Store Result into datalist
	if err != nil {
		return datalist
	} else {
		return datalist
	}
}

func (d *GiftDao) IncrLeftNum(id, num int) (int64, error) {
	r, err := d.engine.ID(id).
		Incr("left_num", num).
		//Where("left_num=?", num).
		Update(&models.LtGift{Id: id})
	return r, err
}

func (d *GiftDao) DecrLeftNum(id, num int) (int64, error) {
	r, err := d.engine.ID(id).
		Decr("left_num", num).
		Where("left_num>=?", num).
		Update(&models.LtGift{Id: id})
	return r, err
}
