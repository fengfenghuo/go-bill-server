package routers

import (
	"strconv"

	"github.com/bill-server/go-bill-server/gin/account"
	"github.com/bill-server/go-bill-server/gin/conf"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type Router struct {
	router *gin.Engine
}

func NewRouter() *Router {
	runMode := conf.AppConfig.String("runmode")
	router := gin.New()

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

	// accountRouter := router.Group("/account")
	// accountRouter.POST("register", controllers.AccountRegister)

	// txRouter := router.Group("/tx")
	// txRouter.POST("/:account", controllers.CreateTx)

	return &Router{router: router}
}

func (router *Router) StartRun() {
	port := conf.AppConfig.DefaultInt("httpport", 38080)
	portStr := strconv.Itoa(port)
	router.router.Run(":" + portStr)
}

func (router *Router) StartRunGRPC() {
	srv := grpc.NewServer()
	account.RegisterAccountServiceServer(srv, &account.AccountService{})
}
