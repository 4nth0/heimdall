package main

import (
	"flag"
	"fmt"
	"strings"
)

func helpCmd() command {
	fs := flag.NewFlagSet("heimdall help", flag.ExitOnError)

	return command{fs, func(args []string) error {
		fs.Parse(args)
		help()
		return nil
	}}
}

func help(messages ...string) {
	message := "Heimdall version " + Version

	if len(messages) > 0 {
		message = fmt.Sprintf("%s\n\n\033[91m%s\033[0m\n", message, strings.Join(messages, "\n"))
	}

	message += `
Usage: heimdall <command> [command flags]

watch command:
  -config string
    Config file path (default golem.yaml)

`

	fmt.Println(message)
}
