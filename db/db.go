package db

import (
	"fmt"
	"sync"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type DBInterface struct {
}

var (
	db   *DBInterface
	once sync.Once
)

func GetDBInstance() (*DBInterface, error) {
	var err error
	once.Do(func() {
		db, err = RegisterDB()
	})
	return db, err
}

func RegisterDB() (*DBInterface, error) {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	// (可选)设置最大空闲连接
	maxIdle := beego.AppConfig.DefaultInt("db::maxIdle", 60)
	// (可选) 设置最大数据库连接 (go >= 1.2)
	maxConn := beego.AppConfig.DefaultInt("db::maxConn", 320)
	db := beego.AppConfig.DefaultString("db::db", "")
	if db == "" {
		return nil, fmt.Errorf("no db link")
	}

	if err := orm.RegisterDataBase("default", "mysql", db, maxIdle, maxConn); err != nil {
		return nil, fmt.Errorf("RegisterDateBase error: " + err.Error())
	}

	return nil, nil
}

func (db *DBInterface) RegisterTable(modules ...interface{}) error {
	orm.RegisterModel(modules...)

	err := orm.RunSyncdb("default", false, true)
	if err != nil {
		return fmt.Errorf("RunSyncdb error: " + err.Error())
	}
	return nil
}

func (db *DBInterface) Insert(data interface{}) (int64, error) {
	o := orm.NewOrm()
	o.Using("default")

	return o.Insert(data)
}
