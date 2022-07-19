package web

import (
	"net/http"
)

func (app *WebApp) Home(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "main.page.tmpl", &templateData{})
}

func (app *WebApp) ShowPriceTable(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "price-table.page.tmpl", &templateData{})
}
