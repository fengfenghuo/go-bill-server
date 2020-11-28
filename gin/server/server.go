package server

import (
	"encoding/json"
	"net/http"
	"reflect"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
)

type GinServer struct{}

func NewGinServer() IServer {
	return &GinServer{}
}

func (s *GinServer) Run() {
	// if !config.C.Debug {
	// 	gin.SetMode(gin.ReleaseMode)
	// }

	r := s.router()

	if err := r.Run(""); err != nil {
		panic(err)
	}
}

func (s *GinServer) router() *gin.Engine {
	r := gin.Default()
	// alarm := api.NewAlarmHandle()
	// r.POST("/api/gateway/alarm", s.handle(alarm.SendAlarm))
	// r.GET("/api/gateway/health", s.handle(alarm.CheckHealth))
	return r
}

func (s *GinServer) handle(rcvr interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		mtype := reflect.TypeOf(rcvr)
		var args []reflect.Value
		if mtype.NumIn() > 0 {
			argType := mtype.In(0)
			argValue := reflect.New(argType.Elem())

			msg := argValue.Interface()
			if err := ctx.ShouldBind(msg); err != nil {
				logs.Error("GinServer: Handle: Bind err: %v", err)
				s.ResponseErr(ctx, err)
				return
			}

			args = append(args, reflect.ValueOf(msg))

			by, _ := json.Marshal(msg)
			logs.Debug("request: %s", string(by))
		}

		reply := reflect.ValueOf(rcvr).Call(args)

		errPos := 0
		var data interface{}

		if len(reply) > 1 {
			errPos = len(reply) - 1
		}

		if !reply[errPos].IsNil() {
			err := reply[errPos].Interface().(error)
			logs.Error("GinServer: Handle: Request: %v", err)
			s.ResponseErr(ctx, err)
			return
		}

		if errPos != 0 {
			data = reply[0].Interface()
		}

		s.ResponseOK(ctx, data)
		return
	}
}

// ResponseErr 错误回复
func (s *GinServer) ResponseErr(ctx *gin.Context, err error) {
	if config.C.Debug {
		logs.Debug("%s: response err: %s", time.Now().Format(time.RFC3339), err.Error())
	}
	ctx.JSON(http.StatusOK, api.MessageResponse{
		Code: http.StatusBadRequest,
		Msg:  err.Error(),
	})
}

// ResponseOK 正常回复
func (s *GinServer) ResponseOK(ctx *gin.Context, data interface{}) {
	if config.C.Debug {
		logs.Debug("%s: response ok", time.Now().Format(time.RFC3339))
	}
	ctx.JSON(http.StatusOK, api.MessageResponse{
		Code: 0,
		Msg:  "ok",
		Data: data,
	})
}
