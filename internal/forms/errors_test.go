package forms

import "testing"

func TestAdd(t *testing.T) {
	errs := errors{}

	errs.Add("123", "321")
}

func TestGet(t *testing.T) {
	errs := errors{}

	errs.Add("123", "321")

	result := errs.Get("123")

	if result == "" {
		t.Error("should get result but empty")
	}

	result = errs.Get("3333")

	if result != "" {
		t.Error("should get empty but have result")
	}
}
