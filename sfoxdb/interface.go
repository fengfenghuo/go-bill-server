package sfoxdb

import (
	"github.com/bill-server/go-bill-server/common"
)

type SFoxDb interface {
	QueryTXByID(id int64) (tx *common.Tx, err error)
	QueryAccountByAddress(address common.AccountID) (account *common.Account, err error)

	InsertAccountData(account *common.Account) (err error)
	InsertTx(tx *common.Tx) (err error)
}
