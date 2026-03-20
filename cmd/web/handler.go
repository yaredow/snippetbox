package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"
	"github.com/yaredow/snippetbox/pkg/forms"
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

	errors := make(map[string]string)

	// validate the title
	if strings.TrimSpace(title) == "" {
		errors["title"] = "this field can not be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "This field is too long (maximum is 100 characters)"
	}

	// validate the content
	if strings.TrimSpace(content) == "" {
		errors["content"] = "This field cannot be blank"
	}

	if strings.TrimSpace(expires) == "" {
		errors["expires"] = "This field is invalid"
	} else if expires != "365" && expires != "7" && expires != "1" {
		errors["expires"] = "This field is invalid"
	}

	if len(errors) > 0 {
		app.render(w, r, "create.page.html", &templateData{
			Form: forms.New(nil),
		})
		fmt.Fprint(w, errors)
		return
	}

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
