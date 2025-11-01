-- name: SelectMailByID :one
SELECT
    *
FROM
    mails
WHERE
    id = $1;
