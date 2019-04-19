package routers

import (
	"github.com/bill-server/go-bill-server-gin/controllers"
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
	router.router.Run("28080")
}
