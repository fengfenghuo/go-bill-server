package tx

import (
	"sync"

	"github.com/bill-server/go-bill-server/node"
	"github.com/bill-server/go-bill-server/rpc"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("tx")
)

const (
	ServiceName    = "tx"
	ServiceVersion = "1.0"
)

type Service struct {
	ctx *node.ServiceContext
	mu  sync.RWMutex
}

func NewService(ctx *node.ServiceContext) (node.Service, error) {
	log.Info("service " + ServiceName + " create success")

	return &Service{
		ctx: ctx}, nil
}

func (srv *Service) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: ServiceName,
			Version:   ServiceVersion,
			Service:   NewPublicTxAPI(srv),
			Public:    true,
		},
	}
}

func (srv *Service) Start() error {
	log.Info("service " + ServiceName + " start success")
	return nil
}

func (srv *Service) Stop() error {
	log.Info("service " + ServiceName + " stop success")
	return nil
}
