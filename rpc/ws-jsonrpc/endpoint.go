package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/bill-server/go-bill-server/rpc"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/op/go-logging"
)

const (
	EndpointName = "WS-JSONRPC"
	contentType  = "application/json"
)

var log = logging.MustGetLogger(EndpointName)

type subscribeData map[string]rpc.ID

type jsonRpcEndpoint struct {
	srv        *rpc.Server
	httpServer *http.Server
	upgrader   *websocket.Upgrader

	clients            map[*websocket.Conn]subscribeData
	subscribesMu       sync.Mutex
	internalSubscribes map[url.URL]subscribeData
	notifier           *rpc.Notifier
}

func NewEndpoint(srv *rpc.Server, addr string) *jsonRpcEndpoint {
	var endpoint *jsonRpcEndpoint

	endpoint = &jsonRpcEndpoint{
		srv:                srv,
		clients:            map[*websocket.Conn]subscribeData{},
		internalSubscribes: map[url.URL]subscribeData{},
		notifier:           rpc.NewNotifier(func(id rpc.ID, data interface{}) error { return endpoint.notifyEvent(id, data) }),
	}

	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		endpoint.callMethod(w, r)
	})

	srv.ForEachMethod(func(s *rpc.ServiceInfo, m *rpc.MethodInfo) {
		log.Debugf("\t" + s.Name + "." + m.Name + " start listen")
	})

	endpoint.upgrader = &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	endpoint.httpServer = &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	return endpoint
}

func (ep *jsonRpcEndpoint) Name() string {
	return EndpointName
}

func (ep *jsonRpcEndpoint) RunLoop() error {
	err := ep.httpServer.ListenAndServe()
	return err
}

func (ep *jsonRpcEndpoint) Shotdown(ctx context.Context) error {
	err := ep.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}

	ep.httpServer = nil

	return nil
}

func (ep *jsonRpcEndpoint) callMethod(w http.ResponseWriter, r *http.Request) {

	ctx := buildContext(r.Context(), r)

	ws, err := ep.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err.Error())
	}

	defer func() {
		ws.Close()
		delete(ep.clients, ws)
	}()

	ep.clients[ws] = nil

	for {
		var res *rpc.RpcResponse

		reqMethodInfo, err := ep.readJson(ws)
		if err != nil {
			log.Error("endpoint readJson fail, error: " + err.Error())
			break
		}

		s, m, err := ep.findMethod(reqMethodInfo)
		if err != nil {
			log.Error("endpoint findMethod fail, error: " + err.Error())
			res = &rpc.RpcResponse{Err: err}
			ep.sendResponse(ws, reqMethodInfo, res)
			continue
		}

		req, err := ep.readRequest(ctx, s, m, reqMethodInfo, ws, r)
		if err != nil {
			log.Error(err.Error())
			res = &rpc.RpcResponse{Err: err}
			ep.sendResponse(ws, reqMethodInfo, res)
			continue
		}

		if strings.Compare("subscribe", m.Name) == 0 {
			_, err := ep.Subscribe(ws, req, reqMethodInfo)
			if err != nil {
				log.Error(err.Error())
				res = &rpc.RpcResponse{Err: err}
				ep.sendResponse(ws, reqMethodInfo, res)
				continue
			}
			res = &rpc.RpcResponse{Data: rpc.ResCommonResult{Status: "success"}}
			ep.sendResponse(ws, reqMethodInfo, res)
			continue
		} else if strings.Compare("unsubscribe", m.Name) == 0 {
			err := ep.UnSubscribe(ws, req, reqMethodInfo)
			if err != nil {
				log.Error(err.Error())
				res = &rpc.RpcResponse{Err: err}
				ep.sendResponse(ws, reqMethodInfo, res)
				continue
			}
			res = &rpc.RpcResponse{Data: rpc.ResCommonResult{Status: "success"}}
			ep.sendResponse(ws, reqMethodInfo, res)
			continue
		}

		res, err = ep.srv.ServeSingleRequest(ctx, req)
		if err != nil {
			log.Error(err.Error())
			res = &rpc.RpcResponse{Err: err}
			ep.sendResponse(ws, reqMethodInfo, res)
			continue
		}

		if res == nil {
			log.Error("no response or error")
			w.WriteHeader(500)
			continue
		}

		ep.sendResponse(ws, reqMethodInfo, res)
	}
}

func (ep *jsonRpcEndpoint) sendResponse(ws *websocket.Conn, reqMethodInfo *rpc.ReqMethodInfo, res *rpc.RpcResponse) {
	resMethodInfo := rpc.ResMethodInfo{ID: reqMethodInfo.ID}
	if res.Err != nil {
		resMethodInfo.Error = res.Err.Error()
	} else if res.Data != nil {
		resMethodInfo.Result = res.Data
	}
	ws.WriteJSON(resMethodInfo)
}

func (ep *jsonRpcEndpoint) readJson(ws *websocket.Conn) (*rpc.ReqMethodInfo, error) {
	var reqMethodInfo rpc.ReqMethodInfo
	err := ws.ReadJSON(&reqMethodInfo)
	if err != nil {
		return nil, err
	}
	return &reqMethodInfo, nil
}

func (ep *jsonRpcEndpoint) findMethod(reqMethodInfo *rpc.ReqMethodInfo) (*rpc.ServiceInfo, *rpc.MethodInfo, error) {
	method := strings.Split(reqMethodInfo.Method, ".")
	s, m := ep.srv.FindMethod(method[0], method[1])
	if s == nil {
		return nil, nil, fmt.Errorf("endpoint readRequest find service %s not exist!", method[0])
	}

	if m == nil {
		return nil, nil, fmt.Errorf("endpoint readRequest find service %s method %s not exist!", method[0], method[1])
	}
	return s, m, nil
}

func (ep *jsonRpcEndpoint) readRequest(ctx context.Context, s *rpc.ServiceInfo, m *rpc.MethodInfo, reqMethodInfo *rpc.ReqMethodInfo, ws *websocket.Conn, r *http.Request) (*rpc.RpcRequest, error) {

	req := rpc.NewRpcRequest(s, m, ctx)

	foreachStructType(req, func(v reflect.Value) {
		structValue := v.Elem()
		argType := structValue.Type()

		for j := 0; j < argType.NumField(); j++ {
			fieldValue := structValue.Field(j).Addr()
			if fieldValue.CanInterface() {
				v := reflect.ValueOf(reqMethodInfo.Params)
				if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
					return
				}

				u := fieldValue.Interface()
				if v.Len() == 0 {
					return
				}
				msg, err := json.Marshal(v.Index(j).Interface())
				if err != nil {
					log.Error(err.Error())
					return
				}

				if err = json.Unmarshal(msg, u); err != nil {
					log.Error(err.Error())
					return
				}
			}
		}
	})
	return req, nil
}

func foreachStructType(req *rpc.RpcRequest, visitor func(v reflect.Value)) {
	argStart := 1
	if req.Method.HasCtx {
		argStart++
	}

	for i := argStart; i < len(req.Args); i++ {
		argType := req.Args[i].Type()
		if argType.Kind() != reflect.Ptr && argType.Elem().Kind() != reflect.Struct {
			continue
		}

		visitor(req.Args[i])
	}
}

func buildContext(ctx context.Context, r *http.Request) context.Context {
	ctx = context.WithValue(ctx, "remote", r.RemoteAddr)
	ctx = context.WithValue(ctx, "scheme", r.Proto)
	ctx = context.WithValue(ctx, "local", r.Host)
	return ctx
}

func (ep *jsonRpcEndpoint) notifyWSEvent(id rpc.ID, data interface{}) error {
	var ws *websocket.Conn
	var method string

	ep.subscribesMu.Lock()
	for checkWs, subscribeTarget := range ep.clients {
		for checkMethod, checkId := range subscribeTarget {
			if checkId == id {
				ws = checkWs
				method = checkMethod
			}
		}
	}
	ep.subscribesMu.Unlock()

	if ws == nil {
		return fmt.Errorf("no subscriber")
	}

	sendData := rpc.SubscribeMethodInfo{Method: method, Result: data}
	jsonData, err := json.Marshal(sendData)
	if err != nil {
		log.Error("endpoint: notifyEvent: json.Marshal fail, " + err.Error())
		return err
	}
	log.Debugf("endpoint notify to %s event %s", ws.LocalAddr().String(), jsonData)

	err = ws.WriteJSON(sendData)
	if err != nil {
		return err
	}

	return nil
}

func (ep *jsonRpcEndpoint) notifyInternalSubscribeEvent(id rpc.ID, data interface{}) error {
	var url *url.URL = nil

	ep.subscribesMu.Lock()
	for checkUrl, subscribeTarget := range ep.internalSubscribes {
		for _, checkId := range subscribeTarget {
			if checkId == id {
				urlBuf := checkUrl
				url = &urlBuf
			}
		}
	}

	ep.subscribesMu.Unlock()

	if url == nil {
		return fmt.Errorf("no subscriber")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Error("RESTFul: notifyEvent: json.Marshal fail, " + err.Error())
		return err
	}

	log.Debugf("endpoint notify to %s event %s", url, jsonData)

	resp, err := http.Post(url.String(), contentType, bytes.NewReader(jsonData))
	if err != nil {
		log.Error(err.Error())
		return err
	}

	defer resp.Body.Close()

	return nil
}

func (ep *jsonRpcEndpoint) notifyEvent(id rpc.ID, data interface{}) error {
	var url *url.URL = nil
	var ws *websocket.Conn = nil

	ep.subscribesMu.Lock()
	for checkUrl, subscribeTarget := range ep.internalSubscribes {
		for _, checkId := range subscribeTarget {
			if checkId == id {
				urlBuf := checkUrl
				url = &urlBuf
			}
		}
	}

	for checkWs, subscribeTarget := range ep.clients {
		for _, checkID := range subscribeTarget {
			if checkID == id {
				ws = checkWs
			}
		}
	}

	if url != nil {
		return ep.notifyInternalSubscribeEvent(id, data)
	} else if ws != nil {
		return ep.notifyWSEvent(id, data)
	} else {
		log.Error("endpoint notify event no subscribe")
	}

	return nil
}

func (ep *jsonRpcEndpoint) InternalSubscribe(svrName string, methodName string, strSubscribeUrl string) (*rpc.ID, error) {
	subscribeUri, err := url.Parse(strSubscribeUrl)
	if err != nil {
		log.Error(svrName + "." + methodName + ": subscribeUri " + strSubscribeUrl + " format error, " + err.Error())
		return nil, err
	}

	ep.subscribesMu.Lock()
	defer ep.subscribesMu.Unlock()

	if subscribeTarget, ok := ep.internalSubscribes[*subscribeUri]; ok {
		if id, ok := subscribeTarget[svrName+"."+methodName]; ok {
			return &id, nil
		}
	}

	ctx := context.Background()
	ctx = rpc.AttachNotifier(ctx, ep.notifier)

	s, m := ep.srv.FindMethod(svrName, methodName)
	if m == nil {
		log.Error(svrName + "." + methodName + ": method not exist")
		return nil, fmt.Errorf(svrName + "." + methodName + ": method not exist")
	}

	req := rpc.NewRpcRequest(s, m, ctx)

	res, err := ep.srv.ServeSingleRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	if res.Data == nil {
		log.Error(svrName + "." + methodName + ": call success, no subscribe")
		return nil, fmt.Errorf(svrName + "." + methodName + ": call success, no subscribe")
	}

	subscription, ok := res.Data.(*rpc.Subscription)
	if !ok || subscription == nil {
		log.Error(s.Name + "." + m.Name + ": call success, subscribe typer error")
		return nil, fmt.Errorf(s.Name + "." + m.Name + ": call success, subscribe typer error")
	}
	ep.notifier.Activate(subscription.ID)

	subscribeTarget, ok := ep.internalSubscribes[*subscribeUri]
	if !ok {
		subscribeTarget = make(subscribeData)
		ep.internalSubscribes[*subscribeUri] = subscribeTarget
	}
	subscribeTarget[s.Name+"."+m.Name] = subscription.ID

	return &subscription.ID, nil
}

func (ep *jsonRpcEndpoint) Subscribe(ws *websocket.Conn, req *rpc.RpcRequest, reqMethodInfo *rpc.ReqMethodInfo) (*rpc.ID, error) {
	ep.subscribesMu.Lock()
	defer ep.subscribesMu.Unlock()

	s, m, err := ep.findMethod(reqMethodInfo)
	if err != nil {
		return nil, err
	}

	if subscribeTarget, ok := ep.clients[ws]; ok {
		if id, ok := subscribeTarget[s.Name+"."+m.Name]; ok {
			return &id, fmt.Errorf(s.Name + "." + m.Name + ": method already subscribe")
		}
	}

	ctx := context.Background()
	ctx = rpc.AttachNotifier(ctx, ep.notifier)

	res, err := ep.srv.ServeSingleRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	if res.Data == nil {
		return nil, fmt.Errorf(s.Name + "." + m.Name + ": call success, no subscribe")
	}

	subscription, ok := res.Data.(*rpc.Subscription)
	if !ok || subscription == nil {
		return nil, fmt.Errorf(s.Name + "." + m.Name + ": call success, subscribe typer error")
	}
	ep.notifier.Activate(subscription.ID)

	subscribeTarget, ok := ep.clients[ws]
	if !ok {
		subscribeTarget = make(subscribeData)
		ep.clients[ws] = subscribeTarget
	}
	subscribeTarget[s.Name+"."+m.Name] = subscription.ID

	return &subscription.ID, nil
}

func (ep *jsonRpcEndpoint) UnSubscribe(ws *websocket.Conn, req *rpc.RpcRequest, reqMethodInfo *rpc.ReqMethodInfo) error {
	ep.subscribesMu.Lock()
	defer ep.subscribesMu.Unlock()

	s, m, err := ep.findMethod(reqMethodInfo)
	if err != nil {
		return err
	}

	var subId rpc.ID
	if subscribeTarget, ok := ep.clients[ws]; ok {
		if id, ok := subscribeTarget[s.Name+".subscribe"]; ok {
			subId = id
		}
	}

	if subId == "" {
		return fmt.Errorf(s.Name + "." + m.Name + ": method not subscribe")
	}

	ctx := context.Background()
	ctx = rpc.AttachNotifier(ctx, ep.notifier)

	req = rpc.NewRpcRequest(s, m, ctx)

	_, err = ep.srv.ServeSingleRequest(ctx, req)
	if err != nil {
		return err
	}
	ep.notifier.Unsubscribe(subId)
	ep.clients[ws] = nil
	return nil
}
