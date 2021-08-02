package core

import "github.com/polpettone/preed/cmd/core/models"

func filterBookingsByYear(bookings []models.Booking, year int) []models.Booking {
	var filteredBookings []models.Booking
	for _,b := range bookings {
		if b.StartDate.Year() == year {
			filteredBookings = append(filteredBookings, b)
		}
	}
	return filteredBookings
}