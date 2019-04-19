package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("WS-tool")
	// addr   = flag.String("addr", "localhost:20020", "http service adress")
	addr   = flag.String("addr", "180.97.69.197:20020", "http service adress")
	Method = map[string]string{
		"o.state":    "online.updateState",
		"o.tx":       "online.updateTx",
		"o.account":  "online.account",
		"t.tx":       "tx.tx",
		"t.transfer": "tx.createTransfer",
		"c.block":    "chain.blocks",
	}
)

const (
	u1 = "0x5e113de49b2324d032632b4b6170cb432aaf1e42"
	u2 = "0xb10f225476c74ccb75fd4d150bbee0bab8ab2903"
)

func main() {
	u := url.URL{Scheme: "ws", Host: *addr, Path: ""}
	var dialer *websocket.Dialer

	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	go waitSendMessage(conn)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}

		fmt.Printf("received: %s\n", message)
	}
}

func waitSendMessage(conn *websocket.Conn) {
	c := make(chan string)
	id := 1000

	go func() {
		defer close(c)

		for {
			select {
			case cmd := <-c:
				method := strings.Split(cmd, " ")

				if m, ok := Method[method[0]]; ok {
					method[0] = m
				} else {
					log.Error("error method %s not find!", method[0])
					continue
				}

				params, err := dealWithParams(method)
				if err != nil {
					log.Error(err.Error())
					continue
				}

				req, err := newRequest(method[0], id, params)
				if err != nil {
					log.Error(err.Error())
					continue
				}

				if err := sendRequest(conn, req); err != nil {
					log.Error(err.Error())
					continue
				}

				id++
			}
		}

	}()

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		c <- line
	}
}

func dealWithParams(method []string) ([]interface{}, error) {
	params := make([]interface{}, 0, 10)
	for i := 1; i < len(method); i++ {
		if strings.Compare(method[i], "u1") == 0 {
			method[i] = u1
		} else if strings.Compare(method[i], "u2") == 0 {
			method[i] = u2
		}

		if strings.Compare(method[i], "seq") == 0 {
			params = append(params, newTxSequence())
			continue
		}

		if strings.Compare(method[i], "true") == 0 {
			params = append(params, true)
			continue
		}

		if strings.Compare(method[i], "false") == 0 {
			params = append(params, false)
			continue
		}

		param, err := strconv.Atoi(method[i])
		if err == nil {
			params = append(params, param)
		} else if strings.Contains(method[i], ".") {
			param, err := strconv.ParseFloat(method[i], 32)
			if err != nil {
				return nil, err
			}
			params = append(params, param)
		} else {
			params = append(params, method[i])
		}
	}
	return params, nil
}

func newRequest(method string, id int, params []interface{}) (string, error) {
	var req string
	paramsStr, err := json.Marshal(params)
	if err != nil {
		return req, err
	}

	req = `{"method": "%s", "id": %d, "params":%v}`
	req = fmt.Sprintf(req, method, id, string(paramsStr[:]))
	fmt.Println(req)
	return req, nil
}

func sendRequest(conn *websocket.Conn, req string) error {
	err := conn.WriteMessage(websocket.TextMessage, []byte(req))
	if err != nil {
		return err
	}
	return nil
}

func newTxSequence() int64 {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Int63()
}

// req := `{"method": "order.cancel", "id": %d, "params":["0x11121244134", 2]}`
// 			req = fmt.Sprintf(req, id)
// 			fmt.Println(req)
// 			conn.WriteMessage(websocket.TextMessage, []byte(req))
// 			id++
// 			time.Sleep(time.Second * 30)
// 			break
