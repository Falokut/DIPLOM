package controller

import (
	"context"
	"encoding/json"

	tgbotapi "dish_as_a_service/bot/api"
	"dish_as_a_service/entity"
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

func (c Order) HandlePayment(ctx context.Context, msg *tgbotapi.Message) tgbotapi.Chattable {
	var payload entity.PaymentPayload
	err := json.Unmarshal([]byte(msg.SuccessfulPayment.InvoicePayload), &payload)
	if err != nil {
		return nil
	}

	err = c.orderService.UpdateOrderStatus(ctx, payload.OrderId, entity.OrderItemStatusSuccess)
	if err != nil {
		return HandleError(msg, err, false)
	}

	order, err := c.orderService.GetOrder(ctx, payload.OrderId)
	if err != nil {
		return HandleError(msg, err, false)
	}
	err = c.userService.NotifySuccessPayment(ctx, order)
	if err != nil {
		return HandleError(msg, err, false)
	}
	return nil
}

func (c Order) HandlePreCheckout(ctx context.Context, query *tgbotapi.PreCheckoutQuery) tgbotapi.Chattable {
	var payload entity.PaymentPayload
	err := json.Unmarshal([]byte(query.InvoicePayload), &payload)
	if err != nil {
		return tgbotapi.PreCheckoutConfig{
			PreCheckoutQueryID: query.ID,
			OK:                 false,
			ErrorMessage:       "invalid payload",
		}
	}

	canceled, err := c.orderService.IsOrderCanceled(ctx, payload.OrderId)
	if err != nil {
		return tgbotapi.PreCheckoutConfig{
			PreCheckoutQueryID: query.ID,
			OK:                 false,
			ErrorMessage:       "internal error",
		}
	}
	if canceled {
		return tgbotapi.PreCheckoutConfig{
			PreCheckoutQueryID: query.ID,
			OK:                 false,
			ErrorMessage:       "order canceled",
		}
	}

	return tgbotapi.PreCheckoutConfig{
		PreCheckoutQueryID: query.ID,
		OK:                 true,
	}
}
