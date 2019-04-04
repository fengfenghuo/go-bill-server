package controllers

import (
	"github.com/astaxie/beego"
	_ "github.com/bill-server/go-bill-server/models"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}
