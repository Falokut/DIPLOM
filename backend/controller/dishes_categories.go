package controller

import (
	"context"
	"dish_as_a_service/domain"
	"errors"
	"net/http"

	"github.com/Falokut/go-kit/http/apierrors"
)

type DishesCategoriesService interface {
	GetDishesCategories(ctx context.Context) ([]domain.DishCategory, error)
	GetAllCategories(ctx context.Context) ([]domain.DishCategory, error)
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
//	@Success	200	{array}		domain.DishCategory
//	@Failure	500	{object}	apierrors.Error
//	@Router		/dishes/all_categories [GET]
func (c DishesCategories) GetAllCategories(ctx context.Context) ([]domain.DishCategory, error) {
	return c.service.GetAllCategories(ctx)
}

// Get categories
//
//	@Tags		dishes_categories
//	@Summary	Получить категории
//	@Accept		json
//	@Produce	json
//	@Success	200	{array}		domain.DishCategory
//	@Failure	500	{object}	apierrors.Error
//	@Router		/dishes/categories [GET]
func (c DishesCategories) GetDishesCategories(ctx context.Context) ([]domain.DishCategory, error) {
	return c.service.GetDishesCategories(ctx)
}

// Get category
//
//	@Param		body	body	domain.GetDishesCategory	true	"request body"
//
//	@Tags		dishes_categories
//	@Summary	Получить категорию
//	@Produce	json
//
//	@Param		id	path		int	true	"Идентификатор категории"
//
//	@Success	200	{object}	domain.DishCategory
//	@Failure	400	{object}	apierrors.Error
//	@Failure	404	{object}	apierrors.Error
//	@Failure	500	{object}	apierrors.Error
//	@Router		/dishes/categories/:id [GET]
func (c DishesCategories) GetCategory(ctx context.Context, req domain.GetDishesCategory) (*domain.DishCategory, error) {
	category, err := c.service.GetCategory(ctx, req.Id)
	switch {
	case errors.Is(err, domain.ErrDishCategoryNotFound):
		return nil, apierrors.New(http.StatusNotFound,
			domain.ErrCodeDishCategoryNotFound,
			domain.ErrDishCategoryNotFound.Error(),
			err,
		)
	case err != nil:
		return nil, err
	default:
		return &category, nil
	}
}

// Add category
//
//	@Param		body		body	domain.AddCategoryRequest	true	"request body"
//
//	@Param		X-USER-ID	header	string						true	"id пользователя"
//
//	@Tags		dishes_categories
//	@Summary	Создать категорию
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	domain.DishCategory
//	@Failure	400	{string}	apierrors.Error
//	@Failure	403	{object}	apierrors.Error
//	@Failure	500	{object}	apierrors.Error
//	@Router		/dishes/categories [POST]
func (c DishesCategories) AddCategory(ctx context.Context, req domain.AddCategoryRequest) (*domain.AddCategoryResponse, error) {
	id, err := c.service.AddCategory(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &domain.AddCategoryResponse{Id: id}, nil
}

// Rename category
//
//	@Param		body		body	domain.RenameCategoryRequest	true	"request body"
//
//	@Param		X-USER-ID	header	string							true	"id пользователя"
//
//	@Tags		dishes_categories
//	@Summary	Переименовать категорию
//	@Accept		json
//	@Produce	json
//
//	@Param		id		path		int								true	"Идентификатор категории"
//
//	@Param		body	body		domain.RenameCategoryRequest	true	"request body"
//	@Success	204		{object}	domain.Empty
//	@Failure	400		{object}	apierrors.Error
//	@Failure	403		{object}	apierrors.Error
//	@Failure	500		{object}	apierrors.Error
//	@Router		/dishes/categories/:id [POST]
func (c DishesCategories) RenameCategory(ctx context.Context, req domain.RenameCategoryRequest) error {
	err := c.service.RenameCategory(ctx, req)
	switch {
	case errors.Is(err, domain.ErrDishCategoryConflict):
		return apierrors.New(
			http.StatusConflict,
			domain.ErrCodeDishCategoryConflict,
			domain.ErrDishCategoryConflict.Error(),
			err,
		)
	case err != nil:
		return err
	default:
		return nil
	}
}

// Delete category
//
//	@Param		body		body	domain.DeleteCategoryRequest	true	"request body"
//
//	@Param		X-USER-ID	header	string							true	"id пользователя"
//
//	@Tags		dishes_categories
//	@Summary	Удалить категорию
//	@Accept		json
//	@Produce	json
//	@Success	204	{object}	domain.Empty
//	@Failure	400	{object}	apierrors.Error
//	@Failure	500	{object}	apierrors.Error
//	@Router		/dishes/categories/:id [DELETE]
func (c DishesCategories) DeleteCategory(ctx context.Context, req domain.DeleteCategoryRequest) error {
	return c.service.DeleteCategory(ctx, req.Id)
}
