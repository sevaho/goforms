-- migrate:up
CREATE TABLE mails (
    id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    created_at timestamp NOT NULL,
    mail_provider varchar(256) NOT NULL,
    success bool NOT NULL,
    mail_from text NOT NULL, -- should be encrypted
    recipients text[] NOT NULL, -- should be encrypted
    subject text NOT NULL, -- should be encrypted
    content text NOT NULL, -- should be encrypted
    error text
);

-- migrate:down
