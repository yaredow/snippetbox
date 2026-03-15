package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/yaredow/snippetbox/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	s, err := app.snippets.Latest()
	if err != nil {
		app.serveError(w, err)
		return
	}

	data := &templateData{Snippets: s}

	app.render(w, r, "home.page.html", data)
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err == models.ErrorNoRecord {
		app.notFound(w)
	} else if err != nil {
		app.serveError(w, err)
		return
	}

	data := &templateData{Snippet: s}

	app.render(w, r, "show.page.html", data)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serveError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.html", nil)
}
