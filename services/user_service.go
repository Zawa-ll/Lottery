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

// User information can be cached (locally or in redis), and when there is an update, the cache can be cleared directly or updated on a case-by-case basis.
var cachedUserList = make(map[int]*models.LtUser) // make() to initialize and allocate memory
var cachedUserLock = sync.Mutex{}

type UserService interface {
	GetAll(page, size int) []models.LtUser
	CountAll() int
	//Search(country string) []models.LtUser
	Get(id int) *models.LtUser
	//Delete(id int) error
	Update(user *models.LtUser, columns []string) error
	Create(user *models.LtUser) error
}

type userService struct {
	dao *dao.UserDao
}

func NewUserService() UserService {
	return &userService{
		dao: dao.NewUserDao(datasource.InstanceDbMaster()),
	}
}

func (s *userService) GetAll(page, size int) []models.LtUser {
	return s.dao.GetAll(page, size)
}

func (s *userService) CountAll() int {
	return s.dao.CountAll()
}

//func (s *userService) Search(country string) []models.LtUser {
//	return s.dao.Search(country)
//}

func (s *userService) Get(id int) *models.LtUser {
	data := s.getByCache(id)
	if data == nil || data.Id <= 0 {
		data = s.dao.Get(id)
		if data == nil || data.Id <= 0 {
			data = &models.LtUser{Id: id}
		}
		s.setByCache(data)
	}
	return data
}

//func (s *userService) Delete(id int) error {
//	return s.dao.Delete(id)
//}

// UpdateOne
func (s *userService) Update(data *models.LtUser, columns []string) error {
	s.updateByCache(data, columns)
	return s.dao.Update(data, columns)
}

func (s *userService) Create(data *models.LtUser) error {
	err := s.dao.Create(data)
	if err == nil {
		s.updateByCache(data, nil)
	}
	return err
}

func (s *userService) getByCache(id int) *models.LtUser {
	key := fmt.Sprintf("info_user_%d", id)
	rds := datasource.InstanceCache()
	dataMap, err := redis.StringMap(rds.Do("HGETALL", key))
	if err != nil {
		log.Println("user_service.getByCache HGETALL key=", key, ", error=", err)
		return nil
	}
	dataId := comm.GetInt64FromStringMap(dataMap, "Id", 0)
	if dataId <= 0 {
		return nil
	}
	data := &models.LtUser{
		Id:         int(dataId),
		Username:   comm.GetStringFromStringMap(dataMap, "Username", ""),
		Blacktime:  int(comm.GetInt64FromStringMap(dataMap, "Blacktime", 0)),
		Realname:   comm.GetStringFromStringMap(dataMap, "Realname", ""),
		Mobile:     comm.GetStringFromStringMap(dataMap, "Mobile", ""),
		Address:    comm.GetStringFromStringMap(dataMap, "Address", ""),
		SysCreated: int(comm.GetInt64FromStringMap(dataMap, "SysCreated", 0)),
		SysUpdated: int(comm.GetInt64FromStringMap(dataMap, "SysUpdated", 0)),
		SysIp:      comm.GetStringFromStringMap(dataMap, "SysIp", ""),
	}
	return data
}

func (s *userService) setByCache(data *models.LtUser) {
	if data == nil || data.Id <= 0 {
		return
	}
	id := data.Id
	key := fmt.Sprintf("info_user_%d", id)
	rds := datasource.InstanceCache()
	params := []interface{}{key}
	params = append(params, "Id", id)
	if data.Username != "" {
		params = append(params, "Username", data.Username)
		params = append(params, "Blacktime", data.Blacktime)
		params = append(params, "Realname", data.Realname)
		params = append(params, "Mobile", data.Mobile)
		params = append(params, "Address", data.Address)
		params = append(params, "SysCreated", data.SysCreated)
		params = append(params, "SysUpdated", data.SysUpdated)
		params = append(params, "SysIp", data.SysIp)
	}
	_, err := rds.Do("HMSET", params...)
	if err != nil {
		log.Println("user_service.setByCache HMSET params=", params, ", error=", err)
	}
}

func (s *userService) updateByCache(data *models.LtUser, columns []string) {
	if data == nil || data.Id <= 0 {
		return
	}
	key := fmt.Sprintf("info_user_%d", data.Id)
	rds := datasource.InstanceCache()
	rds.Do("DEL", key)
}
