package main

import (
	"flag"

	"github.com/4nth0/heimdall/internal/config"
	"github.com/4nth0/heimdall/pkg/gjallarhorn"
	"github.com/4nth0/heimdall/pkg/notifier/mailgun"
	"github.com/4nth0/heimdall/pkg/notifier/slack"
	"github.com/4nth0/heimdall/pkg/watcher"

	log "github.com/sirupsen/logrus"
)

type watchOpts struct {
	configFile string
}

func watchCmd() command {
	fs := flag.NewFlagSet("heimdall watch", flag.ExitOnError)

	opts := &watchOpts{}

	fs.StringVar(&opts.configFile, "config", ConfigPath, "Config File")

	return command{fs, func(args []string) error {
		fs.Parse(args)
		return watch(opts)
	}}
}

func watch(opts *watchOpts) (err error) {
	cfg := config.LoadConfig(opts.configFile)
	frequency := cfg.GetFrequency()

	notifiers := initNotifiers(cfg)
	reporter := gjallarhorn.NewReporter(notifiers)
	watchers := watcher.InitWtachers(cfg.Targets, frequency)
	// Use context ..
	responses, _ := watchers.Watch()

	reporter.Analyze(responses)

	return nil
}

func initNotifiers(cfg *config.Config) []gjallarhorn.Notifier {
	notifiers := []gjallarhorn.Notifier{}

	if len(cfg.Notifiers["mailgun"]) > 0 {
		mg, err := mailgun.New(
			cfg.Notifiers["mailgun"]["domain"],
			cfg.Notifiers["mailgun"]["private_key"],
			cfg.Notifiers["mailgun"]["sender"],
			cfg.Notifiers["mailgun"]["recipient"],
		)
		if err != nil {
			log.Error("Mailgun initialization error: ", err.Error())
		} else {
			notifiers = append(notifiers, mg)
		}
	}

	notifiers = append(notifiers, slack.New())

	return notifiers
}
