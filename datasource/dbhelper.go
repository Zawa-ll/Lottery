package datasource

import (
	"fmt"
	"log"
	"sync"

	"github.com/Zawa-ll/raffle/conf"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var dbLock sync.Mutex
var masterInstance *xorm.Engine

func InstanceDbMaster() *xorm.Engine {
	if masterInstance != nil {
		return masterInstance
	}

	// Lock if no masterInstance exist
	dbLock.Lock()
	defer dbLock.Unlock()

	// return the instance created during lock, to avoid creating instance for multiple times
	if masterInstance != nil {
		return masterInstance
	}

	return NewDbMaster()
}

// generate an xorm database engine each time
func NewDbMaster() *xorm.Engine {
	sourcename := fmt.Sprint("%s:%s@tcp(%s:%d)/%s?charset=utf8",
		conf.DbMaster.User,
		conf.DbMaster.Pwd,
		conf.DbMaster.Host,
		conf.DbMaster.Port,
		conf.DbMaster.Database)
	instance, err := xorm.NewEngine(conf.DriveName, sourcename)
	if err != nil {
		log.Fatal("dbhelper.NewDbMaster NewEngine error ", err)
		return nil
	}

	// For testing: display SQL commands + executeTime
	instance.ShowSQL(true)

	masterInstance = instance
	return instance
}
