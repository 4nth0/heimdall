package slack

import (
	"encoding/json"
	"fmt"

	"github.com/AnthonyCapirchio/heimdall/pkg/reporter"
)

type SlackNotifier struct{}

func New() *SlackNotifier {
	return &SlackNotifier{}
}

func (sn SlackNotifier) Notify(kind string, report *reporter.Report) {
	fmt.Println("")
	fmt.Println("-- SLACK NOTIFICATION --")
	fmt.Println("kind: ", kind)

	b, _ := json.Marshal(report)

	fmt.Println(string(b))
}
