package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	// to use a defined var use:
	// flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	addr := flag.String("addr", ":8080", "HTTP network address")
	flag.Parse()

	// Use the slog.New() function to initialize a new structured logger, which
	// writes to the standard out stream and uses the default settings.
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippets/view", snippetView)
	mux.HandleFunc("/snippets/create", snippetCreate)
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	logger.Info("starting server", slog.String("addr", ":4000"))

	err := http.ListenAndServe(*addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
}
