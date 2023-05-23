package jwt

import (
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	once       sync.Once
)

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// Generate signed token from user's email.
func Generate(userEmail string) (string, error) {
	claims := Claims{
		Email: userEmail,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
			Issuer:    "go-restapi",
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := t.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// Verify signed token.
func Verify(token string) (Claims, error) {
	t, err := jwt.ParseWithClaims(token, &Claims{}, verify)
	if err != nil {
		return Claims{}, err
	}

	if !t.Valid {
		return Claims{}, errors.New("invalid token")
	}

	claims, ok := t.Claims.(*Claims)
	if !ok {
		return Claims{}, errors.New("the claims could not be obtained")
	}
	return *claims, nil
}

func verify(t *jwt.Token) (interface{}, error) {
	return publicKey, nil
}

// LoadFiles read RSA files. Ensures that it can only be executed once.
func LoadFiles(privateFile, publicFile string) error {
	var err error

	once.Do(func() {
		err = loadFiles(privateFile, publicFile)
	})
	return err
}

func loadFiles(private, public string) error {
	// With ioutil.ReadFile(s) only reads the content
	// and you don't worry to close the resource with defer.
	privateBytes, err := ioutil.ReadFile(private)
	if err != nil {
		return err
	}

	publicBytes, err := ioutil.ReadFile(public)
	if err != nil {
		return err
	}

	return parseRSA(privateBytes, publicBytes)
}

func ParseRSA(private, public string) error {
	return parseRSA([]byte(private), []byte(public))
}

func parseRSA(private, public []byte) error {
	var err error

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(private)
	if err != nil {
		return err
	}

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(public)
	if err != nil {
		return err
	}

	return nil
}
