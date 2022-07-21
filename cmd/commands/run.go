package commands

import (
	"github.com/polpettone/preed/cmd/adapter/persistence"
	"github.com/polpettone/preed/cmd/adapter/web/server"
	"github.com/polpettone/preed/cmd/config"
	"github.com/polpettone/preed/cmd/core"
	"github.com/spf13/cobra"
)

func (app *Application) NewRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "",
		Long:  ``,

		Run: func(command *cobra.Command, args []string) {
			app.handleRunCommand()
		},
	}
}

func (app *Application) handleRunCommand() {

	repo := persistence.NewRepo(app.Logging,
		app.DBPort,
		app.DBUser,
		app.DBPassword,
		app.DBName)

	bookingService := &core.BookingService{
		Repo: *repo,
	}

	ledgerService := core.NewLedgerService(*repo)

	app.Logging.InfoLog.Printf("run command")
	app.Logging.InfoLog.Printf("run migrations")

	err := repo.InitMigration()
	if err != nil {
		app.Logging.ErrorLog.Printf("%s", err)
	}

	err = repo.RunMigration()
	if err != nil {
		app.Logging.ErrorLog.Printf("%s", err)
	}

	server.StartWebAppServer(app.Logging, bookingService, ledgerService)
}

func init() {
	logging := config.NewLogging()
	app := NewApplication(logging)
	runCmd := app.NewRunCmd()
	rootCmd.AddCommand(runCmd)
}
