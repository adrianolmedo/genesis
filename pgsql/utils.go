package pgsql

import (
	"database/sql"
	"time"
)

// PtrTime shorcut: return pointer to time.Time.
func PtrTime(t time.Time) *time.Time { return &t }

// ToNullTime convert *time.Time → sql.NullTime.
func ToNullTime(t *time.Time) sql.NullTime {
	if t != nil {
		return sql.NullTime{Time: *t, Valid: true}
	}
	return sql.NullTime{}
}

// ToTimePtr convert sql.NullTime → *time.Time.
func ToTimePtr(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}
