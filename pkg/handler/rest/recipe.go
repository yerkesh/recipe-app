package rest

import (
	"net/http"
	"recipe-app/pkg/handler"
	"recipe-app/pkg/util/writer"
)

type RecipeRest struct {
	ctx *handler.Ctx
}

func NewRecipeRest() *RecipeRest {
	return &RecipeRest{}
}

func (r *RecipeRest) Hello(res http.ResponseWriter, req *http.Request) {
	writer.HTTPResponseWriter(res,nil,r.ctx.RecipeService.Hello())
}