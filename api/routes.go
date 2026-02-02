package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// Register the relevant methods, URL patterns and handler functions
	// for testing purposes only, not required
	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/persons", app.createPersonHandler)
	router.HandlerFunc(http.MethodGet, "/persons", app.listPersonsHandler)
	// catches /persons/:id, /persons/color/:id
	router.HandlerFunc(http.MethodGet, "/persons/*path", app.pathHandler)

	return app.recoverPanic(router)
}
