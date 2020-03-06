package gjallarhorn

import (
	"errors"
	"sync"
	"time"

	"github.com/4nth0/heimdall/pkg/store"
	"github.com/4nth0/heimdall/pkg/watcher"
)

type Reporter struct {
	Store     *store.Database
	Notifiers []Notifier
	Reports   map[string]*Report
	mu        sync.Mutex
}

type Notifier interface {
	Notify(string, *Report)
}

type Report struct {
	Time     time.Time
	Count    int
	Target   string
	Kind     string
	Expected interface{}
	Current  interface{}
	Error    error
}

func NewReporter(notifiers []Notifier) *Reporter {
	return &Reporter{
		Store:     store.New("./reports.db", true),
		Notifiers: notifiers,
		Reports:   map[string]*Report{},
	}
}

func (r *Reporter) Analyze(response chan watcher.Response) {
	for {
		select {
		case resp := <-response:
			report, valid := isValidResponse(resp)
			r.Store.Push(report)
			if valid == false {
				r.reportInvalidResponse(report, resp.Target)
			} else if exists := r.ErrorExists(resp.TargetID); exists == true {
				r.closeError(resp.TargetID)
			}
		}
	}
}

func (r *Reporter) closeError(target string) {
	r.notify("restore", r.Reports[target])
	delete(r.Reports, target)
}

func (r *Reporter) ErrorExists(target string /* maybe use Kind to store error by kind */) bool {
	if _, ok := r.Reports[target]; ok == true {
		return true
	} else {
		return false
	}
}

func (r *Reporter) reportInvalidResponse(report *Report, target *watcher.Config) {
	if r.ErrorExists(report.Target) {
		r.Reports[report.Target].Count++
	} else {
		r.Reports[report.Target] = report
	}
	if r.Reports[report.Target].Count == target.Threshold {
		r.notify("error", r.Reports[report.Target])
	}
}

func (r *Reporter) notify(kind string, report *Report) {
	for _, notifier := range r.Notifiers {
		go notifier.Notify(kind, report)
	}
}

func isValidResponse(response watcher.Response) (*Report, bool) {
	// Check error
	if response.Error != nil {
		return &Report{
			Time:   time.Now(),
			Kind:   "process",
			Target: response.TargetID,
			Error:  response.Error,
		}, false
	}

	// Check code
	if response.Status != response.Target.Status {
		return &Report{
			Time:     time.Now(),
			Kind:     "status",
			Target:   response.TargetID,
			Error:    errors.New("Invalid status code."),
			Expected: response.Target.Status,
			Current:  response.Status,
		}, false
	}

	// Check latency
	if response.Latency > response.Target.Latency {
		return &Report{
			Time:     time.Now(),
			Kind:     "latency",
			Target:   response.TargetID,
			Error:    errors.New("Latency exceeded."),
			Expected: response.Target.Latency,
			Current:  response.Latency,
		}, false
	}

	return nil, true
}
