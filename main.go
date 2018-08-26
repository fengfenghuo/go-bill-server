package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"strings"

	"github.com/bill-server/go-bill-server/account"
	"github.com/bill-server/go-bill-server/config"
	"github.com/bill-server/go-bill-server/node"
	"github.com/bill-server/go-bill-server/redis"
	"github.com/bill-server/go-bill-server/rpc"
	"github.com/bill-server/go-bill-server/rpc/restful"
	"github.com/bill-server/go-bill-server/rpc/ws-jsonrpc"
	"github.com/bill-server/go-bill-server/sfoxdb"
	// "github.com/bill-server/go-bill-server/sfoxdb/cluster"
	"github.com/bill-server/go-bill-server/tx"
	"github.com/op/go-logging"
	"gopkg.in/urfave/cli.v1"
)

var (
	log        = logging.MustGetLogger("BillServer")
	app        = cli.NewApp()
	configName = ""
)

func main() {
	app.Name = filepath.Base(os.Args[0])
	app.Author = ""
	app.Email = ""
	app.Version = ""
	app.Action = RunNode
	app.HideVersion = true // we have a command to print the version
	app.Copyright = "Copyright 2013-2018 The SecuritFox Authors"
	app.Commands = []cli.Command{}

	if len(os.Args) > 1 {
		configName = os.Args[1]
	}

	// log.Debugf("XXXXXX config name: %s", configName)
	// config.SetNodeConfigName(configName)
	// config.ReadNodeConfig()

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func RunNode(ctx *cli.Context) error {
	node, err := makeFullNode(ctx)
	if err != nil {
		return err
	}

	err = node.Start()
	if err != nil {
		log.Error("RunNode: start fail, " + err.Error())
		return err
	}

	restfulEndpoint := restful.NewEndpoint(nil, node.RpcServer(), ":20020")
	node.AddRpcEndpoint(restfulEndpoint)
	// newSubscrib(fullNode, restfulEndpoint)

	wsEndpoint := jsonrpc.NewEndpoint(node.RpcServer(), ":20022")
	node.AddRpcEndpoint(wsEndpoint)

	// if fullNode := config.ResolveNodeConfig("RESTFul"); fullNode != nil {
	// 	if err := startListen(node, fullNode, "RESTFul"); err != nil {
	// 		log.Error(err.Error())
	// 		return err
	// 	}
	// }

	// if fullNode := config.ResolveNodeConfig("WS"); fullNode != nil {
	// 	if err := startListen(node, fullNode, "WS"); err != nil {
	// 		log.Error(err.Error())
	// 		return err
	// 	}
	// }

	go func() {
		log.Debugf(http.ListenAndServe(":6060", nil).Error())
	}()

	node.Wait()

	return nil
}

func startListen(node *node.Node, fullNode interface{}, rpcType string) error {
	// listen, ok := config.GetParameFrom(fullNode, "listen")
	// if !ok {
	// 	err := fmt.Errorf("RunNode start fail, listen not exist!")
	// 	return err
	// }

	if strings.Compare(rpcType, "RESTFul") == 0 {
		restfulEndpoint := restful.NewEndpoint(nil, node.RpcServer(), ":20020")
		node.AddRpcEndpoint(restfulEndpoint)
		newSubscrib(fullNode, restfulEndpoint)
	} else if strings.Compare(rpcType, "WS") == 0 {
		wsEndpoint := jsonrpc.NewEndpoint(node.RpcServer(), ":20022")
		node.AddRpcEndpoint(wsEndpoint)
		// newSubscrib(fullNode, wsEndpoint)
	}

	return nil
}

func makeFullNode(ctx *cli.Context) (*node.Node, error) {
	db, err := newDB()
	if err != nil {
		return nil, err
	}

	redis, err := newRedis()
	if err != nil {
		return nil, err
	}
	n, err := node.NewNode(db, redis)
	if err != nil {
		return nil, err
	}

	n.Register(account.NewService)
	n.Register(tx.NewService)

	return n, nil
}

func newDB() (sfoxdb.SFoxDb, error) {
	// if clusterdb := config.ResolveNodeConfig("cluster"); clusterdb != nil {
	// 	result, ok := config.GetParameFrom(clusterdb, "db")
	// 	if !ok {
	// 		return nil, fmt.Errorf("cluster db not exist, fail")
	// 	}

	// 	url, ok := result.(string)
	// 	if !ok {
	// 		return nil, fmt.Errorf("cluster db url not string")
	// 	}

	// 	db, err := cluster.NewSFoxCluster(url)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return db, nil
	// }
	return nil, nil
}

func newSubscrib(node interface{}, endpoint rpc.Endpoint) error {
	subscrib, ok := config.GetParameFrom(node, "autoSubscrib")
	if !ok {
		return nil
	}

	su, ok := subscrib.([]interface{})
	if !ok {
		return fmt.Errorf("newSubscrib config fmt error")
	}
	for _, temp := range su {
		scope, ok := config.GetParameFrom(temp, "scope")
		if !ok {
			continue
		}

		method, ok := config.GetParameFrom(temp, "method")
		if !ok {
			continue
		}

		url, ok := config.GetParameFrom(temp, "url")
		if !ok {
			continue
		}

		endpoint.InternalSubscribe(scope.(string), method.(string), url.(string))
	}
	return nil
}

func newRedis() (redis.DataRedis, error) {
	// redis, err := redis.NewDataRedis(":6379")
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}
