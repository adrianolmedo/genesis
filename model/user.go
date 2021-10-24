package model

import (
	"database/sql"
	"time"
)

// User model.
type User struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
	ID        int64     `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
}

type Users []*User

// helpers...
type scanner interface {
	Scan(dest ...interface{}) error
}

// ScanRowUser return nulled fields of User parsed.
func ScanRowUser(s scanner) (*User, error) {
	var updatedAtNull, deletedAtNull sql.NullTime
	user := &User{}

	err := s.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&updatedAtNull,
		&deletedAtNull,
	)
	if err != nil {
		return &User{}, err
	}

	user.UpdatedAt = updatedAtNull.Time
	user.DeletedAt = updatedAtNull.Time

	return user, nil
}

// TimeToNull helper to try empty time fields.
func TimeToNull(t time.Time) sql.NullTime {
	null := sql.NullTime{Time: t}

	if !null.Time.IsZero() {
		null.Valid = true
	}
	return null
}
