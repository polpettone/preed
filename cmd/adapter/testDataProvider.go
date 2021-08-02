package adapter

import (
	"github.com/Rhymond/go-money"
	"github.com/polpettone/preed/cmd/core/models"
	"time"
)

func TestBooking(notes string, nameAnschrift string, providerName string, itemName string) *models.Booking {
	pricePerDay := *money.New(5500, "EUR")
	cleaningPrice := *money.New(2500, "EUR")
	provision := *money.New(1000, "EUR")

	booking := models.NewBooking()

	booking.Customer = models.NewCustomer(nameAnschrift)
	booking.Provider = providerName
	booking.Item = itemName
	booking.Notes = notes
	booking.PricePerDay = pricePerDay
	booking.CleaningPrice = cleaningPrice
	booking.Provision = provision

	return booking
}

func ParseDate(date, layout string) *time.Time  {
	parsed, err := time.Parse(layout, date)
	if err != nil {
		panic(err)
	}
	return &parsed
}
