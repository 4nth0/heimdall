package main

import (
	"os"

	"github.com/4nth0/heimdall/internal/config"
	"github.com/4nth0/heimdall/pkg/gjallarhorn"
	"github.com/4nth0/heimdall/pkg/notifier/mailgun"
	"github.com/4nth0/heimdall/pkg/notifier/slack"
	"github.com/4nth0/heimdall/pkg/watcher"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	cfg := config.LoadConfig("targets.yaml")
	frequency := cfg.GetFrequency()

	notifiers := initNotifiers(cfg)
	reporter := gjallarhorn.NewReporter(notifiers)
	watchers := watcher.InitWtachers(cfg.Targets, frequency)
	// Use context ..
	responses, _ := watchers.Watch()

	reporter.Analyze(responses)
}

func initNotifiers(cfg *config.Config) []gjallarhorn.Notifier {
	return []gjallarhorn.Notifier{
		mailgun.New(
			cfg.Notifiers["mailgun"]["domain"],
			cfg.Notifiers["mailgun"]["private_key"],
			cfg.Notifiers["mailgun"]["sender"],
			cfg.Notifiers["mailgun"]["recipient"],
		),
		slack.New(),
	}
}
