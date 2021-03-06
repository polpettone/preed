package models

import (
	"github.com/Rhymond/go-money"
)

type BookingStatistics struct {
	Bookings []Booking
}


func (o BookingStatistics) Count() int {
	count := 0
	for _, b := range o.Bookings {
		if !b.Canceled {
			count += 1
		}
	}
	return count
}

func (o BookingStatistics) CountCancellation() int {
	count := 0
	for _, b := range o.Bookings {
		if b.Canceled {
			count += 1
		}
	}
	return count
}

func (o BookingStatistics) TotalAllocationDays() int {
	total := 0
	for _, b := range o.Bookings {
		if !b.Canceled {
			days := b.EndDate.YearDay() - b.StartDate.YearDay()
			total += days
		}
	}
	return total
}

func (o BookingStatistics) TotalIncome() *money.Money {
	total := money.New(0, "EUR")
	for _, b := range o.Bookings {
		if !b.Canceled {
			days := b.EndDate.YearDay() - b.StartDate.YearDay()
			p := b.PricePerDay.Multiply(int64(days))
			total, _ = total.Add(p)
		}
	}
	return total
}

func (o BookingStatistics) TotalProvision() *money.Money {
	total := money.New(0, "EUR")
	for _, b := range o.Bookings {
		if !b.Canceled {
			total, _ = total.Add(&b.Provision)
		}
	}
	return total
}
