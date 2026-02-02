package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"assecor.assessment.test/internal/mock"
)

func newTestApp(_ *testing.T) *application {
	return &application{
		logger: log.New(io.Discard, "", 0),
		models: mock.NewTestModels(),
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(_ *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := rs.Body.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body
}

func (ts *testServer) post(t *testing.T, urlPath string, body []byte) (int, http.Header,
	[]byte) {
	rs, err := ts.Client().Post(ts.URL+urlPath, "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := rs.Body.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	body, err = io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body
}

func readJSON(t *testing.T, body []byte, dst interface{}) {
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		switch {
		case errors.As(err, &syntaxError):
			t.Fatal(fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset))
		case errors.Is(err, io.ErrUnexpectedEOF):
			t.Fatal(errors.New("body contains badly-formed JSON"))
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				t.Fatal(fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field))
			}
			t.Fatal(fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset))
		case errors.Is(err, io.EOF):
			t.Fatal(errors.New("body must not be empty"))
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			t.Fatal(fmt.Errorf("body contains unknown key %s", fieldName))
		case errors.As(err, &invalidUnmarshalError):
			t.Fatal(err)
		default:
			t.Fatal(err)
		}
	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		t.Fatal(errors.New("body must only contain a single JSON value"))
	}
}

func writeJSON(t *testing.T, data interface{}) []byte {
	if _, ok := data.([]interface{}); ok {
		// return an empty array instead of null
		data = make([]string, 0)
	}
	js, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	return js
}
