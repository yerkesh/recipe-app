package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"recipe-app/pkg/domain"
	"recipe-app/pkg/repository/database"
)

type RecipeRepoer interface {
	database.Beginner
	GetRecipe(
		reqCtx context.Context,
		tx pgx.Tx,
		id uint64,
	) (r *domain.RecipeView, err error)
	GetRecipeSteps(
		reqCtx context.Context,
		tx pgx.Tx,
		recipeID uint64,
	) (steps []*domain.Step, err error)
	GetRecipeReview(
		reqCtx context.Context,
		tx pgx.Tx,
		recipeID uint64,
	) (reviews []*domain.Review, err error)
	LeaveReview(
		reqCtx context.Context,
		tx pgx.Tx,
		review *domain.ReviewCreate,
	) (rID uint64, err error)
	AddToFavourite(
		reqCtx context.Context,
		tx pgx.Tx,
		userID, recipeID uint64,
	) (favID uint64, err error)
	GetUserFavourite(
		reqCtx context.Context,
		tx pgx.Tx,
		userID uint64,
	) (fs []*domain.UserFavourite, err error)
	RemoveFavourite(reqCtx context.Context, tx pgx.Tx, userID, recipeID uint64) (err error)
}
