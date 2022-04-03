package service

import (
	"context"
	"recipe-app/pkg/domain"
)

type RecipeServicer interface {
	Recipe(reqCtx context.Context, id uint64) (r *domain.RecipeView, err error)
	RecipeSteps(reqCtx context.Context, recipeID uint64) (steps []*domain.Step, err error)
	RecipeReview(reqCtx context.Context, recipeID uint64) (reviews []*domain.Review, err error)
	LeaveReview(reqCtx context.Context, r *domain.ReviewCreate) (crv *domain.CreatedObjectView, err error)
}
