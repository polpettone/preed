package models

import (
	"github.com/Rhymond/go-money"
	"time"
)

func NewPayment() Payment {
	return Payment{
		FinalAmount:   *money.New(0, "EUR"),
		PartialAmount: *money.New(0, "EUR"),
	}
}

type Payment struct {
	FinalAmount   money.Money
	PartialAmount money.Money

	FinalAmountPaid   bool
	PartialAmountPaid bool

	FinalAmountPaidAt   time.Time
	PartialAmountPaidAt time.Time

	PartialAmountDueDate time.Time
	FinalAmountDueDate   time.Time
}
