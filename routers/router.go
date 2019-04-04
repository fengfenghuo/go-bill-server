package routers

import (
	"github.com/astaxie/beego"
	"github.com/bill-server/go-bill-server/controllers"
)

func init() {
	beego.Router("/account/register", &controllers.AccountController{}, "post:Register")
}
