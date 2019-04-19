package node

import (
	"reflect"

	"github.com/bill-server/go-bill-server/node/event"
	"github.com/bill-server/go-bill-server/node/redis"
	"github.com/bill-server/go-bill-server/node/rpc"
	"github.com/bill-server/go-bill-server/node/sfoxdb"
)

type serviceContainer map[reflect.Type]Service

type ServiceContext struct {
	DB       sfoxdb.SFoxDb
	Redis    redis.DataRedis
	ec       *event.EventCenter
	services serviceContainer
}

type ServiceConstructor func(ctx *ServiceContext) (Service, error)

type Service interface {
	APIs() []rpc.API
	Start() error
	Stop() error
}

func (ctx *ServiceContext) NotifyEvent(evt interface{}) {
	ctx.ec.NotifyEvent(evt)
}

func (ctx *ServiceContext) Subscribe(types []reflect.Type, c chan interface{}) *event.Subscription {
	return ctx.ec.Subscribe(types, c)
}

func (ctx *ServiceContext) FindService(srvType reflect.Type) Service {
	if srv, ok := ctx.services[srvType]; ok {
		return srv
	}
	return nil
}
