package pgsql

import (
	"database/sql"
	"time"
)

// TimePtr returns a pointer to the given time.Time.
func TimePtr(t time.Time) *time.Time { return &t }

// NullTimeFromPtr converts *time.Time to sql.NullTime.
func NullTimeFromPtr(t *time.Time) sql.NullTime {
	if t != nil {
		return sql.NullTime{Time: *t, Valid: true}
	}
	return sql.NullTime{}
}

// PtrFromNullTime converts sql.NullTime to *time.Time.
func PtrFromNullTime(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}
