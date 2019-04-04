package main

import (
	_ "github.com/bill-server/go-bill-server/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
