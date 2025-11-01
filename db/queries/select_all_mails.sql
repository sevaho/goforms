-- name: SelectAllMails :many
SELECT
    *
FROM
    mails
ORDER BY
    created_at DESC
LIMIT $1 offset $2;

-- name: CountAllMails :one
SELECT
    count(*)
FROM
    mails;
-- name: FilterMailsOnCreatedAt :many
SELECT
    *
FROM
    mails
WHERE
    created_at > sqlc.narg ('created_at_lt')
    AND created_at < sqlc.narg ('created_at_gt')
ORDER BY
    created_at DESC
LIMIT $1;
