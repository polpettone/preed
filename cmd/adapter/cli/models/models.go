package models

import (
	money2 "github.com/Rhymond/go-money"
	coreModels "github.com/polpettone/preed/cmd/core/models"
	"time"
)

type BookingEditView struct {
	StartDate time.Time
	EndDate   time.Time
	Notes     string

	NameAnschrift string
	ItemName      string
	ProviderName  string

	PricePerDay   MonetaryAmount
	CleaningPrice MonetaryAmount
	Provision     MonetaryAmount
}

func NewBookingEditView() BookingEditView {
	return BookingEditView{
		PricePerDay:   MonetaryAmount{}.New(0, "EUR"),
		CleaningPrice: MonetaryAmount{}.New(0, "EUR"),
		Provision:     MonetaryAmount{}.New(0, "EUR"),
	}
}

func (b BookingEditView) Equal(a BookingEditView) bool {

	if b.StartDate != a.StartDate {
		return false
	}
	if b.EndDate != a.EndDate {
		return false
	}
	if b.Notes != a.Notes {
		return false
	}
	if b.NameAnschrift != a.NameAnschrift {
		return false
	}
	if b.ProviderName != a.ProviderName {
		return false
	}
	if b.ItemName != a.ItemName {
		return false
	}

	if b.PricePerDay != a.PricePerDay {
		return false
	}
	if b.CleaningPrice != a.CleaningPrice {
		return false
	}
	if b.Provision != a.Provision {
		return false
	}

	return true
}

func (b *BookingEditView) From(booking coreModels.Booking) {
	b.StartDate = booking.StartDate
	b.EndDate = booking.EndDate
	b.Notes = booking.Notes
	b.NameAnschrift = booking.Customer.NameAnschrift
	b.ItemName = booking.Item
	b.ProviderName = booking.Provider

	b.PricePerDay.Currency = booking.PricePerDay.Currency().Code
	b.PricePerDay.Amount = booking.PricePerDay.Amount()

	b.CleaningPrice.Currency = booking.CleaningPrice.Currency().Code
	b.CleaningPrice.Amount = booking.CleaningPrice.Amount()

	b.Provision.Currency = booking.Provision.Currency().Code
	b.Provision.Amount = booking.Provision.Amount()
}

func (b *BookingEditView) To(booking coreModels.Booking) (coreModels.Booking, error) {
	booking.StartDate = b.StartDate
	booking.EndDate = b.EndDate
	booking.Notes = b.Notes
	booking.Customer.NameAnschrift = b.NameAnschrift
	booking.Item = b.ItemName
	booking.Provider = b.ProviderName

	booking.Provision = *money2.New(b.Provision.Amount, b.Provision.Currency)
	booking.PricePerDay = *money2.New(b.PricePerDay.Amount, b.PricePerDay.Currency)
	booking.CleaningPrice = *money2.New(b.CleaningPrice.Amount, b.CleaningPrice.Currency)

	return booking, nil
}

type MonetaryAmount struct {
	Amount   int64
	Currency string
}

func (m MonetaryAmount) New(amount int64, currencyCode string) MonetaryAmount {
	return MonetaryAmount{
		Amount:   amount,
		Currency: currencyCode,
	}
}
