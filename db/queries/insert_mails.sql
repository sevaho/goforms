-- name: InsertMail :one
INSERT INTO mails (
    -- COLUMS --
    created_at, --
    mail_provider, --
    mail_from, --
    success, --
    recipients, --
    subject, --
    content, --
    error --
)
    VALUES (
        -- VALUES --
        $1, --
        $2, --
        $3, --
        $4, --
        $5, --
        $6, --
        $7, --
        $8 --
)
RETURNING
    id;
