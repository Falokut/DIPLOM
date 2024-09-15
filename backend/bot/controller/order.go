package controller

import (
	"context"
	"encoding/json"

	"dish_as_a_service/domain"
	"dish_as_a_service/entity"

	"github.com/Falokut/go-kit/telegram_bot"
	"github.com/Falokut/go-kit/telegram_bot/apierrors"
)

type OrderService interface {
	GetOrder(ctx context.Context, orderId string) (*entity.Order, error)
	UpdateOrderStatus(ctx context.Context, orderId string, newStatus string) error
	IsOrderCanceled(ctx context.Context, orderId string) (bool, error)
}

type OrderUserService interface {
	NotifySuccessPayment(ctx context.Context, req *entity.Order) error
}

type Order struct {
	orderService OrderService
	userService  OrderUserService
}

func NewOrder(service OrderService, userService OrderUserService) Order {
	return Order{
		orderService: service,
		userService:  userService,
	}
}

func (c Order) HandlePayment(ctx context.Context, update telegram_bot.Update) (telegram_bot.Chattable, error) {
	msg := update.Message
	if msg.SuccessfulPayment == nil {
		return nil, nil //nolint:nilnil
	}
	var payload entity.PaymentPayload
	err := json.Unmarshal([]byte(msg.SuccessfulPayment.InvoicePayload), &payload)
	if err != nil {
		return nil, apierrors.NewBusinessError(domain.ErrCodeInvalidArgument, "invalid payment payload", err)
	}

	err = c.orderService.UpdateOrderStatus(ctx, payload.OrderId, entity.OrderItemStatusSuccess)
	if err != nil {
		return nil, err
	}

	order, err := c.orderService.GetOrder(ctx, payload.OrderId)
	if err != nil {
		return nil, err
	}
	err = c.userService.NotifySuccessPayment(ctx, order)
	if err != nil {
		return nil, err
	}
	return nil, nil // nolint:nilnil
}

// nolint:nilerr
func (c Order) HandlePreCheckout(ctx context.Context, update telegram_bot.Update) (telegram_bot.Chattable, error) {
	query := update.PreCheckoutQuery
	var payload entity.PaymentPayload
	err := json.Unmarshal([]byte(query.InvoicePayload), &payload)
	if err != nil {
		return telegram_bot.PreCheckoutConfig{
			PreCheckoutQueryID: query.Id,
			OK:                 false,
			ErrorMessage:       "invalid payload",
		}, nil
	}

	canceled, err := c.orderService.IsOrderCanceled(ctx, payload.OrderId)
	if err != nil {
		return nil, apierrors.NewInternalServiceError(err)
	}
	if canceled {
		return telegram_bot.PreCheckoutConfig{
			PreCheckoutQueryID: query.Id,
			OK:                 false,
			ErrorMessage:       "order canceled",
		}, nil
	}

	return telegram_bot.PreCheckoutConfig{
		PreCheckoutQueryID: query.Id,
		OK:                 true,
	}, nil
}
