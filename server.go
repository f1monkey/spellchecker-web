package web

import (
	"context"

	"github.com/f1monkey/spellchecker-web/internal/routes"
	"github.com/f1monkey/spellchecker-web/internal/spellchecker"
	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/rest/web"
	swgui "github.com/swaggest/swgui/v5emb"
)

func NewServer(appCtx context.Context, registry *spellchecker.Registry) *web.Service {
	s := web.NewService(openapi31.NewReflector())

	s.OpenAPISchema().SetTitle("Spellchecker")
	s.OpenAPISchema().SetDescription("To fix words")
	s.OpenAPISchema().SetVersion("v1")

	s.Route("/v1", routes.Routes(registry))

	// Swagger UI endpoint at /docs.
	s.Docs("/docs", swgui.New)

	return s
}
