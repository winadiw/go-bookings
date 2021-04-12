package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/winadiw/go-bookings/internal/config"
	"github.com/winadiw/go-bookings/internal/handlers"
	"github.com/winadiw/go-bookings/internal/models"
	"github.com/winadiw/go-bookings/internal/render"
)

const portNumber = "localhost:8080"

var app config.AppConfig
var session *scs.SessionManager

// main is the main application function
func main() {
	// what am i going to put in the session
	gob.Register(models.Reservation{})

	// change this to true when production
	app.InProduction = false

	// Adjust Session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal("Cannot create template cache", err)
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Println("Starting application on port: " + portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()

	log.Fatal(err)
}
