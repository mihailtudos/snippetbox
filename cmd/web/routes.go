package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippets/view", app.snippetView)
	mux.HandleFunc("/snippets/create", app.snippetCreate)
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return secureHeaders(mux)
}
