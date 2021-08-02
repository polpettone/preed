package cli

import (
	"github.com/polpettone/preed/cmd/adapter/cli/models"
	"testing"
	"time"
)

func Test_should_convert_bookingView_to_yaml(t *testing.T) {

	bookingView := createBookingView()
	yaml, err := ConvertBookingViewToYaml(bookingView)
	if err != nil {
		t.Errorf("%v", err)
	}

	result, err := ConvertYamlToBookingView(yaml)
	if err != nil {
		t.Errorf("%v", err)
	}

	if !bookingView.Equal(*result) {
		t.Errorf("not equal")
	}
}

func Test_should_convert_yaml_to_bookingView(t *testing.T) {
	bookingView, err := ConvertYamlToBookingView(bookingViewYaml)
	if err != nil {
		t.Errorf("%v", err)
	}

	expectedBookingView := createBookingView()

	if !bookingView.Equal(expectedBookingView) {
		t.Errorf("not equal")
	}
}

func createBookingView() models.BookingEditView {

	layout := "2006-01-02T15:04:05Z"

	startDate, _ := time.Parse(layout, "2001-01-01T00:00:00Z")
	endDate, _ := time.Parse(layout, "2001-01-09T00:00:00Z")

	return models.BookingEditView{
		StartDate:     startDate,
		EndDate:       endDate,
		Notes:         "notes 0",
		NameAnschrift: "customer 0",
		ItemName:      "item 0",
		ProviderName:  "provider 0",
		PricePerDay:   models.MonetaryAmount{}.New(2500, "EUR"),
		CleaningPrice: models.MonetaryAmount{}.New(2500, "EUR"),
		Provision:     models.MonetaryAmount{}.New(2500, "EUR"),
	}
}

var bookingViewYaml = `
startdate: 2001-01-01T00:00:00Z
enddate: 2001-01-09T00:00:00Z
notes: notes 0
nameanschrift: customer 0
itemname: item 0
providername: provider 0
priceperday:
  amount: 2500
  currency: "EUR"
cleaningprice:
  amount: 2500
  currency: "EUR"
provision:
  amount: 2500
  currency: "EUR"
`
