package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"assecor.assessment.test/internal/data"
)

func (app *application) readIDParam(param string) (int64, error) {
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

// Define a writeJSON() helper for sending responses.
func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
	if v, ok := data.([]interface{}); ok {
		if reflect.ValueOf(v).Len() == 0 {
			// return an empty array instead of null
			data = make([]string, 0)
		}
	}
	// Encode the data to JSON, returning the error if there was one.
	js, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	// Append a newline to make it easier to view in terminal applications.
	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		// If the JSON contains a field which cannot be mapped to the target destination
		// then Decode() will now return an error message in the format "json: unknown
		// field "<name>"". We check for this, extract the field name from the error,
		// and interpolate it into our custom error message. Note that there's an open
		// issue at https://github.com/golang/go/issues/29035 regarding turning this
		// into a distinct error type in the future.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}
	// Call Decode() again, using a pointer to an empty anonymous struct as the
	// destination. If the request body only contained a single JSON value this will
	// return an io.EOF error. So if we get anything else, we know that there is
	// additional data in the request body and we return our own custom error message.
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}

// Convert color id to string
func (app *application) formatPerson(persons *data.Person) interface{} {
	return struct {
		ID       int64  `json:"id"`
		Name     string `json:"name"`
		Lastname string `json:"lastname"`
		Zipcode  string `json:"zipcode"`
		City     string `json:"city"`
		Color    string `json:"color"`
	}{
		ID:       persons.ID,
		Name:     persons.Name,
		Lastname: persons.Lastname,
		Zipcode:  persons.Zipcode,
		City:     persons.City,
		Color:    data.Color(persons.Color).String(),
	}
}

// Convert color id to string
func (app *application) formatPersonArray(persons []*data.Person) interface{} {
	var formated []interface{}
	for _, person := range persons {
		formated = append(formated, struct {
			ID       int64  `json:"id"`
			Name     string `json:"name"`
			Lastname string `json:"lastname"`
			Zipcode  string `json:"zipcode"`
			City     string `json:"city"`
			Color    string `json:"color"`
		}{
			ID:       person.ID,
			Name:     person.Name,
			Lastname: person.Lastname,
			Zipcode:  person.Zipcode,
			City:     person.City,
			Color:    data.Color(person.Color).String(),
		})
	}
	return formated
}
