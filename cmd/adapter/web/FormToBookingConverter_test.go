package web

import (
	"github.com/polpettone/preed/cmd/core/models"
	"github.com/polpettone/preed/pkg/forms"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestFormToBookingConverter_convertFormToBooking(t *testing.T) {
	type fields struct {
		TimeFormat string
	}
	type args struct {
		form forms.Form
		booking models.Booking
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Booking
		wantErr bool
	}{

		{
			name:    "Valid Dates",
			fields:  struct{ TimeFormat string }{TimeFormat: "02.01.2006"},
			args:    struct{
				form forms.Form
				booking models.Booking
			}{
				form: validDateForm(),
				booking: *models.NewBooking(),
			},

			want:    validDateBooking(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converter := FormToBookingConverter{
				TimeFormat: tt.fields.TimeFormat,
			}
			got, err := converter.convertFormToBooking(tt.args.form, tt.args.booking)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertFormToBooking() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertFormToBooking() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func validDateBooking() *models.Booking {
	b := models.NewBooking()
	b.StartDate, _ = time.Parse("02.01.2006", "01.01.2022")
	b.EndDate, _ = time.Parse("02.01.2006", "01.02.2022")

	b.Notes = "some notes"

	return b
}

func validDateForm() forms.Form {
	data := url.Values{}
	form := forms.New(data)
	form.Set("startDate", "01.01.2022")
	form.Set("endDate", "01.02.2022")
	form.Set("notes", "some notes")
	return *form
}
