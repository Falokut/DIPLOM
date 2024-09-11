package domain

type DishCategory struct {
	Id   int32
	Name string
}

type AddCategoryRequest struct {
	Name string
}

type AddCategoryResponse struct {
	Id int32
}

type RenameCategoryRequest struct {
	Id   int32
	Name string
}

type DeleteCategoryRequest struct {
	Id int32
}

type GetDishesCategory struct {
	Id int32
}
