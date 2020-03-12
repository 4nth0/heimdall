package watcher

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func Call(id string, target *Config) Response {

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
