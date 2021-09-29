package model

import "errors"

// ErrResourceCantBeEmpty = "the resource can't be empty"
var ErrResourceCantBeEmpty = errors.New("the resource can't be empty")

// ErrResourceDoesNotExist = "resource does not exist".
var ErrResourceDoesNotExist = errors.New("resource does not exist")
