package models

import ()

type Account struct {
	ID        int64  `json:"id" orm:"column(id);auto"`
	AccountID string `json:"account_id" orm:"column(account_id);unique"`
	Email     string `json:"email" orm:"column(email)"`
}

func NewAccount(accountID string, email string) *Account {
	return &Account{AccountID: accountID, Email: email}
}
