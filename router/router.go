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
	rst := rest.NewRecipeRest(h)
	r.Use(middleware.Logger)
	r.Route("/", func(r chi.Router) {
		r.Get("/recipe/review/{recipeID}", rst.RecipeReview)
		r.Post("/leave/review", rst.LeaveReview)
		r.Get("/recipe/{recipeID}", rst.GetRecipe)
		r.Get("/recipe/steps/{recipeID}", rst.RecipeSteps)
		r.Get("/user/favourite/{userID}", rst.GetUserFavourites)
		r.Post("/user/favourite", rst.AddToFavourites)
		r.Delete("/user/favourite/{userID}/recipe/{recipeID}", rst.RemoveFavourite)
	})

	return r
}
