package services

import (
	"fmt"
	"log"
	"sync"

	"github.com/Zawa-ll/raffle/comm"
	"github.com/Zawa-ll/raffle/dao"
	"github.com/Zawa-ll/raffle/datasource"
	"github.com/Zawa-ll/raffle/models"
	"github.com/gomodule/redigo/redis"
)

// IP information, which can be cached (locally or in Redis).
// When there are updates, the cache should be refreshed accordingly based on the specific situation.
var cachedBlackipLock = sync.Mutex{}

type BlackipService interface {
	GetAll(page, size int) []models.LtBlackip
	CountAll() int64
	Search(ip string) []models.LtBlackip
	Get(id int) *models.LtBlackip
	//Delete(id int) error
	Update(user *models.LtBlackip, columns []string) error
	Create(user *models.LtBlackip) error
	GetByIp(ip string) *models.LtBlackip
}

type blackipService struct {
	dao *dao.BlackipDao
}

func NewBlackipService() BlackipService {
	return &blackipService{
		dao: dao.NewBlackipDao(datasource.InstanceDbMaster()),
	}
}

func (s *blackipService) GetAll(page, size int) []models.LtBlackip {
	return s.dao.GetAll(page, size)
}

func (s *blackipService) CountAll() int64 {
	return s.dao.CountAll()
}

func (s *blackipService) Search(ip string) []models.LtBlackip {
	return s.dao.Search(ip)
}

func (s *blackipService) Get(id int) *models.LtBlackip {
	return s.dao.Get(id)
}

//func (s *blackipService) Delete(id int) error {
//	return s.dao.Delete(id)
//}

func (s *blackipService) Update(data *models.LtBlackip, columns []string) error {
	s.updateByCache(data, columns)
	return s.dao.Update(data, columns)
}

func (s *blackipService) Create(data *models.LtBlackip) error {
	return s.dao.Create(data)
}

func (s *blackipService) GetByIp(ip string) *models.LtBlackip {

	data := s.getByCache(ip)

	if data == nil || data.Ip == "" {
		data = s.dao.GetByIp(ip)
		if data == nil || data.Ip == "" {
			data = &models.LtBlackip{Ip: ip}
		}
		s.setByCache(data)
	}
	return data
}

func (s *blackipService) getByCache(ip string) *models.LtBlackip {
	key := fmt.Sprintf("info_blackip_%s", ip) // IP
	rds := datasource.InstanceCache()         // Redis Instance
	dataMap, err := redis.StringMap(rds.Do("HGETALL", key))
	if err != nil {
		log.Println("blackip_service.getByCache HGETALL key=", key, ", error=", err)
		return nil
	}
	dataIp := comm.GetStringFromStringMap(dataMap, "Ip", "")
	if dataIp == "" {
		return nil
	}
	data := &models.LtBlackip{
		Id:         int(comm.GetInt64FromStringMap(dataMap, "Id", 0)),
		Ip:         dataIp,
		Blacktime:  int(comm.GetInt64FromStringMap(dataMap, "Blacktime", 0)),
		SysCreated: int(comm.GetInt64FromStringMap(dataMap, "SysCreated", 0)),
		SysUpdated: int(comm.GetInt64FromStringMap(dataMap, "SysUpdated", 0)),
	}
	return data
}

func (s *blackipService) setByCache(data *models.LtBlackip) {
	if data == nil || data.Ip == "" {
		return
	}
	key := fmt.Sprintf("info_blackip_%s", data.Ip)
	rds := datasource.InstanceCache()

	// Initialize a slice of interface{} with a single elemtent 'key'
	// For dynamically build a list of arguments for a Redis command
	params := []interface{}{key}

	params = append(params, "Ip", data.Ip)
	if data.Id > 0 {
		params = append(params, "Blacktime", data.Blacktime)
		params = append(params, "SysCreated", data.SysCreated)
		params = append(params, "SysUpdated", data.SysUpdated)
	}
	_, err := rds.Do("HMSET", params...) // HMSET sets multiple hash fields to multiple values
	if err != nil {
		log.Println("blackip_service.setByCache HMSET params=", params, ", error=", err)
	}
}

func (s *blackipService) updateByCache(data *models.LtBlackip, columns []string) {
	if data == nil || data.Ip == "" {
		return
	}

	key := fmt.Sprintf("info_blackip_%s", data.Ip)
	rds := datasource.InstanceCache()

	rds.Do("DEL", key)
}
