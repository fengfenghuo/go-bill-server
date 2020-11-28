package routers

import (
	"net"

	"github.com/astaxie/beego/logs"
	"github.com/bill-server/go-bill-server/gin/account"
	"github.com/bill-server/go-bill-server/gin/conf"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type Router struct {
	router *gin.Engine
	grpc   net.Listener
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) StartRun() {
	runMode := conf.AppConfig.String("runmode")
	r.router = gin.New()

	if runMode != "" {
		var mode string
		if runMode == "dev" || runMode == "debug" {
			mode = gin.DebugMode
		} else if runMode == "test" {
			mode = gin.TestMode
		} else if runMode == "prod" || runMode == "release" {
			mode = gin.ReleaseMode
		} else {
			mode = gin.DebugMode
		}
		gin.SetMode(mode)
	}

	// accountRouter := r.router.Group("/account")
	// accountRouter.POST("register", controllers.AccountRegister)

	port := conf.AppConfig.DefaultString("httpport", "38080")
	r.router.Run(":" + port)
}

func (r *Router) StartRunGRPC() {
	srv := grpc.NewServer()
	account.RegisterAccountServiceServer(srv, &account.AccountService{})

	port := conf.AppConfig.DefaultString("httpport", "38080")
	var err error
	r.grpc, err = net.Listen("tcp", ":"+port)
	if err != nil {
		logs.Error("StartRunGRPC err: %v", err)
	}
	srv.Serve(r.grpc)
}
