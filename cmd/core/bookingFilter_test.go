package core

import (
	"github.com/polpettone/preed/cmd/adapter"
	"github.com/polpettone/preed/cmd/core/models"
	"reflect"
	"testing"
)

func Test_filterBookingsByYear(t *testing.T) {
	type args struct {
		bookings []models.Booking
		year     int
	}
	tests := []struct {
		name string
		args args
		want []models.Booking
	}{

		{
			name: "2021",
			args: struct {
				bookings []models.Booking
				year int
			} {
				bookings: createBookings(),
				year: 2021,
			},
			want: createWanted(),
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterBookingsByYear(tt.args.bookings, tt.args.year); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterBookingsByYear() = %v, want %v", got, tt.want)
			}
		})
	}
}


func createWanted() []models.Booking {
	b0 := models.NewBooking()
	b0.StartDate = *adapter.ParseDate("01.10.2021", "02.01.2006")

	var bookings []models.Booking
	bookings = append(bookings, *b0)
	return bookings
}

func createBookings() []models.Booking {
	b0 := models.NewBooking()
	b0.StartDate = *adapter.ParseDate("01.10.2021", "02.01.2006")

	b1 := models.NewBooking()
	b1.StartDate = *adapter.ParseDate("01.10.2020", "02.01.2006")

	b2 := models.NewBooking()
	b2.StartDate = *adapter.ParseDate("01.10.2019", "02.01.2006")

	var bookings []models.Booking
	bookings = append(bookings, *b0)
	bookings = append(bookings, *b1)
	bookings = append(bookings, *b2)

	return bookings
}
