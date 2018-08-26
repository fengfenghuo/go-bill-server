package config

import (
	// "bytes"
	"encoding/json"
	"fmt"
	"github.com/op/go-logging"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	// "reflect"
	"strings"
)

var (
	DefaultConfigName = "test-config.json"
)

var (
	log = logging.MustGetLogger("config")
)

var NodeConfig map[string]interface{}

func SetNodeConfigName(name string) {
	if len(name) == 0 {
		return
	}
	DefaultConfigName = name
}

func ReadNodeConfig() error {
	home := os.Getenv("HOME")
	if home == "" {
		if user, err := user.Current(); err == nil {
			home = user.HomeDir
		}
	}

	defaultDatasetDir := filepath.Join(home, ".sfox")

	file := filepath.Join(defaultDatasetDir, DefaultConfigName)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("read node config fail, error: " + err.Error())
	}

	err = json.Unmarshal(data, &NodeConfig)
	if err != nil {
		return fmt.Errorf("read node config fail, error: " + err.Error())
	}
	// for key, value := range NodeConfig {
	// 	log.Debugf("key: %s, value: %v", key, value)
	// }
	return nil
}

func GetParameFrom(config interface{}, parame string) (interface{}, bool) {
	if value, ok := config.(map[string]interface{}); ok {
		for k, v := range value {
			if strings.Compare(k, parame) == 0 {
				return v, true
			} else {
				result, ok := GetParameFrom(v, parame)
				if ok {
					return result, ok
				}
			}
		}
		// if k, ok := GetParameFrom(value, parame); ok {
		// 	return k, ok
		// }
	} else if value, ok := config.([]interface{}); ok {
		for _, result := range value {
			k, ok := GetParameFrom(result, parame)
			if ok {
				return k, ok
			}
		}
	}
	return nil, false
}

func ResolveNodeConfig(parame string) interface{} {
	result, ok := GetParameFrom(NodeConfig, parame)
	if ok {
		return result
	}
	return nil
}

/*
func GetParameFrom(config interface{}, parame string) (interface{}, bool) {
	re := reflect.TypeOf(config)
	v := reflect.ValueOf(config)
	switch re.Kind() {
	case reflect.Map:
		for _, key := range v.MapKeys() {
			log.Debugf("XXXXXXmap: %s", key.String())
			if strings.Compare(key.String(), parame) == 0 {
				return v.MapIndex(key), true
			} else {
				log.Error("XXXXXX%v", v.MapIndex(key))
				result, ok := GetParameFrom(v.MapIndex(key), parame)
				if ok {
					return result, ok
				}
			}
		}
		return nil, false
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			result, ok := GetParameFrom(v.Index(i), parame)
			if ok {
				return result, ok
			}
		}
		return nil, false
	case reflect.Interface:
		if v.IsNil() {
			return nil, false
		}
		return GetParameFrom(v.Elem(), parame)
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).CanInterface() {
				log.Debug("XXXXXX %s, XXXXXX %v", v.Type().Field(i).Name, v.Field(i).Interface())
				if strings.Compare(v.Type().Field(i).Name, parame) == 0 {
					return v.Field(i).Interface(), true
				}

				result, ok := GetParameFrom(v.Field(i).Interface(), parame)
				if ok {
					return result, ok
				}
			}
		}
		return nil, false
	default:
		log.Debugf("error type %s", re.Kind().String())
		return nil, false
	}
}

func ResolveNodeConfig(parame string) interface{} {
	result, ok := GetParameFrom(NodeConfig, parame)
	if ok {
		return result
	}
	return nil
}
*/
