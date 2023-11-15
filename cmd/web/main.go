package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

// application struct is used to hold the application-wide dependencies
type application struct {
	logger *slog.Logger
}

func main() {
	// to use a defined var use:
	// flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	addr := flag.String("addr", ":8080", "HTTP network address")
	flag.Parse()

	// Use the slog.New() function to initialize a new structured logger, which
	// writes to the standard out stream and uses the default settings.
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	app := &application{
		logger: logger,
	}

	logger.Info("starting server", slog.String("addr", *addr))

	err := http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
