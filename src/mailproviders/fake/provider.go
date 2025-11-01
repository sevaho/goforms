package fake

import (
	"github.com/sevaho/goforms/src/internal/models"
)

type Fake struct {
	name string
}

func New() *Fake {
	return &Fake{name: "fake"}
}

func (m *Fake) Name() string { return m.name }

func (m *Fake) Mail(
	html string,
	text string,
	subject string,
	from models.Recipient,
	to []models.Recipient,
) error {
	return nil
}
