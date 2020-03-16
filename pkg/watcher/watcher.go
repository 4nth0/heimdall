package watcher

import (
	"time"
)

type Config struct {
	URL       string            `yaml:"url,omitempty"`
	Method    string            `yaml:"method,omitempty"`
	Status    int               `yaml:"status,omitempty"`
	Latency   int64             `yaml:"latency,omitempty"`
	Threshold int               `yaml:"threshold,omitempty"`
	Matchs    map[string]string `yaml:"matchs,omitempty"`
}

type Watcher struct {
	Targets   map[string]Config
	Frequency time.Duration
}

type Response struct {
	Target      *Config
	TargetID    string
	Timestanp   string
	Site        string
	Status      int
	Latency     int64
	BodySize    int
	BodyContent string
	Error       error
}

func InitWtachers(targets map[string]Config, frequency time.Duration) *Watcher {
	return &Watcher{
		Targets:   targets,
		Frequency: frequency,
	}
}

func (w *Watcher) Watch() (chan Response, chan bool) {
	resp := make(chan Response)
	ticker := time.NewTicker(w.Frequency)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				for id, target := range w.Targets {
					go func(id string, target Config) {
						resp <- Call(id, &target)
					}(id, target)
				}
			}
		}
	}()

	return resp, done
}
