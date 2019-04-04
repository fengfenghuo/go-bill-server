package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/bill-server/go-bill-server/models"
)

type AccountController struct {
	beego.Controller
}

func (c *AccountController) Register() {
	log.Debug("xxxxxxRegister")
	defer c.ServeJSON()

	var ac models.Account
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &ac)
	if err != nil {
		log.Error("Register: Unmarshal error: " + err.Error())
		c.Data["json"] = NewErrorMsg("register decode error", ErrorCodeSystemDecodeError)
		return
	}

	acc := models.NewAccount(ac.AccountID, ac.Password, ac.Email)
	err = acc.Register()
	if err != nil {
		log.Error("Register: Register error: " + err.Error())
		c.Data["json"] = NewErrorMsg("register instert error", ErrorCodeAccountInsertError)
		return
	}
	c.Data["json"] = NewErrorMsg("", ErrorNo)
}
