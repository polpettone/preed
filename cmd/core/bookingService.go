package core

import (
	"github.com/polpettone/preed/cmd/adapter/persistence"
	"github.com/polpettone/preed/cmd/core/models"
	"time"
)

type BookingService struct {
	Repo persistence.Repo
}

func NewBookingService(repo persistence.Repo) BookingService {
	return BookingService{Repo : repo}
}

func (bookingService *BookingService) Overview() (*models.BookingStatistics, error) {
	bookings, err := bookingService.Repo.FindAllBookings()

	if err != nil {
		return nil, err
	}

	bookingOverView := &models.BookingStatistics{Bookings: bookings}
	return bookingOverView, nil
}

func (bookingService *BookingService) GetAllBookings() ([]models.Booking, error) {
	bookings, err := bookingService.Repo.FindAllBookings()
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

func (bookingService *BookingService) GetAllBookingsForYear(year int) ([]models.Booking, error) {
	bookings, err := bookingService.Repo.FindAllBookings()
	if err != nil {
		return nil, err
	}
	return filterBookingsByYear(bookings, year), nil
}



func (bookingService *BookingService) GetBookingById(id int64) (*models.Booking, error) {
	return bookingService.Repo.FindBookingById(id)
}

func (bookingService *BookingService) SaveBooking(booking *models.Booking) error {
	booking.ModifiedAt = time.Now()
	return bookingService.Repo.SaveBooking(booking)
}

func (bookingService *BookingService) CreateBooking(booking *models.Booking) error {
	booking.CreatedAt  = time.Now()
	booking.ModifiedAt = time.Now()
	return bookingService.Repo.SaveBooking(booking)
}

func (bookingService *BookingService) DeleteBooking(booking *models.Booking) error {
	return bookingService.Repo.DeleteBooking(booking)
}

func (bookingService *BookingService) CancelBooking(booking *models.Booking) error {
	booking.Cancel()
	return bookingService.Repo.SaveBooking(booking)
}

func (bookingService *BookingService) ResetCancellationOfBooking(booking *models.Booking) error {
	booking.ResetCancellation()
	return bookingService.Repo.SaveBooking(booking)
}
