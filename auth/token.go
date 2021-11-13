package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(data *LoginRequest) (string, error) {
	claim := Claim{
		Email: data.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
			Issuer:    "adrianolmedo",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)
	signedToken, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// ValidateToken is for validation token.
func ValidateToken(t string) (Claim, error) {
	token, err := jwt.ParseWithClaims(t, &Claim{}, verify)
	if err != nil {
		return Claim{}, err
	}
	if !token.Valid {
		return Claim{}, errors.New("invalid token")
	}

	claim, ok := token.Claims.(*Claim)
	if !ok {
		return Claim{}, errors.New("the claims could not be obtained")
	}
	return *claim, nil
}

func verify(t *jwt.Token) (interface{}, error) {
	return verifiKey, nil
}
