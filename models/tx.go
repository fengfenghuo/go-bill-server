package models

type TxReason uint8

type Tx struct {
	ID        int64    `json:"id" orm:"column(id);auto"`
	AccountID string   `json:"account_id" orm:"column(account_id)"`
	Amount    int64    `json:"amount" orm:"column(amount)"`
	Reason    TxReason `json:"reason" orm:"column(reason)"`
	Remarks   string   `json:"remarks" orm:"column(remarks)"`
}

func NewTx(accountID string, amount int64, reason TxReason, remarks string) *Tx {
	return &Tx{AccountID: accountID, Amount: amount, Reason: reason, Remarks: remarks}
}