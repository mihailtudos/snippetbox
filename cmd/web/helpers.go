package main

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
)

// serverError helper writes a Error log then sends a generic 500 response
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	// creates an error level log entry with the method nad uri
	app.logger.Error(err.Error(), slog.String("method", method), slog.String("uri", uri), slog.String("trace", trace))

	// sends an HTTP error to the user
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError helper sends a specific status code and corresponding description
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// notFound helper is a convenience wrapper around clientError which sends a 404
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// render helper renders the templates from the cache
func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	// Retrieve the appropriate template set from the cache
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	// If an error ocurred when parsing writing the templates
	// the user would still get 200 OK and receive a broken (half generated page)
	/*
		To fix this we need to make the template render a two-stage process. First, we should make
		a ‘trial’ render by writing the template into a buffer. If this fails, we can respond to the user
		with an error message. But if it works, we can then write the contents of the buffer to our
		http.ResponseWriter.
	*/

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// If the template is written to the buffer without any errors, we are safe
	// to go ahead and write the HTTP status code to http.ResponseWriter.
	// Write out the provided HTTP status code
	w.WriteHeader(status)

	// Execute the template set and write the response body
	buf.WriteTo(w)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear: time.Now().Year(),
		Flash:       app.sessionManager.PopString(r.Context(), "flash"),
	}
}

// decodePostForm parses and decords a HTML form, it returns an error if one ocurred
func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// formDecode's Decode method will fill the form struct with the relevant values from the HTML form
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// If we try to use an invalid target destination, the Decode() method
		// will return an error with the type *form.InvalidDecoderError.We use
		// errors.As() to check for this and raise a panic rather than returning
		// the error.
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		// For all other errors, we return them as normal.
		return err
	}

	return nil
}
