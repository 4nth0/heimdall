package main

import (
	"time"

	"github.com/AnthonyCapirchio/heimdall/pkg/notifier/slack"
	"github.com/AnthonyCapirchio/heimdall/pkg/reporter"
	"github.com/AnthonyCapirchio/heimdall/pkg/watcher"
)

var Targets = map[string]watcher.Config{
	"video-ref": watcher.Config{
		URL:       "https://wizads.val1.p.plop.fr/iptv-sfr?zipcode=94550&video=TF103TF581369",
		Method:    "GET",
		Status:    301,
		Latency:   500,
		Threshold: 5,
	},
	"segment": watcher.Config{
		URL:     "http://127.0.0.1:8080/web-tf1",
		Method:  "GET",
		Status:  200,
		Latency: 500,
	},
	"freewheel": watcher.Config{
		URL:       "https://tf1pub.drawapi.com/qa/smart-xml",
		Method:    "GET",
		Status:    200,
		Latency:   500,
		Threshold: 15,
	},
}

var Frequency = 500 * time.Millisecond

func main() {

	notifiers := []reporter.Notifier{slack.New()}
	reporter := reporter.NewReporter(notifiers)
	watchers := watcher.InitWtachers(Targets, Frequency)
	// Use context ..
	responses, ok := watchers.Watch()

	reporter.Analyze(responses)

	<-ok
}
