package render

import (
	"net/http"
	"testing"

	"github.com/winadiw/go-bookings/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()

	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "123")

	result := AddDefaultData(&td, r)

	if result != nil && result.Flash != "123" {
		t.Error("flash of value of 123 not found in session")
	} else {
		t.Error("nil result")
	}
}

func getSession() (*http.Request, error) {

	r, err := http.NewRequest("GET", "/some-url", nil)

	if err != nil {
		return nil, err
	}

	ctx := r.Context()

	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))

	r = r.WithContext(ctx)

	return r, nil
}
