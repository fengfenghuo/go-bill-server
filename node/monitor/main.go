package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-gomail/gomail"
	"github.com/op/go-logging"
)

var (
	log     = logging.MustGetLogger("monitor")
	relNode = "https://hkbt.shadowfox.cc" //"52.74.32.119"
)

func main() {
	if err := newCheckUpdateState(); err != nil {
		log.Error(err.Error())
	}

	for {
		time.Sleep(15 * 60 * time.Second)
		startCheckTimer()
	}
}

func startCheckTimer() {
	err1 := newCheckUpdateState()
	if err1 != nil {
		log.Error(err1.Error())
	}

	err2 := newCheckUpdateTx()
	if err2 != nil {
		log.Error(err2.Error())
	}

	err3 := newCheckTransferTx()
	if err3 != nil {
		log.Error(err3.Error())
	}

	err4 := newOnlineAccount()
	if err4 != nil {
		log.Error(err4.Error())
	}

	if err := newMail(err1.Error() + err2.Error() + err3.Error() + err4.Error()); err != nil {
		log.Error(err.Error())
	}
}

func newCheckUpdateState() error {
	req := `{"state": true}`
	req_new := bytes.NewBuffer([]byte(req))

	url := fmt.Sprintf("%s/account/0x94ee60bda85cf37aa4a2b163996274b7254eac13/online/state", relNode)
	request, err := http.NewRequest("POST", url, req_new)
	if err != nil {
		return err
	}

	if err = doRequestCheck(request); err != nil {
		return err
	}

	return nil
}

func newCheckUpdateTx() error {
	if err := newCheckUpdateState(); err != nil {
		return err
	}

	seq := NewTxSequence()
	req := `{"user": "0x94ee60bda85cf37aa4a2b163996274b7254eac13", "NU":1, "state": false, "sequence": %d }`
	req = fmt.Sprintf(req, seq)
	req_new := bytes.NewBuffer([]byte(req))

	url := fmt.Sprintf("%s/account/0xB06bf630aB574E2930D83DE416Ea65bce7D15073/online/tx/", relNode)
	request, err := http.NewRequest("POST", url, req_new)
	if err != nil {
		return err
	}

	if err = doRequestCheck(request); err != nil {
		return err
	}

	return nil
}

func newCheckTransferTx() error {
	seq := NewTxSequence()

	req := fmt.Sprintf(`{"from": "0xB06bf630aB574E2930D83DE416Ea65bce7D15073", "to": "0x94ee60bda85cf37aa4a2b163996274b7254eac13", "SFOX": 0.1, "sequence": %d }`, seq)
	req_new := bytes.NewBuffer([]byte(req))

	url := fmt.Sprintf("%s/tx/transfer/", relNode)
	request, err := http.NewRequest("POST", url, req_new)
	if err != nil {
		return err
	}

	if err = doRequestCheck(request); err != nil {
		return err
	}

	return nil
}

func newOnlineAccount() error {
	url := fmt.Sprintf("%s/account/0x94ee60bda85cf37aa4a2b163996274b7254eac13/online/", relNode)
	res, err := http.Get(url)
	if err != nil {
		if err := newMail(url + "\n" + err.Error()); err != nil {
			log.Error(err.Error())
		}
		return fmt.Errorf(err.Error())
	}

	if res.StatusCode != 200 {
		body, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf(url + "\n" + string(body[:]))
	} else if res.StatusCode == 200 {
		body, _ := ioutil.ReadAll(res.Body)
		log.Info(string(body[:]))
	}
	return nil
}

func doRequestCheck(request *http.Request) error {
	request.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		if err := newMail(err.Error()); err != nil {
			log.Error(err.Error())
		}
		return fmt.Errorf(err.Error())
	}

	if response.StatusCode != 200 {
		body, _ := ioutil.ReadAll(response.Body)
		return fmt.Errorf(string(body[:]))
	} else if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		log.Info(string(body[:]))
	}

	return nil
}

func newMail(msg string) error {
	m := gomail.NewMessage()

	m.SetAddressHeader("From", "3188953334@qq.com" /*"发件人地址"*/, "OnlineNode") // 发件人

	m.SetHeader("To",
		m.FormatAddress("guyikang@adups.com", "Tester")) // 收件人
	// m.SetHeader("Cc",
	// 	m.FormatAddress("xxxx@foxmail.com", "收件人")) //抄送

	m.SetHeader("Subject", "NodeError!!!") // 主题

	//m.SetBody("text/html",xxxxx ") // 可以放html..还有其他的
	m.SetBody("text/html", "Error Msg: \n"+msg) // 正文

	//m.Attach("我是附件") //添加附件

	d := gomail.NewPlainDialer("smtp.qq.com", 465, "3188953334@qq.com", "bnieluisodyidcef") // 发送邮件服务器、端口、发件人账号、发件人密码
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	log.Info("done.发送成功")
	return nil
}

func NewTxSequence() int64 {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Int63()
}

// req := `{"height": 955 }`
// req_new := bytes.NewBuffer([]byte(req))
// fmt.Println(req_new)
// request, _ := http.NewRequest("POST", "http://127.0.0.1:20000/chain/blocks/", req_new)

// request, _ := http.NewRequest("POST", "http://127.0.0.1:20000/account/0xb10f225476c74ccb75fd4d150bbee0bab8ab2903/online/using/", nil)

// request, _ := http.NewRequest("POST", "http://127.0.0.1:20000/account/0xb10f225476c74ccb75fd4d150bbee0bab8ab2903/online/provide/", nil)

// req := "{\"uri\":\"http://180.97.69.197:32201/climb/bc/txDetail/\"}"
// req_new := bytes.NewBuffer([]byte(req))
// request, _ := http.NewRequest("POST", "http://180.97.69.197:20000/tx/subscribe/", req_new)
