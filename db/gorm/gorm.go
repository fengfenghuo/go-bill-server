package gorm

import (
	"fmt"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type GormInterface struct {
	gormDB *gorm.DB
}

func RegisterDB() (*GormInterface, error) {
	// (可选)设置最大空闲连接
	maxIdle := beego.AppConfig.DefaultInt("db::maxIdle", 60)
	// (可选) 设置最大数据库连接 (go >= 1.2)
	maxConn := beego.AppConfig.DefaultInt("db::maxConn", 320)
	dbLink := beego.AppConfig.DefaultString("db::db", "")
	if dbLink == "" {
		return nil, fmt.Errorf("no db link")
	}

	db, err := gorm.Open("mysql", dbLink)
	if err != nil {
		return nil, fmt.Errorf("RegisterDateBase error: " + err.Error())
	}

	db.DB().SetMaxIdleConns(maxIdle)
	db.DB().SetMaxOpenConns(maxConn)

	return &GormInterface{gormDB: db}, nil
}

func (db *GormInterface) RegisterTable(modules ...interface{}) error {
	db.gormDB.SingularTable(true)
	for _, module := range modules {
		if db.gormDB.HasTable(module) {
			continue
		}

		db.gormDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(module)
	}
	return nil
}

func (db *GormInterface) Insert(data interface{}) (int64, error) {
	dbTemp := db.gormDB.Create(data)
	if dbTemp.Error != nil {
		return 0, dbTemp.Error
	}
	return dbTemp.RowsAffected, nil
}

func (db *GormInterface) QueryAccount(account string, data interface{}) error {
	dbTemp := db.gormDB.Where("account_id = ?", account).Find(data)
	if dbTemp.Error != nil {
		return dbTemp.Error
	}
	return nil
}

func (db *GormInterface) QueryTxesByAccount(account string, count, offset int, data interface{}) error {
	dbTemp := db.gormDB.Where("account_id = ?", account).Offset(offset).Count(count).Find(data)
	if dbTemp.Error != nil {
		return dbTemp.Error
	}
	return nil
}
