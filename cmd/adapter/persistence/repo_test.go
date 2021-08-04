package persistence

import (
	"github.com/polpettone/preed/cmd/adapter"
	"testing"
)

func checkError(t *testing.T, err error, contextText string) {
	if err != nil {
		t.Errorf("%s %v", contextText, err)
	}
}

func Test_should_delete_booking(t *testing.T) {

	repo, postgres := startDB(t)
	defer stopDB(postgres, t)

	booking0 := adapter.TestBooking("n0", "c0", "p0", "i0")
	err := repo.SaveBooking(booking0)
	checkError(t, err, "Save booking error")
	foundBooking0, err := repo.FindBookingById(booking0.ID)
	if foundBooking0 == nil {
		t.Errorf("wanted %v but got nil", booking0)
	}

	err = repo.DeleteBooking(booking0)
	checkError(t, err, "Delete booking error")

	notExistingBooking, err := repo.FindBookingById(booking0.ID)
	if notExistingBooking != nil {
		t.Errorf("wanted nil but got %v", notExistingBooking)
	}

}

func Test_should_save_booking_with_same_nameAnschrift(t *testing.T) {

	repo, postgres := startDB(t)
	defer stopDB(postgres, t)

	booking0 := adapter.TestBooking("n0", "c0", "p0", "i0")
	booking1 := adapter.TestBooking("n1", "c0", "p1", "i1")

	err := repo.SaveBooking(booking0)
	checkError(t, err, "Save booking error")
	err = repo.SaveBooking(booking1)
	checkError(t, err, "Save booking error")
}

func Test_should_save_booking(t *testing.T) {

	repo, postgres := startDB(t)
	defer stopDB(postgres, t)

	// BOOKING 0 ------------

	booking0 := adapter.TestBooking("n0", "c0", "p0", "i0")

	err := repo.SaveBooking(booking0)
	checkError(t, err, "Save booking error")

	foundBooking, err := repo.FindBookingById(booking0.ID)
	checkError(t, err, "no error wanted but got")

	expectedNote := "n0"
	if foundBooking.Notes != expectedNote {
		t.Errorf("wanted %s but got %s", expectedNote, foundBooking.Notes)
	}

	expectedCustomerNameAnschrift := "c0"
	if foundBooking.Customer.NameAnschrift != expectedCustomerNameAnschrift {
		t.Errorf("wanted [%s] but got [%s]", expectedCustomerNameAnschrift, foundBooking.Customer.NameAnschrift)
	}
	// BOOKING 1 ------------

	booking1 := adapter.TestBooking("n1", "c1", "p1", "i1")

	err = repo.SaveBooking(booking1)
	checkError(t, err, "Save booking error")

	foundBooking1, err := repo.FindBookingById(booking1.ID)
	checkError(t, err, "no error wanted but got")

	expectedNote1 := "n1"
	if foundBooking1.Notes != expectedNote1 {
		t.Errorf("wanted %s but got %s", expectedNote1, foundBooking1.Notes)
	}

	expectedCustomerNameAnschrift1 := "c1"
	if foundBooking1.Customer.NameAnschrift != expectedCustomerNameAnschrift1 {
		t.Errorf("wanted [%s] but got [%s]", expectedCustomerNameAnschrift1, foundBooking1.Customer.NameAnschrift)
	}

	// BOOKING 2 ------------

	booking2 := adapter.TestBooking("n2", "c0", "p0", "i0")

	err = repo.SaveBooking(booking2)
	checkError(t, err, "Save booking error")

	foundBooking2, err := repo.FindBookingById(booking2.ID)
	checkError(t, err, "no error wanted but got")

	expectedNote2 := "n2"
	if foundBooking2.Notes != expectedNote2 {
		t.Errorf("wanted %s but got %s", expectedNote2, foundBooking2.Notes)
	}

	expectedCustomerNameAnschrift2 := "c0"
	if foundBooking2.Customer.NameAnschrift != expectedCustomerNameAnschrift2 {
		t.Errorf("wanted [%s] but got [%s]", expectedCustomerNameAnschrift2, foundBooking2.Customer.NameAnschrift)
	}

	customers, _ := repo.findAllCustomers()

	if len(customers) != 3 {
		t.Errorf("wanted 3 customers got %d", len(customers))
	}
}

func Test_findAllBookings(t *testing.T) {

	repo, postgres := startDB(t)
	defer stopDB(postgres, t)

	booking0 := adapter.TestBooking("n0", "c0", "p0", "i0")
	booking1 := adapter.TestBooking("n1", "c1", "p1", "i1")

	err := repo.SaveBooking(booking0)
	checkError(t, err, "no error wanted but got")

	err = repo.SaveBooking(booking1)
	checkError(t, err, "no error wanted but got")

	bookings, err := repo.FindAllBookings()
	checkError(t, err, "no error wanted but got")

	count := len(bookings)
	if count != 2 {
		t.Errorf("wanted %d got %d", 2, count)
	}

	customer0 := bookings[0].Customer
	if customer0 == nil {
		t.Errorf("Customer 0 should not be nil")
	}

	customer1 := bookings[1].Customer
	if customer1 == nil {
		t.Errorf("Customer 1 should not be nil")
	}
}

func Test_shouldUpdateBookings(t *testing.T) {
	repo, postgres := startDB(t)
	defer stopDB(postgres, t)

	booking0 := adapter.TestBooking("n0", "c0", "p0", "i0")
	err := repo.SaveBooking(booking0)
	checkError(t, err, "save booking")

	foundBooking0, err := repo.FindBookingById(booking0.ID)
	checkError(t, err, "no error wanted but got")

	notesText := "notes added"
	foundBooking0.Notes = notesText

	err = repo.SaveBooking(foundBooking0)
	checkError(t, err, "")

	againFoundBooking0, err := repo.FindBookingById(booking0.ID)
	checkError(t, err, "no error wanted but got")

	if againFoundBooking0.Notes != notesText {
		t.Errorf("wanted %s got %s", notesText, againFoundBooking0.Notes)
	}

	allBookings, err := repo.FindAllBookings()
	checkError(t, err, "no error wanted but got")

	if len(allBookings) != 1 {
		t.Errorf("wanted 1 booking got %d", len(allBookings))
	}
}

func Test_shouldUpdateBookingsCustomer(t *testing.T) {
	repo, postgres := startDB(t)
	defer stopDB(postgres, t)

	booking0 := adapter.TestBooking("n0", "c0", "p0", "i0")
	err := repo.SaveBooking(booking0)
	checkError(t, err, "")

	foundBooking0, err := repo.FindBookingById(booking0.ID)
	checkError(t, err, "no error wanted but got")
	foundCustomer, err := repo.FindCustomerById(foundBooking0.CustomerID)
	checkError(t, err, "no error wanted but got")

	foundCustomer.NameAnschrift = "changed"
	err = repo.SaveCustomer(foundCustomer)
	checkError(t, err, "no error wanted but got")

	againFoundBooking0, err := repo.FindBookingById(booking0.ID)
	checkError(t, err, "no error wanted but got")

	if againFoundBooking0.Customer.NameAnschrift != "changed" {
		t.Errorf("wanted %s got %s", "changed", againFoundBooking0.Customer.NameAnschrift)
	}
}

func Test_shouldUpdateBookingsCustomer1(t *testing.T) {
	repo, postgres := startDB(t)
	defer stopDB(postgres, t)

	booking0 := adapter.TestBooking("n0", "c0", "p0", "i0")
	err := repo.SaveBooking(booking0)
	checkError(t, err, "")

	foundBooking0, err := repo.FindBookingById(booking0.ID)
	checkError(t, err, "no error wanted but got")

	foundBooking0.Customer.NameAnschrift = "changed"
	err = repo.SaveBooking(foundBooking0)
	checkError(t, err, "no error wanted but got")

	againFoundBooking0, err := repo.FindBookingById(booking0.ID)
	checkError(t, err, "no error wanted but got")

	if againFoundBooking0.Customer.NameAnschrift != "changed" {
		t.Errorf("wanted %s got %s", "changed", againFoundBooking0.Customer.NameAnschrift)
	}
}
