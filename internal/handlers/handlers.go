package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/winadiw/go-bookings/internal/config"
	"github.com/winadiw/go-bookings/internal/driver"
	"github.com/winadiw/go-bookings/internal/forms"
	"github.com/winadiw/go-bookings/internal/helpers"
	"github.com/winadiw/go-bookings/internal/models"
	"github.com/winadiw/go-bookings/internal/render"
	"github.com/winadiw/go-bookings/internal/repository"
	"github.com/winadiw/go-bookings/internal/repository/dbrepo"
)

// Repo the repository used by the handlers
var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(rw http.ResponseWriter, r *http.Request) {
	render.Template(rw, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(rw http.ResponseWriter, r *http.Request) {
	render.Template(rw, r, "about.page.tmpl", &models.TemplateData{})
}

// Generals is the generals page handler
func (m *Repository) Generals(rw http.ResponseWriter, r *http.Request) {
	render.Template(rw, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors is the majors-suite page handler
func (m *Repository) Majors(rw http.ResponseWriter, r *http.Request) {
	render.Template(rw, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability is the search page handler
func (m *Repository) Availability(rw http.ResponseWriter, r *http.Request) {
	render.Template(rw, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability is the search page handler
func (m *Repository) PostAvailability(rw http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)

	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	if len(rooms) == 0 {
		// no availability
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(rw, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})

	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(rw, r, "choose-room.page.tmpl", &models.TemplateData{Data: data})
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
		helpers.ServerError(rw, err)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(out)
}

// Contact is the contact page handler
func (m *Repository) Contact(rw http.ResponseWriter, r *http.Request) {
	render.Template(rw, r, "contact.page.tmpl", &models.TemplateData{})
}

// MakeReservation is the make reservation page handler
func (m *Repository) MakeReservation(rw http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.Template(rw, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	// https://www.pauladamsmith.com/blog/2011/05/go_time.html
	// 01/02 03:04:05PM '06 -0700
	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}

	forms := forms.New(r.PostForm)

	forms.Required("first_name", "last_name", "email")
	forms.MinLength("first_name", 3)
	forms.IsEmail("email")

	if !forms.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(rw, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: forms,
			Data: data,
		})
		return
	}

	newReservationId, err := m.DB.InsertReservation(reservation)

	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: newReservationId,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestrictions(restriction)

	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(rw, r, "/reservation-summary", http.StatusSeeOther)
}

func (m *Repository) ReservationSummary(rw http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		m.App.ErrorLog.Println("Can't get reservation from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})

	data["reservation"] = reservation

	render.Template(rw, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) ChooseRoom(rw http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		helpers.ServerError(rw, err)
		return
	}

	res.RoomID = roomID

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(rw, r, "/make-reservation", http.StatusSeeOther)
}
