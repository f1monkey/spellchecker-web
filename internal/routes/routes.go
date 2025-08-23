package routes

import (
	"net/http"
	"regexp"

	"github.com/f1monkey/spellchecker-web/internal/spellchecker"
	"github.com/go-chi/chi/v5"
	"github.com/swaggest/rest/nethttp"
)

type Empty struct{}

func Routes(registry *spellchecker.Registry, splitter *regexp.Regexp) func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/dictionaries", dictionaryRoutes(registry, splitter))
		r.Route("/aliases", aliasRoutes(registry))
	}
}

func dictionaryRoutes(registry *spellchecker.Registry, splitter *regexp.Regexp) func(r chi.Router) {
	return func(r chi.Router) {
		r.Method(http.MethodGet, "/", nethttp.NewHandler(
			dictionaryList(registry),
		))

		r.Method(http.MethodPost, "/{code}", nethttp.NewHandler(
			dictionaryCreate(registry),
		))

		r.Method(http.MethodDelete, "/{code}", nethttp.NewHandler(
			dictionaryDelete(registry),
		))

		r.Method(http.MethodPost, "/{code}/save", nethttp.NewHandler(
			dictionarySave(registry),
		))

		r.Method(http.MethodPost, "/{code}/add", nethttp.NewHandler(
			dictionaryItemAdd(registry, splitter),
		))

		r.Method(http.MethodPost, "/{code}/fix", nethttp.NewHandler(
			dictionaryFix(registry, splitter),
		))
	}
}

func aliasRoutes(registry *spellchecker.Registry) func(r chi.Router) {
	return func(r chi.Router) {
		r.Method(http.MethodGet, "/", nethttp.NewHandler(
			aliasList(registry),
		))

		r.Method(http.MethodGet, "/{alias}", nethttp.NewHandler(
			aliasGet(registry),
		))

		r.Method(http.MethodPut, "/{alias}", nethttp.NewHandler(
			aliasSet(registry),
		))

		r.Method(http.MethodDelete, "/{alias}", nethttp.NewHandler(
			aliasDelete(registry),
		))
	}
}
