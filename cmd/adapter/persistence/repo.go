package persistence

import (
	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/polpettone/preed/cmd/config"
	models2 "github.com/polpettone/preed/cmd/core/models"
)

type Repo struct {
	Logging   *config.Logging
	DBOptions *pg.Options
}

func NewRepo(logging *config.Logging, addr string, user string, password string, database string) *Repo {
	return &Repo{
		Logging: logging,
		DBOptions: &pg.Options{
			Addr:     addr,
			User:     user,
			Password: password,
			Database: database,
		},
	}
}

func (repo *Repo) CreateSchema() error {
	db := pg.Connect(repo.DBOptions)
	defer db.Close()
	models := []interface{}{
		(*models2.Customer)(nil),
		(*models2.Booking)(nil),
	}
	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp: false,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *Repo) SaveBooking(booking *models2.Booking) error {
	db := pg.Connect(repo.DBOptions)
	defer db.Close()

	err := repo.SaveCustomer(booking.Customer)
	if err != nil {
		return err
	}

	booking.CustomerID = booking.Customer.ID

	_, err = db.Model(booking).
		OnConflict("(id) DO UPDATE").
		Set("start_date = EXCLUDED.start_date").
		Set("end_date = EXCLUDED.end_date").
		Set("notes = EXCLUDED.notes").
		Set("price_per_day = EXCLUDED.price_per_day").
		Set("cleaning_price = EXCLUDED.cleaning_price").
		Set("provision = EXCLUDED.provision").
		Set("provider = EXCLUDED.provider").
		Set("item = EXCLUDED.item").
		Set("canceled = EXCLUDED.canceled").
		Insert()
	if err != nil {
		return err
	}

	return nil
}

func (repo *Repo) DeleteBooking(booking *models2.Booking) error {
	db := pg.Connect(repo.DBOptions)
	defer db.Close()
	_, err := db.Model(booking).Where("id=?0", booking.ID).Delete()
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repo) DeleteLedgerEntry(entry *models2.LedgerEntry) error {
	db := pg.Connect(repo.DBOptions)
	defer db.Close()
	_, err := db.Model(entry).Where("id=?0", entry.ID).Delete()
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repo) SaveLedgerEntry(ledgerEntry *models2.LedgerEntry) error {
	db := pg.Connect(repo.DBOptions)
	defer db.Close()
	result, err := db.Model(ledgerEntry).
		OnConflict("(id) DO UPDATE").
		Set("item = EXCLUDED.item").
		Set("receiver = EXCLUDED.receiver").
		Set("amount = EXCLUDED.amount").
		Set("due_date = EXCLUDED.due_date").
		Set("paid_date = EXCLUDED.paid_date").
		Set("notes = EXCLUDED.notes").
		Insert()
	if err != nil {
		return err
	}
	repo.Logging.DebugLog.Printf("%v", result)
	return nil
}

func (repo *Repo) FindLedgerEntryByID(id string) (*models2.LedgerEntry, error) {
	db := pg.Connect(repo.DBOptions)
	defer db.Close()

	ledgerEntry := &models2.LedgerEntry{ID: id}
	err := db.Model(ledgerEntry).
		WherePK().
		Select()

	if err != nil {
		return nil, err
	}

	return ledgerEntry, nil
}

func (repo *Repo) FindAllLedgerEntries() ([]models2.LedgerEntry, error) {
	db := pg.Connect(repo.DBOptions)
	defer db.Close()
	var entries []models2.LedgerEntry
	err := db.Model(&entries).Select()
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func (repo *Repo) SaveCustomer(customer *models2.Customer) error {
	db := pg.Connect(repo.DBOptions)
	defer db.Close()
	result, err := db.Model(customer).OnConflict("(id) DO UPDATE").Set("name_anschrift = EXCLUDED.name_anschrift").Insert()
	if err != nil {
		return err
	}
	repo.Logging.DebugLog.Printf("%v", result)
	return nil
}

func (repo *Repo) findAllCustomers() ([]models2.Customer, error) {
	db := pg.Connect(repo.DBOptions)
	defer db.Close()
	var customers []models2.Customer
	err := db.Model(&customers).Select()
	if err != nil {
		return nil, err
	}
	return customers, nil
}

func (repo *Repo) FindAllBookings() ([]models2.Booking, error) {
	db := pg.Connect(repo.DBOptions)
	defer db.Close()
	var bookings []models2.Booking
	err := db.Model(&bookings).Select()
	if err != nil {
		return nil, err
	}

	var enrichedBookings []models2.Booking
	for _, b := range bookings {

		booking := &models2.Booking{ID: b.ID}
		err = db.Model(booking).
			Relation("Customer").
			WherePK().
			Select()
		if err != nil {
			return nil, err
		}
		enrichedBookings = append(enrichedBookings, *booking)
	}

	return enrichedBookings, nil
}

func (repo *Repo) FindBookingById(id int64) (*models2.Booking, error) {
	db := pg.Connect(repo.DBOptions)
	defer db.Close()

	booking := &models2.Booking{ID: id}
	err := db.Model(booking).
		Relation("Customer").
		WherePK().
		Select()

	if err != nil {
		return nil, err
	}

	return booking, nil
}

func (repo *Repo) FindCustomerById(id int64) (*models2.Customer, error) {
	db := pg.Connect(repo.DBOptions)
	defer db.Close()

	customer := &models2.Customer{ID: id}
	err := db.Model(customer).
		WherePK().
		Select()

	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (repo *Repo) InitMigration() error {
	db := pg.Connect(repo.DBOptions)
	defer db.Close()

	oldVersion, newVersion, err := migrations.Run(db, "init")

	if err != nil {
		return err
	}

	repo.Logging.InfoLog.Printf("DB Init Migrations: old %d -> new %d", oldVersion, newVersion)

	return nil
}

func (repo *Repo) RunMigration() error {

	db := pg.Connect(repo.DBOptions)
	defer db.Close()

	oldVersion, newVersion, err := migrations.Run(db, "up")

	if err != nil {
		return err
	}

	repo.Logging.InfoLog.Printf("DB Migration: old %d -> new %d", oldVersion, newVersion)
	return nil
}
