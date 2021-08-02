package importer

import (
	"github.com/Rhymond/go-money"
	"github.com/polpettone/preed/cmd/core/models"
)

func ConvertRawBookingsToBooking(rawBookings []*RawBooking) ([]*models.Booking, error) {
	var bookings []*models.Booking
	for _, b := range rawBookings {
		booking, err := convertRawBookingToBooking(*b)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}
	return bookings, nil
}

func convertRawBookingToBooking(rawBooking RawBooking) (*models.Booking, error) {
	customer := &models.Customer{
		NameAnschrift: rawBooking.NameAnschrift,
	}

	booking := &models.Booking{
		StartDate:     rawBooking.StartDate.Time,
		EndDate:       rawBooking.EndDate.Time,
		Customer:      customer,
		Provider:      rawBooking.Provider,
		Item:          rawBooking.Item,
		PricePerDay:   rawBooking.PricePerDay.Money,
		CleaningPrice: rawBooking.CleaningIncome.Money,
		Provision:     rawBooking.Provision.Money,
		NumberOfGuests: 1,

		Payment: models.Payment{
			FinalAmount:   *money.New(0, "EUR"),
			PartialAmount: *money.New(0, "EUR"),
		},
	}
	return booking, nil
}
