package account

import (
	"github.com/bill-server/go-bill-server/node"
	"github.com/bill-server/go-bill-server/rpc"
	"github.com/op/go-logging"
	"reflect"
	"sync"
)

const (
	ServiceName    = "account"
	ServiceVersion = "1.0"
)

var (
	ServiceType = reflect.TypeOf(Service{})
	log         = logging.MustGetLogger("account")
)

type Service struct {
	ctx *node.ServiceContext
	mu  sync.RWMutex
}

func NewService(ctx *node.ServiceContext) (node.Service, error) {
	log.Info("service " + ServiceName + " create success")

	return &Service{
		ctx: ctx,
	}, nil
}

func (srv *Service) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: ServiceName,
			Version:   ServiceVersion,
			Service:   NewPublicAccountAPI(srv),
			Public:    true,
		},
	}
}

func (srv *Service) Start() error {
	log.Info("service " + ServiceName + " start success")

	return nil
}

func (srv *Service) Stop() error {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	log.Info("service " + ServiceName + " stop success")
	return nil
}
