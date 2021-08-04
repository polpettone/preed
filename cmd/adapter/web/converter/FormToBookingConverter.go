package converter

import (
	money2 "github.com/Rhymond/go-money"
	"github.com/polpettone/preed/cmd/core/models"
	"github.com/polpettone/preed/pkg/forms"
	"net/url"
	"strconv"
	"time"
)

func NewFormToBookingConverter() FormToBookingConverter {
	return FormToBookingConverter{
		TimeFormat: "02.01.2006",
	}
}

type FormToBookingConverter struct {
	TimeFormat string
}

func (converter FormToBookingConverter) ConvertFormToBooking(form forms.Form, booking models.Booking) (*models.Booking, error) {


	booking.Customer.NameAnschrift = form.Get("nameAnschrift")

	pricePerDayInEuro, _ := strconv.Atoi(form.Get("pricePerDay"))
	cleaningPriceInEuro, _ := strconv.Atoi(form.Get("cleaningPrice"))
	provisionInEuro, _ := strconv.Atoi(form.Get("provision"))
	booking.PricePerDay = *money2.New(int64(pricePerDayInEuro * 100), "EUR")
	booking.Provision= *money2.New(int64(provisionInEuro * 100), "EUR")
	booking.CleaningPrice= *money2.New(int64(cleaningPriceInEuro * 100), "EUR")

	startDate, err := time.Parse(converter.TimeFormat, form.Get("startDate"))
	if err != nil {
		return nil, err
	}
	booking.StartDate = startDate

	endDate, err := time.Parse(converter.TimeFormat, form.Get("endDate"))
	if err != nil {
		return nil, err
	}
	booking.EndDate = endDate

	booking.Notes = form.Get("notes")
	booking.Provider = form.Get("provider")

	numberOfGuests, _ := strconv.Atoi(form.Get("numberOfGuests"))
	booking.NumberOfGuests = numberOfGuests

	return &booking, nil
}

func (converter FormToBookingConverter) ConvertBookingToForm(booking models.Booking) forms.Form {
	data := url.Values{}
	form := forms.New(data)

	form.Set("id", strconv.Itoa(int(booking.ID)))
	form.Set("startDate", booking.StartDate.Format(converter.TimeFormat))
	form.Set("endDate", booking.EndDate.Format(converter.TimeFormat))
	form.Set("notes", booking.Notes)
	form.Set("provider", booking.Provider)

	form.Set("nameAnschrift", booking.Customer.NameAnschrift)

	form.Set("pricePerDay", strconv.Itoa(int(booking.PricePerDay.Amount()/ 100)))
	form.Set("cleaningPrice", strconv.Itoa(int(booking.CleaningPrice.Amount() / 100)))
	form.Set("provision", strconv.Itoa(int(booking.Provision.Amount() / 100)))
	form.Set("numberOfGuests", strconv.Itoa(booking.NumberOfGuests))

	return *form
}
