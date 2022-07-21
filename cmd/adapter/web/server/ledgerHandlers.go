package server

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/polpettone/preed/cmd/adapter/web/server/templates"

	converter2 "github.com/polpettone/preed/cmd/adapter/web/converter"
	"github.com/polpettone/preed/cmd/core/models"
	"github.com/polpettone/preed/pkg/forms"
)

func (app *WebApp) CreateLedgerEntryForm(w http.ResponseWriter, r *http.Request) {
	app.Render(w, r, "createLedgerEntry.page.tmpl", &templates.TemplateData{
		Form: forms.New(nil),
	})
}

func (app *WebApp) CreateLedgerEntry(w http.ResponseWriter, r *http.Request) {

	app.InfoLog.Printf("%s", "Create Ledger Entry")

	err := ParseForm(w, r, app)
	if err != nil {
		return
	}

	form := forms.New(r.PostForm)
	form.Required("item",
		"receiver",
		"amount")

	form.ValidDateFormat("dueDate")
	form.ValidDateFormat("paidDate")

	form.ValidMoneyFormat("amount")

	form.MaxLength("notes", 100)

	if !form.Valid() {
		app.Render(w, r, "createLedgerEntry.page.tmpl", &templates.TemplateData{Form: form})
		return
	}

	ledgerEntry, err := converter2.ConvertFormToLedgerEntry(*form, *models.NewLedgerEntry())
	if err != nil {
		app.ErrorLog.Printf("%v", err)
	}

	app.InfoLog.Printf("%v", ledgerEntry)

	err = app.LedgerService.SaveLedgerEntry(ledgerEntry)
	if err != nil {
		app.ErrorLog.Printf("%v", err)
	}

	app.Session.Put(r, "flash", "Ledger Entry erfolgreich angelegt")

	http.Redirect(w, r, fmt.Sprintf("/ledger"), http.StatusSeeOther)
}

func (app *WebApp) DeleteLedgerEntryForm(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":id")

	app.InfoLog.Printf("DeleteLedgerEntry Form ID: %s", id)

	e, err := getLedgerEntryByID(w, id, app)
	if err != nil {
		return
	}

	app.InfoLog.Printf("Found DeleteLedgerEntry to delete : %v", e)
	data := url.Values{}
	form := forms.New(data)
	form.Set("id", e.ID)

	app.Render(w, r, "deleteLedgerEntry.page.tmpl", &templates.TemplateData{
		Form: form,
	})
}

func (app *WebApp) DeleteLedgerEntry(w http.ResponseWriter, r *http.Request) {

	err := ParseForm(w, r, app)
	if err != nil {
		return
	}

	form := forms.New(r.PostForm)
	id := form.Get("id")

	app.InfoLog.Printf("DeleteLedgerEntry ID: %s", id)

	e, err := getLedgerEntryByID(w, id, app)
	if err != nil {
		return
	}

	err = app.LedgerService.DeleteLedgerEntry(e)
	if err != nil {
		app.ServerError(w, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/ledger"), http.StatusSeeOther)
}

func (app *WebApp) EditLedgerEntryForm(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":id")

	e, err := getLedgerEntryByID(w, id, app)
	if err != nil {
		return
	}

	form := converter2.ConvertLedgerEntryToForm(*e)

	app.Render(w, r, "editLedgerEntry.page.tmpl", &templates.TemplateData{
		Form: &form,
	})
}

func (app *WebApp) EditLedgerEntry(w http.ResponseWriter, r *http.Request) {

	err := ParseForm(w, r, app)
	if err != nil {
		return
	}

	form := forms.New(r.PostForm)
	form.Required("item",
		"receiver",
		"amount")

	form.ValidDateFormat("dueDate")
	form.ValidDateFormat("paidDate")
	form.ValidMoneyFormat("amount")
	form.MaxLength("notes", 100)

	if !form.Valid() {
		app.Render(w, r, "editLedgerEntry.page.tmpl", &templates.TemplateData{Form: form})
		return
	}

	app.InfoLog.Printf("%v", form)

	app.Session.Put(r, "flash", "Ledger Entry erfolgreich geaendert")

	id := form.Get("id")

	entry, err := getLedgerEntryByID(w, id, app)
	if err != nil {
		app.ErrorLog.Printf("%v", err)
		return
	}

	entry, err = converter2.ConvertFormToLedgerEntry(*form, *entry)
	if err != nil {
		app.ErrorLog.Printf("%v", err)
	}

	err = app.LedgerService.SaveLedgerEntry(entry)
	if err != nil {
		app.ErrorLog.Printf("%v", err)
	}

	http.Redirect(w, r, fmt.Sprintf("/ledger"), http.StatusSeeOther)
}

func (app *WebApp) ShowLedger(w http.ResponseWriter, r *http.Request) {

	ledgerEntries, err := app.LedgerService.GetAllLedgerEntries()
	if err != nil {
		return
	}

	app.Render(w, r, "ledgerEntries.page.tmpl", &templates.TemplateData{
		LedgerEntries: ledgerEntries,
	})
}

func getLedgerEntryByID(w http.ResponseWriter, id string, app *WebApp) (*models.LedgerEntry, error) {
	b, err := app.LedgerService.FindLedgerEntryById(id)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			app.NotFound(w)
		} else {
			app.ServerError(w, err)
		}
		return nil, err
	}
	return b, nil
}
