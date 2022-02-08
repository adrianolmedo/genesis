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

func TestAddProduct(t *testing.T) {
	tt := []struct {
		name         string
		inputForm    []byte
		wantResponse string
		wantHTTPCode int
	}{
		{
			name: "wrong-body-json",
			inputForm: []byte(`{
				"name": Coca-Cola,
				"observations": "",
				"price": 3
			}`),
			wantResponse: `{"message_error":{"code":"ER002","content":"the JSON structure is not correct"}}`,
			wantHTTPCode: http.StatusBadRequest,
		},
		{
			name:         "empty-body-json",
			inputForm:    []byte(``),
			wantResponse: `{"message_error":{"code":"ER004","content":"the product has no name"}}`,
			wantHTTPCode: http.StatusInternalServerError,
		},
		{
			name: "right-body-json",
			inputForm: []byte(`{
				"name": "Coca-Cola",
				"price": 3,
				"observations": ""
			}`),
			wantResponse: `{"message_ok":{"code":"OK002","content":"product added"},"data":{"name":"Coca-Cola","observations":"","price":3}}`,
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

		err = reqMeth(c, addProduct(*s))
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

func TestFindProduct(t *testing.T) {
	tt := []struct {
		input        string
		wantHTTPCode int
		wantResponse string
	}{
		{
			input:        "1",
			wantHTTPCode: http.StatusOK,
			wantResponse: `{"message_ok":{"code":"OK002","content":""},"data":{"id":1,"name":"Coca-Cola","observations":"","price":3}}`,
		},
		{
			input:        "3",
			wantHTTPCode: http.StatusNotFound,
			wantResponse: `{"message_error":{"code":"ER007","content":"product not found"}}`,
		},
		{
			input:        "0",
			wantHTTPCode: http.StatusNotFound,
			wantResponse: `{"message_error":{"code":"ER007","content":"product not found"}}`,
		},
		{
			input:        "-1",
			wantHTTPCode: http.StatusBadRequest,
			wantResponse: `{"message_error":{"code":"ER005","content":"positive number expected for ID product"}}`,
		},
	}

	s, err := service.New(mock.StorageOk{})
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range tt {
		e := echo.New()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c := e.NewContext(r, w)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues(tc.input)

		err = reqMeth(c, findProduct(*s))
		if err != nil {
			t.Fatalf("input %s: %v", tc.input, err)
		}

		if tc.wantResponse != strings.TrimRight(w.Body.String(), "\n") {
			t.Errorf("input %s: wrong response body: want %s, got %s", tc.input, tc.wantResponse, w.Body.String())
		}

		if tc.wantHTTPCode != w.Code {
			t.Errorf("input %s: http code: want %d, got %d", tc.input, tc.wantHTTPCode, w.Code)
		}
	}
}

func TestListProducts(t *testing.T) {
	rp := mock.StorageOk{}
	s, err := service.New(rp)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c := e.NewContext(r, w)

	err = reqMeth(c, listProducts(*s))
	if err != nil {
		t.Fatal(err)
	}

	// Check HTTP status code.
	if http.StatusOK != w.Code {
		t.Errorf("http code: want %d, got %d", http.StatusOK, w.Code)
	}

	t.Logf("response body: %v", w.Body.String())

	/*if wantResponse != strings.TrimRight(w.Body.String(), "\n") {
		t.Errorf("wrong response body: want %s, got %s", wantResponse, w.Body.String())
	}*/
}
