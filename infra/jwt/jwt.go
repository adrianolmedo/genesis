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
	signKey   *rsa.PrivateKey
	verifiKey *rsa.PublicKey
	once      sync.Once
)

// JWTClaims is for JSON struct of JWT.
type JWTClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// New generate signed token.
func New(userEmail string) (string, error) {
	claims := JWTClaims{
		Email: userEmail,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
			Issuer:    "adrianolmedo",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedToken, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// Verify is for validation token.
func Verify(strToken string) (JWTClaims, error) {
	token, err := jwt.ParseWithClaims(strToken, &JWTClaims{}, verify)
	if err != nil {
		return JWTClaims{}, err
	}

	if !token.Valid {
		return JWTClaims{}, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return JWTClaims{}, errors.New("the claims could not be obtained")
	}
	return *claims, nil
}

func verify(token *jwt.Token) (interface{}, error) {
	return verifiKey, nil
}

// LoadFiles read SRA files.
// Se asegura poderse ejecutar una única vez.
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

func parseRSA(privateBytes, publicBytes []byte) error {
	var err error

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		return err
	}

	verifiKey, err = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		return err
	}

	return nil
}
