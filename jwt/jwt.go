package jwt

import (
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	once       sync.Once
)

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// Generate signed token from user's email.
func Generate(userEmail string) (string, error) {
	claims := Claims{
		Email: userEmail,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
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

// LoadFiles read RSA files. Se asegura poderse ejecutar una única vez.
func LoadFiles(privateFile, publicFile string) error {
	var err error

	once.Do(func() {
		err = loadFiles(privateFile, publicFile)
	})
	return err
}

func loadFiles(privateFile, publicFile string) error {
	// Con ioutil.ReadFile(s) solo lees el contenido
	// y no tienes que procuparte por cerrar el recurso con defer.
	privateBytes, err := ioutil.ReadFile(privateFile)
	if err != nil {
		return err
	}

	publicBytes, err := ioutil.ReadFile(publicFile)
	if err != nil {
		return err
	}

	return parseRSA(privateBytes, publicBytes)
}

func ParseRSA(privateRSA, publicRSA string) error {
	return parseRSA([]byte(privateRSA), []byte(publicRSA))
}

func parseRSA(privateBytes, publicBytes []byte) error {
	var err error

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		return err
	}

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		return err
	}

	return nil
}
