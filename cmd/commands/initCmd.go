package commands

import (
	"github.com/polpettone/preed/cmd/adapter/persistence"
	"github.com/polpettone/preed/cmd/config"
	"github.com/spf13/cobra"
)

func (app *Application) NewInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "",
		Long:  ``,

		Run: func(command *cobra.Command, args []string) {
			app.handleInitCommand()
		},
	}
}

func (app *Application) handleInitCommand() {

	app.Logging.Stdout.Printf("init")

	repo := persistence.NewRepo(app.Logging,
		app.DBPort,
		app.DBUser,
		app.DBPassword,
		app.DBName)

		app.Logging.Stdout.Printf("Create Schema")
		err := repo.CreateSchema()
		if err != nil {
			app.Logging.Stdout.Printf("%v", err)
		}
}

func init() {
	logging := config.NewLogging()
	app := NewApplication(logging)
	initCmd := app.NewInitCmd()
	rootCmd.AddCommand(initCmd)
}
