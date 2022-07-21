package server

import (
	"github.com/polpettone/preed/cmd/adapter/web"
	"net/http"
)

func (app *WebApp) Home(w http.ResponseWriter, r *http.Request) {
	app.Render(w, r, "main.page.tmpl", &web.TemplateData{})
}
