package models

import "fmt"

type Customer struct {
	ID            int64 `pg:",unique"`
	NameAnschrift string
}

func NewCustomer(nameAnschrift string) *Customer {
	return &Customer{
		NameAnschrift: nameAnschrift,
	}
}

func (c Customer) String() string {
	return fmt.Sprintf("Customer<%d %s>", c.ID, c.NameAnschrift)
}
