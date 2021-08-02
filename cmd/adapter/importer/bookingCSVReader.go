package importer

import (
	"encoding/csv"
	"github.com/Rhymond/go-money"
	"github.com/gocarina/gocsv"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type RawBooking struct {
	NameAnschrift       string         `csv:"NameAnschrift"`
	StartDate           DateTime       `csv:"von"`
	EndDate             DateTime       `csv:"bis"`
	Days                int            `csv:"AnzahlTage"`
	Item                string         `csv:"Item"`
	Provider            string         `csv:"Provider"`
	PricePerDay         MonetaryAmount `csv:"PreisProTag"`
	IntermediateSum     MonetaryAmount `csv:"Zwischensumme"`
	CleaningIncome      MonetaryAmount `csv:"Endreinigungseinnahme"`
	Total               MonetaryAmount `csv:"Gesamtbetrag"`
	Provision           MonetaryAmount `csv:"Provision"`
	TotalMinusProvision MonetaryAmount `csv:"GesamtMinusProvision"`
	CleaningCost        MonetaryAmount `csv:"Endreinigungskosten"`
	BookingNumber       string         `csv:"Buchungsnummer"`
	InvoiceNumber       string         `csv:"Rechnungsnummer"`
	CashTransferDate    DateTime       `csv:"Überweisungsdatum"`
	CleaningDate        DateTime       `csv:"Reinigungsdatum"`
}

type DateTime struct {
	time.Time
}

type MonetaryAmount struct {
	money.Money
}

func (date *DateTime) UnmarshalCSV(csv string) (err error) {
	//TODO: time location not set, default is utc
	if csv != "" {
		date.Time, err = time.ParseInLocation("02.01.2006", csv, time.Local)
		return err
	}
	return nil
}

func (m *MonetaryAmount) UnmarshalCSV(csv string) (err error) {
	s := strings.Replace(csv, "€", "", -1)
	s = strings.TrimSpace(s)
	s = strings.Replace(s, ",", "", -1)

	value, err := strconv.Atoi(s)

	if err != nil {
		return err
	}

	m.Money = *money.New(int64(value), "EUR")
	return err
}

func ReadCSV(file string) ([]*RawBooking, error) {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var rawBookings []*RawBooking

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = ';'
		return r
	})

	if err := gocsv.UnmarshalFile(f, &rawBookings); err != nil {
		return nil, err
	}
	return rawBookings, nil
}
