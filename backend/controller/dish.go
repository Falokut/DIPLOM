package controller

import (
	"context"
	"net/http"
	"strconv"

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
// @Tags dishes
// @Summary dish
// @Description возвращает список блюд
// @Param query ids string false "список индентификаторов блюд через запятую"
// @Param query limit int false "максимальное количество блюд"
// @Param query offset int false "смещение"
// @Accept  json
// @Produce json
// @Success 200 {object} []domain.Dish
// @Failure 500 {string} string
// @Router /dishes [GET]

func (c Dish) List(ctx echo.Context) error {
	ids, _ := stringToIntSlice(ctx.QueryParam("ids"))
	var resp []domain.Dish
	var err error
	if len(ids) > 0 {
		resp, err = c.service.GetByIds(ctx.Request().Context(), ids)
	} else {
		limit, _ := strconv.ParseInt(ctx.QueryParam("limit"), 10, 32)
		if limit < 0 {
			return ctx.String(http.StatusBadRequest, "invalid limit value")
		}
		offset, _ := strconv.ParseInt(ctx.QueryParam("offset"), 10, 32)
		if offset < 0 {
			return ctx.String(http.StatusBadRequest, "invalid offset value")
		}
		resp, err = c.service.List(ctx.Request().Context(), int32(limit), int32(offset))
	}

	if err != nil {
		ctx.Error(err)
		return err
	}
	return ctx.JSON(http.StatusOK, resp)
}

// Add dish
// @Tags dishes
// @Summary Add Dish
// @Param body body domain.AddDishRequest true "request body"
// @Accept  json
// @Produce json
// @Success 200 {string} string
// @Failure 500 {string} string
// @Router /dishes [POST]

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
