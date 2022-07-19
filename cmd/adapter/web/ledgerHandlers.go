package web

import (
	"errors"
	"fmt"
	"net/http"

	converter2 "github.com/polpettone/preed/cmd/adapter/web/converter"
	"github.com/polpettone/preed/cmd/core/models"
	"github.com/polpettone/preed/pkg/forms"
)

func (app *WebApp) CreateLedgerEntryForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "createLedgerEntry.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *WebApp) CreateLedgerEntry(w http.ResponseWriter, r *http.Request) {

	app.InfoLog.Printf("%s", "Create Ledger Entry")

	err := parseForm(w, r, app)
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
		app.render(w, r, "createLedgerEntry.page.tmpl", &templateData{Form: form})
		return
	}

	ledgerEntry, err := converter2.ConvertFormToLedgerEntry(*form, *models.NewLedgerEntry())
	if err != nil {
		app.ErrorLog.Printf("%v", err)
	}

	app.InfoLog.Printf("%v", ledgerEntry)

	err = app.BookingService.Repo.SaveLedgerEntry(ledgerEntry)
	if err != nil {
		app.ErrorLog.Printf("%v", err)
	}

	app.Session.Put(r, "flash", "Ledger Entry erfolgreich angelegt")

	http.Redirect(w, r, fmt.Sprintf("/ledger"), http.StatusSeeOther)
}

func (app *WebApp) EditLedgerEntryForm(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":id")

	e, err := getLedgerEntryByID(w, id, app)
	if err != nil {
		return
	}

	form := converter2.ConvertLedgerEntryToForm(*e)

	app.render(w, r, "editLedgerEntry.page.tmpl", &templateData{
		Form: &form,
	})
}

func (app *WebApp) EditLedgerEntry(w http.ResponseWriter, r *http.Request) {

	err := parseForm(w, r, app)
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
		app.render(w, r, "editLedgerEntry.page.tmpl", &templateData{Form: form})
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

	err = app.BookingService.Repo.SaveLedgerEntry(entry)
	if err != nil {
		app.ErrorLog.Printf("%v", err)
	}

	http.Redirect(w, r, fmt.Sprintf("/ledger"), http.StatusSeeOther)
}

func (app *WebApp) ShowLedger(w http.ResponseWriter, r *http.Request) {

	ledgerEntries, err := app.BookingService.GetAllLedgerEntries()
	if err != nil {
		return
	}

	app.render(w, r, "ledgerEntries.page.tmpl", &templateData{
		LedgerEntries: ledgerEntries,
	})
}


func getLedgerEntryByID(w http.ResponseWriter, id string, app *WebApp) (*models.LedgerEntry, error) {
	b, err := app.BookingService.Repo.FindLedgerEntryByID(id)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return nil, err
	}
	return b, nil
}
