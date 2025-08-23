package web

import (
	"context"
	"net/http"

	"github.com/f1monkey/spellchecker"
	"github.com/f1monkey/spellchecker-web/internal/routes"
	"github.com/go-chi/chi/v5"
	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/rest/nethttp"
	"github.com/swaggest/rest/web"
	swgui "github.com/swaggest/swgui/v5emb"
)

func NewServer(appCtx context.Context, sc *spellchecker.Spellchecker) *web.Service {
	s := web.NewService(openapi31.NewReflector())

	s.OpenAPISchema().SetTitle("Spellchecker")
	s.OpenAPISchema().SetDescription("To fix words")
	s.OpenAPISchema().SetVersion("v1")

	s.Route("/v1", func(r chi.Router) {
		r.Method(http.MethodPost, "/spellchecker/fix", nethttp.NewHandler(
			routes.SpellcheckerFix(sc)),
		)
		r.Method(http.MethodPost, "/dictionary/add", nethttp.NewHandler(
			routes.DictionaryAdd(sc)),
		)
	})

	// Swagger UI endpoint at /docs.
	s.Docs("/docs", swgui.New)

	return s
}
