package sfoxdb

import (
	"github.com/bill-server/go-bill-server/common"
)

type Account struct {
	Address common.AccountID
}

type Tx struct {
	TxID    common.TxID
	Address common.AccountID
}

type SFoxDb interface {
	QueryTXByID(id int64) (tx *Tx, err error)
	QueryAccountByAddress(address common.AccountID) (account *Account, err error)
}
