package web

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *WebApp) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	dynamicMiddleware := alice.New(app.Session.Enable, noSurf)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.Home))

	mux.Get("/bookings", dynamicMiddleware.ThenFunc(app.BookingOverview))

	mux.Get("/booking/create", dynamicMiddleware.ThenFunc(app.CreateBookingForm))
	mux.Post("/booking/create", dynamicMiddleware.ThenFunc(app.CreateBooking))

	mux.Get("/booking/edit/:id", dynamicMiddleware.ThenFunc(app.EditBookingForm))
	mux.Post("/booking/edit", dynamicMiddleware.ThenFunc(app.EditBooking))

	mux.Get("/booking/:id", dynamicMiddleware.ThenFunc(app.ShowBooking))

	mux.Get("/booking/delete/:id", dynamicMiddleware.ThenFunc(app.DeleteBookingForm))
	mux.Post("/booking/delete", dynamicMiddleware.ThenFunc(app.DeleteBooking))

	mux.Get("/booking/cancel/:id", dynamicMiddleware.ThenFunc(app.CancelBookingForm))
	mux.Post("/booking/cancel", dynamicMiddleware.ThenFunc(app.CancelBooking))

	mux.Get("/booking/reset-cancellation/:id", dynamicMiddleware.ThenFunc(app.ResetBookingCancellationForm))
	mux.Post("/booking/reset-cancellation", dynamicMiddleware.ThenFunc(app.ResetBookingCancellation))

	mux.Get("/upload", dynamicMiddleware.ThenFunc(app.UploadFileForBookingForm))
	mux.Post("/upload", dynamicMiddleware.ThenFunc(app.UploadFileForBooking))

	mux.Get("/ledger", dynamicMiddleware.ThenFunc(app.ShowLedger))

	mux.Get("/ledgerEntry/create", dynamicMiddleware.ThenFunc(app.CreateLedgerEntryForm))
	mux.Post("/ledgerEntry/create", dynamicMiddleware.ThenFunc(app.CreateLedgerEntry))

	mux.Get("/ledgerEntry/edit/:id", dynamicMiddleware.ThenFunc(app.EditLedgerEntryForm))
	mux.Post("/ledgerEntry/edit", dynamicMiddleware.ThenFunc(app.EditLedgerEntry))

	mux.Get("/price-table", dynamicMiddleware.ThenFunc(app.ShowPriceTable))
	mux.Get("/statistics", dynamicMiddleware.ThenFunc(app.ShowStatistics))

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
