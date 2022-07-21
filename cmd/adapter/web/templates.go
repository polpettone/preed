package web

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/Rhymond/go-money"
	coreModels "github.com/polpettone/preed/cmd/core/models"
	"github.com/polpettone/preed/pkg/forms"
)

type TemplateData struct {
	CSRFToken         string
	CurrentYear       int
	Flash             string
	Form              *forms.Form
	IsAuthenticated   bool
	Bookings          []coreModels.Booking
	Booking           coreModels.Booking
	BookingStatistics coreModels.BookingStatistics
	LedgerEntries     []coreModels.LedgerEntry
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02.01.2006")
}

func displayMoney(money money.Money) string {
	return money.Display()
}

var functions = template.FuncMap{
	"humanDate":    humanDate,
	"displayMoney": displayMoney,
}

func NewTemplateCache(dir string) (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
