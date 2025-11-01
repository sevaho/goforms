package mailproviders

import (
	"github.com/sevaho/goforms/src/internal/models"
)

type MailProvider interface {
	Name() string
	Mail(
		html string,
		text string,
		subject string,
		from models.Recipient,
		to []models.Recipient,
	) error
}

type MailProviders struct {
	providers map[string]MailProvider
}

func New(providers ...MailProvider) *MailProviders {
	m := MailProviders{providers: make(map[string]MailProvider, len(providers))}

	for _, v := range providers {
		m.providers[v.Name()] = v
	}

	return &m

}

func (m *MailProviders) Get(s string) MailProvider {
	return m.providers[s]
}
