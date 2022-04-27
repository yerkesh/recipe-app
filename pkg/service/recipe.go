package service

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"recipe-app/pkg/domain"
	"recipe-app/pkg/domain/constant"
	"recipe-app/pkg/repository"
	"recipe-app/pkg/util/fault"
	"recipe-app/pkg/util/validator"
	"recipe-app/pkg/util/writer"
)

type RecipeService struct {
	repo repository.RecipeRepoer
}

func NewRecipeService(repo repository.RecipeRepoer) *RecipeService {
	return &RecipeService{repo: repo}
}

func (svc *RecipeService) Recipe(reqCtx context.Context, id uint64) (r *domain.RecipeView, err error) {
	if err = svc.repo.Begin(reqCtx, func(tx pgx.Tx) error {
		r, err = svc.repo.GetRecipe(reqCtx, tx, id)
		if err != nil {
			return fmt.Errorf("couldn't get recipe err: %w", err)
		}

		return nil
	}); err != nil {
		return nil, fault.SanitizeServiceError(err)
	}

	return r, nil
}

func (svc *RecipeService) RecipeSteps(reqCtx context.Context, recipeID uint64) (steps []*domain.Step, err error) {
	if err = svc.repo.Begin(reqCtx, func(tx pgx.Tx) error {
		steps, err = svc.repo.GetRecipeSteps(reqCtx, tx, recipeID)
		if err != nil {
			return fmt.Errorf("couldn't get recipe steps err: %w", err)
		}

		return nil
	}); err != nil {
		return nil, fault.SanitizeServiceError(err)
	}

	for _, step := range steps {
		if step.ImageURL.Valid {
			step.Image = step.ImageURL.String
		} else {
			step.Image = ""
		}
	}
	return steps, nil
}

func (svc *RecipeService) RecipeReview(reqCtx context.Context, recipeID uint64) (reviews []*domain.Review, err error) {
	if err = svc.repo.Begin(reqCtx, func(tx pgx.Tx) error {
		reviews, err = svc.repo.GetRecipeReview(reqCtx, tx, recipeID)
		if err != nil {
			return fmt.Errorf("couldn't get recipe reviews err: %w", err)
		}

		return nil
	}); err != nil {
		return nil, fault.SanitizeServiceError(err)
	}

	return reviews, nil
}

func (svc *RecipeService) LeaveReview(reqCtx context.Context, r *domain.ReviewCreate) (crv *domain.CreatedObjectView, err error) {
	var cr domain.CreatedObjectView
	var rID uint64
	//if msg, vmap := vld.Validate(&r); msg != nil {
	//	return 0, fault.WhsValidateError(*msg, vmap)
	//}

	if err = svc.repo.Begin(reqCtx, func(tx pgx.Tx) error {
		rID, err = svc.repo.LeaveReview(reqCtx, tx, r)
		if err != nil {
			return fmt.Errorf("couldn't create recipe review err: %w", err)
		}

		return nil
	}); err != nil {
		return nil, fault.SanitizeServiceError(err)
	}

	cr.ID = rID
	cr.ServiceResponse = writer.ServiceResponseCreated(constant.MsgCreated)

	return &cr, nil
}

func (svc *RecipeService) AddToFavourite(reqCtx context.Context, fav *domain.UserFavouriteCreate) (crv *domain.CreatedObjectView, err error) {
	var cr domain.CreatedObjectView
	var favID uint64
	if msg, vmap := validator.Validate(fav); msg != nil {
		return nil, fault.WhsValidateError(*msg, vmap)
	}

	if err = svc.repo.Begin(reqCtx, func(tx pgx.Tx) error {
		if favID, err = svc.repo.AddToFavourite(reqCtx, tx, fav.UserID, fav.RecipeID); err != nil {
			return fmt.Errorf("couldn't add user favourites err %w", err)
		}

		return nil
	}); err != nil {
		return nil, fault.SanitizeServiceError(err)
	}

	cr.ID = favID
	cr.ServiceResponse = writer.ServiceResponseCreated(constant.MsgCreated)

	return &cr, nil
}

func (svc *RecipeService) UserFavourites(reqCtx context.Context, userID uint64) (fs []*domain.UserFavourite, err error) {
	if err = svc.repo.Begin(reqCtx, func(tx pgx.Tx) error {
		if fs, err = svc.repo.GetUserFavourite(reqCtx, tx, userID); err != nil {
			return fmt.Errorf("couldn't get user favourites err %w", err)
		}

		return nil
	}); err != nil {
		return nil, fault.SanitizeServiceError(err)
	}

	return fs, nil
}

func (svc *RecipeService) RemoveUserFavourite(reqCtx context.Context, userID, recipeID uint64) (res writer.ServiceResponse, err error) {
	if err = svc.repo.Begin(reqCtx, func(tx pgx.Tx) error {
		if err = svc.repo.RemoveFavourite(reqCtx, tx, userID, recipeID); err != nil {
			return fmt.Errorf("couldn't remove err %w", err)
		}

		return nil
	}); err != nil {
		return res, fault.SanitizeServiceError(err)
	}

	res = writer.ServiceResponseOk(constant.MsgDeleted)

	return res, nil
}
