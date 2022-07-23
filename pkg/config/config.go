package config

import (
	"gojek/web-server-gin/pkg/handleError"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var DBpass string
var RedisPass string

type Config struct {
	Db    string `yaml:"db"`
	Cache string `yaml:"cache"`
}

func Init() {
	config := &Config{}

	cfgFile, err := ioutil.ReadFile("config.yaml")
	handleError.Check(err)

	err = yaml.Unmarshal(cfgFile, config)
	handleError.Check(err)

	RedisPass = config.Cache
	DBpass = config.Db
}

func TestInit() {
	config := &Config{}

	cfgFile, err := ioutil.ReadFile("../config.yaml")
	handleError.Check(err)

	err = yaml.Unmarshal(cfgFile, config)
	handleError.Check(err)

	RedisPass = config.Cache
	DBpass = config.Db
}
