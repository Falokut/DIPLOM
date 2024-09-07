package controller

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"dish_as_a_service/domain"
)

type DishService interface {
	List(ctx context.Context, limit, offset int32) ([]domain.Dish, error)
	GetByIds(ctx context.Context, ids []int32) ([]domain.Dish, error)
	AddDish(ctx context.Context, req *domain.AddDishRequest) error
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
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		domain.Dish
//	@Failure		500	{string}	string
//	@Router			/dishes [GET]
func (c Dish) List(ctx echo.Context) error {
	req, err := bindRequest[domain.GetDishesRequest](ctx)
	if err != nil {
		return err
	}
	ids, _ := stringToIntSlice(req.Ids)
	var resp []domain.Dish
	if len(ids) > 0 {
		resp, err = c.service.GetByIds(ctx.Request().Context(), ids)
	} else {
		resp, err = c.service.List(ctx.Request().Context(), req.Limit, req.Offset)
	}

	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, resp)
}

// Add dish
//
//	@Tags		dishes
//	@Summary	Add Dish
//	@Param		body	body	domain.AddDishRequest	true	"request body"
//
// @Param X-USER-ID header string true "id пользователя"
//
//	@Accept		json
//	@Produce	json
//	@Success	200	{string}	string
//	@Failure	403	{string}	string
//	@Failure	500	{string}	string
//	@Router		/dishes [POST]
func (c Dish) AddDish(ctx echo.Context) error {
	req, err := bindRequest[domain.AddDishRequest](ctx)
	if err != nil {
		return err
	}

	err = c.service.AddDish(ctx.Request().Context(), &req)
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusOK)
}
