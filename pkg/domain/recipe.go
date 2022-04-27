package domain

import (
	"database/sql"
	"time"
)

type RecipeView struct {
	RecipeID    uint64       `json:"recipe_id"`
	RecipeName  string       `json:"recipe_name"`
	Description string       `json:"description"`
	ImageURL    string       `json:"image_url"`
	Rate        float64      `json:"rate"`
	Calorie     uint64       `json:"calorie"`
	CookingTime uint64       `json:"cooking_time"`
	Ingredients []Ingredient `json:"ingredients"`
}

type Ingredient struct {
	IngredientID        uint64  `json:"ingredient_id"`
	IngredientName      string  `json:"ingredient_name"`
	IngredientImageURL  string  `json:"ingredient_image_url"`
	Quantity            float64 `json:"quantity"`
	UnitOfMeasurementID string  `json:"unit_of_measurement"`
}

type ReviewCreate struct {
	UserID      uint64 `json:"user_id"`
	RecipeID    uint64 `json:"recipe_id"`
	CommentText string `json:"comment_text"`
	Star        uint64 `json:"star"`
}

type Step struct {
	StepNumber  uint64         `json:"step_number"`
	Description string         `json:"text"`
	Duration    uint64         `json:"duration"`
	ImageURL    sql.NullString `json:"-"`
	Image       string         `json:"image"`
}

func (s *Step) ScanFields() []interface{} {
	return []interface{}{
		&s.StepNumber,
		&s.Duration,
		&s.Description,
		&s.ImageURL,
	}
}

type Review struct {
	CommentID       uint64         `json:"comment_id"`
	UserName        string         `json:"username"`
	Star            uint64         `json:"star"`
	CommentText     string         `json:"comment_text"`
	CommentNullable sql.NullString `json:"-"`
	CreatedDate     time.Time      `json:"created_date"`
}

func (r *Review) ScanFields() []interface{} {
	return []interface{}{
		&r.CommentID,
		&r.CommentNullable,
		&r.Star,
		&r.CreatedDate,
		&r.UserName,
	}
}

type UserFavouriteCreate struct {
	UserID   uint64 `json:"user_id" validate:"required"`
	RecipeID uint64 `json:"recipe_id" validate:"required"`
}

type UserFavourite struct {
	RecipeId    uint64     `json:"id"`
	Name        string     `json:"name"`
	Complexity  Complexity `json:"complexity"`
	Category    Category   `json:"category"`
	Rate        uint64     `json:"rate"`
	Duration    uint64     `json:"cookingTime"`
	Calorie     uint64     `json:"calorie"`
	ImageURL    string     `json:"image"`
	Description string     `json:"description"`
}

type Complexity struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type Category struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

func (u *UserFavourite) ScanFields() []interface{} {
	return []interface{}{
		&u.RecipeId,
		&u.Name,
		&u.Description,
		&u.Duration,
		&u.Calorie,
		&u.ImageURL,
		&u.Rate,
		&u.Complexity.ID,
		&u.Complexity.Name,
		&u.Category.ID,
		&u.Category.Name,
		&u.Category.Image,
	}
}
