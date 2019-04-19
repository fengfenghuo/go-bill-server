package main

import (
	"bytes"
	"fmt"
	"github.com/op/go-logging"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

var (
	log       = logging.MustGetLogger("monitor")
	localNode = "http://localhost:20000"
	// localNode = "https://bt.shadowfox.cc"
	// localNode = "http://180.97.69.197:20000"
)

type tx struct {
	from string
	to   string
	SFox float64
	NU   float64
}

const (
	from   = "0xB06bf630aB574E2930D83DE416Ea65bce7D15073"
	seeder = "0x94ee60bda85cf37aa4a2b163996274b7254eac13"
	test1  = "0x5e113de49b2324d032632b4b6170cb432aaf1e42"
	test2  = "0xb10f225476c74ccb75fd4d150bbee0bab8ab2903"
)

func main() {
	// update state
	stateSeq := NewTxSequence()
	if err := newCheckUpdateState(seeder, stateSeq, true); err != nil {
		log.Error(err.Error())
	}

	// update tx nu
	txSeq := NewTxSequence()
	if err := newCheckUpdateTx(test1, seeder, 0.1, txSeq, true); err != nil {
		log.Error(err.Error())
	}

	time.Sleep(1 * time.Second)
	if err := newCheckUpdateTx(test1, seeder, 0.2, txSeq, true); err != nil {
		log.Error(err.Error())
	}

	time.Sleep(1 * time.Second)
	if err := newCheckUpdateTx(test1, seeder, 0.3, NewTxSequence(), true); err != nil {
		log.Error(err.Error())
	}

	// // update transfer
	// if err := newCheckTransferTx(from, seeder, 1, 1); err != nil {
	// 	log.Error(err.Error())
	// }

	// // query block
	// if err := getBlockQuery(222); err != nil {
	// 	log.Error(err.Error())
	// }
}

func newCheckUpdateState(address string, seq int64, isOnline bool) error {
	req := fmt.Sprintf(`{"state": %v, "sequence": %d }`, isOnline, seq)
	req_new := bytes.NewBuffer([]byte(req))
	url := fmt.Sprintf("%s/account/%s/online/state", localNode, address)
	request, err := http.NewRequest("POST", url, req_new)
	if err != nil {
		return err
	}

	if err = doRequestCheck(request); err != nil {
		return err
	}

	return nil
}

func newCheckUpdateTx(address string, user string, nu float64, seq int64, isOnline bool) error {
	req := fmt.Sprintf(`{"user": "%s", "NU": %f, "state": %v, "sequence": %d }`, user, nu, isOnline, seq)
	req_new := bytes.NewBuffer([]byte(req))
	url := fmt.Sprintf("%s/account/%s/online/tx", localNode, address)
	request, err := http.NewRequest("POST", url, req_new)
	if err != nil {
		return err
	}

	if err = doRequestCheck(request); err != nil {
		return err
	}

	return nil
}

func newCheckTransferTx(from string, to string, sfox float64, nu float64) error {
	seq := NewTxSequence()

	req := fmt.Sprintf(`{"from": "%s", "to": "%s", "SFOX": %f, "NU": %f, "sequence": %d }`, from, to, sfox, nu, seq)
	req_new := bytes.NewBuffer([]byte(req))
	fmt.Println(req)

	url := fmt.Sprintf("%s/tx/transfer", localNode)
	fmt.Println(url)
	request, err := http.NewRequest("POST", url, req_new)
	if err != nil {
		return err
	}

	if err = doRequestCheck(request); err != nil {
		return err
	}

	return nil
}

func getBlockQuery(height int64) error {
	url := fmt.Sprintf("%s/chain/blocks/%d", localNode, height)
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	if res.StatusCode == 200 {
		body, _ := ioutil.ReadAll(res.Body)
		log.Info(string(body[:]))
	}
	return nil
}

func newOnlineAccount(address string) error {
	url := fmt.Sprintf("%s/account/%s/online", localNode, address)
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	if res.StatusCode == 200 {
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
		return fmt.Errorf(err.Error())
	}

	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		log.Info(string(body[:]))
	}

	return nil
}

func NewTxSequence() int64 {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Int63()
}

// txes := []tx{
// 	{from: "0xB06BF630AB574E2930D83DE416EA65BCE7D15073", to: "0xa677995a0ea616e0dcff324a83e3e0cf65ba74f5", SFox: 0.1, NU: 500},
// 	{from: "0xB06BF630AB574E2930D83DE416EA65BCE7D15073", to: "0xa677995a0ea616e0dcff324a83e3e0cf65ba74f5", SFox: 5, NU: 500},
// 	{from: "0xB06BF630AB574E2930D83DE416EA65BCE7D15073", to: "0x99e4dbl5513d5df68e62e863080ae24f2c4elc6b", SFox: 0.1, NU: 500},
// 	{from: "0xB06BF630AB574E2930D83DE416EA65BCE7D15073", to: "0x03b39691742009361030f7d76463c809c7a7671f", SFox: 0.1, NU: 500},
// 	{from: "0xB06BF630AB574E2930D83DE416EA65BCE7D15073", to: "0xa828efbc839eaedb7f9d9fb7952c2d4f407c767e", SFox: 0.1, NU: 500},
// 	{from: "0xB06BF630AB574E2930D83DE416EA65BCE7D15073", to: "0x2cf57f8ef062eb4c2fbffld364626df5081eab3f", SFox: 0.1, NU: 500},
// }

/*
req := `{"height": 955 }`
req_new := bytes.NewBuffer([]byte(req))
fmt.Println(req_new)
request, _ := http.NewRequest("POST", "http://127.0.0.1:20000/chain/blocks/", req_new)

request, _ := http.NewRequest("POST", "http://127.0.0.1:20000/account/0xb10f225476c74ccb75fd4d150bbee0bab8ab2903/online/using/", nil)

request, _ := http.NewRequest("POST", "http://127.0.0.1:20000/account/0xb10f225476c74ccb75fd4d150bbee0bab8ab2903/online/provide/", nil)

req := "{\"uri\":\"http://180.97.69.197:32201/climb/bc/txDetail/\"}"
req_new := bytes.NewBuffer([]byte(req))
request, _ := http.NewRequest("POST", "http://180.97.69.197:20000/tx/subscribe/", req_new)

func main() {
	req := `{"state": true}`
	req_new := bytes.NewBuffer([]byte(req))
	request, _ := http.NewRequest("POST", "http://127.0.0.1:20000/account/0x5e113de49b2324d032632b4b6170cb432aaf1e42/online/state", req_new)
	request, _ := http.NewRequest("POST", "http://180.97.69.197:20000/account/0xb10f225476c74ccb75fd4d150bbee0bab8ab2910/online/state", req_new)
	request, _ := http.NewRequest("POST", "http://52.74.32.119:20000/account/0xb10f225476c74ccb75fd4d150bbee0bab8ab2910/online/state", req_new)
	request, _ := http.NewRequest("POST", "https://bt.shadowfox.cc/account/0xb10f225476c74ccb75fd4d150bbee0bab8ab0000/online/state", req_new)

	req := `{"user": "0x94ee60bda85cf37aa4a2b163996274b7254eac13", "NU":1, "state": false, "sequence": 11233414131 }`
	req_new := bytes.NewBuffer([]byte(req))
	request, _ := http.NewRequest("POST", "http://localhost:20000/account/0xb10f225476c74ccb75fd4d150bbee0bab8ab2902/online/tx/", req_new)

	req := `{"height": 955 }`
	req_new := bytes.NewBuffer([]byte(req))
	fmt.Println(req_new)
	request, _ := http.NewRequest("POST", "http://127.0.0.1:20000/chain/blocks/", req_new)

	req := `{"from": "0xB06bf630aB574E2930D83DE416Ea65bce7D15010", "to": "0xb10f225476c74ccb75fd4d150bbee0bab8ab2902", "SFOX": 0, "sequence": 231231 }`
	req_new := bytes.NewBuffer([]byte(req))
	request, _ := http.NewRequest("POST", "http://127.0.0.1:20000/tx/transfer/", req_new)
	request, _ := http.NewRequest("POST", "https://bt.shadowfox.cc/tx/transfer", req_new)

	request, _ := http.NewRequest("GET", "http://127.0.0.1:20000/account/0xb10f225476c74ccb75fd4d150bbee0bab8ab2910/online", nil)

	request, _ := http.NewRequest("POST", "http://127.0.0.1:20000/account/0xb10f225476c74ccb75fd4d150bbee0bab8ab2903/online/using/", nil)

	request, _ := http.NewRequest("POST", "http://127.0.0.1:20000/account/0xb10f225476c74ccb75fd4d150bbee0bab8ab2903/online/provide/", nil)

	req := "{\"uri\":\"http://180.97.69.197:32201/climb/bc/txDetail/\"}"
	req_new := bytes.NewBuffer([]byte(req))
	request, _ := http.NewRequest("POST", "http://180.97.69.197:20000/tx/subscribe/", req_new)

	request.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	response, _ := client.Do(request)
	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(body))
	}
}

https://bt.shadowfox.cc/account/0x94ee60bda85cf37aa4a2b163996274b7254eac13/online/
*/
