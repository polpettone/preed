package cli

import (
	"github.com/polpettone/preed/cmd/adapter/cli/models"
	"gopkg.in/yaml.v2"
)

func ConvertBookingViewToYaml(bookingView models.BookingEditView) (string, error) {
	bookingYaml, err := yaml.Marshal(bookingView)
	if err != nil {
		return "", err
	}
	return string(bookingYaml), nil
}

func ConvertYamlToBookingView(bookingYaml string) (*models.BookingEditView, error) {
	bookingView := models.BookingEditView{}
	err := yaml.Unmarshal([]byte(bookingYaml), &bookingView)
	if err != nil {
		return nil, err
	}
	return &bookingView, nil
}
