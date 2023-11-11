package main

import (
	"fmt"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "View snippet page")
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Create snippet page")
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/snippet/view", snippetView)
	http.HandleFunc("/snippet/create", snippetCreate)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
