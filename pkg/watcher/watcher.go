package watcher

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	URL       string `yaml:"url,omitempty"`
	Method    string `yaml:"method,omitempty"`
	Status    int    `yaml:"status,omitempty"`
	Latency   int64  `yaml:"latency,omitempty"`
	Threshold int    `yaml:"threshold,omitempty"`
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
						resp <- w.Call(id, &target)
					}(id, target)
				}
			}
		}
	}()

	return resp, done
}

func (w Watcher) Call(id string, target *Config) Response {

	log.WithFields(
		log.Fields{
			"method": target.Method,
			"url":    target.URL,
		}).Info("Make new request.")

	response := Response{
		Target:   target,
		TargetID: id,
	}
	startRequest := time.Now()

	c := &http.Client{}

	req, err := http.NewRequest(target.Method, target.URL, nil)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := c.Do(req)
	if err != nil {
		response.Error = err
		log.WithFields(
			log.Fields{
				"method": target.Method,
				"url":    target.URL,
				"error":  err,
			}).Debug("Request failed.")
		return response
	}

	response.Status = resp.StatusCode
	response.Latency = time.Now().Sub(startRequest).Milliseconds()

	log.WithFields(
		log.Fields{
			"method":  target.Method,
			"url":     target.URL,
			"status":  response.Status,
			"latency": response.Latency,
		}).Debug("Request finished.")

	return response
}
