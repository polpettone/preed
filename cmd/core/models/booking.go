package models

import (
	"fmt"
	"github.com/Rhymond/go-money"
	"math"
	"time"
)

func NewBooking() *Booking {
	return &Booking{
		Customer: &Customer{},

		PricePerDay:   *money.New(0, "EUR"),
		CleaningPrice: *money.New(0, "EUR"),
		Provision:     *money.New(0, "EUR"),

		Payment: NewPayment(),
	}
}

type Booking struct {
	ID         int64
	StartDate  time.Time
	EndDate    time.Time
	Notes      string
	Customer   *Customer `pg:"rel:has-one"`
	CustomerID int64
	Item       string
	Provider   string

	NumberOfGuests int

	PricePerDay   money.Money
	CleaningPrice money.Money
	Provision     money.Money

	Payment Payment

	CreatedAt time.Time
	ModifiedAt time.Time
}

func (b *Booking) String() string {
	return fmt.Sprintf("Booking<%d %v %v %s %s>", b.ID, b.StartDate, b.EndDate, b.Customer.NameAnschrift, b.Notes)
}

func (b *Booking) Days() int {
	return int(b.EndDate.Sub(b.StartDate).Hours() / 24)
}

func (b *Booking) Total() money.Money {
	return *b.PricePerDay.Multiply(int64(b.Days()))
}

func (b *Booking) ProvisionInPercentOfTotal() float64 {
	total := b.Total()
	if total.Amount() == 0 {
		return 0
	}

	result := ((float64(b.Provision.Amount())) / float64(total.Amount())) * 100

	return math.Round(result * 100) / 100
}
