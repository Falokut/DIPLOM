package controller

import (
	"context"
	"dish_as_a_service/domain"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type DishesCategoriesService interface {
	GetCategories(ctx context.Context) ([]domain.DishCategory, error)
	GetCategory(ctx context.Context, id int32) (domain.DishCategory, error)
	AddCategory(ctx context.Context, category string) (int32, error)
	RenameCategory(ctx context.Context, req domain.RenameCategoryRequest) error
	DeleteCategory(ctx context.Context, id int32) error
}
type DishesCategories struct {
	service DishesCategoriesService
}

func NewDishesCategories(service DishesCategoriesService) DishesCategories {
	return DishesCategories{service: service}
}

// Get categories
//
//	@Tags		dishes_categories
//	@Summary	Получить категории
//	@Accept		json
//	@Produce	json
//	@Success	200	{array}	domain.DishCategory
//	@Failure	500	{string}	string
//	@Router		/dishes/categories [GET]
func (c DishesCategories) GetCategories(ctx echo.Context) error {
	categories, err := c.service.GetCategories(ctx.Request().Context())
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, categories)
}

// Get category
//
//	@Param		body	body	domain.GetDishesCategory	true	"request body"
//
//	@Tags		dishes_categories
//	@Summary	Получить категорию
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	domain.DishCategory
//	@Failure	400	{string}	string
//	@Failure	404	{string}	string
//	@Failure	500	{string}	string
//	@Router		/dishes/categories/:id [GET]
func (c DishesCategories) GetCategory(ctx echo.Context) error {
	req, err := bindRequest[domain.GetDishesCategory](ctx)
	if err != nil {
		return err
	}
	category, err := c.service.GetCategory(ctx.Request().Context(), req.Id)
	switch {
	case errors.Is(err, domain.ErrDishCategoryNotFound):
		return ctx.String(http.StatusNotFound, domain.ErrDishCategoryNotFound.Error())
	case err != nil:
		return err
	default:
		return ctx.JSON(http.StatusOK, category)
	}
}

// Add category
//
//	@Param		body	body	domain.AddCategoryRequest	true	"request body"
//
// @Param X-USER-ID header string true "id пользователя"
//
//	@Tags		dishes_categories
//	@Summary	Создать категорию
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	domain.DishCategory
//	@Failure	400	{string}	string
//	@Failure	403	{string}	string
//	@Failure	500	{string}	string
//	@Router		/dishes/categories [POST]
func (c DishesCategories) AddCategory(ctx echo.Context) error {
	req, err := bindRequest[domain.AddCategoryRequest](ctx)
	if err != nil {
		return err
	}
	id, err := c.service.AddCategory(ctx.Request().Context(), req.Name)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, domain.AddCategoryResponse{Id: id})
}

// Rename category
//
//	@Param		body	body	domain.RenameCategoryRequest	true	"request body"
//
// @Param X-USER-ID header string true "id пользователя"
//
//	@Tags		dishes_categories
//	@Summary	Переименовать категорию
//	@Accept		json
//	@Produce	json
//	@Success	204	{object}	domain.Empty
//	@Failure	400	{string}	string
//	@Failure	403	{string}	string
//	@Failure	500	{string}	string
//	@Router		/dishes/categories/:id [POST]
func (c DishesCategories) RenameCategory(ctx echo.Context) error {
	req, err := bindRequest[domain.RenameCategoryRequest](ctx)
	if err != nil {
		return err
	}
	err = c.service.RenameCategory(ctx.Request().Context(), req)
	switch {
	case errors.Is(err, domain.ErrDishCategoryNotFoundOrConflict):
		return ctx.String(http.StatusBadRequest, domain.ErrDishCategoryNotFoundOrConflict.Error())
	case err != nil:
		return err
	default:
		return ctx.NoContent(http.StatusNoContent)
	}
}

// Delete category
//
//	@Param		body	body	domain.DeleteCategoryRequest	true	"request body"
//
// @Param X-USER-ID header string true "id пользователя"
//
//	@Tags		dishes_categories
//	@Summary	Удалить категорию
//	@Accept		json
//	@Produce	json
//	@Success	204	{object}	domain.Empty
//	@Failure	400	{string}	string
//	@Failure	500	{string}	string
//	@Router		/dishes/categories/:id [DELETE]
func (c DishesCategories) DeleteCategory(ctx echo.Context) error {
	req, err := bindRequest[domain.DeleteCategoryRequest](ctx)
	if err != nil {
		return err
	}
	err = c.service.DeleteCategory(ctx.Request().Context(), req.Id)
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}
