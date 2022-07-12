package models

import (
	"time"

	"github.com/Rhymond/go-money"
)

type LedgerEntry struct {
	ID       string
	Item     string
	Receiver string
	Amount   money.Money
	DueDate  time.Time
	PaidDate time.Time
	Notes    string
}
