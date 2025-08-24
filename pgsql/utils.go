package pgsql

import (
	"database/sql"
	"time"
)

// TimeToPtr returns a pointer to the given time.Time.
func TimeToPtr(t time.Time) *time.Time { return &t }

// TimePtrToNull converts *time.Time to sql.NullTime.
func TimePtrToNull(t *time.Time) sql.NullTime {
	if t != nil {
		return sql.NullTime{Time: *t, Valid: true}
	}
	return sql.NullTime{}
}

// NullTimeToPtr converts sql.NullTime to *time.Time.
func NullTimeToPtr(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}
