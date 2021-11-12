package user

import "errors"

// ErrResourceCantBeEmpty = "the resource can't be empty"
var ErrResourceCantBeEmpty = errors.New("the resource can't be empty")

// ErrResourceDoesNotExist = "resource does not exist".
var ErrResourceDoesNotExist = errors.New("resource does not exist")

type Request struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Response struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
