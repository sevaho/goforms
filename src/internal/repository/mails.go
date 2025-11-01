package repository

import (
	"context"
	"errors"
	"time"

	"github.com/k3a/html2text"
	"github.com/jackc/pgx/v5"
	"github.com/sevaho/goforms/src/db"
	"github.com/sevaho/goforms/src/pkg/encryption"
	"github.com/sevaho/goforms/src/pkg/logger"
	"github.com/sevaho/goforms/src/internal/models"
)

var ErrMailNotFound = errors.New("mail not found")

type Repository struct {
	encryptor *encryption.Encryptor
	db        db.Querier
}

func New(database db.Querier, encryptionKey string) *Repository {
	encryptor, err := encryption.NewEncryptor(encryptionKey)
	if err != nil {
		panic("Failed to initialize encryptor: " + err.Error())
	}
	return &Repository{encryptor: encryptor, db: database}
}

func (r *Repository) GetMails(offset int, limit int) ([]models.DecryptedMail, int, error) {
	params := db.SelectAllMailsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	mails, err := r.db.SelectAllMails(context.Background(), params)
	if err != nil {
		return nil, 0, err
	}

	decryptedMails := make([]models.DecryptedMail, len(mails))
	for i, mail := range mails {
		decryptedMail, err := r.decryptMailWithoutContent(mail)
		if err != nil {
			return nil, 0, err
		}
		decryptedMails[i] = decryptedMail
	}

	count, err := r.db.CountAllMails(context.Background())
	if err != nil {
		return nil, 0, err
	}

	return decryptedMails, int(count), nil
}

func (r *Repository) GetMailByID(id int) (models.DecryptedMailWithContent, error) {
	mail, err := r.db.SelectMailByID(context.Background(), int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.DecryptedMailWithContent{}, ErrMailNotFound
		}
		return models.DecryptedMailWithContent{}, err
	}

	decryptedMail, err := r.decryptMailWithoutContent(mail)
	if err != nil {
		return models.DecryptedMailWithContent{}, err
	}

	decryptedContent, err := r.encryptor.Decrypt(mail.Content)
	if err != nil {
		return models.DecryptedMailWithContent{}, err
	}

	return models.DecryptedMailWithContent{
		DecryptedMail: decryptedMail,
		Content:       decryptedContent,
	}, nil
}

func (r *Repository) Store(
	mailProvider string,
	subject string,
	content string,
	mail_from string,
	recipients []models.Recipient,
	error error,
) {
	encryptedSubject, err := r.encryptor.Encrypt(subject)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("Failed to encrypt subject")
		panic(err)
	}

	encryptedContent, err := r.encryptor.Encrypt(content)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("Failed to encrypt content")
		panic(err)
	}

	encryptedMailFrom, err := r.encryptor.Encrypt(mail_from)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("Failed to encrypt mail_from")
		panic(err)
	}

	recipient_as_string := []string{}

	for _, recepient := range recipients {
		recipient_as_string = append(recipient_as_string, recepient.Email+" "+recepient.Name)
	}

	encryptedRecipients, err := r.encryptor.EncryptStringSlice(recipient_as_string)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("Failed to encrypt recipients")
		panic(err)
	}

	params := db.InsertMailParams{
		CreatedAt:    db.TimeToPGTimestamp(time.Now().UTC()),
		Subject:      encryptedSubject,
		Content:      encryptedContent,
		MailFrom:     encryptedMailFrom,
		MailProvider: mailProvider,
		Recipients:   encryptedRecipients,
	}

	if error != nil {
		params.Error = db.StringtoPGText(error.Error())
		params.Success = false
	} else {
		params.Success = true
	}

	id, err := r.db.InsertMail(context.Background(), params)

	if err != nil {
		panic(err)
	}
	logger.Logger.Info().Msgf("Mail stored with id: %d", id)
}

func (r *Repository) decryptMailWithoutContent(mail db.Mail) (models.DecryptedMail, error) {
	decryptedSubject, err := r.encryptor.Decrypt(mail.Subject)
	if err != nil {
		return models.DecryptedMail{}, err
	}

	decryptedMailFrom, err := r.encryptor.Decrypt(mail.MailFrom)
	if err != nil {
		return models.DecryptedMail{}, err
	}

	mailFromPlanText := html2text.HTML2Text(decryptedMailFrom)

	decryptedRecipients, err := r.encryptor.DecryptStringSlice(mail.Recipients)
	if err != nil {
		return models.DecryptedMail{}, err
	}

	return models.DecryptedMail{
		ID:           mail.ID,
		CreatedAt:    mail.CreatedAt.Time,
		MailProvider: mail.MailProvider,
		Success:      mail.Success,
		MailFrom:     decryptedMailFrom,
		MailFromPlainText:			mailFromPlanText,
		Recipients:   decryptedRecipients,
		Subject:      decryptedSubject,
		Error:        mail.Error.String,
	}, nil
}
