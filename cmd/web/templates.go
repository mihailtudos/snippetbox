package main

import "github.com/mihailtudos/snippetbox/internals/models"

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}
