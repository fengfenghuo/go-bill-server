package controllers

import (
	"github.com/gin-gonic/gin"
	// "encoding/json"
	// "github.com/bill-server/go-bill-server-gin/models"
)

type Account struct {
}

func (a *Account) AccountRegister(c *gin.Context) {
	// defer c.ServeJSON()

	// var ac models.Account
	// err := json.Unmarshal(c.Ctx.Input.RequestBody, &ac)
	// if err != nil {
	// 	log.Error("Register: Unmarshal error: " + err.Error())
	// 	c.Data["json"] = NewErrorMsg("register decode error", ErrorCodeSystemDecodeError)
	// 	return
	// }

	// acc := models.NewAccount(ac.AccountID, ac.Password, ac.Email)
	// err = acc.Register()
	// if err != nil {
	// 	log.Error("Register: Register error: " + err.Error())
	// 	c.Data["json"] = NewErrorMsg("register instert error", ErrorCodeAccountInsertError)
	// 	return
	// }
	// c.Data["json"] = NewErrorMsg("", ErrorNo)
}
