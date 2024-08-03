package controller

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"

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
// @Tags order
// @Summary Заказать
// @Accept  json
// @Produce json
// @Success 200 {object} domain.ProcessOrderResponse
// @Failure 500 {string} string
// @Router /orders [POST]

func (c Order) ProcessOrder(ctx echo.Context) error {
	req, err := bindRequest[domain.ProcessOrderRequest](ctx)
	if err != nil {
		return err
	}

	url, err := c.service.ProcessOrder(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, domain.ProcessOrderResponse{PaymentUrl: url})
}
