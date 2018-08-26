package tx

import (
	"context"
	"fmt"
	"reflect"

	"github.com/bill-server/go-bill-server/common"
	"github.com/bill-server/go-bill-server/rpc"
)

type PublicTxAPI struct {
	srv *Service
}

func NewPublicTxAPI(srv *Service) *PublicTxAPI {
	return &PublicTxAPI{
		srv: srv,
	}
}

type ApiCreateTxRequest struct {
	Account common.AccountID
}

type ApiCreateTxResult struct {
	Account common.AccountID
}

func (api *PublicTxAPI) Createtx(ctx context.Context, req *ApiCreateTxRequest) (*ApiCreateTxResult, error) {
	api.srv.mu.RLock()
	defer api.srv.mu.RUnlock()

	result := ApiCreateTxResult{}
	return &result, nil
}

func (api *PublicTxAPI) Subscribe(ctx context.Context) (*rpc.Subscription, error) {
	notifier := rpc.NotifierFromContext(ctx)
	if notifier == nil {
		return nil, fmt.Errorf("%s.SubscribeTx: can`t find filter from context!", ServiceName)
	}

	subscription := notifier.CreateSubscription()
	if subscription == nil {
		return nil, fmt.Errorf("%s.SubscribeTx: create subscription fail!", ServiceName)
	}

	c := make(chan interface{})
	apiSub := api.srv.ctx.Subscribe([]reflect.Type{reflect.TypeOf(EvtTxUpdated{})}, c)

	go func() {
		defer func() {
			apiSub.Unsubscribe()
			close(c)
		}()

		for {
			select {
			case v := <-c:
				if _, ok := v.(EvtTxUpdated); ok {
					account := "HelloWorld"
					notifier.Notify(subscription.ID, account)
				}
			case <-subscription.Err():
			}
		}
	}()

	return subscription, nil
}
