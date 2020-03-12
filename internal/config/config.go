package config

import (
	"fmt"
	"io/ioutil"

	"github.com/4nth0/heimdall/pkg/watcher"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	path      string
	Targets   map[string]watcher.Config    `yaml:"targets,omitempty"`
	Notifiers map[string]map[string]string `yaml:"notifiers,omitempty"`
}

// LoadConfig load configuration yaml file content from the specified path
func LoadConfig(path string) *Config {
	t := Config{
		path: path,
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Err: ", err)
	}

	err = yaml.Unmarshal(data, &t)
	if err != nil {
		fmt.Println("error: %v", err)
	}

	return &t
}

func InitConfig(path string) *Config {
	cfg := Config{
		path: path,
	}

	return &cfg
}
