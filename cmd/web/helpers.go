package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serveError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	app.errorLog.Println(trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]

	if !ok {
		app.serveError(w, fmt.Errorf("the template %s doesn't exist", name))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, ts)
	if err != nil {
		app.serveError(w, err)
		return
	}

	buf.WriteTo(w)
}
