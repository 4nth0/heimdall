package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
)

type command struct {
	fs *flag.FlagSet
	fn func(args []string) error
}

var Version string
var ConfigPath string = "./targets.yaml"

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)

	commands := map[string]command{
		// "init": initCmd(),
		// "help": helpCmd(),
		// "targets": helpCmd(),
		"watch": watchCmd(),
	}

	fs := flag.NewFlagSet("heimdall", flag.ExitOnError)
	fs.Parse(os.Args[1:])
	args := fs.Args()

	if len(args) == 0 {
		help("No command provided.")
		return
	}

	if cmd, ok := commands[args[0]]; !ok {
		log.Fatalf("Unknown command: %s", args[0])
	} else if err := cmd.fn(args[1:]); err != nil {
		help("Unknown command: %s", args[0])
		log.Print(err)
	}
}
