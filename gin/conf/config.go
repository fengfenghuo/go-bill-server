package conf

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type ConfigEngine struct {
	Data map[string]interface{}
	Path string
}

func init() {

}

func NewConfigEngine(path string) *ConfigEngine {
	return &ConfigEngine{Path: path, Data: make(map[string]interface{})}
}

func (conf *ConfigEngine) loadFromYaml() error {
	file, err := ioutil.ReadFile(conf.Path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, &conf.Data)
	if err != nil {
		return err
	}
	return nil
}
