package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// Valid returns true if there are no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{data, errors(map[string][]string{})}
}

// Required is used to check if data exists
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Has checks if form field is in post and not empty
func (f *Form) Has(field string) bool {
	x := f.Get(field)

	return x != ""
}

// MinLength is used to check the length of the field
func (f *Form) MinLength(field string, length int) bool {
	x := f.Get(field)

	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}

	return true
}

func (f *Form) IsEmail(field string) {
	x := f.Get(field)
	if !govalidator.IsEmail(x) {
		f.Errors.Add(field, "Invalid email address")
	}
}
