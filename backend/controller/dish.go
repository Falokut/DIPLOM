package controller

import (
	"context"

	"dish_as_a_service/domain"

	"github.com/Falokut/go-kit/http/apierrors"
	_ "github.com/Falokut/go-kit/http/apierrors"
)

type DishService interface {
	List(ctx context.Context, limit, offset int32) ([]domain.Dish, error)
	GetByIds(ctx context.Context, ids []int32) ([]domain.Dish, error)
	AddDish(ctx context.Context, req domain.AddDishRequest) error
}

type Dish struct {
	service DishService
}

func NewDish(service DishService) Dish {
	return Dish{
		service: service,
	}
}

// List
//
//	@Tags			dishes
//	@Summary		dish
//	@Description	возвращает список блюд
//	@Param			ids		query	string	false	"список идентификаторов блюд через запятую"
//	@Param			limit	query	int		false	"максимальное количество блюд"
//	@Param			offset	query	int		false	"смещение"
//	@Produce		json
//	@Success		200	{array}		domain.Dish
//	@Failure		400	{object}	apierrors.Error
//	@Failure		500	{object}	apierrors.Error
//	@Router			/dishes [GET]
func (c Dish) List(ctx context.Context, req domain.GetDishesRequest) ([]domain.Dish, error) {
	ids, err := stringToIntSlice(req.Ids)
	if err != nil {
		return nil, apierrors.NewBusinessError(domain.ErrCodeInvalidArgument, "invalid ids", err)
	}

	var dishes []domain.Dish
	if len(ids) > 0 {
		dishes, err = c.service.GetByIds(ctx, ids)
	} else {
		dishes, err = c.service.List(ctx, req.Limit, req.Offset)
	}

	if err != nil {
		return nil, err
	}
	return dishes, nil
}

// Add dish
//
//	@Tags		dishes
//	@Summary	Add Dish
//	@Param		body		body	domain.AddDishRequest	true	"request body"
//
//	@Param		X-USER-ID	header	string					true	"id пользователя"
//
//	@Accept		json
//	@Success	200	{object}	domain.Empty
//	@Failure	403	{object}	apierrors.Error
//	@Failure	500	{object}	apierrors.Error
//	@Router		/dishes [POST]
func (c Dish) AddDish(ctx context.Context, req domain.AddDishRequest) error {
	err := c.service.AddDish(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
