package main

import (
	"github.com/astaxie/beego"
	_ "github.com/bill-server/go-bill-server/node/beego/routers"
)

func main() {
	beego.Run()
}
