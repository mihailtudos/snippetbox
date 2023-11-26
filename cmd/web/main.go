package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mihailtudos/snippetbox/internal/models"
)

var (
	addr   string
	dsn    string
	logger *slog.Logger
)

// application struct is used to hold the application-wide dependencies
type application struct {
	logger        *slog.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
}

func init() {
	flag.StringVar(&addr, "addr", ":8080", "HTTP network address")
	flag.StringVar(&dsn, "dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	// The slog.New() function to initialize a new structured logger, which
	// writes to the standard out stream and uses the default settings.
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	flag.Parse()

	db, err := openDB(dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	templates, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

	app := &application{
		logger:        logger,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templates,
		formDecoder:   formDecoder,
	}

	logger.Info("starting server", slog.String("addr", addr))

	err = http.ListenAndServe(addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
