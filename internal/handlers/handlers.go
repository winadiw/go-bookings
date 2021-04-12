package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/winadiw/go-bookings/internal/config"
	"github.com/winadiw/go-bookings/internal/forms"
	"github.com/winadiw/go-bookings/internal/models"
	"github.com/winadiw/go-bookings/internal/render"
)

// Repo the repository used by the handlers
var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(rw http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(rw, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(rw http.ResponseWriter, r *http.Request) {

	// perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	// send the data to the template
	render.RenderTemplate(rw, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Generals is the generals page handler
func (m *Repository) Generals(rw http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(rw, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors is the majors-suite page handler
func (m *Repository) Majors(rw http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(rw, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability is the search page handler
func (m *Repository) Availability(rw http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(rw, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability is the search page handler
func (m *Repository) PostAvailability(rw http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	rw.Write([]byte(start + end))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles request for availability and send JSON response
func (m *Repository) AvailabilityJSON(rw http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "    ")

	if err != nil {
		log.Fatal("err")
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(out)
}

// Contact is the contact page handler
func (m *Repository) Contact(rw http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(rw, r, "contact.page.tmpl", &models.TemplateData{})
}

// MakeReservation is the make reservation page handler
func (m *Repository) MakeReservation(rw http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(rw, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	forms := forms.New(r.PostForm)

	forms.Has("first_name", r)

	if !forms.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplate(rw, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: forms,
			Data: data,
		})
		return
	}
}
