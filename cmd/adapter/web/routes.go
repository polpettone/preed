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
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))

	mux.Get("/bookings", dynamicMiddleware.ThenFunc(app.bookingOverview))
	mux.Get("/booking/create", dynamicMiddleware.ThenFunc(app.createBookingForm))
	mux.Post("/booking/create", dynamicMiddleware.ThenFunc(app.createBooking))
	mux.Get("/booking/edit/:id", dynamicMiddleware.ThenFunc(app.editBookingForm))
	mux.Post("/booking/edit", dynamicMiddleware.ThenFunc(app.editBooking))
	mux.Get("/booking/:id", dynamicMiddleware.ThenFunc(app.showBooking))

	mux.Get("/booking/delete/:id", dynamicMiddleware.ThenFunc(app.deleteBookingForm))
	mux.Post("/booking/delete", dynamicMiddleware.ThenFunc(app.deleteBooking))

	mux.Get("/upload", dynamicMiddleware.ThenFunc(app.uploadFileForBookingForm))
	mux.Post("/upload", dynamicMiddleware.ThenFunc(app.uploadFileForBooking))


	mux.Get("/ledger", dynamicMiddleware.ThenFunc(app.showLedger))
	mux.Get("/price-table", dynamicMiddleware.ThenFunc(app.showPriceTable))
	mux.Get("/statistics", dynamicMiddleware.ThenFunc(app.showStatistics))



	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
