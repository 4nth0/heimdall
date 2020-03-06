package watcher

import (
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	URL       string
	Status    int
	Latency   int64
	Method    string
	Threshold int // ??
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
		return response
	}

	response.Status = resp.StatusCode
	response.Latency = time.Now().Sub(startRequest).Milliseconds()

	return response
}
