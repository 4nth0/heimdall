package main

import (
	"flag"
	"fmt"

	"github.com/4nth0/heimdall/internal/config"
	"github.com/4nth0/heimdall/pkg/watcher"
)

type initOpts struct {
	configFile string
}

var initMessages map[string]string = map[string]string{
	"success": "\n\033[32mConfiguration file has been created.\033[0m\nConfiguration file location: %s\n\n",
	"exists":  "\n\033[31mConfiguration file already exists.\033[0m\nConfiguration file location: %s\n\n",
}

// This command create a new configuration file
// heimdall init
func initCmd() command {
	fs := flag.NewFlagSet("heimdall init", flag.ExitOnError)

	return command{fs, func(args []string) error {
		fs.Parse(args)
		return initHeimdall()
	}}
}

func initHeimdall() (err error) {
	if config.Exists(ConfigPath) {
		fmt.Printf(initMessages["exists"], ConfigPath)
		return nil
	}

	cfg := config.InitConfig(ConfigPath)

	cfg.Frequency.Unit = "s"
	cfg.Frequency.Value = 60

	cfg.Targets = map[string]watcher.Config{
		"httpbin": watcher.Config{
			URL:     "https://httpbin.org/get",
			Method:  "GET",
			Status:  200,
			Latency: 500,
		},
	}

	cfg.Save()

	fmt.Printf(initMessages["success"], ConfigPath)

	return nil
}
