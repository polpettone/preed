package commands

import (
	"github.com/olekukonko/tablewriter"
	"github.com/polpettone/preed/cmd/adapter/persistence"
	"github.com/polpettone/preed/cmd/config"
	"github.com/polpettone/preed/cmd/core"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

func (app *Application) NewOverviewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "overview",
		Short: "",
		Long:  ``,

		Run: func(command *cobra.Command, args []string) {
			app.handleOverviewCommand()
		},
	}
}

func (app *Application) handleOverviewCommand() {

	repo := persistence.NewRepo(app.Logging,
		app.DBPort,
		app.DBUser,
		app.DBPassword,
		app.DBName)

	bookingService := &core.BookingService{
		Repo: *repo,
	}

	overview, err := bookingService.Overview()

	if err != nil {
		app.Logging.ErrorLog.Printf("%v", err)
		return
	}

	app.Logging.Stdout.Printf("Anzahl Buchungen: %d", overview.Count())
	app.Logging.Stdout.Printf("Belegungstage   : %d", overview.TotalAllocationDays())
	app.Logging.Stdout.Printf("Total Income (Days * Price per Day) : %s", overview.TotalIncome().Display())
	app.Logging.Stdout.Printf("Total Provision   : %s", overview.TotalProvision().Display())

	var data [][]string
	for _, booking := range overview.Bookings {
		total := booking.Total()
		data = append(data, []string{
			strconv.Itoa(int(booking.ID)),
			booking.StartDate.Format("02.01.2006"),
			booking.EndDate.Format("02.01.2006"),
			strconv.Itoa(booking.Days()),
			booking.Customer.NameAnschrift,
			booking.PricePerDay.Display(),
			booking.Provider,
			booking.Provision.Display(),
			booking.CleaningPrice.Display(),
			total.Display(),
			booking.Payment.PartialAmount.Display(),
			strconv.FormatBool(booking.Payment.PartialAmountPaid),
			booking.Payment.FinalAmount.Display(),
			strconv.FormatBool(booking.Payment.FinalAmountPaid),
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Nr",
		"Start",
		"Ende",
		"Tage",
		"Name Anschrift",
		"Preis pro Tag",
		"Provider",
		"Provision",
		"Reinigungskosten",
		"Gesamt",
		"Anzahlung",
		"Anzahlung bezahlt",
		"Endbetrag",
		"Endbetrag bezahlt",
	})
	for _, d := range data {
		table.Append(d)
	}

	table.Render()

}

func init() {
	logging := config.NewLogging()
	app := NewApplication(logging)
	overviewCmd := app.NewOverviewCmd()
	rootCmd.AddCommand(overviewCmd)
}
