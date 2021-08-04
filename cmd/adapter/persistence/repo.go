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

func (repo *Repo) RunMigration() error {
	db := pg.Connect(repo.DBOptions)
	defer db.Close()

	oldVersion, newVersion, err := migrations.Run(db, "version")

	if err != nil {
		return err
	}

	repo.Logging.Stdout.Printf("old %d", oldVersion)
	repo.Logging.Stdout.Printf("new %d", newVersion)

	return nil
}