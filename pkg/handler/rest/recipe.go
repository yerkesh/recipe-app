package rest

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"recipe-app/pkg/domain"
	"recipe-app/pkg/domain/constant"
	"recipe-app/pkg/handler"
	"recipe-app/pkg/util"
	"recipe-app/pkg/util/fault"
	"recipe-app/pkg/util/writer"
)

const (
	recipeIDCtxKey domain.RestCtxKey = "recipeID"
	userIDCtxKey   domain.RestCtxKey = "userID"
)

type RecipeRest struct {
	ctx *handler.Ctx
}

func NewRecipeRest(ctx *handler.Ctx) *RecipeRest {
	return &RecipeRest{ctx: ctx}
}

func (r *RecipeRest) GetRecipe(res http.ResponseWriter, req *http.Request) {
	var parsedID uint64
	var err error
	var rew *domain.RecipeView

	if idStr := chi.URLParam(req, recipeIDCtxKey.String()); idStr != "" {
		parsedID, err = util.ParseUint64(idStr)
		if err != nil {
			writer.HTTPResponseWriter(res, fault.Whs400Error(err.Error(), constant.MsgRequestBodyErr), nil)

			return
		}

		rew, err = r.ctx.RecipeService.Recipe(req.Context(), parsedID)
		if err != nil {
			writer.HTTPResponseWriter(res, err, nil)

			return
		}
	} else {
		writer.HTTPResponseWriter(res, fault.Whs404Error(err.Error(), constant.MsgNotFoundErr), nil)

		return
	}

	writer.HTTPResponseWriter(res, nil, rew)
}

func (r *RecipeRest) RecipeSteps(res http.ResponseWriter, req *http.Request) {
	var parsedID uint64
	var err error
	var steps []*domain.Step

	if idStr := chi.URLParam(req, recipeIDCtxKey.String()); idStr != "" {
		parsedID, err = util.ParseUint64(idStr)
		if err != nil {
			writer.HTTPResponseWriter(res, fault.Whs400Error(err.Error(), constant.MsgRequestBodyErr), nil)

			return
		}

		steps, err = r.ctx.RecipeService.RecipeSteps(req.Context(), parsedID)
		if err != nil {
			writer.HTTPResponseWriter(res, err, nil)

			return
		}
	} else {
		writer.HTTPResponseWriter(res, fault.Whs404Error(err.Error(), constant.MsgNotFoundErr), nil)

		return
	}

	writer.HTTPResponseWriter(res, nil, steps)
}

func (r *RecipeRest) RecipeReview(res http.ResponseWriter, req *http.Request) {
	var parsedID uint64
	var err error
	var reviews []*domain.Review

	if idStr := chi.URLParam(req, recipeIDCtxKey.String()); idStr != "" {
		parsedID, err = util.ParseUint64(idStr)
		if err != nil {
			writer.HTTPResponseWriter(res, fault.Whs400Error(err.Error(), constant.MsgRequestBodyErr), nil)

			return
		}

		reviews, err = r.ctx.RecipeService.RecipeReview(req.Context(), parsedID)
		if err != nil {
			writer.HTTPResponseWriter(res, err, nil)

			return
		}
	} else {
		writer.HTTPResponseWriter(res, fault.Whs404Error(err.Error(), constant.MsgNotFoundErr), nil)

		return
	}

	writer.HTTPResponseWriter(res, nil, reviews)
}

func (r *RecipeRest) LeaveReview(res http.ResponseWriter, req *http.Request) {
	var rew domain.ReviewCreate

	if err := json.NewDecoder(req.Body).Decode(&rew); err != nil {
		writer.HTTPResponseWriter(res, fault.Whs400Error(err.Error(), constant.MsgRequestBodyErr), nil)

		return
	}

	result, err := r.ctx.RecipeService.LeaveReview(req.Context(), &rew)
	if err != nil {
		writer.HTTPResponseWriter(res, err, nil)

		return
	}

	res.WriteHeader(http.StatusCreated)
	writer.HTTPResponseWriter(res, nil, result)
}

func (r *RecipeRest) AddToFavourites(res http.ResponseWriter, req *http.Request) {
	var f domain.UserFavouriteCreate
	var err error

	if err = json.NewDecoder(req.Body).Decode(&f); err != nil {
		writer.HTTPResponseWriter(res, fault.Whs400Error(err.Error(), constant.MsgRequestBodyErr), nil)

		return
	}

	result, err := r.ctx.RecipeService.AddToFavourite(req.Context(), &f)
	if err != nil {
		writer.HTTPResponseWriter(res, err, nil)

		return
	}

	res.WriteHeader(http.StatusCreated)
	writer.HTTPResponseWriter(res, nil, result)
}

func (r *RecipeRest) GetUserFavourites(res http.ResponseWriter, req *http.Request) {
	var userID uint64
	var err error
	var favs []*domain.UserFavourite
	if idStr := chi.URLParam(req, userIDCtxKey.String()); idStr != "" {
		if userID, err = util.ParseUint64(idStr); err != nil {
			writer.HTTPResponseWriter(res, fault.Whs400Error(err.Error(), constant.MsgRequestBodyErr), nil)

			return
		}

		favs, err = r.ctx.RecipeService.UserFavourites(req.Context(), userID)
		if err != nil {
			writer.HTTPResponseWriter(res, err, nil)

			return
		}
	} else {
		writer.HTTPResponseWriter(res, fault.Whs404Error(err.Error(), constant.MsgNotFoundErr), nil)

		return
	}

	writer.HTTPResponseWriter(res, nil, favs)
}
