package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippets/view", snippetView)
	mux.HandleFunc("/snippets/create", snippetCreate)
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
