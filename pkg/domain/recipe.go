package domain

import "time"

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
	StepNumber  uint64 `json:"step_number"`
	Description string `json:"text"`
	Duration    uint64 `json:"duration"`
}

func (s *Step) ScanFields() []interface{} {
	return []interface{}{
		&s.StepNumber,
		&s.Duration,
		&s.Description,
	}
}

type Review struct {
	CommentID   uint64    `json:"comment_id"`
	UserName    string    `json:"username"`
	Star        uint64    `json:"star"`
	CommentText string    `json:"comment_text"`
	CreatedDate time.Time `json:"created_date"`
}

func (r *Review) ScanFields() []interface{} {
	return []interface{}{
		&r.CommentID,
		&r.CommentText,
		&r.Star,
		&r.CreatedDate,
		&r.UserName,
	}
}
