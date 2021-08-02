package commands

import (
	"github.com/polpettone/preed/cmd/adapter/persistence"
	"github.com/polpettone/preed/cmd/adapter/web"
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

	app.Logging.InfoLog.Printf("run command")

	web.StartWebAppServer(app.Logging, bookingService)
}

func init() {
	logging := config.NewLogging()
	app := NewApplication(logging)
	runCmd := app.NewRunCmd()
	rootCmd.AddCommand(runCmd)
}
