package server

import (
	"bytes"
	"fmt"
	"github.com/polpettone/preed/cmd/adapter/web/server/templates"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/justinas/nosurf"
)

func (app *WebApp) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *WebApp) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *WebApp) NotFound(w http.ResponseWriter) {
	app.ClientError(w, http.StatusNotFound)
}

func (app *WebApp) addDefaultData(td *templates.TemplateData, r *http.Request) *templates.TemplateData {
	if td == nil {
		td = &templates.TemplateData{}
	}
	td.CSRFToken = nosurf.Token(r)
	td.CurrentYear = time.Now().Year()
	td.Flash = app.Session.PopString(r, "flash")
	td.IsAuthenticated = app.IsAuthenticated(r)
	return td
}

func (app *WebApp) Render(w http.ResponseWriter,
	r *http.Request,
	name string,
	td *templates.TemplateData) {
	ts, ok := app.TemplateCache[name]
	if !ok {
		app.ServerError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.ServerError(w, err)
	}

	buf.WriteTo(w)
}

func (app *WebApp) IsAuthenticated(r *http.Request) bool {

	isAuthenticated, ok := r.Context().Value(contextKeyIsAuthenticated).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}

func ParseForm(w http.ResponseWriter, r *http.Request, app *WebApp) error {
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return err
	}
	return nil
}
