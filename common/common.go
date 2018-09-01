package common

import (
	"strings"
	"time"
)

type AccountID string
type TxID uint64

type Account struct {
	Address AccountID
}

type Tx struct {
	ID      TxID
	Address AccountID
}

func TimeFormat(cur_time time.Time) string {
	str := cur_time.Format("2006-01-02 15:04:05Z07:00")
	if strings.HasSuffix(str, "Z") {
		str = strings.Replace(str, "Z", "+00:00", -1)
	}
	return str
}
