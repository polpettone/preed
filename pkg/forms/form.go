package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d characters)", d))
	}
}

func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("This field is too short (minimum ist %d characters)", d))
	}
}

func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value) {
		f.Errors.Add(field, "This field is invalid")
	}
}


func (f *Form) IsNumber(field string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	_, err := strconv.Atoi(value)
	if err != nil {
		f.Errors.Add(field, "This field is invalid")
	}
}

func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid")
}

func (f *Form) ValidDateFormat(field string) {
	_, err := parseDate(f.Get(field), "02.01.2006")
	if err != nil {
		f.Errors.Add(field, fmt.Sprintf("%s", err))
	}
}

func (f *Form) ValidMoneyFormat(field string) {
	value := f.Get(field)
	_, err := strconv.Atoi(value)
	if err != nil {
		f.Errors.Add(field, fmt.Sprintf("%s", err))
	}
}

func parseDate(date, layout string) (*time.Time, error) {
	parsed, err := time.Parse(layout, date)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
