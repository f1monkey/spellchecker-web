package routes

import (
	"net/http"

	"github.com/f1monkey/spellchecker-web/internal/spellchecker"
	"github.com/go-chi/chi/v5"
	"github.com/swaggest/rest/nethttp"
)

type EmptyResponse struct{}

func Routes(registry *spellchecker.Registry) func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/dictionaries", dictionaryRoutes(registry))
	}
}

func dictionaryRoutes(registry *spellchecker.Registry) func(r chi.Router) {
	return func(r chi.Router) {
		r.Method(http.MethodPost, "/{code}", nethttp.NewHandler(
			dictionaryCreate(registry)),
		)

		r.Method(http.MethodDelete, "/{code}", nethttp.NewHandler(
			dictionaryDelete(registry)),
		)

		r.Method(http.MethodPost, "/{code}/save", nethttp.NewHandler(
			dictionarySave(registry)),
		)
	}
}
