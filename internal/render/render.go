package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/winadiw/go-bookings/internal/config"
	"github.com/winadiw/go-bookings/internal/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

var pathToTemplates = "./templates"

// NewRenderer sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// AddDefaultData sets the config for default data
// pointers used to modify data here
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")

	td.CSRFToken = nosurf.Token(r)
	return td
}

// Template renders templates using html/template
func Template(rw http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	// Get the template cache from the app config
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// Get template cache by its name
	t, ok := tc[tmpl]

	// If template not found, throw fatal for now
	if !ok {
		log.Println("Missing template", tmpl)
		return errors.New("can't get template from cache")
	}

	// Create temp buffer
	buff := new(bytes.Buffer)

	// Execute AddDefaultData if needed
	td = AddDefaultData(td, r)

	// Execute buffer with given data
	_ = t.Execute(buff, td)

	// Write buff to http.ResponseWriter
	_, err := buff.WriteTo(rw)

	if err != nil {
		fmt.Println("Error writing template to browser", err)
		return err
	}

	return nil
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))

	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		// name ex: about.page.tmpl
		name := filepath.Base(page)

		// create template with given name and parse the files
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)

		if err != nil {
			return myCache, err
		}

		//Get all layout templates
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))

		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			// ts.ParseGlob to obtain all template layouts
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))

			if err != nil {
				return myCache, err
			}
		}

		// Add to map by given name + template combo
		myCache[name] = ts
	}

	return myCache, nil
}
