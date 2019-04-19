package account

import (
	"github.com/bill-server/go-bill-server/node/common"
)

type EvtAccountUpdated struct {
	account common.AccountID
}

func (srv *Service) NotifyTxUpdated(account common.AccountID) {
	result := EvtAccountUpdated{account: account}
	srv.ctx.NotifyEvent(result)
}
