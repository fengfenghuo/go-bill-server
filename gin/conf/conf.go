package conf

import (
	"github.com/go-sfox-lib/sfox/config"
	_ "github.com/go-sfox-lib/sfox/config/yaml"
)

// AppConfig ...
var AppConfig config.Configer

func init() {
	AppConfig, _ = config.NewConfig("yaml", "./conf/conf.yaml")
}
