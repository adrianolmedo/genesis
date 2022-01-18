package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/adrianolmedo/go-restapi-practice/internal/mock"
	"github.com/adrianolmedo/go-restapi-practice/internal/service"
	"github.com/adrianolmedo/go-restapi-practice/jwt"

	"github.com/labstack/echo/v4"
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

// TestLoginUser successful attemp.
func TestLoginUser(t *testing.T) {
	inputForm := []byte(`{
		"email": "qwerty@hotmail.com",
		"password": "1234567a"
	}`)

	// Read mocked credentials.
	err := jwt.ParseRSA([]byte(privateSRA), []byte(publicSRA))
	if err != nil {
		t.Fatal(err)
	}

	s, err := service.New(mock.StorageOk{})
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(inputForm))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c := e.NewContext(r, w)

	err = reqMeth(c, loginUser(*s))
	if err != nil {
		t.Fatal(err)
	}

	// Check HTTP status code.
	if http.StatusCreated != w.Code {
		t.Errorf("http code: want %d, got %d", http.StatusCreated, w.Code)
	}

	got := w.Body.String()
	gotResponse := response{messageOK: &messageOK{}, messageError: &messageError{}}
	err = json.NewDecoder(w.Body).Decode(&gotResponse)
	if err != nil {
		t.Fatal(err)
	}

	wantResponse := newResponse(MsgOK, "OK004", "logged", gotResponse.Data)
	want, err := json.Marshal(wantResponse)
	if err != nil {
		t.Fatal(err)
	}

	// Check body JSON whit the same token.
	if string(want) != strings.TrimRight(got, "\n") {
		t.Errorf("response body: want %s, got %s", want, got)
	}
}

// TestLoginUserFailure attempts.
func TestLoginUserFailure(t *testing.T) {
	tt := []struct {
		name         string
		inputForm    []byte
		wantResponse string
		wantHTTPCode int
	}{
		{
			name: "wrong-body-json",
			inputForm: []byte(`{
				"email": qwerty@hotmail.com",
				"password": "1234567a"
			}`),
			wantResponse: `{"message_error":{"code":"ER002","content":"the JSON structure is not correct"}}`,
			wantHTTPCode: http.StatusBadRequest,
		},
		{
			name: "email-not-valid",
			inputForm: []byte(`{
				"email": "qwertyhotmail.com",
				"password": "1234567a"
			}`),
			wantResponse: `{"message_error":{"code":"ER009","content":"email not valid"}}`,
			wantHTTPCode: http.StatusBadRequest,
		},
		{
			name: "not-found",
			inputForm: []byte(`{
				"email": "ytrewq@hotmail.com",
				"password": "1234567a"
			}`),
			wantResponse: `{"message_error":{"code":"ER007","content":"user not found"}}`,
			wantHTTPCode: http.StatusUnauthorized,
		},
	}

	// Read mocked credentials.
	err := jwt.ParseRSA([]byte(privateSRA), []byte(publicSRA))
	if err != nil {
		t.Fatal(err)
	}

	s, err := service.New(mock.StorageOk{})
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range tt {
		e := echo.New()
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(tc.inputForm))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c := e.NewContext(r, w)

		err = reqMeth(c, loginUser(*s))
		if err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}

		// Check HTTP status code.
		if tc.wantHTTPCode != w.Code {
			t.Errorf("%s: http code: want %d, got %d", tc.name, tc.wantHTTPCode, w.Code)
		}

		// Check body JSON reponse.
		if tc.wantResponse != strings.TrimRight(w.Body.String(), "\n") {
			t.Errorf("%s: response body: want %s, got %s", tc.name, tc.wantResponse, w.Body.String())
		}
	}
}
