package db_test

import (
	"fmt"
	"github.com/bill-server/go-bill-server/db"
	"testing"
)

func TestInitDB(t *testing.T) {
	err := db.RegisterDB()
	if err != nil {
		fmt.Printf("db register error: %v", err)
	}
}
