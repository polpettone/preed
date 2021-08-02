package commands

import (
	"github.com/polpettone/preed/cmd/adapter/cli"
	"github.com/polpettone/preed/cmd/adapter/cli/models"
	"github.com/polpettone/preed/cmd/adapter/persistence"
	"github.com/polpettone/preed/cmd/config"
	"github.com/polpettone/preed/cmd/core"
	models2 "github.com/polpettone/preed/cmd/core/models"
	"github.com/polpettone/preed/pkg"
	"github.com/spf13/cobra"
)

func (app *Application) CreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "",
		Long:  ``,

		Run: func(command *cobra.Command, args []string) {
			app.handleCreateCommand(args)
		},
	}
}

func initBookingService(app *Application) *core.BookingService {
	repo := persistence.NewRepo(app.Logging,
		app.DBPort,
		app.DBUser,
		app.DBPassword,
		app.DBName)
	bookingService := &core.BookingService{
		Repo: *repo,
	}
	return bookingService
}

func (app *Application) handleCreateCommand(args []string) {

	bookingService := initBookingService(app)

	bookingView := models.NewBookingEditView()
	bookingYaml, err := cli.ConvertBookingViewToYaml(bookingView)
	if err != nil {
		app.Logging.ErrorLog.Printf("%s", err)
		return
	}

	createdBookingYaml, err := pkg.CaptureInputFromEditor(bookingYaml)
	if err != nil {
		app.Logging.ErrorLog.Printf("%s", err)
		return
	}

	createdBookingView, err := cli.ConvertYamlToBookingView(createdBookingYaml)
	if err != nil {
		app.Logging.ErrorLog.Printf("%s", err)
		return
	}

	booking := models2.NewBooking()
	createdBooking, err := createdBookingView.To(*booking)
	if err != nil {
		app.Logging.ErrorLog.Printf("%s", err)
		return
	}

	err = bookingService.CreateBooking(&createdBooking)
	if err != nil {
		app.Logging.ErrorLog.Printf("%s", err)
		return
	}
}

func init() {
	logging := config.NewLogging()
	app := NewApplication(logging)
	createCmd := app.CreateCmd()
	rootCmd.AddCommand(createCmd)
}
