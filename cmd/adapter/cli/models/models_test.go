package models

import "testing"

func Test_MonetaryAmount_should_be_equal(t *testing.T) {
	moneyA := MonetaryAmount{}.New(100, "EUR")
	moneyB := MonetaryAmount{}.New(100, "EUR")
	if moneyA != moneyB {
		t.Errorf("not equal")
	}
}

func Test_MonetaryAmount_should_be_not_equal(t *testing.T) {
	moneyA := MonetaryAmount{}.New(100, "EUR")
	moneyB := MonetaryAmount{}.New(200, "EUR")
	if moneyA == moneyB {
		t.Errorf("equal")
	}
}
