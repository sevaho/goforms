package mailersend

import (
	"context"
	"io"
	"time"

	"github.com/mailersend/mailersend-go"
	"github.com/sevaho/goforms/src/pkg/logger"
	"github.com/sevaho/goforms/src/internal/models"
)

type MailerSend struct {
	apiKey string
	name   string
}

func New(apiKey string) *MailerSend {
	return &MailerSend{apiKey: apiKey, name: "mailersend"}
}

func (m *MailerSend) Name() string { return m.name }

func (m *MailerSend) Mail(
	html string,
	text string,
	subject string,
	from models.Recipient,
	to []models.Recipient,
) error {
	ms := mailersend.NewMailersend(m.apiKey)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	mail_from := mailersend.From{
		Name:  from.Name,
		Email: from.Email,
	}

	recipients := []mailersend.Recipient{}
	for _, x := range to {
		recipients = append(recipients, mailersend.Recipient{
			Name:  x.Name,
			Email: x.Email,
		})
	}

	// Send in 5 minute
	// sendAt := time.Now().Add(time.Minute * 5).Unix()

	// tags := []string{"foo", "bar"}

	message := ms.Email.NewMessage()

	message.SetFrom(mail_from)
	message.SetRecipients(recipients)
	message.SetSubject(subject)
	message.SetHTML(html)
	message.SetText(text)
	// message.SetTags(tags)
	// message.SetSendAt(sendAt)
	// message.SetInReplyTo("client-id")

	res, err := ms.Email.Send(ctx, message)

	if err != nil {
		return err
	}

	response, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	//jsonmessage, _ := json.Marshal(body)

	logger.Logger.Info().Msgf("Email send with mailersend. [message-id=%s, to=%v, response_status_code=%d, response=%s].", res.Header.Get("X-Message-Id"), to, res.StatusCode, string(response))

	return nil
}
