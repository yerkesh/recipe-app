package database

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"recipe-app/pkg/domain"
	"recipe-app/pkg/domain/constant"
	"recipe-app/pkg/repository/database/repository"
	"recipe-app/pkg/util"
	"recipe-app/pkg/util/fault"
	"recipe-app/pkg/util/sql"
)

type RecipeRepo struct {
	table constant.Table
	*repository.Base
}

func NewRecipeApp(pool *pgxpool.Pool) *RecipeRepo {
	return &RecipeRepo{
		Base:  repository.New(pool),
		table: constant.TblRecipe,
	}
}

func (repo *RecipeRepo) GetRecipe(
	reqCtx context.Context,
	tx pgx.Tx,
	id uint64,
) (r *domain.RecipeView, err error) {
	var recipe domain.RecipeView
	ingredients, _, _ := sql.SB().Select(
		`json_agg(json_build_object(
					'ingredient_id', ing.id,
					'ingredient_name', ing.name, 
					'unit_of_measurement', uom.name, 
					'quantity', ir.quantity))`).
		From(constant.TblIngredient.String() + " ing").
		Join(constant.TblUnitOfMeasurement.String() + " uom on uom.id=ing.unit_of_measurement_id").
		Join(constant.TblIngredientRecipe.String() + " ir on ir.ingredient_id=ing.id").
		Where("ir.recipe_id=r.id").
		Prefix("(").Suffix(")").ToSql()

	qs, args, err := sql.SB().Select(
		ingredients,
		"r.id",
		"r.name",
		"r.cooking_time",
		"r.calorie",
		"r.description",
		"r.image",
	).From(constant.TblRecipe.As("r")).Where(sq.Eq{"r.id": id}).ToSql()
	if err != nil {
		log.Printf("sql compose err: %v", err)

		return nil, fault.SanitizeDBError(err, qs, args)
	}

	err = tx.QueryRow(reqCtx, qs, args...).Scan(
		&recipe.Ingredients,
		&recipe.RecipeID,
		&recipe.RecipeName,
		&recipe.CookingTime,
		&recipe.Calorie,
		&recipe.Description,
		&recipe.ImageURL,
	)
	if err != nil {
		log.Printf("sql scan err: %v", err)

		return nil, fault.SanitizeDBError(err, qs, args)
	}

	return &recipe, nil
}

func (repo *RecipeRepo) GetRecipeSteps(
	reqCtx context.Context,
	tx pgx.Tx,
	recipeID uint64,
) (steps []*domain.Step, err error) {
	qs, args, err := sql.SB().Select("number", "duration", "description").
		From(constant.TblRecipeStep.String()).Where(sq.Eq{"recipe_id": recipeID}).ToSql()
	if err != nil {
		log.Printf("sql scan err: %v", err)

		return nil, fault.SanitizeDBError(err, qs, args)
	}

	s := new(domain.Step)
	scans := s.ScanFields()
	if _, err = tx.QueryFunc(reqCtx, qs, args, scans, func(row pgx.QueryFuncRow) error {
		curStep := *s
		steps = append(steps, &curStep)

		return nil
	}); err != nil {
		log.Printf("sql scan err: %v", err)

		return nil, fault.SanitizeDBError(err, qs, args)
	}

	return steps, nil
}

func (repo *RecipeRepo) GetRecipeReview(
	reqCtx context.Context,
	tx pgx.Tx,
	recipeID uint64,
) (reviews []*domain.Review, err error) {
	qs, args, err := sql.SB().Select("c.id", "c.text", "c.star", "c.created_date", "u.username").
		From(constant.TblComment.As("c")).Join(constant.TblUsers.String() + " u on u.id=c.user_id").
		Where(sq.Eq{"c.recipe_id": recipeID}).ToSql()
	if err != nil {
		log.Printf("sql scan err: %v", err)

		return nil, fault.SanitizeDBError(err, qs, args)
	}

	c := new(domain.Review)
	scans := c.ScanFields()
	if _, err = tx.QueryFunc(reqCtx, qs, args, scans, func(row pgx.QueryFuncRow) error {
		curRew := *c
		reviews = append(reviews, &curRew)

		return nil
	}); err != nil {
		log.Printf("sql scan err: %v", err)

		return nil, fault.SanitizeDBError(err, qs, args)
	}

	return reviews, nil
}

func (repo *RecipeRepo) LeaveReview(
	reqCtx context.Context,
	tx pgx.Tx,
	review *domain.ReviewCreate,
) (rID uint64, err error) {
	qs, args, err := sql.SB().Insert(constant.TblComment.String()).
		Columns(
			"user_id",
			"recipe_id",
			"text",
			"star",
			"created_date").
		Values(
			review.UserID,
			review.RecipeID,
			review.CommentText,
			review.Star,
			util.CurTime()).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		log.Printf("sql scan err: %v", err)

		return 0, fault.SanitizeDBError(err, qs, args)
	}

	if err = tx.QueryRow(reqCtx, qs, args...).Scan(&rID); err != nil {
		log.Printf("sql scan err: %v", err)

		return 0, fault.SanitizeDBError(err, qs, args)
	}

	return rID, nil
}

func (repo *RecipeRepo) updateRecipeRate(reqCtx context.Context, tx pgx.Tx, recipeID uint64) error {
	calculatedStar := sql.SB().Select("round(sum(star::numeric)/count(id), 1").
		From(constant.TblComment.String()).
		Where(sq.Eq{"recipe_id": recipeID})
	qs, args, err := sql.SB().Update(constant.TblRecipe.String()).Set("rate", calculatedStar).ToSql()
	if err != nil {
		log.Printf("sql scan err: %v", err)

		return fault.SanitizeDBError(err, qs, args)
	}

	if _, err = tx.Exec(reqCtx, qs, args...); err != nil {
		log.Printf("sql exec err: %v", err)

		return fault.SanitizeDBError(err, qs, args)
	}

	return nil
}
