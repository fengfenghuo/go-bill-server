package db

import (
	"sync"

	"github.com/astaxie/beego"
	"github.com/bill-server/go-bill-server/node/db/bee"
	"github.com/bill-server/go-bill-server/node/db/gorm"
)

type DBInterface interface {
	RegisterTable(modules ...interface{}) error
	Insert(data interface{}) (int64, error)
	QueryAccount(account string, data interface{}) error
	QueryTxesByAccount(account string, count, offset int, data interface{}) error
	QueryTxByID(txID int64, data interface{}) error
	Update(data interface{}, newData interface{}) error
	DeleteTxByID(data interface{}) error
}

var (
	db   DBInterface
	once sync.Once
)

func GetDBInstance() (DBInterface, error) {
	var err error
	DBName := beego.AppConfig.DefaultString("db::module", "bee")
	once.Do(func() {
		switch DBName {
		case "bee":
			db, err = beeorm.RegisterDB()
		case "gorm":
			db, err = gorm.RegisterDB()
		}
	})
	return db, err
}
