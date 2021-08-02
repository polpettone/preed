package commands

import "github.com/polpettone/preed/cmd/config"

type Application struct {
	Logging    *config.Logging
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func NewApplication(logging *config.Logging) *Application {
	return &Application{
		Logging:    logging,
		DBPort:     ":5432",
		DBUser:     "preed",
		DBPassword: "preed",
		DBName:     "preed",
	}
}
