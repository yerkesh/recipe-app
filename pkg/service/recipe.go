package service

type RecipeService struct {

}

func NewRecipeService() *RecipeService {
	return &RecipeService{}
}

func (svc *RecipeService) Hello() string {
	return "HELLO WORLD......\n REST is initialized"
}