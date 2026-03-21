package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/justinas/alice"
)

func adaptHandler(h http.Handler) http.HandlerFunc {
	return h.ServeHTTP
}

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable)

	r.Get("/", adaptHandler(dynamicMiddleware.ThenFunc(app.home)))
	r.Post("/snippet/create", adaptHandler(dynamicMiddleware.ThenFunc(app.createSnippetForm)))
	r.Get("/snippet/create", adaptHandler(dynamicMiddleware.ThenFunc(app.createSnippet)))
	r.Get("/snippet/{id}", adaptHandler(dynamicMiddleware.ThenFunc(app.showSnippet)))

	fileServer := http.FileServer(http.Dir("./ui/static"))
	r.Handle("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(r)
}
