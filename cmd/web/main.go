package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// to use a defined var use:
	// flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	addr := flag.String("addr", ":8080", "HTTP network address")
	flag.Parse()

	fmt.Println(*addr)
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippets/view", snippetView)
	mux.HandleFunc("/snippets/create", snippetCreate)
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Printf("starting server on %s", *addr)

	log.Fatal(http.ListenAndServe(*addr, mux))
}
