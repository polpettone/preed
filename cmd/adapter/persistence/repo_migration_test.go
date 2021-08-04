package persistence

import (
	"github.com/polpettone/preed/cmd/adapter"
	"testing"
)

func Test_should_run_migration(t *testing.T) {
	repo, postgres := startDB(t)
	defer stopDB(postgres, t)

	booking0 := adapter.TestBooking("n0", "c0", "p0", "i0")
	err := repo.SaveBooking(booking0)
	checkError(t, err, "Save booking 0 error")

	err = repo.InitMigration()
	checkError(t, err, "init migration")

	err = repo.RunMigration()
	checkError(t, err, "run migration version")

	booking1 := adapter.TestBooking("n1", "c0", "p1", "i1")
	booking1.Canceled = true
	err = repo.SaveBooking(booking1)
	checkError(t, err, "Save booking 1 error")

	found, err := repo.FindBookingById(booking1.ID)
	checkError(t, err, "Find Booking By ID")

	if found.Canceled != true {
		t.Errorf("Wanted Canceled == true but it is false")
	}

}