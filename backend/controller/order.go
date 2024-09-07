package controller

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"dish_as_a_service/domain"
)

type OrderService interface {
	ProcessOrder(ctx context.Context, req domain.ProcessOrderRequest) (string, error)
}

type Order struct {
	service OrderService
}

func NewOrder(service OrderService) Order {
	return Order{
		service: service,
	}
}

// Process order
//
//	@Tags		order
//	@Summary	Заказать
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	domain.ProcessOrderResponse
//	@Failure	400	{string}	string
//	@Failure	404	{string}	string
//	@Failure	500	{string}	string
//	@Router		/orders [POST]
func (c Order) ProcessOrder(ctx echo.Context) error {
	req, err := bindRequest[domain.ProcessOrderRequest](ctx)
	if err != nil {
		return err
	}

	url, err := c.service.ProcessOrder(ctx.Request().Context(), req)
	switch {
	case errors.Is(err, domain.ErrDishNotFound):
		return ctx.String(http.StatusNotFound, domain.ErrDishNotFound.Error())
	case errors.Is(err, domain.ErrInvalidDishCount):
		return ctx.String(http.StatusBadRequest, domain.ErrInvalidDishCount.Error())
	case err != nil:
		return err
	default:
		return ctx.JSON(http.StatusOK, domain.ProcessOrderResponse{PaymentUrl: url})
	}
}
