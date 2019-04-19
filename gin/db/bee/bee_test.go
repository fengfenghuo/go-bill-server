package beeorm_test

import (
	"fmt"
	"github.com/bill-server/go-bill-server/node/db/bee"
	"testing"
)

func TestInitDB(t *testing.T) {
	_, err := beeorm.RegisterDB()
	if err != nil {
		fmt.Printf("db register error: %v", err)
	}
}
