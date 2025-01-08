package service

import (
	"context"
	"strings"

	"github.com/Falokut/go-kit/log"
	"github.com/google/uuid"

	"dish_as_a_service/domain"
	"dish_as_a_service/entity"

	"github.com/pkg/errors"
)

type DishRepo interface {
	List(ctx context.Context, limit, offset int32) ([]entity.Dish, error)
	GetDishesByIds(ctx context.Context, ids []int32) ([]entity.Dish, error)
	GetDishesByCategories(ctx context.Context, limit int32, offset int32, ids []int32) ([]entity.Dish, error)
	AddDish(ctx context.Context, dish *entity.AddDishRequest) error
	EditDish(ctx context.Context, dish *entity.EditDishRequest) error
	DeleteDish(ctx context.Context, id int32) error
}

type ImagesRepo interface {
	UploadImage(ctx context.Context, req entity.UploadFileRequest) error
	DeleteImage(ctx context.Context, category, imageId string) error
	GetImageUrl(category, imageId string) string
}

const dishImageCategory = "image-dish"

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
		imageId = uuid.NewString()
		uploadImageReq := entity.UploadFileRequest{
			Category: dishImageCategory,
			Filename: imageId,
			Content:  req.Image,
		}
		err = s.imagesRepo.UploadImage(ctx, uploadImageReq)
		if err != nil {
			return errors.WithMessage(err, "upload image")
		}
	}

	err = s.repo.AddDish(ctx, &entity.AddDishRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Categories:  req.Categories,
		ImageId:     imageId,
	})
	if err != nil {
		if len(req.Image) > 0 {
			delErr := s.imagesRepo.DeleteImage(ctx, dishImageCategory, imageId)
			if delErr != nil {
				s.logger.Error(ctx, "delete image",
					log.Any("imageId", imageId),
					log.Any("category", dishImageCategory),
					log.Any("error", delErr),
				)
			}
		}
		return errors.WithMessage(err, "add dish")
	}
	return nil
}

func (s Dish) EditDish(ctx context.Context, req domain.EditDishRequest) error {
	dishes, err := s.repo.GetDishesByIds(ctx, []int32{req.Id})
	if err != nil {
		return errors.WithMessage(err, "get dishes by ids")
	}
	if len(dishes) == 0 {
		return domain.ErrDishNotFound
	}
	if dishes[0].ImageId != "" {
		err = s.imagesRepo.DeleteImage(ctx, dishImageCategory, dishes[0].ImageId)
		if err != nil {
			s.logger.Warn(ctx, "delete image",
				log.Any("imageId", dishes[0].ImageId),
				log.Any("category", dishImageCategory),
				log.Any("error", err),
			)
			return errors.WithMessage(err, "delete dish image")
		}
	}
	var imageId string
	if len(req.Image) > 0 {
		imageId = uuid.NewString()
		uploadImageReq := entity.UploadFileRequest{
			Category: dishImageCategory,
			Filename: imageId,
			Content:  req.Image,
		}
		err = s.imagesRepo.UploadImage(ctx, uploadImageReq)
		if err != nil {
			return errors.WithMessage(err, "upload image")
		}
	}
	err = s.repo.EditDish(ctx, &entity.EditDishRequest{
		Id:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		ImageId:     imageId,
		Categories:  req.Categories,
		Price:       req.Price,
	})
	if err != nil {
		return errors.WithMessage(err, "edit dish")
	}
	return nil
}

func (s Dish) DeleteDish(ctx context.Context, id int32) error {
	dishes, err := s.repo.GetDishesByIds(ctx, []int32{id})
	if err != nil {
		return errors.WithMessage(err, "get dishes by ids")
	}
	if len(dishes) == 0 {
		return domain.ErrDishNotFound
	}
	if dishes[0].ImageId != "" {
		err = s.imagesRepo.DeleteImage(ctx, dishImageCategory, dishes[0].ImageId)
		if err != nil {
			s.logger.Warn(ctx, "delete image",
				log.Any("imageId", dishes[0].ImageId),
				log.Any("category", dishImageCategory),
				log.Any("error", err),
			)
			return errors.WithMessage(err, "delete dish image")
		}
	}
	err = s.repo.DeleteDish(ctx, id)
	if err != nil {
		return errors.WithMessage(err, "delete dish")
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
