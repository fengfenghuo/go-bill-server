package routers

import (
	"github.com/astaxie/beego"
	"github.com/bill-server/go-bill-server/node/beego/controllers"
)

func init() {
	beego.Router("/account/register", &controllers.AccountController{}, "post:Register")
	beego.Router("/tx/:account", &controllers.TxController{}, "post:Tx")
}
