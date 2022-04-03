package service

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"recipe-app/pkg/domain"
	"recipe-app/pkg/domain/constant"
	"recipe-app/pkg/repository"
	"recipe-app/pkg/util/fault"
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
