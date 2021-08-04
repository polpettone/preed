package web

import (
	"errors"
	"fmt"
	"github.com/polpettone/preed/cmd/core/models"
	"github.com/polpettone/preed/pkg/forms"
	"net/http"
	"net/url"
	"sort"
	"strconv"
)

func (app *WebApp) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "main.page.tmpl", &templateData{
	})
}

func (app *WebApp) showLedger(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "ledger.page.tmpl", &templateData{
	})
}

func (app *WebApp) showPriceTable(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "price-table.page.tmpl", &templateData{
	})
}

func (app *WebApp) showStatistics(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "statistics.page.tmpl", &templateData{
	})
}

func (app *WebApp) bookingOverview(w http.ResponseWriter, r *http.Request) {

	var bookings []models.Booking

	year := r.URL.Query().Get("year")
	if year == "" {
		foundBookings, err := app.BookingService.GetAllBookings()
		if err != nil {
			app.serverError(w, err)
			return
		}
		bookings = foundBookings
	} else {
		yearInt, err := strconv.Atoi(year)
		if err != nil {
			app.InfoLog.Printf("", err)
		} else {
			foundBookings, err := app.BookingService.GetAllBookingsForYear(yearInt)
			if err != nil {
				app.serverError(w, err)
				return
			}
			bookings = foundBookings
		}
	}

	sort.Slice(bookings, func(i, j int) bool {
		return bookings[i].StartDate.Before(bookings[j].StartDate)
	})

	app.render(w, r, "bookings.page.tmpl", &templateData{
		Bookings: bookings,
	})
}


func (app *WebApp) uploadFileForBooking(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		app.ErrorLog.Println("Error Retrieving the File")
		app.ErrorLog.Println(err)
		return
	}
	defer file.Close()
	app.InfoLog.Printf("Uploaded File: %+v\n", handler.Filename)
	app.InfoLog.Printf("File Size: %+v\n", handler.Size)
	app.InfoLog.Printf("MIME Header: %+v\n", handler.Header)
}

func (app *WebApp) uploadFileForBookingForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "upload.page.tmpl", &templateData{
	})
}

func (app *WebApp) deleteBooking(w http.ResponseWriter, r *http.Request) {

	err := parseForm(w, r, app)
	if err != nil {
		return
	}

	bookingID, err := getBookingIDFromForm(w, r, app)
	if err != nil {
		return
	}

	b, err := getBookingByID(w, *bookingID, app)
	if err != nil {
		return
	}

	err = app.BookingService.DeleteBooking(b)
	if err != nil {
		app.serverError(w, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/bookings?year=2021"), http.StatusSeeOther)
}

func (app *WebApp) deleteBookingForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	b, err := getBookingByID(w, int64(id), app)
	if err != nil {
		return
	}

	data := url.Values{}
	form := forms.New(data)
	form.Set("id", strconv.Itoa(int(b.ID)))

	app.render(w, r, "deleteBooking.page.tmpl", &templateData{
		Form: form,
	})
}

func (app *WebApp) cancelBooking(w http.ResponseWriter, r *http.Request) {

	err := parseForm(w, r, app)
	if err != nil {
		return
	}

	bookingID, err := getBookingIDFromForm(w, r, app)
	if err != nil {
		return
	}

	b, err := getBookingByID(w, *bookingID, app)
	if err != nil {
		return
	}

	err = app.BookingService.CancelBooking(b)
	if err != nil {
		app.serverError(w, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/bookings?year=2021"), http.StatusSeeOther)
}

func (app *WebApp) cancelBookingForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	b, err := getBookingByID(w, int64(id), app)
	if err != nil {
		return
	}

	data := url.Values{}
	form := forms.New(data)
	form.Set("id", strconv.Itoa(int(b.ID)))

	app.render(w, r, "cancelBooking.page.tmpl", &templateData{
		Form: form,
	})
}

func (app *WebApp) showBooking(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	b, err := getBookingByID(w, int64(id), app)
	if err != nil {
		return
	}

	app.render(w, r, "booking.page.tmpl", &templateData{
		Booking: *b,
	})
}

func (app *WebApp) createBookingForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "createBooking.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *WebApp) createBooking(w http.ResponseWriter, r *http.Request) {

	err := parseForm(w, r, app)
	if err != nil {
		return
	}

	form := forms.New(r.PostForm)
	form.Required("startDate",
		"endDate",
		"nameAnschrift",
		"provider",
		"pricePerDay",
		"provision",
		"cleaningPrice",
		"numberOfGuests")
	form.ValidDateFormat("startDate")
	form.ValidDateFormat("endDate")
	form.ValidMoneyFormat("pricePerDay")
	form.ValidMoneyFormat("provision")
	form.ValidMoneyFormat("cleaningPrice")
	form.MaxLength("notes", 100)
	form.IsNumber("numberOfGuests")

	if !form.Valid() {
		app.render(w, r, "createBooking.page.tmpl", &templateData{Form: form})
		return
	}

	app.InfoLog.Printf("Von: %s", form.Get("startDate"))
	app.InfoLog.Printf("Bis: %s", form.Get("endDate"))

	app.Session.Put(r, "flash", "Buchung erfolgreich angelegt")

	converter := NewFormToBookingConverter()

	booking, err := converter.convertFormToBooking(*form, *models.NewBooking())
	if err != nil {
		app.ErrorLog.Printf("%v", err)
	}

	err = app.BookingService.CreateBooking(booking)
	if err != nil {
		app.ErrorLog.Printf("%v", err)
	}

	http.Redirect(w, r, fmt.Sprintf("/bookings"), http.StatusSeeOther)
}

func (app *WebApp) editBookingForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	b, err := getBookingByID(w, int64(id), app)
	if err != nil {
		return
	}

	converter := NewFormToBookingConverter()
	form := converter.convertBookingToForm(*b)

	app.InfoLog.Printf("%v", form.Values)

	app.render(w, r, "editBooking.page.tmpl", &templateData{
		Form: &form,
	})
}

func (app *WebApp) editBooking(w http.ResponseWriter, r *http.Request) {

	err := parseForm(w, r, app)
	if err != nil {
		return
	}

	form := forms.New(r.PostForm)
	form.Required("startDate",
		"endDate",
		"nameAnschrift",
		"pricePerDay",
		"provision",
		"cleaningPrice",
		"numberOfGuests")
	form.ValidDateFormat("startDate")
	form.ValidDateFormat("endDate")
	form.ValidMoneyFormat("pricePerDay")
	form.ValidMoneyFormat("provision")
	form.ValidMoneyFormat("cleaningPrice")
	form.MaxLength("notes", 5000)
	form.IsNumber("numberOfGuests")

	if !form.Valid() {
		app.render(w, r, "editBooking.page.tmpl", &templateData{Form: form})
		return
	}

	app.InfoLog.Printf("%v", form)

	app.Session.Put(r, "flash", "Buchung erfolgreich geaendert")

	converter := NewFormToBookingConverter()

	bookingID, err := getBookingIDFromForm(w, r, app)
	if err != nil {
		return
	}

	b, err := getBookingByID(w, *bookingID, app)
	if err != nil {
		return
	}
	booking, err := converter.convertFormToBooking(*form, *b)
	if err != nil {
		app.ErrorLog.Printf("%v", err)
	}

	err = app.BookingService.SaveBooking(booking)
	if err != nil {
		app.ErrorLog.Printf("%v", err)
	}

	http.Redirect(w, r, fmt.Sprintf("/bookings"), http.StatusSeeOther)
}


var (
	ErrNoRecord           = errors.New("models: no matching record found")
)


func parseForm(w http.ResponseWriter, r *http.Request, app *WebApp) error {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return err
	}
	return nil
}

func getBookingByID(w http.ResponseWriter, bookingID int64, app *WebApp) (*models.Booking, error) {
	b, err := app.BookingService.GetBookingById(bookingID)
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

func getBookingIDFromForm(w http.ResponseWriter, r *http.Request, app *WebApp) (*int64, error) {
	form := forms.New(r.PostForm)
	var bookingId int64
	if form.Get("id") != "" {
		id, err := strconv.Atoi(form.Get("id"))
		if err != nil {
			app.notFound(w)
			return nil, err
		}
		bookingId = int64(id)
	}

	return &bookingId, nil
}









