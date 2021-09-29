package model

import "github.com/dgrijalva/jwt-go"

// Login is for JSON struct of login data.
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Claims is for JSON struct of JWT.
type Claim struct {
	Email string `json:"email"`
	jwt.StandardClaims
}
