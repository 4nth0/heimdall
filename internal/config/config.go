package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/4nth0/heimdall/pkg/watcher"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	path      string
	Frequency struct {
		Unit  string `yaml:"unit,omitempty"`
		Value int    `yaml:"value,omitempty"`
	}
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

func (c Config) GetFrequency() time.Duration {
	switch c.Frequency.Unit {
	case "sec", "secconde":
		return time.Duration(c.Frequency.Value) * time.Second
	case "min", "minute":
		return time.Duration(c.Frequency.Value) * time.Minute
	}

	return 60 * time.Minute
}

func InitConfig(path string) *Config {
	cfg := Config{
		path: path,
	}

	return &cfg
}

func Exists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (c Config) Save() error {

	b, _ := yaml.Marshal(c)
	err := ioutil.WriteFile(c.path, b, 0644)
	if err != nil {
		fmt.Println("Err: ", err)
	}

	return nil
}
