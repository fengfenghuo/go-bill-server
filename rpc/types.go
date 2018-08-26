package rpc

import (
	"context"
	"fmt"
	"reflect"
)

type API struct {
	Namespace string      // namespace under which the rpc methods of Service are exposed
	Version   string      // api version for DApp's
	Service   interface{} // receiver instance which holds the methods
	Public    bool        // indication if the methods must be considered safe for public use
}

type ServiceInfo struct {
	Name string
}

type MethodInfo struct {
	Name        string
	Receiver    reflect.Value
	Method      reflect.Method
	ArgTypes    []reflect.Type
	HasCtx      bool
	ErrPos      int
	IsSubscribe bool
}

type ReqMethodInfo struct {
	Method string      `json:"method"`
	ID     uint32      `json:"id"`
	Params interface{} `json:"params"`
}

type ResCommonResult struct {
	Status string `json:"status"`
}

type ResMethodInfo struct {
	Error  interface{} `json:"error"`
	ID     uint32      `json:"id"`
	Result interface{} `json:"result"`
}

type SubscribeMethodInfo struct {
	Method string      `json:"method"`
	ID     uint32      `json:"id"`
	Result interface{} `json:"result"`
}

type RpcRequest struct {
	Service *ServiceInfo
	Method  *MethodInfo
	Args    []reflect.Value
}

type RpcResponse struct {
	Data interface{}
	Err  error
}

type Endpoint interface {
	Name() string
	RunLoop() error
	Shotdown(ctx context.Context) error
	InternalSubscribe(svrName string, methodName string, strSubscribeUrl string) (*ID, error)
}

func NewRpcRequest(s *ServiceInfo, m *MethodInfo, ctx context.Context) *RpcRequest {
	req := RpcRequest{
		Service: s,
		Method:  m,
		Args:    []reflect.Value{m.Receiver},
	}

	if m.HasCtx {
		req.Args = append(req.Args, reflect.ValueOf(ctx))
	}

	argCount := len(m.ArgTypes)
	if argCount > 0 {
		args := make([]reflect.Value, argCount)
		for i := 0; i < argCount; i++ {
			argType := m.ArgTypes[i]

			if argType.Kind() == reflect.Ptr {
				argVal := reflect.New(argType.Elem())
				args[i] = argVal
			} else {
				argVal := reflect.New(argType)
				args[i] = argVal.Elem()
			}
		}
		req.Args = append(req.Args, args...)
	}

	return &req
}

func (req *RpcRequest) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%s.%s(", req.Service.Name, req.Method.Name)

	argPos := 1
	if req.Method.HasCtx {
		argPos += 1
	}

	for ; argPos < len(req.Args); argPos++ {
		arg := req.Args[argPos]
		fmt.Fprintf(s, "%v", arg)
	}

	fmt.Fprintf(s, ")")
}
