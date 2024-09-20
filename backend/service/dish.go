package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Falokut/go-kit/log"

	"dish_as_a_service/domain"
	"dish_as_a_service/entity"
	"github.com/pkg/errors"
)

type DishRepo interface {
	List(ctx context.Context, limit, offset int32) ([]entity.Dish, error)
	GetDishesByIds(ctx context.Context, ids []int32) ([]entity.Dish, error)
	GetDishesByCategories(ctx context.Context, limit int32, offset int32, ids []int32) ([]entity.Dish, error)
	AddDish(ctx context.Context, dish *entity.AddDishRequest) error
}

type ImagesRepo interface {
	UploadImage(ctx context.Context, category string, image []byte) (string, error)
	DeleteImage(ctx context.Context, category, imageId string) error
	GetImageUrl(category, imageId string) string
}

const dishImageCategory = "dish"

type Dish struct {
	repo       DishRepo
	imagesRepo ImagesRepo
	logger     log.Logger
}

func NewDish(repo DishRepo, imagesRepo ImagesRepo, logger log.Logger) Dish {
	return Dish{
		repo:       repo,
		imagesRepo: imagesRepo,
		logger:     logger,
	}
}

func (s Dish) List(ctx context.Context, limit, offset int32) ([]domain.Dish, error) {
	if limit == 0 {
		limit = 30
	}
	dish, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, errors.WithMessage(err, "dish list")
	}
	converted := make([]domain.Dish, len(dish))
	for i, f := range dish {
		converted[i] = s.dishFromEntity(f)
	}
	return converted, nil
}

func (s Dish) GetByIds(ctx context.Context, ids []int32) ([]domain.Dish, error) {
	dish, err := s.repo.GetDishesByIds(ctx, ids)
	if err != nil {
		return nil, errors.WithMessage(err, "dish list by ids")
	}
	converted := make([]domain.Dish, len(dish))
	for i, f := range dish {
		converted[i] = s.dishFromEntity(f)
	}
	return converted, nil
}

func (s Dish) GetByCategories(ctx context.Context, limit, offset int32, ids []int32) ([]domain.Dish, error) {
	dish, err := s.repo.GetDishesByCategories(ctx, limit, offset, ids)
	if err != nil {
		return nil, errors.WithMessage(err, "dish list by ids")
	}
	converted := make([]domain.Dish, len(dish))
	for i, f := range dish {
		converted[i] = s.dishFromEntity(f)
	}
	return converted, nil
}

func (s Dish) AddDish(ctx context.Context, req domain.AddDishRequest) error {
	var imageId string
	var err error
	if len(req.Image) > 0 {
		imageId, err = s.imagesRepo.UploadImage(ctx, dishImageCategory, req.Image)
		if err != nil {
			return errors.WithMessage(err, "upload image")
		}
	}

	dish := &entity.AddDishRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Categories:  req.Categories,
		ImageId:     imageId,
	}
	err = s.repo.AddDish(ctx, dish)
	if err != nil {
		if len(req.Image) > 0 {
			delErr := s.imagesRepo.DeleteImage(ctx, dishImageCategory, imageId)
			if delErr != nil {
				s.logger.Error(ctx, fmt.Sprintf("error while deleting image with id=%s category=%s err=%v",
					imageId, dishImageCategory, err))
			}
		}
		return errors.WithMessage(err, "add dish")
	}
	return nil
}

func (s Dish) dishFromEntity(dish entity.Dish) domain.Dish {
	categories := []string{}
	if dish.Categories != "" {
		categories = strings.Split(dish.Categories, ",")
	}
	return domain.Dish{
		Id:          dish.Id,
		Name:        dish.Name,
		Description: dish.Description,
		Price:       dish.Price,
		Url:         s.imagesRepo.GetImageUrl(dishImageCategory, dish.ImageId),
		Categories:  categories,
	}
}
