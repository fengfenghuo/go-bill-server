package routers

import (
	"strconv"

	"github.com/bill-server/go-bill-server/gin/conf"
	"github.com/bill-server/go-bill-server/gin/controllers"
	"github.com/gin-gonic/gin"
)

type Router struct {
	router *gin.Engine
}

func NewRouter() *Router {
	router := gin.Default()
	accountRouter := router.Group("/account")
	accountRouter.POST("register", controllers.AccountRegister)

	txRouter := router.Group("/tx")
	txRouter.POST("/:account", controllers.CreateTx)

	return &Router{router: router}
}

func (router *Router) StartRun() {
	port := conf.AppConfig.DefaultInt("httpport", 28080)
	portStr := strconv.Itoa(port)
	router.router.Run(":" + portStr)
}
