package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"gopkg.in/asaskevich/govalidator.v9"
)

// Form creates a custom form struct and embeds url.values object
type Form struct {
	url.Values
	Errors errors
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Has checks if form field is inpost and not empty
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		return false
	}

	return true
}

// Valid returns true if there are no errors otherwise it returns false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// Required checks for required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be empty")
		}
	}
}

// MinLength checks for string minimum lenght
func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	x := r.Form.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}

	return true
}

func (f *Form) IsEmail(field string) {

	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "invalid email adress")
	}
}
