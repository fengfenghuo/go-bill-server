package models

import (
	"strings"
	"time"

	"github.com/bill-server/go-bill-server/db"
	"github.com/go-sfox-lib/sfox/log"
)

var log = logger.NewLogInstance("models", "debug", "", "")

func init() {
	db, err := db.GetDBInstance()
	if err != nil {
		log.Error(err.Error())
	}
	db.RegisterTable(
		new(Account),
		new(Tx),
	)
}

func TimeFormat(cur_time time.Time) string {
	str := cur_time.Format("2006-01-02 15:04:05Z07:00")
	if strings.HasSuffix(str, "Z") {
		str = strings.Replace(str, "Z", "+00:00", -1)
	}
	return str
}
