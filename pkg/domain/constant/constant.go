package constant

const (
	MsgCreated  = "Объект создан"
	MsgUpdated  = "Успешно обновлен"
	MsgDeleted  = "Успешно удален"
	MsgPatched  = "Успешно частично обновлен"
	MsgSuccess  = "Успешно"
	MsgRequired = "required"
)

// Request errors.
const (
	MsgNotFoundErr    = "Не найдена запись в БД"
	MsgRequiredErr    = "Не отправлены обязательные поля"
	MsgUnhandledErr   = "Непредвиденная ошибка"
	MsgRequestBodyErr = "Переданы некорректные данные"
	MsgAuthorizeErr   = "Ошибка авторизации"
	MsgAlreadyExists  = "Такая запись уже существует в БД"
)

// Tablenames.
type Table string

const (
	TblCategory                 Table = "category"
	TblCategoryUser             Table = "category_user"
	TblComplexity               Table = "complexity"
	TblCuisine                  Table = "cuisine"
	TblCuisineUser              Table = "cuisine_user"
	TblDailyRecommendation      Table = "daily_recommendation"
	TblIngredient               Table = "ingredient"
	TblIngredientRecipe         Table = "ingredient_recipe"
	TblRecipe                   Table = "recipe"
	TblRecipeUserRecommendation Table = "recipe_user_recommendation"
	TblRefreshToken             Table = "refreshtoken"
	TblSelectedRecipeByUser     Table = "selected_recipe_by_user"
	TblUnitOfMeasurement        Table = "unit_of_measurement"
	TblUsers                    Table = "users"
	TblRecipeStep               Table = "recipe_step"
	TblComment                  Table = "comment"
	TblRate                     Table = "rate"
	TblUserFavourite            Table = "user_favourite"
)

func (t Table) As(as ...string) string {
	if len(as) == 0 {
		return t.String()
	}

	return t.String() + " " + as[0]
}

func (t Table) String() string {
	return string(t)
}

type junction struct {
	left  Table
	right Table
}

func (j junction) swap() junction {
	return junction{left: j.right, right: j.left}
}

type SQLAction string

const (
	Insert        SQLAction = "insert"
	Update        SQLAction = "update"
	PartialUpdate SQLAction = "patch"
	Delete        SQLAction = "delete"
	Replace       SQLAction = "replace"
	Select        SQLAction = "select"
)

func (a SQLAction) String() string {
	return string(a)
}
