package core

import (
	"github.com/Rhymond/go-money"
	"github.com/polpettone/preed/cmd/core/models"
	"testing"
	"time"
)

func Test_totalAllocationDays(t *testing.T) {
	b1 := models.Booking{
		StartDate: date(2020, 1, 1),
		EndDate:   date(2020, 1, 3),
	}

	b2 := models.Booking{
		StartDate: date(2020, 2, 1),
		EndDate:   date(2020, 2, 3),
	}

	bookings := []models.Booking{b1, b2}
	overview := models.BookingStatistics{Bookings: bookings}

	totalAllocationDays := overview.TotalAllocationDays()

	if totalAllocationDays != 4 {
		t.Errorf("Wanted %d got %d", 4, totalAllocationDays)
	}
}

func date(year int, month int, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
}

func Test_totalIncome(t *testing.T) {

	b1 := booking(date(2020, 1, 1), date(2020, 1, 3), *money.New(2000, "EUR"))
	b2 := booking(date(2020, 5, 1), date(2020, 5, 3), *money.New(5000, "EUR"))

	bookings := []models.Booking{b1, b2}
	overview := models.BookingStatistics{Bookings: bookings}

	totalIncome := overview.TotalIncome()

	equals, _ := totalIncome.Equals(money.New(14000, "EUR"))

	if !equals {
		t.Errorf("error")
	}

}

func booking(
	startDate time.Time,
	endDate time.Time,
	pricePerDay money.Money) models.Booking {
	return models.Booking{
		PricePerDay: pricePerDay,
		StartDate:   startDate,
		EndDate:     endDate,
	}
}
