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
	if err := registerPost(); err != nil {
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
