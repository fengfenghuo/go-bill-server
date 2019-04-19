package account

import (
	"context"
	"fmt"
	"github.com/bill-server/go-bill-server/node/common"
	"github.com/bill-server/go-bill-server/node/rpc"
	"reflect"
)

type PublicAccountAPI struct {
	srv *Service
}

func NewPublicAccountAPI(srv *Service) *PublicAccountAPI {
	return &PublicAccountAPI{
		srv: srv,
	}
}

type ApiCreateTxRequest struct {
	Account common.AccountID
}

type ApiCreateTxResult struct {
	Account common.AccountID
}

func (api *PublicAccountAPI) Createtx(ctx context.Context, req *ApiCreateTxRequest) (*ApiCreateTxResult, error) {
	api.srv.mu.RLock()
	defer api.srv.mu.RUnlock()

	result := ApiCreateTxResult{}
	return &result, nil
}

/*订阅通知 */
func (api *PublicAccountAPI) Subscribe(ctx context.Context) (*rpc.Subscription, error) {
	notifier := rpc.NotifierFromContext(ctx)
	if notifier == nil {
		return nil, fmt.Errorf("%s.SubscribeTx: can`t find filter from context!", ServiceName)
	}

	subscription := notifier.CreateSubscription()
	if subscription == nil {
		return nil, fmt.Errorf("%s.SubscribeTx: create subscription fail!", ServiceName)
	}

	c := make(chan interface{})
	apiSub := api.srv.ctx.Subscribe([]reflect.Type{reflect.TypeOf(EvtAccountUpdated{})}, c)

	go func() {
		defer func() {
			apiSub.Unsubscribe()
			close(c)
		}()

		for {
			select {
			case v := <-c:
				if _, ok := v.(EvtAccountUpdated); ok {
					// notifier.Notify(subscription.ID, ApiBlock{BlockHeight: evt.blockHeight, BlockId: evt.blockId, Time: common.TimeFormat(evt.time), Txes: txes})
				}
			case <-subscription.Err():
			}
		}
	}()
	return subscription, nil
}
