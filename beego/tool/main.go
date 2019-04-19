package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	url = "http://127.0.0.1:38080/"
)

func main() {
	// if err := registerPost(); err != nil {
	// 	log.Printf("%v", err)
	// }

	if err := txPost(); err != nil {
		log.Printf("%v", err)
	}
}

func registerPost() error {
	data := `{"account_id": "guyikang@adups.com", "password":"5656688974", "email":"guyikang@adups.com"}`
	req_new := bytes.NewBuffer([]byte(data))
	request, err := http.NewRequest("POST", url+"account/register", req_new)
	if err != nil {
		return err
	}

	if err = doRequestCheck(request); err != nil {
		return err
	}
	return nil
}

// type Tx struct {
// 	ID         int64     `json:"id" orm:"column(id);auto"`
// 	AccountID  string    `json:"account_id" orm:"column(account_id)"`
// 	Amount     int64     `json:"amount" orm:"column(amount)"`
// 	TxType     TxType    `json:"tx_type" orm:"column(tx_type)"`
// 	Reason     TxReason  `json:"reason" orm:"column(reason)"`
// 	Remarks    string    `json:"remarks" orm:"column(remarks);null"`
// 	CreateTime time.Time `json:"create_time" orm:"auto_now_add;column(create_time)"`
// }

func txPost() error {
	data := `{"amount": 100, "tx_type": 1, "reason": 2, "remarks": "haha"}`
	req_new := bytes.NewBuffer([]byte(data))
	request, err := http.NewRequest("POST", url+"tx/guyikang@adups.com", req_new)
	if err != nil {
		return err
	}

	if err = doRequestCheck(request); err != nil {
		return err
	}
	return nil
}

func doRequestCheck(request *http.Request) error {
	request.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		log.Println(string(body[:]))
	}

	return nil
}
