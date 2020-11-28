package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/bill-server/go-bill-server/gin/routers"
)

func main() {
	go func() {
		log.Printf(http.ListenAndServe(":6060", nil).Error())
	}()
	router := routers.NewRouter()
	router.StartRunGRPC()
}
