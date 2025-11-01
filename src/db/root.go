package db

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/avast/retry-go"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sevaho/goforms/src/pkg/logger"
)

var DB_CONNECTION_RETRIES int = 3
var DB_CONNECTION_RETRY_DELAY time.Duration = 1

// function to make sure we can connect to database, test with retries
func (db *Queries) TestConnection(dsn string) error {
	err := retry.Do(
		func() (err error) {
			ctx := context.Background()
			_, err = db.db.Exec(ctx, "SELECT NOW()")
			if err != nil {
				logger.Logger.Warn().Err(err).Msgf("Unable to connect to db %s", strings.Split(dsn, "@")[1])
			}
			return
		},
		retry.Attempts(uint(DB_CONNECTION_RETRIES)),
		retry.Delay(DB_CONNECTION_RETRY_DELAY),
	)
	if err != nil {
		msg := fmt.Sprintf("Unable to connect to db %s", strings.Split(dsn, "@")[1])
		return errors.New(msg)
	}

	logger.Logger.Info().Msgf("[DB] connected to %s", strings.Split(dsn, "@")[1])
	return nil
}

func NewDB(dsn string, withTransaction bool) (*Queries, *pgx.Tx) {
	var transaction *pgx.Tx
	var db *Queries

	conn, err := pgxpool.New(context.Background(), dsn)

	if err != nil {
		panic(err)
	}

	if withTransaction {
		tx, err := conn.Begin(context.Background())

		transaction = &tx
		if err != nil {
			panic(err)
		}
		db = New(tx)
	} else {
		db = New(conn)
	}

	db.TestConnection(dsn)
	return db, transaction
}
