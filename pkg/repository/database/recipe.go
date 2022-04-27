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
	qs, args, err := sql.SB().Select("number", "duration", "description", "image").
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
		if curRew.CommentNullable.Valid {
			curRew.CommentText = curRew.CommentNullable.String
		}

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

func (repo *RecipeRepo) AddToFavourite(reqCtx context.Context, tx pgx.Tx, userID, recipeID uint64) (favID uint64, err error) {
	qs, args, err := sql.SB().Insert(constant.TblUserFavourite.String()).
		Columns("user_id", "recipe_id").Values(userID, recipeID).Suffix("returning id").ToSql()
	if err != nil {
		log.Printf("sql compose err %s", err)

		return 0, fault.SanitizeDBError(err, qs, args)
	}

	if err = tx.QueryRow(reqCtx, qs, args...).Scan(&favID); err != nil {
		log.Printf("sql scan err %s", err)

		return 0, fault.SanitizeServiceError(err)
	}

	return favID, nil
}

func (repo *RecipeRepo) GetUserFavourite(
	reqCtx context.Context,
	tx pgx.Tx,
	userID uint64,
) (fs []*domain.UserFavourite, err error) {
	qs, args, err := sql.SB().Select(
		"rec.id",
		"rec.name",
		"rec.Description",
		"rec.cooking_time",
		"rec.calorie",
		"rec.image",
		"rec.rate",
		"cplx.id",
		"cplx.name",
		"cat.id",
		"cat.name",
		"cat.image").From(constant.TblRecipe.As("rec")).
		Join(constant.TblUserFavourite.As("f on f.recipe_id=rec.id")).
		Join(constant.TblComplexity.As("cplx on cplx.id=rec.complexity_id")).
		Join(constant.TblCategory.As("cat on cat.id=rec.category_id")).
		Where(sq.Eq{"f.user_id": userID}).ToSql()
	if err != nil {
		log.Printf("sql compose err %s", err)

		return
	}

	f := new(domain.UserFavourite)
	if _, err = tx.QueryFunc(reqCtx, qs, args, f.ScanFields(), func(row pgx.QueryFuncRow) error {
		curr := *f
		fs = append(fs, &curr)

		return nil
	}); err != nil {
		log.Printf("sql scan err %s", err)

		return nil, fault.SanitizeDBError(err, qs, args)
	}

	return fs, nil
}

func (repo *RecipeRepo) RemoveFavourite(reqCtx context.Context, tx pgx.Tx, userID, recipeID uint64) (err error) {
	qs, args, err := sql.SB().Delete(constant.TblUserFavourite.String()).Where(
		sq.And{sq.Eq{"user_id": userID}, sq.Eq{"recipe_id": recipeID}}).ToSql()
	if err != nil {
		log.Printf("sql compose err %s", err)

		return fault.SanitizeDBError(err, qs, args)
	}

	if _, err = tx.Exec(reqCtx, qs, args...); err != nil {
		log.Printf("sql exec err %s", err)

		return fault.SanitizeDBError(err, qs, args)
	}

	return nil
}
