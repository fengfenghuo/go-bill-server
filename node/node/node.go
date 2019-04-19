package node

import (
	"errors"
	"reflect"
	"sync"

	"github.com/bill-server/go-bill-server/node/event"
	"github.com/bill-server/go-bill-server/node/redis"
	"github.com/bill-server/go-bill-server/node/rpc"
	"github.com/bill-server/go-bill-server/node/sfoxdb"
	"github.com/op/go-logging"
)

var (
	ErrNodeRunning    = errors.New("node already running")
	ErrNodeNotRunning = errors.New("node not running")
	log               = logging.MustGetLogger("node")
)

type Node struct {
	db           sfoxdb.SFoxDb
	redis        redis.DataRedis
	ec           *event.EventCenter
	runing       bool
	serviceFuncs []ServiceConstructor
	services     map[reflect.Type]Service
	rpcServer    *rpc.Server
	rpcAPIs      []rpc.API
	rpcEndpoints []rpc.Endpoint
	wg           *sync.WaitGroup
	lock         sync.RWMutex
}

func NewNode(db sfoxdb.SFoxDb, redis redis.DataRedis) (*Node, error) {
	return &Node{
		db:           db,
		redis:        redis,
		ec:           event.NewEventCenter(),
		serviceFuncs: []ServiceConstructor{},
		wg:           &sync.WaitGroup{},
	}, nil
}

func (n *Node) Register(constructor ServiceConstructor) error {
	n.lock.Lock()
	defer n.lock.Unlock()

	if n.runing {
		return ErrNodeRunning
	}

	n.serviceFuncs = append(n.serviceFuncs, constructor)
	return nil
}

func (n *Node) Start() error {
	n.lock.Lock()
	defer n.lock.Unlock()

	if n.runing {
		return ErrNodeRunning
	}

	/*创建services */
	services := make(serviceContainer)
	for _, constructor := range n.serviceFuncs {
		ctx := &ServiceContext{
			DB:       n.db,
			Redis:    n.redis,
			ec:       n.ec,
			services: serviceContainer{},
		}

		for kind, s := range services {
			ctx.services[kind] = s
		}

		service, err := constructor(ctx)
		if err != nil {
			return err
		}

		kind := reflect.TypeOf(service)
		if kind.Kind() == reflect.Ptr {
			kind = kind.Elem()
		}

		if _, exists := services[kind]; exists {
			return &DuplicateServiceError{Kind: kind}
		}
		services[kind] = service
	}

	/*启动所有Service */
	started := []reflect.Type{}
	for kind, service := range services {
		if err := service.Start(); err != nil {
			for _, kind := range started {
				services[kind].Stop()
			}
			return err
		}

		started = append(started, kind)
	}

	/*创建rpcServer*/
	rpcServer := rpc.NewServer()
	apis := n.basicApis()
	for _, service := range services {
		apis = append(apis, service.APIs()...)
	}
	for _, api := range apis {
		if err := rpcServer.RegisterName(api.Namespace, api.Service); err != nil {
			return err
		}
		log.Debug("InProc registered", "service", api.Service, "namespace", api.Namespace)
	}

	/*成功 */
	n.services = services
	n.rpcServer = rpcServer
	n.runing = true

	return nil
}

func (n *Node) Cancel() error {
	n.lock.Lock()
	defer n.lock.Unlock()

	if !n.runing {
		return ErrNodeNotRunning
	}

	//TODO:
	return nil
}

func (n *Node) Wait() {
	n.lock.RLock()
	if !n.runing {
		n.lock.RUnlock()
		return
	}
	n.lock.RUnlock()

	n.wg.Wait()
}

func (n *Node) AddRpcEndpoint(endpoint rpc.Endpoint) error {
	n.lock.Lock()
	defer n.lock.Unlock()

	if !n.runing {
		return ErrNodeNotRunning
	}

	n.wg.Add(1)
	go func() {
		defer n.wg.Done()
		log.Info("endpoint " + endpoint.Name() + " start")
		err := endpoint.RunLoop()
		if err != nil {
			log.Error("endpoint "+endpoint.Name()+" stop: ", err.Error())
		} else {
			log.Info("endpoint " + endpoint.Name() + " stop")
		}
	}()

	n.rpcEndpoints = append(n.rpcEndpoints, endpoint)
	return nil
}

func (n *Node) basicApis() []rpc.API {
	return nil
}

func (n *Node) RpcServer() *rpc.Server {
	return n.rpcServer
}
