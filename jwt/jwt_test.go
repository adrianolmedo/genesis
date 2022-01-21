package jwt_test

import (
	"testing"

	"github.com/adrianolmedo/go-restapi-practice/jwt"
)

const (
	publicSRA = `
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC81YoM084MeqWmv8FLev1T84SG
l0wsz1fLXmZxhgnCEW+dfMTnqfDGYMTup1ca0QtgNpGl2VM3M50SWCXnAUGQ7n6D
kGirT+BDKKHB7gMreyIZEafUjmaEi5oxnCjYkNDxKfNXQkxTyYB2c8yZHnG0ESxb
RAj9cFFEQ1WNTWJBkQIDAQAB
-----END PUBLIC KEY-----
`

	privateSRA = `
-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQC81YoM084MeqWmv8FLev1T84SGl0wsz1fLXmZxhgnCEW+dfMTn
qfDGYMTup1ca0QtgNpGl2VM3M50SWCXnAUGQ7n6DkGirT+BDKKHB7gMreyIZEafU
jmaEi5oxnCjYkNDxKfNXQkxTyYB2c8yZHnG0ESxbRAj9cFFEQ1WNTWJBkQIDAQAB
AoGAW8a9LbbTciU51W1lGLZR4Td9tZxbHXw4g1MCHzKyE2w9/yDg4mcp6oCltggG
wbXP/ZcH+r9BPpcLRBsrcLafkmXE+CG0s7JDoA8Lru4PhnJYMOWhLqmX7lgMu+Ic
ZDaex90nTbq8wp3C18sHXMyF5+ofZZMTaRpp4Ust5QuhH0ECQQD4eyO+Lrfy9sSd
7rZJwYz84+pwmycC8qI5IaptQffiyyIB7ZxJcq4Er1phEeKpTo6/0FUMTCQc4hSJ
9hVjx3cdAkEAwoxYNm1YWaunWWKoarDAx9/U86eqFR/phgtytrpVQkS0R9L5/85y
6+Xdfyx0x1oVmMykwVN/CalIDVoQzKdGBQJBAMtsXG21Z6kENyEorZmiWB8tI+A+
VOjH5OEq25CI4jyMmnHDqiBDP43cVPyFHPAIvTrxfr8LksEGoVP037wJL00CQQCn
HnoUXv+7H7pFDXvREn63863xlECFnwEyNZlYIF5m66/V1wUMWmLcA3yu5xh1uwu8
U2bf74K8YN9VIN43fyWlAkEAyyu3S6Fve1SaB0/L4HkYRnpy8lSIEHwVu5JQo7OP
DIhdpkWF5HTL9xPMKy4aUHcCorx8QC2sdDlGq6QyoMEiXA==
-----END RSA PRIVATE KEY-----
`
)

func TestJWTClaims(t *testing.T) {
	input := "example@gmail.com"

	loadKeys(t)
	token := genToken(t, input)
	verifyClaims(t, token)
}

// loadKeys read mocked credentials.
func loadKeys(t *testing.T) {
	err := jwt.ParseRSA(privateSRA, publicSRA)
	if err != nil {
		t.Fatal(err)
	}
}

// genToken generate token from an email input.
func genToken(t *testing.T, email string) (token string) {
	token, err := jwt.Generate(email)
	if err != nil {
		t.Fatal(err)
	}
	return
}

func verifyClaims(t *testing.T, token string) {
	_, err := jwt.Verify(token)
	if err != nil {
		t.Fatal(err)
	}
}
