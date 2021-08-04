package persistence

import (
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/polpettone/preed/cmd/config"
	"testing"
)

const addr = ":16432"
const user = "kaufen"
const password = "kaufen"
const db = "preed"

func startDB(t *testing.T) (*Repo, *embeddedpostgres.EmbeddedPostgres) {
	logging := config.NewLogging()
	repo := NewRepo(logging, addr, user, password, db)
	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Username(repo.DBOptions.User).
		Password(repo.DBOptions.Password).
		Database(repo.DBOptions.Database).
		Port(16432))

	err := postgres.Start()
	if err != nil {
		t.Errorf("Error %v start embedded db", err)
	}

	err = repo.CreateSchema()
	if err != nil {
		t.Errorf("Create schema error %v", err)
	}

	return repo, postgres
}

func stopDB(postgres *embeddedpostgres.EmbeddedPostgres, t *testing.T) {
	err := postgres.Stop()
	if err != nil {
		t.Errorf("Error %v stop embedded db", err)
	}
}

func checkError(t *testing.T, err error, contextText string) {
	if err != nil {
		t.Errorf("%s %v", contextText, err)
	}
}
