package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/mihailtudos/snippetbox/internals/models"
)

type templateData struct {
	Snippet     models.Snippet
	Snippets    []models.Snippet
	CurrentYear int
	Form        any
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache
	cache := make(map[string]*template.Template, 10)

	// The filepath.Globe returns a slice of filepaths of
	// all matching paths
	pages, err := filepath.Glob("./ui/html/pages/*.gohtml")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// get the file name (e.g. home.gohtml)
		name := filepath.Base(page)

		// Parse the base template file into a template set.
		// Register the emplate.FuncMap to the templates
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.gohtml")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() *on this template set* to add any partials.
		ts, err = ts.ParseGlob("./ui/html/partials/*.gohtml")
		if err != nil {
			return nil, err
		}

		// Call ParseFiles() *on this template set* to add the page template.
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map, using the name of the page
		cache[name] = ts
	}

	return cache, nil
}

// humanDate returns a formatted string representation of time.Time
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// store template.FuncMap in a new object
var functions = template.FuncMap{
	"humanDate": humanDate,
}
