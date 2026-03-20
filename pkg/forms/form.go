package forms

import (
	"fmt"
	"net/url"
	"slices"
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

func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)

	if value == "" {
		return
	}

	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (Max: %d)", d))
	}
}

func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)

	if value == "" {
		return
	}

	if slices.Contains(opts, value) {
		f.Errors.Add(field, "This field is invalid")
	}
}

func (f *Form) valid() bool {
	return len(f.Errors) == 0
}
