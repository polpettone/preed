package models

import (
	"time"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
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

func NewLedgerEntry() *LedgerEntry {
	return &LedgerEntry{
		ID: uuid.New().String(),
	}
}
