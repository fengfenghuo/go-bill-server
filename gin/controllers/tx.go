package controllers

import (
	"github.com/gin-gonic/gin"
	// "encoding/json"
	// "github.com/bill-server/go-bill-server-gin/models"
)

func CreateTx(c *gin.Context) {
	// defer tx.ServeJSON()

	// var t models.Tx
	// err := json.Unmarshal(tx.Ctx.Input.RequestBody, &t)
	// if err != nil {
	// 	log.Error("Tx: Unmarshal error: " + err.Error())
	// 	tx.Data["json"] = NewErrorMsg("tx decode error", ErrorCodeSystemDecodeError)
	// 	return
	// }

	// t.AccountID = tx.GetString(":account")
	// newT := models.NewTx(t.AccountID, t.Amount, t.TxType, t.Reason, t.Remarks)
	// err = newT.CreateTx()
	// if err != nil {
	// 	log.Error("Tx: CreateTx error: " + err.Error())
	// 	tx.Data["json"] = NewErrorMsg("tx create error", ErrorCodeTxInsertError)
	// 	return
	// }

	// tx.Data["json"] = NewErrorMsg("", ErrorNo)
}
