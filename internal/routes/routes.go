package routes

import (
	"net/http"

	"github.com/f1monkey/spellchecker"
	"github.com/go-chi/chi/v5"
	"github.com/swaggest/rest/nethttp"
)

func Routes(sc *spellchecker.Spellchecker) func(r chi.Router) {
	return func(r chi.Router) {
		r.Method(http.MethodPost, "/spellchecker/fix", nethttp.NewHandler(
			SpellcheckerFix(sc)),
		)
		r.Method(http.MethodPost, "/dictionary/add", nethttp.NewHandler(
			DictionaryAdd(sc)),
		)
	}
}
