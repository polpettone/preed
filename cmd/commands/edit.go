package commands

import (
	"github.com/polpettone/preed/cmd/adapter/cli"
	"github.com/polpettone/preed/cmd/adapter/cli/models"
	"github.com/polpettone/preed/cmd/adapter/persistence"
	"github.com/polpettone/preed/cmd/config"
	"github.com/polpettone/preed/cmd/core"
	"github.com/polpettone/preed/pkg"
	"github.com/spf13/cobra"
	"strconv"
)

func (app *Application) EditCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "edit",
		Short: "",
		Long:  ``,

		Run: func(command *cobra.Command, args []string) {
			app.handleEditCommand(args)
		},
	}
}

func (app *Application) handleEditCommand(args []string) {

	if len(args) < 1 {
		app.Logging.Stdout.Printf("No Booking Id given")
		return
	}

	repo := persistence.NewRepo(app.Logging,
		app.DBPort,
		app.DBUser,
		app.DBPassword,
		app.DBName)

	bookingService := &core.BookingService{
		Repo: *repo,
	}

	bookingID, err := strconv.Atoi(args[0])

	if err != nil {
		app.Logging.Stdout.Printf("%s", err)
		return
	}

	booking, err := bookingService.GetBookingById(int64(bookingID))
	if err != nil {
		app.Logging.Stdout.Printf("Could not find Booking with ID: %d", bookingID)
		return
	}

	bookingView := models.BookingEditView{}
	bookingView.From(*booking)

	bookingYaml, err := cli.ConvertBookingViewToYaml(bookingView)
	if err != nil {
		app.Logging.ErrorLog.Printf("%s", err)
		return
	}

	editedBookingYaml, err := pkg.CaptureInputFromEditor(bookingYaml)
	if err != nil {
		app.Logging.ErrorLog.Printf("%s", err)
		return
	}

	editedBookingView, err := cli.ConvertYamlToBookingView(editedBookingYaml)
	if err != nil {
		app.Logging.ErrorLog.Printf("%s", err)
		return
	}

	editedBooking, err := editedBookingView.To(*booking)
	if err != nil {
		app.Logging.ErrorLog.Printf("%s", err)
		return
	}

	err = bookingService.SaveBooking(&editedBooking)
	if err != nil {
		app.Logging.ErrorLog.Printf("%s", err)
		return
	}
}

func init() {
	logging := config.NewLogging()
	app := NewApplication(logging)
	editCmd := app.EditCmd()
	rootCmd.AddCommand(editCmd)
}
