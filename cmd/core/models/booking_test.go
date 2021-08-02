package models

import (
	"github.com/Rhymond/go-money"
	"reflect"
	"testing"
	"time"
)

func TestBooking(t *testing.T) {
	b0 := NewBooking()

	startDate, _ := time.Parse("02.01.2006", "01.01.2020")
	b0.StartDate = startDate

	endDate, _ := time.Parse("02.01.2006", "02.01.2020")
	b0.EndDate = endDate

	b1 := NewBooking()

	startDate1, _ := time.Parse("02.01.2006", "01.01.2020")
	b1.StartDate = startDate1

	endDate1, _ := time.Parse("02.01.2006", "02.01.2020")
	b1.EndDate = endDate1

	if !reflect.DeepEqual(b0, b1) {
		t.Errorf("%v should be equal %v", b0, b1)
	}
}

func TestTotal(t *testing.T) {
	b := NewBooking()

	b.PricePerDay = *money.New(5000, "EUR")
	startDate, _ := time.Parse("02.01.2006", "01.10.2020")
	endDate, _ := time.Parse("02.01.2006", "03.10.2020")

	b.StartDate = startDate
	b.EndDate = endDate

	total := b.Total()
	days := b.Days()

	if days != 2 {
		t.Errorf("wanted %d got %d days", 2, days)
	}

	var wanted int64
	wanted = 10000

	if total.Amount() != wanted {
		t.Errorf("wanted %d got %d", wanted, total.Amount())
	}
}

func TestEqualMoney(t *testing.T) {
	if !reflect.DeepEqual(*money.New(0, "EUR"), *money.New(0, "EUR")) {
		t.Errorf("Money in go sucks")
	}
}

func TestBooking_ProvisionInPercentOfTotal(t *testing.T) {

	tests := []struct{
		name string
		booking Booking
		wanted float64
	} {

		{
			name: "correct calculation",
			booking: createBooking(10, 10000, 6000),
			wanted: 16.67,
		},

		{
			name: "total is zero percentage should be zero",
			booking: createBooking(10, 10000, 0),
			wanted: 0,
		},

		{
			name: "provision is zero percentage should be zero",
			booking: createBooking(10, 0, 6000),
			wanted: 0,
		},

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			percentage := tt.booking.ProvisionInPercentOfTotal()
			if percentage != tt.wanted {
				t.Errorf("want %f got %f", tt.wanted, percentage)
			}
		})
	}

}

func createBooking(days int, provision int64, pricePerDay int64) Booking  {
	booking := *NewBooking()
	booking.Provision = *money.New(provision, "EUR")
	booking.PricePerDay = *money.New(pricePerDay, "EUR")
	booking.StartDate.AddDate(2020, 1, 1)
	booking.EndDate = booking.StartDate.AddDate(0,0,days)
	return booking
}