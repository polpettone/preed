package web

import (
	"crypto/tls"
	"github.com/golangcollege/sessions"
	"github.com/polpettone/preed/cmd/config"
	"github.com/polpettone/preed/cmd/core"
	"html/template"
	"log"
	"net/http"
	"time"
)

type contextKey string

const contextKeyIsAuthenticated = contextKey("isAuthenticated")

type WebApp struct {
	Session        *sessions.Session
	TemplateCache  map[string]*template.Template
	ErrorLog       *log.Logger
	InfoLog        *log.Logger
	BookingService *core.BookingService
}

func StartWebAppServer(logging *config.Logging, bookingService *core.BookingService) {

	secret := "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge"

	session := sessions.New([]byte(secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		logging.ErrorLog.Fatal(err)
	}

	webApp := &WebApp{
		ErrorLog:       logging.ErrorLog,
		InfoLog:        logging.InfoLog,
		Session:        session,
		TemplateCache:  templateCache,
		BookingService: bookingService,
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         ":4000",
		ErrorLog:     webApp.ErrorLog,
		Handler:      webApp.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logging.InfoLog.Printf("Starting server on %s", ":4000")
	//TODO: change tls location and make configurable
	err = srv.ListenAndServeTLS("./cmd/adapter/web/tls/cert.pem", "./cmd/adapter/web/tls/key.pem")
	logging.ErrorLog.Fatal(err)
}
