package routers_test

import (
	"testing"

	"github.com/bill-server/go-bill-server/gin/routers"
)

func TestRouter(t *testing.T) {
	router := routers.NewRouter()
	router.StartRun()
}
