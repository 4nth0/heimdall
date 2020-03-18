package mailgun

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/4nth0/heimdall/pkg/gjallarhorn"
	mailgun "github.com/mailgun/mailgun-go/v4"
)

type MailgunNotifier struct {
	Domain     string
	PrivateKey string
	Sender     string
	Recipient  string
}

func New(domain, privateKey, sender, recipient string) (*MailgunNotifier, error) {
	errorKeys := []string{}
	if domain == "" {
		errorKeys = append(errorKeys, "Domain is not provided")
	}
	if privateKey == "" {
		errorKeys = append(errorKeys, "Private Key is not provided")
	}
	if sender == "" {
		errorKeys = append(errorKeys, "Sender is not provided")
	}
	if recipient == "" {
		errorKeys = append(errorKeys, "Recipient is not provided")
	}

	if len(errorKeys) > 0 {
		return nil, errors.New(strings.Join(errorKeys, ", "))
	}

	return &MailgunNotifier{
		Domain:     domain,
		PrivateKey: privateKey,
		Sender:     sender,
		Recipient:  recipient,
	}, nil
}

func (m MailgunNotifier) Notify(kind string, report *gjallarhorn.Report) {
	m.sendSimpleMessage()
}

func (m MailgunNotifier) sendSimpleMessage() {
	mg := mailgun.NewMailgun(m.Domain, m.PrivateKey)

	subject := "Heimdall error report"
	body := "Hi! "

	message := mg.NewMessage(m.Sender, subject, body, m.Recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)
}
