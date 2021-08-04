package commands

import (
	"github.com/polpettone/preed/cmd/adapter/importer"
	"github.com/polpettone/preed/cmd/adapter/persistence"
	"github.com/polpettone/preed/cmd/config"
	"github.com/polpettone/preed/cmd/core"
	"github.com/spf13/cobra"
)

func (app *Application) NewImportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "import",
		Short: "",
		Long:  ``,

		Run: func(command *cobra.Command, args []string) {
			app.handleImportCommand(command)
		},
	}
}

func (app *Application) handleImportCommand(cobraCommand *cobra.Command) {
	initial, _ := cobraCommand.Flags().GetBool("initial")

	app.Logging.InfoLog.Printf("import")

	repo := persistence.NewRepo(app.Logging,
		app.DBPort,
		app.DBUser,
		app.DBPassword,
		app.DBName)

	if initial {
		app.Logging.InfoLog.Printf("Create Schema")
		err := repo.CreateSchema()
		if err != nil {
			app.Logging.ErrorLog.Printf("%v", err)
		}
	}

	pathToCSV := "add path to csv"
	rawBookings, err := importer.ReadCSV(pathToCSV)
	if err != nil {
		app.Logging.ErrorLog.Printf("%v", err)
	}
	bookingService := core.NewBookingService(*repo)
	bookings, _ := importer.ConvertRawBookingsToBooking(rawBookings)
	for _, b := range bookings {
		err = bookingService.CreateBooking(b)
		if err != nil {
			app.Logging.ErrorLog.Printf("%v", err)
		}
	}
}

func init() {
	logging := config.NewLogging()
	app := NewApplication(logging)
	importCmd := app.NewImportCmd()

	importCmd.Flags().BoolP(
		"initial",
		"i",
		false,
		"Indicates an initial csv import and will try to create the db schema",
	)

	rootCmd.AddCommand(importCmd)
}
