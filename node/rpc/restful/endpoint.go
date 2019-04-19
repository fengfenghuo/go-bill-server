package restful

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/bill-server/go-bill-server/node/common"
	"github.com/bill-server/go-bill-server/node/rpc"
	"github.com/julienschmidt/httprouter"
	"github.com/op/go-logging"
	"github.com/rs/cors"
)

const (
	EndpointName            = "RESTFul"
	contentType             = "application/json"
	maxRequestContentLength = 1024 * 128
)

var log = logging.MustGetLogger(EndpointName)

type subscribeData map[string]rpc.ID

type httpRestFulEndpoint struct {
	srv        *rpc.Server
	httpServer *http.Server

	subscribesMu sync.Mutex
	subscribes   map[url.URL]subscribeData
	notifier     *rpc.Notifier
}

func NewEndpoint(allowedOrigins []string, srv *rpc.Server, addr string) *httpRestFulEndpoint {
	var endpoint *httpRestFulEndpoint

	endpoint = &httpRestFulEndpoint{
		srv:        srv,
		subscribes: map[url.URL]subscribeData{},
		notifier:   rpc.NewNotifier(func(id rpc.ID, data interface{}) error { return endpoint.notifyEvent(id, data) }),
	}

	router := httprouter.New()
	srv.ForEachMethod(func(s *rpc.ServiceInfo, m *rpc.MethodInfo) {
		hasAccountID := false

		for i := 0; i < len(m.ArgTypes); i++ {
			argType := m.ArgTypes[i]
			if argType.Kind() == reflect.Ptr {
				argType = argType.Elem()
			}

			if argType.Kind() != reflect.Struct {
				continue
			}

			if haveAccountIDField(argType) {
				hasAccountID = true
				break
			}
		}

		mountPoint := ""
		if hasAccountID {
			mountPoint += "/account/:AccountID/" + s.Name + "/"
		} else {
			mountPoint += "/" + s.Name + "/"
		}

		if strings.HasPrefix(m.Name, "update") {
			mountPoint += formatEndpointName(m.Name, s.Name, "update", hasAccountID)
			log.Debugf("    mount %s post", mountPoint)
			router.POST(mountPoint, func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
				endpoint.callMethod(w, r, ps, s, m, hasAccountID)
			})
		} else if strings.HasPrefix(m.Name, "create") {
			mountPoint += formatEndpointName(m.Name, s.Name, "create", hasAccountID)
			log.Debugf("    mount %s post", mountPoint)
			router.POST(mountPoint, func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
				endpoint.callMethod(w, r, ps, s, m, hasAccountID)
			})
		} else if strings.HasPrefix(m.Name, "subscribe") {
			if !m.IsSubscribe {
				log.Error(s.Name + "." + m.Name + ": is not subscribe function")
				return
			}

			mountPoint += "subscribe/" + formatEndpointName(m.Name, s.Name, "subscribe", hasAccountID)
			log.Debugf("    mount %s post", mountPoint)
			router.POST(mountPoint, func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
				endpoint.callMethod(w, r, ps, s, m, hasAccountID)
			})
		} else if strings.HasPrefix(m.Name, "unSubscribe") {
			mountPoint += "unsubscribe/" + formatEndpointName(m.Name, s.Name, "unSubscribe", hasAccountID)
			log.Debugf("    mount %s post", mountPoint)
			router.POST(mountPoint, func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
				endpoint.callMethod(w, r, ps, s, m, hasAccountID)
			})
		} else {
			mountPoint += formatEndpointName(m.Name, s.Name, "", hasAccountID)
			mountPoint += formatEndpointNameArgs(m)

			log.Debugf("    mount %s get", mountPoint)
			router.GET(mountPoint, func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
				endpoint.callMethod(w, r, ps, s, m, hasAccountID)
			})
		}
	})

	var handler http.Handler = router
	if len(allowedOrigins) >= 0 {
		c := cors.New(cors.Options{
			AllowedOrigins: allowedOrigins,
			AllowedMethods: []string{http.MethodPost, http.MethodGet},
			MaxAge:         600,
			AllowedHeaders: []string{"*"},
		})

		handler = c.Handler(handler)
	}

	endpoint.httpServer = &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return endpoint
}

func (ep *httpRestFulEndpoint) InternalSubscribe(svrName string, methodName string, strSubscribeUrl string) (*rpc.ID, error) {
	subscribeUri, err := url.Parse(strSubscribeUrl)
	if err != nil {
		log.Error(svrName + "." + methodName + ": subscribeUri " + strSubscribeUrl + " format error, " + err.Error())
		return nil, err
	}

	ep.subscribesMu.Lock()
	defer ep.subscribesMu.Unlock()

	if subscribeTarget, ok := ep.subscribes[*subscribeUri]; ok {
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

	subscribeTarget, ok := ep.subscribes[*subscribeUri]
	if !ok {
		subscribeTarget = make(subscribeData)
		ep.subscribes[*subscribeUri] = subscribeTarget
	}
	subscribeTarget[s.Name+"."+m.Name] = subscription.ID

	return &subscription.ID, nil
}

func (ep *httpRestFulEndpoint) Name() string {
	return EndpointName
}

func (ep *httpRestFulEndpoint) RunLoop() error {
	err := ep.httpServer.ListenAndServe()
	return err
}

func (ep *httpRestFulEndpoint) Shotdown(ctx context.Context) error {
	err := ep.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}

	ep.httpServer = nil

	return nil
}

func (ep *httpRestFulEndpoint) callMethod(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params,
	s *rpc.ServiceInfo, m *rpc.MethodInfo, hasAccountID bool) {

	ctx := buildContext(r.Context(), r)

	var subscribeUri *url.URL = nil
	if m.IsSubscribe {
		uri, err := readSubscribeUri(s, m, r)
		if err != nil {
			ep.sendInternalResponse(w, 400, err)
			return
		}
		subscribeUri, err = url.Parse(uri)
		if err != nil {
			log.Error(s.Name + "." + m.Name + ": subscribeUri " + uri + " format error, " + err.Error())
			ep.sendInternalResponse(w, 400, err)
			return
		}

		ep.subscribesMu.Lock()
		if subscribeTarget, ok := ep.subscribes[*subscribeUri]; ok {
			if id, ok := subscribeTarget[s.Name+"."+m.Name]; ok {
				log.Debugf(s.Name + "." + m.Name + ": subscribeUri " + uri + " already created")
				ep.subscribesMu.Unlock()
				ep.sendResponse(w, r, ps, &rpc.RpcResponse{Data: subscribeRes{ID: string(id)}})
				return
			}
		}
		defer func() {
			ep.subscribesMu.Unlock()
		}()

		if rpc.NotifierFromContext(ctx) == nil {
			ctx = rpc.AttachNotifier(ctx, ep.notifier)
		}
	}

	req, err := ep.readRequest(ctx, r, ps, s, m, hasAccountID)
	if err != nil {
		ep.sendInternalResponse(w, 400, err)
		return
	}

	res, err := ep.srv.ServeSingleRequest(ctx, req)
	if err != nil {
		ep.sendInternalResponse(w, 500, err)
		return
	}

	if m.IsSubscribe {
		if res.Data == nil {
			log.Error(s.Name + "." + m.Name + ": call success, no subscribe")
			ep.sendInternalResponse(w, 500, fmt.Errorf("internal error"))
			return
		}

		subscription, ok := res.Data.(*rpc.Subscription)
		if !ok || subscription == nil {
			log.Error(s.Name + "." + m.Name + ": call success, subscribe typer error")
			ep.sendInternalResponse(w, 500, fmt.Errorf("internal error"))
			return
		}
		defer ep.notifier.Activate(subscription.ID)

		res.Data = subscribeRes{ID: string(subscription.ID)}
		subscribeTarget, ok := ep.subscribes[*subscribeUri]
		if !ok {
			subscribeTarget = make(subscribeData)
			ep.subscribes[*subscribeUri] = subscribeTarget
		}
		subscribeTarget[s.Name+"."+m.Name] = subscription.ID
	}

	if res == nil {
		log.Error(s.Name + "." + m.Name + ": no response or error")
		w.WriteHeader(500)
		return
	}

	ep.sendResponse(w, r, ps, res)
}

func haveAccountIDField(structType reflect.Type) bool {
	//log.Info("process " + structType.Name())

	for j := 0; j < structType.NumField(); j++ {
		field := structType.Field(j)
		//log.Info("    " + field.Name + ": " + field.Type.Name())

		if field.Name == "AccountID" {
			return true
		}

		if field.Type.Kind() == reflect.Struct && field.Name == field.Type.Name() {
			if haveAccountIDField(field.Type) {
				return true
			}
		}
	}

	return false
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

func (ep *httpRestFulEndpoint) readRequest(ctx context.Context, r *http.Request, ps httprouter.Params,
	s *rpc.ServiceInfo, m *rpc.MethodInfo, hasAccountID bool) (*rpc.RpcRequest, error) {
	req := rpc.NewRpcRequest(s, m, ctx)

	var rError error = nil

	if rError == nil && r.Method == http.MethodPost {
		var incomingMsg *json.RawMessage = nil

		foreachStructType(req, func(v reflect.Value) {
			if incomingMsg == nil {
				body := io.LimitReader(r.Body, maxRequestContentLength)

				var incomingMsgBuf json.RawMessage
				if err := json.NewDecoder(body).Decode(&incomingMsgBuf); err != nil {
					log.Error(s.Name + "." + m.Name + ": decode msg fail")
					rError = &rpc.InvalidParamsError{
						Msg: err.Error(),
					}
					return
				}
				incomingMsg = &incomingMsgBuf
			}

			if err := json.Unmarshal(*incomingMsg, v.Interface()); err != nil {
				log.Error(s.Name + "." + m.Name + ": parseRequest: read " + v.Elem().Type().Name() + " from json fail" +
					", json=" + string(*incomingMsg) +
					", error=" + err.Error())
				rError = &rpc.InvalidParamsError{err.Error()}
				return
			}
		})
	}

	if rError == nil && hasAccountID {
		foreachStructType(req, func(v reflect.Value) {
			structValue := v.Elem()
			argType := structValue.Type()
			for j := 0; j < argType.NumField(); j++ {
				if argType.Field(j).Name == "account" {
					AccountID := ps.ByName("account")
					if AccountID == "" {
						log.Error(s.Name + "." + m.Name + ": no AccountID parament")
						rError = &rpc.InvalidParamsError{"no AccountID parament"}
						return
					} else {
						fieldValue := structValue.Field(j).Addr()
						if fieldValue.CanInterface() {
							if u, ok := fieldValue.Interface().(*common.AccountID); ok {
								if err := json.Unmarshal([]byte(AccountID), u); err != nil {
									log.Error(
										fmt.Sprintf("%s.%s: field %s.%s set AccountID %s fail, %s",
											s.Name, m.Name, argType.Name(), argType.Field(j).Name, AccountID, err.Error()))
									rError = &rpc.InvalidParamsError{fmt.Sprintf("AccountID %s format error", AccountID)}
									return
								} else {
									return
								}
							}
						}

						log.Error(
							fmt.Sprintf("%s.%s: field %s.%s type %s not support set AccountID",
								s.Name, m.Name, argType.Name(), argType.Field(j).Name, fieldValue.Type().Name()))
						rError = &rpc.InvalidParamsError{fmt.Sprintf("type %s not support set AccountID", fieldValue.Type().Name())}
						return
					}
				}
			}
		})
	}

	foreachStructType(req, func(v reflect.Value) {
		structValue := v.Elem()
		argType := structValue.Type()
		for j := 0; j < argType.NumField(); j++ {
			arg := ps.ByName(argType.Field(j).Name)
			if arg != "" {
				fieldValue := structValue.Field(j).Addr()
				if fieldValue.CanInterface() {
					u := fieldValue.Interface()
					json.Unmarshal([]byte(arg), u)
					return
				}

				log.Error(
					fmt.Sprintf("%s.%s: field %s.%s type %s not support set height",
						s.Name, m.Name, argType.Name(), argType.Field(j).Name, fieldValue.Type().Name()))
				rError = &rpc.InvalidParamsError{fmt.Sprintf("type %s not support set height", fieldValue.Type().Name())}
				return
			}
		}
	})

	if rError != nil {
		log.Error(rError.Error())
		return nil, rError
	} else {
		return req, nil
	}
}

func (ep *httpRestFulEndpoint) sendResponse(w http.ResponseWriter, r *http.Request, ps httprouter.Params, res *rpc.RpcResponse) {
	if res.Err != nil {
		resJson, err := json.Marshal(res.Data)
		if err != nil {
			ep.sendInternalResponse(w, 500, err)
		}

		w.Header().Set("content-type", contentType)
		w.Write(resJson)
	} else if res.Data != nil {
		resJson, err := json.Marshal(res.Data)
		if err != nil {
			ep.sendInternalResponse(w, 500, err)
			return
		}

		w.Header().Set("content-type", contentType)
		w.Write(resJson)
	}
}

func (ep *httpRestFulEndpoint) sendInternalResponse(w http.ResponseWriter, httpErrno int, err error) {
	w.WriteHeader(500)
	w.Write([]byte(err.Error()))
}

type subscribeReq struct {
	Uri string `json:"uri"`
}

type subscribeRes struct {
	ID string `json:"id"`
}

func readSubscribeUri(s *rpc.ServiceInfo, m *rpc.MethodInfo, r *http.Request) (string, error) {
	body := io.LimitReader(r.Body, maxRequestContentLength)
	var incomingMsg json.RawMessage
	if err := json.NewDecoder(body).Decode(&incomingMsg); err != nil {
		log.Error(s.Name + "." + m.Name + ": decode msg fail")
		return "", &rpc.InvalidParamsError{
			Msg: err.Error(),
		}
	}

	req := subscribeReq{}
	if err := json.Unmarshal(incomingMsg, &req); err != nil {
		log.Error(s.Name + "." + m.Name + ": decode uri fail")
		return "", &rpc.InvalidParamsError{
			Msg: err.Error(),
		}
	}

	if req.Uri == "" {
		log.Error(s.Name + "." + m.Name + ": no uri!")
		return "", &rpc.InvalidParamsError{
			Msg: "no uri",
		}
	}

	return req.Uri, nil
}

func buildContext(ctx context.Context, r *http.Request) context.Context {
	ctx = context.WithValue(ctx, "remote", r.RemoteAddr)
	ctx = context.WithValue(ctx, "scheme", r.Proto)
	ctx = context.WithValue(ctx, "local", r.Host)
	return ctx
}

func formatEndpointName(name string, module string, prefix string, hasAccountID bool) string {
	ret := []rune(name)
	ret = ret[len(prefix):]

	parts := []string{}

	if len(ret) > 0 {
		ret[0] = unicode.ToLower(ret[0])
	}

	for i := 0; i < len(ret); i++ {
		if unicode.IsUpper(ret[i]) {
			if i > 0 {
				parts = append(parts, string(ret[:i]))
			}

			ret = ret[i:]
			ret[0] = unicode.ToLower(ret[0])
			i = 0
		}
	}

	if len(ret) > 0 {
		parts = append(parts, string(ret))
	}

	if len(parts) > 0 && parts[0] == module {
		parts = parts[1:]
	}

	if len(parts) > 0 && hasAccountID && parts[0] == "account" {
		parts = parts[1:]
	}

	if len(parts) > 0 {
		return strings.Join(parts, "/") + "/"
	} else {
		return ""
	}
}

func formatEndpointNameArgsOneType(structType reflect.Type) string {
	//log.Info("process " + structType.Name())
	appendMsg := ""

	for j := 0; j < structType.NumField(); j++ {
		field := structType.Field(j)
		//log.Info("    " + field.Name + ": " + field.Type.Name())

		if field.Name == "AccountID" {
			continue
		}

		if field.Type.Kind() == reflect.Struct {
			if field.Name == field.Type.Name() {
				appendMsg += formatEndpointNameArgsOneType(field.Type)
			}
		} else {
			appendMsg += ":" + field.Name + "/"
		}
	}

	return appendMsg
}

func formatEndpointNameArgs(m *rpc.MethodInfo) string {
	appendMsg := ""
	for i := 0; i < len(m.ArgTypes); i++ {
		argType := m.ArgTypes[i]
		if argType.Kind() == reflect.Ptr {
			argType = argType.Elem()
		}

		if argType.Kind() != reflect.Struct {
			continue
		}

		appendMsg += formatEndpointNameArgsOneType(argType)
	}

	return appendMsg
}

func (ep *httpRestFulEndpoint) notifyEvent(id rpc.ID, data interface{}) error {
	var url *url.URL = nil

	ep.subscribesMu.Lock()
	for checkUrl, subscribeTarget := range ep.subscribes {
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
