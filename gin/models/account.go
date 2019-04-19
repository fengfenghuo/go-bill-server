package models

import (
	"fmt"
	"time"

	"github.com/bill-server/go-bill-server/node/db"
)

type Account struct {
	ID         int64     `json:"id" orm:"column(id);auto" gorm:"primary_key"`
	AccountID  string    `json:"account_id" orm:"column(account_id);unique"`
	Password   string    `json:"password" orm:"column(password);null"`
	Email      string    `json:"email" orm:"column(email);null"`
	CreateTime time.Time `json:"create_time" orm:"auto_now_add;column(create_time)"`
}

func NewAccount(accountID string, password string, email string) *Account {
	return &Account{AccountID: accountID, Password: password, Email: email}
}

func (ac *Account) Register() error {
	d, err := db.GetDBInstance()
	if err != nil {
		return fmt.Errorf("GetDBInstance error: " + err.Error())
	}
	ac.CreateTime = time.Now()
	_, err = d.Insert(ac)
	if err != nil {
		return fmt.Errorf("Insert error: " + err.Error())
	}
	return nil
}
