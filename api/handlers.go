package main

import (
	"errors"
	"net/http"
	"strings"

	"assecor.assessment.test/internal/data"
	"assecor.assessment.test/internal/validator"
)

// "GET /healthcheck" endpoint
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "available",
		"version": version,
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// "POST /persons" endpoint
func (app *application) createPersonHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Lastname string `json:"lastname"`
		Zipcode  string `json:"zipcode"`
		City     string `json:"city"`
		Color    int    `json:"color"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()

	person := data.Person{
		Name:     input.Name,
		Lastname: input.Lastname,
		Zipcode:  input.Zipcode,
		City:     input.City,
		Color:    input.Color,
	}

	if data.ValidatePerson(v, &person); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Persons.Insert(&person)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, person, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// "GET /persons" endpoint
func (app *application) listPersonsHandler(w http.ResponseWriter, r *http.Request) {
	persons, err := app.models.Persons.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, app.formatPersonArray(persons), nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// "GET /persons/*path" endpoint
func (app *application) pathHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	count := len(parts)
	if count == 3 {
		app.showPersonHandler(w, r, parts[2])
	} else if count == 4 && parts[2] == "color" {
		app.listPersonsByFavoriteColorHandler(w, r, parts[3])
	} else {
		app.notFoundResponse(w, r)
	}
}

// "GET /persons/:id" endpoint
func (app *application) showPersonHandler(w http.ResponseWriter, r *http.Request, param string) {
	id, err := app.readIDParam(param)
	if err != nil {
		app.failedValidationResponse(w, r, map[string]string{"personID": err.Error()})
		return
	}
	person, err := app.models.Persons.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, app.formatPerson(person), nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// "GET /persons/color/:id" endpoint
func (app *application) listPersonsByFavoriteColorHandler(w http.ResponseWriter, r *http.Request, param string) {
	id, err := app.readIDParam(param)
	if err != nil {
		app.failedValidationResponse(w, r, map[string]string{"colorID": err.Error()})
		return
	}
	if id >= int64(data.LastColorIndex) {
		app.failedValidationResponse(w, r, map[string]string{"colorID": "out of range"})
		return
	}
	persons, err := app.models.Persons.GetAllByColor(data.Color(id))
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, app.formatPersonArray(persons), nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
