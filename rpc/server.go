package rpc

import (
	"context"
	"fmt"
	"github.com/op/go-logging"
	"reflect"
	"runtime"
	"sync/atomic"
)

const (
	MetadataApi = "rpc"
)

var log = logging.MustGetLogger("rpc")

// callback is a method callback which was registered in the server
type callback struct {
	rcvr        reflect.Value  // receiver of method
	method      reflect.Method // callback
	argTypes    []reflect.Type // input argument types
	hasCtx      bool           // method's first argument is a context (not included in argTypes)
	errPos      int            // err return idx, of -1 when method cannot return error
	isSubscribe bool           // indication if the callback is a subscription
}

type callbacks map[string]*callback

type service struct {
	name      string
	typ       reflect.Type
	callbacks callbacks
}

type serviceRegistry map[string]*service

type Server struct {
	services serviceRegistry
	run      int32
}

func NewServer() *Server {
	server := &Server{
		services: serviceRegistry{},
		run:      1,
	}

	rpcService := &RPCService{server}
	server.RegisterName(MetadataApi, rpcService)

	return server
}

func (s *Server) ForEachMethod(visitor func(service *ServiceInfo, method *MethodInfo)) {
	for _, service := range s.services {
		serviceInfo := ServiceInfo{Name: service.name}

		for name, callback := range service.callbacks {
			methodInfo := MethodInfo{
				Name:        name,
				Method:      callback.method,
				Receiver:    callback.rcvr,
				ArgTypes:    callback.argTypes,
				HasCtx:      callback.hasCtx,
				ErrPos:      callback.errPos,
				IsSubscribe: callback.isSubscribe,
			}
			visitor(&serviceInfo, &methodInfo)
		}
	}
}

func (s *Server) FindMethod(sName string, mName string) (*ServiceInfo, *MethodInfo) {
	for _, service := range s.services {
		if service.name != sName {
			continue
		}

		for name, callback := range service.callbacks {
			if name != mName {
				continue
			}

			return &ServiceInfo{
					Name: service.name,
				},
				&MethodInfo{
					Name:        name,
					Method:      callback.method,
					Receiver:    callback.rcvr,
					ArgTypes:    callback.argTypes,
					HasCtx:      callback.hasCtx,
					ErrPos:      callback.errPos,
					IsSubscribe: callback.isSubscribe,
				}
		}
	}

	return nil, nil
}

func (s *Server) ServeSingleRequest(ctx context.Context, req *RpcRequest) (*RpcResponse, error) {
	reqs := make([]*RpcRequest, 1)
	reqs[0] = req
	responses, err := s.serveRequest(ctx, reqs)

	if len(responses) > 0 {
		return responses[0], err
	} else {
		return nil, err
	}
}

func (s *Server) RegisterName(name string, rcvr interface{}) error {
	svc := new(service)
	svc.typ = reflect.TypeOf(rcvr)
	rcvrVal := reflect.ValueOf(rcvr)

	if name == "" {
		return fmt.Errorf("no service name for type %s", svc.typ.String())
	}

	if !isExported(reflect.Indirect(rcvrVal).Type().Name()) {
		return fmt.Errorf("%s is not exported", reflect.Indirect(rcvrVal).Type().Name())
	}

	methods := suitableCallbacks(rcvrVal, svc.typ)

	if regsvc, present := s.services[name]; present {
		if len(methods) == 0 {
			return fmt.Errorf("Service %T doesn't have any suitable methods to expose", rcvr)
		}
		for _, m := range methods {
			regsvc.callbacks[formatName(m.method.Name)] = m
		}
		return nil
	}

	svc.name = name
	svc.callbacks = methods

	if len(svc.callbacks) == 0 {
		return fmt.Errorf("Service %T doesn't have any suitable methods to expose", rcvr)
	}

	s.services[svc.name] = svc
	return nil
}

func (s *Server) serveRequest(ctx context.Context, reqs []*RpcRequest) ([]*RpcResponse, error) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Error(string(buf))
		}
	}()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if atomic.LoadInt32(&s.run) != 1 { // server stopped
		return nil, &shutdownError{}
	}

	reses := make([]*RpcResponse, len(reqs))

	for i := 0; i < len(reqs); i++ {
		res, err := s.handle(ctx, reqs[i])
		if err != nil {
			return nil, err
		}
		reses[i] = res
	}

	return reses, nil
}

func (s *Server) handle(ctx context.Context, req *RpcRequest) (*RpcResponse, error) {
	if len(req.Args) != req.Method.Method.Type.NumIn() {
		log.Error(req.Service.Name+"."+req.Method.Name+" request arg count mismatch!",
			"require", req.Method.Method.Type.NumIn(), "input", len(req.Args))
		return nil, &InvalidParamsError{}
	}

	if req.Method.HasCtx {
		req.Args[1] = reflect.ValueOf(ctx)
	}

	reply := req.Method.Method.Func.Call(req.Args)

	var data interface{} = nil
	if req.Method.ErrPos >= 0 {
		if !reply[req.Method.ErrPos].IsNil() { // test if method returned an error
			err := reply[req.Method.ErrPos].Interface().(error)
			log.Debugf("%v ==> %s", req, err.Error())
			return &RpcResponse{
				Err: err,
			}, nil
		} else {
			for i := 0; i < req.Method.Method.Type.NumOut(); i++ {
				if !reply[i].IsNil() {
					data = reply[i].Interface()
					break
				}
			}
		}
	} else {
		data = reply[0].Interface()
	}

	return &RpcResponse{
		Data: data,
	}, nil
}

type RPCService struct {
	server *Server
}

func (s *RPCService) Modules() map[string]string {
	modules := make(map[string]string)
	for name := range s.server.services {
		modules[name] = "1.0"
	}
	return modules
}
