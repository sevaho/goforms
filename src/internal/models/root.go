package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Recipient struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// TODO:  <26-04-25, Sebastiaan Van Hoecke> // Add validation logic, fields should be checked because required
type FormTemplate struct {
	ID          uuid.UUID   `json:"id"`
	Skipcaptcha bool        `json:"skipcaptcha"`
	Recipients  []Recipient `json:"recipients"`
	Subject     string      `json:"subject"`
	Name        string      `json:"name"`
	Sender      Recipient   `json:"sender"`
	Provider    string      `json:"provider"`
}

type FormsConfig struct {
	Forms []FormTemplate `json:"forms"`
}

func (c *FormsConfig) Check() bool {
	return len(c.Forms) != 0
}

func (c *FormsConfig) Get(id uuid.UUID) (*FormTemplate, error) {
	for _, v := range c.Forms {
		if v.ID == id {
			return &v, nil
		}
	}
	return nil, errors.New("No form found with ID: " + id.String())
}

type DecryptedMail struct {
	ID                int32
	CreatedAt         time.Time
	MailProvider      string
	Success           bool
	MailFrom          string
	MailFromPlainText string
	Recipients        []string
	Subject           string
	Error             string
}

type DecryptedMailWithContent struct {
	DecryptedMail
	Content string
}
