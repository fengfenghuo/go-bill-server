package tx

import (
	"github.com/bill-server/go-bill-server/node/common"
)

type EvtTxUpdated struct {
	txID common.TxID
}

func (srv *Service) NotifyTxUpdated(txID common.TxID) {
	result := EvtTxUpdated{txID: txID}
	srv.ctx.NotifyEvent(result)
}
