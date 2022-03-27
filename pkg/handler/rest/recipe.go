package rest

import (
	"log"
	"net/http"
	"recipe-app/pkg/handler"
	"recipe-app/pkg/util/writer"
)

type RecipeRest struct {
	ctx *handler.Ctx
}

func NewRecipeRest(ctx *handler.Ctx) *RecipeRest {
	return &RecipeRest{ctx: ctx}
}

func (r *RecipeRest) Hello(res http.ResponseWriter, req *http.Request) {
	log.Println("get req")
	writer.HTTPResponseWriter(res, nil, r.ctx.RecipeService.Hello())
}
