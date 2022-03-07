package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"recipe-app/pkg/handler"
	"recipe-app/pkg/handler/rest"
)

func Router(h *handler.Ctx) chi.Router {
	log.Println("Router is initialized")

	r := chi.NewRouter()
	rst :=rest.NewRecipeRest(h)
	r.Use(middleware.Logger)
	r.Route("/", func(r chi.Router) {
		r.Get("/hello", rst.Hello)
	})
	return r
}
