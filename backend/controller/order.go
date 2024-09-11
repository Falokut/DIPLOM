package controller

import (
	"context"
	"errors"
	"net/http"

	"github.com/Falokut/go-kit/http/apierrors"

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
//	@Failure	400	{object}	apierrors.Error
//	@Failure	404	{object}	apierrors.Error
//	@Failure	500	{object}	apierrors.Error
//	@Router		/orders [POST]
func (c Order) ProcessOrder(ctx context.Context, req domain.ProcessOrderRequest) (*domain.ProcessOrderResponse, error) {
	url, err := c.service.ProcessOrder(ctx, req)
	switch {
	case errors.Is(err, domain.ErrDishNotFound):
		return nil, apierrors.New(http.StatusNotFound, domain.ErrCodeDishNotFound, domain.ErrDishNotFound.Error(), err)
	case errors.Is(err, domain.ErrInvalidDishCount):
		return nil, apierrors.NewBusinessError(domain.ErrCodeInvalidDishCount, domain.ErrInvalidDishCount.Error(), err)
	case err != nil:
		return nil, err
	default:
		return &domain.ProcessOrderResponse{PaymentUrl: url}, nil
	}
}
