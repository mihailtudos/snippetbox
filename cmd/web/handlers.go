package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/base.gohtml",
		"./ui/html/partials/nav.gohtml",
		"./ui/html/pages/home.gohtml",
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "View snippet page %d", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Cannot handle your request", http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("creating new snippet"))
}
