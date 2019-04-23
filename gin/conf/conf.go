package conf

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-sfox-lib/sfox/config"
	_ "github.com/go-sfox-lib/sfox/config/yaml"
)

// AppConfig ...
var AppConfig config.Configer

func init() {
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(workPath)
	appConfigPath := filepath.Join(workPath, "conf", "conf.yaml")
	AppConfig, _ = config.NewConfig("yaml", appConfigPath)
	fmt.Println(AppConfig.Int("httpport"))
}
