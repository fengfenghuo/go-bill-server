package beeorm

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type BeeOrmInterface struct {
}

func RegisterDB() (*BeeOrmInterface, error) {
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

	return &BeeOrmInterface{}, nil
}

func (db *BeeOrmInterface) RegisterTable(modules ...interface{}) error {
	orm.RegisterModel(modules...)

	err := orm.RunSyncdb("default", false, true)
	if err != nil {
		return fmt.Errorf("RunSyncdb error: " + err.Error())
	}
	return nil
}

func (db *BeeOrmInterface) Insert(data interface{}) (int64, error) {
	o := orm.NewOrm()
	o.Using("default")

	return o.Insert(data)
}

func (db *BeeOrmInterface) QueryAccount(account string, data interface{}) error {
	return nil
}

func (db *BeeOrmInterface) QueryTxesByAccount(account string, count, offset int, data interface{}) error {
	return nil
}

func (db *BeeOrmInterface) QueryTxByID(txID int64, data interface{}) error {
	return nil
}

func (db *BeeOrmInterface) Update(data interface{}, newData interface{}) error {
	return nil
}

func (db *BeeOrmInterface) DeleteTxByID(data interface{}) error {
	return nil
}
