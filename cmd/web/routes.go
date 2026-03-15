package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	r.Get("/", app.home)
	r.Get("/snippet", app.showSnippet)
	r.Post("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	r.Handle("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(r)
}
