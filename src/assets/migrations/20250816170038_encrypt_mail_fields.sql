-- migrate:up
-- Update comments to indicate these fields are encrypted
COMMENT ON COLUMN mails.mail_from IS 'Encrypted sender email address';
COMMENT ON COLUMN mails.recipients IS 'Encrypted array of recipient email addresses';
COMMENT ON COLUMN mails.subject IS 'Encrypted email subject';
COMMENT ON COLUMN mails.content IS 'Encrypted email content';

-- migrate:down
COMMENT ON COLUMN mails.mail_from IS NULL;
COMMENT ON COLUMN mails.recipients IS NULL;
COMMENT ON COLUMN mails.subject IS NULL;
COMMENT ON COLUMN mails.content IS NULL;
