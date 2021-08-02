package forms

import (
	"testing"
)

func TestParseDate(t *testing.T) {
	_, err := parseDate("01.10.2021", "02.01.2006")
	if err != nil {
		t.Errorf("%s", err)
	}
}
