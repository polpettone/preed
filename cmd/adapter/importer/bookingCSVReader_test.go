package importer

import (
	"testing"
	"time"
)

func Test_should_read_csv(t *testing.T) {
	fileToRead := "testBookings.csv"
	rawBookings, err := ReadCSV(fileToRead)

	if err != nil {
		t.Errorf("No error expected but got %v", err)
	}

	if len(rawBookings) != 2 {
		t.Errorf("wanted %d but got %d", 1, len(rawBookings))
	}

	expectedNameAnschrift := "Familie X, 10119 Berlin, Teststraße 7"
	if rawBookings[0].NameAnschrift != expectedNameAnschrift {
		t.Errorf("wanted %s got %s", expectedNameAnschrift, rawBookings[0].NameAnschrift)
	}

	expectedItem := "Item"
	if rawBookings[0].Item != "Reeddach" {
		t.Errorf("wanted %s but got %s", expectedItem, rawBookings[0].Item)
	}

	expectedDays := 3
	if rawBookings[0].Days != expectedDays {
		t.Errorf("wanted %d but got %d", expectedDays, rawBookings[0].Days)
	}

	expectedMoneyAmount := int64(5500)
	expectedCurrency := "EUR"

	if rawBookings[0].PricePerDay.Amount() != expectedMoneyAmount {
		t.Errorf("%d not equal to %d", expectedMoneyAmount, rawBookings[0].PricePerDay.Amount())
	}

	if rawBookings[0].PricePerDay.Currency().Code != expectedCurrency {
		t.Errorf("%s not equal to %s", expectedCurrency, rawBookings[0].PricePerDay.Currency().Code)
	}

	if rawBookings[0].Provider != "provider0" {
		t.Errorf("%s not equal to %s", "MV", rawBookings[0].Provider)
	}

	if rawBookings[1].Provider != "provider1" {
		t.Errorf("%s not equal to %s", "Bellvilla", rawBookings[1].Provider)
	}

	if rawBookings[0].BookingNumber != "160181" {
		t.Errorf("%s not equal to %s", "160181", rawBookings[0].BookingNumber)
	}

	if rawBookings[0].InvoiceNumber != "" {
		t.Errorf("%s not equal to %s", "", rawBookings[0].InvoiceNumber)
	}

	expectedCleaningDate := time.Date(2020, 06, 11, 00, 00, 00, 00, time.Local)
	if rawBookings[0].CleaningDate.Time != expectedCleaningDate {
		t.Errorf("%s not equal to %s", expectedCleaningDate, rawBookings[0].CleaningDate)
	}

	expectedCashTransferDate := time.Date(2020, 06, 05, 00, 00, 00, 00, time.Local)
	if rawBookings[0].CashTransferDate.Time != expectedCashTransferDate {
		t.Errorf("%s not equal to %s", expectedCashTransferDate, rawBookings[0].CashTransferDate)
	}

	expectedCleaningCosts := int64(3000)
	if rawBookings[0].CleaningCost.Amount() != expectedCleaningCosts {
		t.Errorf("%d not equal to %d", expectedCleaningCosts, rawBookings[0].CleaningCost.Amount())
	}

}
