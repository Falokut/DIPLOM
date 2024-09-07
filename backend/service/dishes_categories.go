package service

import (
	"context"
	"dish_as_a_service/domain"
	"dish_as_a_service/entity"

	"github.com/pkg/errors"
)

type DishesCategoriesRepo interface {
	GetCategories(ctx context.Context) ([]entity.DishCategory, error)
	GetCategory(ctx context.Context, id int32) (entity.DishCategory, error)
	AddCategory(ctx context.Context, category string) (int32, error)
	RenameCategory(ctx context.Context, id int32, newName string) error
	DeleteCategory(ctx context.Context, id int32) error
}

type DishesCategories struct {
	repo DishesCategoriesRepo
}

func NewDishesCategories(repo DishesCategoriesRepo) DishesCategories {
	return DishesCategories{repo: repo}
}

func (s DishesCategories) GetCategories(ctx context.Context) ([]domain.DishCategory, error) {
	categories, err := s.repo.GetCategories(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "get categories")
	}
	domainCategories := make([]domain.DishCategory, len(categories))
	for i, category := range categories {
		domainCategories[i] = domain.DishCategory{
			Id:   category.Id,
			Name: category.Name,
		}
	}
	return domainCategories, nil
}

func (s DishesCategories) GetCategory(ctx context.Context, id int32) (domain.DishCategory, error) {
	category, err := s.repo.GetCategory(ctx, id)
	if err != nil {
		return domain.DishCategory{}, errors.WithMessage(err, "get category")
	}
	return domain.DishCategory{
		Id:   category.Id,
		Name: category.Name,
	}, nil
}
func (s DishesCategories) AddCategory(ctx context.Context, category string) (int32, error) {
	id, err := s.repo.AddCategory(ctx, category)
	if err != nil {
		return 0, errors.WithMessage(err, "add category")
	}
	return id, nil
}

func (s DishesCategories) RenameCategory(ctx context.Context, req domain.RenameCategoryRequest) error {
	err := s.repo.RenameCategory(ctx, req.Id, req.Name)
	if err != nil {
		return errors.WithMessage(err, "rename category")
	}
	return nil
}

func (s DishesCategories) DeleteCategory(ctx context.Context, id int32) error {
	err := s.repo.DeleteCategory(ctx, id)
	if err != nil {
		return errors.WithMessage(err, "delete category")
	}
	return nil
}
