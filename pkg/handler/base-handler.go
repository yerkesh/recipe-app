package handler

import (
	"recipe-app/pkg/service"

	"github.com/gorilla/schema"
)

type Ctx struct {
	RecipeService service.RecipeServicer
	queryDecoder        *schema.Decoder
}


func NewHandlerCtx(opts ...Option) *Ctx {
	var h Ctx
	h.queryDecoder = schema.NewDecoder()

	for _, opt := range opts {
		opt(&h)
	}

	return &h
}

type Option func(ctx *Ctx)

func WithRecipeService(svc service.RecipeServicer) Option {
	return func(ctx *Ctx) {
		ctx.RecipeService = svc
	}
}