package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/go-sfox-lib/sfox/log"

	"github.com/astaxie/beego"
	"github.com/bill-server/go-bill-server/models"
)

var log = logger.NewLogInstance("controllers", "debug", "", "")

type AccountController struct {
	beego.Controller
}

func (c *AccountController) Register() {
	log.Debug("xxxxxx")
	defer c.ServeJSON()

	var ac models.Account
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &ac)
	if err != nil {
		log.Error("Register: Unmarshal error: " + err.Error())
		c.Data["json"] = NewErrorMsg("register decode", ErrorCodeSystemDecodeError)
		return
	}
	fmt.Println("xxxxx%v", ac)
	acc := models.NewAccount(ac.AccountID, ac.Password, ac.Email)
	err = acc.Register()
	if err != nil {
		log.Error("Register: Register error: " + err.Error())
		c.Data["json"] = NewErrorMsg("register instert", ErrorCodeAccountInsertError)
		return
	}
	c.Data["json"] = "register success"
}
