package main

import (
	"os"
	"time"

	"github.com/4nth0/heimdall/internal/config"
	"github.com/4nth0/heimdall/pkg/gjallarhorn"
	"github.com/4nth0/heimdall/pkg/notifier/slack"
	"github.com/4nth0/heimdall/pkg/watcher"
	log "github.com/sirupsen/logrus"
)

var Frequency = 5000 * time.Millisecond

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	cfg := config.LoadConfig("targets.yaml")

	notifiers := []gjallarhorn.Notifier{slack.New()}
	reporter := gjallarhorn.NewReporter(notifiers)
	watchers := watcher.InitWtachers(cfg.Targets, Frequency)
	// Use context ..
	responses, _ := watchers.Watch()

	reporter.Analyze(responses)
}
