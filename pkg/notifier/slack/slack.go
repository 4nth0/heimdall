package slack

import (
	"encoding/json"
	"fmt"

	"github.com/4nth0/heimdall/pkg/gjallarhorn"
)

type SlackNotifier struct{}

func New() *SlackNotifier {
	return &SlackNotifier{}
}

func (sn SlackNotifier) Notify(kind string, report *gjallarhorn.Report) {
	fmt.Println("")
	fmt.Println("-- SLACK NOTIFICATION --")
	fmt.Println("kind: ", kind)

	b, _ := json.Marshal(report)

	fmt.Println(string(b))
}
