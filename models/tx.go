package models

import (
	"fmt"
	"time"

	"github.com/bill-server/go-bill-server/db"
)

type TxReason uint8
type TxType uint8

const (
	TxTypeUnknow TxType = iota
	TxTypeIncome TxType = iota
	TxTypeOutgo  TxType = iota
)

const (
	TxReasonUnknow  TxReason = iota
	TxReasonBuy     TxReason = iota
	TxReasonTraffic TxReason = iota
)

type Tx struct {
	ID         int64     `json:"id" orm:"column(id);auto"`
	AccountID  string    `json:"account_id" orm:"column(account_id)"`
	Amount     int64     `json:"amount" orm:"column(amount)"`
	TxType     TxType    `json:"tx_type" orm:"column(tx_type)"`
	Reason     TxReason  `json:"reason" orm:"column(reason)"`
	Remarks    string    `json:"remarks" orm:"column(remarks);null"`
	CreateTime time.Time `json:"create_time" orm:"auto_now_add;column(create_time)"`
}

func NewTx(accountID string, amount int64, txType TxType, reason TxReason, remarks string) *Tx {
	return &Tx{AccountID: accountID, Amount: amount, TxType: txType, Reason: reason, Remarks: remarks}
}

func (tx *Tx) CreateTx() error {
	d, err := db.GetDBInstance()
	if err != nil {
		return fmt.Errorf("GetDBInstance error: " + err.Error())
	}

	tx.CreateTime = time.Now()

	_, err = d.Insert(tx)
	if err != nil {
		return fmt.Errorf("Insert error: " + err.Error())
	}
	return nil
}
