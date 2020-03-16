package mailgun

import (
	"context"
	"fmt"
	"log"
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
	if domain == "" {
	}
	if privateKey == "" {
	}
	if sender == "" {
	}
	if recipient == "" {
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

	subject := "Fancy subject!"
	body := "Hello from Mailgun Go!"

	message := mg.NewMessage(m.Sender, subject, body, m.Recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)
}
