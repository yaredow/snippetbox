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
	r.Post("/snippet/create", app.createSnippet)
	r.Get("/snippet/create", app.createSnippetForm)
	r.Get("/snippet/{id}", app.showSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	r.Handle("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(r)
}
