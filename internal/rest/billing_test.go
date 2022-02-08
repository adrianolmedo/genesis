package rest

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/adrianolmedo/go-restapi/internal/mock"
	"github.com/adrianolmedo/go-restapi/internal/service"

	"github.com/labstack/echo/v4"
)

func TestGenerateInvoice(t *testing.T) {
	tt := []struct {
		name         string
		inputForm    []byte
		wantResponse string
		wantHTTPCode int
	}{
		{
			name: "wrong-body-json",
			inputForm: []byte(`{
				"header":{
					"client_id":1,
				},
				"items":[
					{
						"product_id":"1"
					},
					{
						"product_id":2
					}
				]
			}`),
			wantResponse: `{"message_error":{"code":"ER002","content":"the JSON structure is not correct"}}`,
			wantHTTPCode: http.StatusBadRequest,
		},
		{
			name: "not-header",
			inputForm: []byte(`{
				"items":[
					{
						"product_id":1
					},
					{
						"product_id":2
					}
				]
			}`),
			wantResponse: `{"message_error":{"code":"ER004","content":"user not found"}}`,
			wantHTTPCode: http.StatusNotFound,
		},
		{
			name: "not-items",
			inputForm: []byte(`{
				"header":{
					"client_id":1
				}
			}`),
			wantResponse: `{"message_error":{"code":"ER004","content":"item list can't be empty"}}`,
			wantHTTPCode: http.StatusInternalServerError,
		},
		{
			name: "right-body-json",
			inputForm: []byte(`{
				"header":{
					"client_id":1
				},
				"items":[
					{
						"product_id":1
					},
					{
						"product_id":2
					}
				]
			}`),
			wantResponse: `{"message_ok":{"code":"OK002","content":"invoice generated"},"data":{"header":{"client_id":1},"items":[{"product_id":1},{"product_id":2}]}}`,
			wantHTTPCode: http.StatusCreated,
		},
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

		err = reqMeth(c, generateInvoice(*s))
		if err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}

		// Check body JSON reponse. To fix linebreak in w.Body.String(), add strings.TrimRight and cut if right off the string.
		// https://stackoverflow.com/a/45275479/3408901
		if tc.wantResponse != strings.TrimRight(w.Body.String(), "\n") {
			t.Errorf("%s: wrong response body: want %s, got %s", tc.name, tc.wantResponse, w.Body.String())
		}

		// w.Body es lo que devuelve c.JSON, para comprobarlo, lo mostramos w.Body.String() en un t.Logf
		//t.Logf("%s: response body: %v\n", tc.name, w.Body.String())

		// Check HTTP status code.
		if tc.wantHTTPCode != w.Code {
			t.Errorf("%s: http code: want %d, got %d", tc.name, tc.wantHTTPCode, w.Code)
		}
	}
}
