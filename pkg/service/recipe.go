package service

type RecipeService struct {
}

func NewRecipeService() *RecipeService {
	return &RecipeService{}
}

func (svc *RecipeService) Hello() string {
	return "Daria keep calm and study well"
}
