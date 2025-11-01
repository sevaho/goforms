package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// See https://github.com/emicklei/pgtalk/blob/v1.3.0/convert/convert.go
func BooltoPGBool(b bool) pgtype.Bool {
	return pgtype.Bool{Bool: b, Valid: true}
}

func TimeToPGDate(t time.Time) pgtype.Date {
	return pgtype.Date{Time: t, Valid: true}
}

func TimeToPGTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func TimeToPGTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{Time: t.UTC(), Valid: true}
}

func UUIDToPGUUID(u uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: u, Valid: true}
}

func StringtoPGText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: true}
}
