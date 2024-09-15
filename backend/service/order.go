package service

import (
	"context"
	"strconv"
	"time"

	"dish_as_a_service/domain"
	"dish_as_a_service/entity"

	"github.com/Falokut/go-kit/http/apierrors"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/exp/maps"
)

type DishesRepo interface {
	GetDishesByIds(ctx context.Context, ids []int32) ([]entity.Dish, error)
}

type PaymentService interface {
	IsPaymentMethodValid(method string) bool
	Process(ctx context.Context, order *entity.Order, method string) (string, error)
}

type OrderRepo interface {
	GetOrder(ctx context.Context, orderId string) (*entity.Order, error)
	ProcessOrder(ctx context.Context, order *entity.Order) error
	UpdateOrderStatus(ctx context.Context, orderId string, newStatus string) error
	IsOrderCanceled(ctx context.Context, orderId string) (bool, error)
	SetOrderingAllowed(ctx context.Context, isAllowed bool) error
	IsOrderingAllowed(ctx context.Context) (bool, error)
}

type Order struct {
	paymentService PaymentService
	orderRepo      OrderRepo
	dishesRepo     DishesRepo
}

func NewOrder(paymentService PaymentService, orderRepo OrderRepo, dishesRepo DishesRepo) Order {
	return Order{
		paymentService: paymentService,
		orderRepo:      orderRepo,
		dishesRepo:     dishesRepo,
	}
}

func (s Order) UpdateOrderStatus(ctx context.Context, orderId, newStatus string) error {
	err := s.orderRepo.UpdateOrderStatus(ctx, orderId, newStatus)
	if err != nil {
		return errors.WithMessage(err, "update order status")
	}

	return nil
}

func (s Order) GetOrder(ctx context.Context, orderId string) (*entity.Order, error) {
	order, err := s.orderRepo.GetOrder(ctx, orderId)
	if err != nil {
		return nil, errors.WithMessage(err, "get order")
	}
	return order, nil
}

func (s Order) SetOrderingAllowed(ctx context.Context, isAllowed bool) error {
	err := s.orderRepo.SetOrderingAllowed(ctx, isAllowed)
	if err != nil {
		return errors.WithMessage(err, "set ordering allowed")
	}
	return nil
}

func (s Order) IsOrderingAllowed(ctx context.Context) (bool, error) {
	allowed, err := s.orderRepo.IsOrderingAllowed(ctx)
	if err != nil {
		return false, errors.WithMessage(err, "get ordering allowed")
	}
	return allowed, nil
}

func (s Order) ProcessOrder(ctx context.Context, req domain.ProcessOrderRequest) (string, error) {
	allowed, err := s.orderRepo.IsOrderingAllowed(ctx)
	if err != nil {
		return "", errors.WithMessage(err, "is ordering allowed")
	}
	if !allowed {
		return "", apierrors.NewBusinessError(
			domain.ErrCodeOrderingForbidden,
			domain.ErrOrderingForbidden.Error(),
			domain.ErrOrderingForbidden,
		)
	}
	if !s.paymentService.IsPaymentMethodValid(req.PaymentMethod) {
		return "", apierrors.NewBusinessError(domain.ErrCodeInvalidArgument, "invalid payment method", errors.New("invalid payment method"))
	}

	items, err := convertMapStringToInt(req.Items)
	if err != nil {
		return "", apierrors.NewBusinessError(domain.ErrCodeInvalidArgument, "invalid order items", err)
	}
	for _, count := range items {
		if count <= 0 {
			return "", domain.ErrInvalidDishCount
		}
	}

	dishes, err := s.dishesRepo.GetDishesByIds(ctx, maps.Keys(items))
	if err != nil {
		return "", errors.WithMessage(err, "get prices")
	}
	if len(dishes) != len(items) {
		return "", domain.ErrDishNotFound
	}

	dishesMap := make(map[int32]entity.Dish)
	for i := range dishes {
		dishesMap[dishes[i].Id] = dishes[i]
	}

	var total int32
	orderItems := make([]entity.OrderItem, 0, len(dishes))
	for id, count := range items {
		orderItems = append(orderItems, entity.OrderItem{
			DishId: id,
			Count:  count,
			Name:   dishesMap[id].Name,
			Price:  count * dishesMap[id].Price,
			Status: entity.OrderItemStatusProcess,
		})
		total += dishesMap[id].Price * count
	}
	order := &entity.Order{
		Id:            uuid.NewString(),
		PaymentMethod: req.PaymentMethod,
		Items:         orderItems,
		UserId:        req.UserId,
		Total:         total,
		Wishes:        req.Wishes,
		CreatedAt:     time.Now().UTC(),
	}

	err = s.orderRepo.ProcessOrder(ctx, order)
	if err != nil {
		return "", errors.WithMessage(err, "process order")
	}

	url, err := s.paymentService.Process(ctx, order, req.PaymentMethod)
	if err != nil {
		return "", errors.WithMessagef(err, "process payment, orderId=%v", order.Id)
	}
	return url, nil
}

func (s Order) IsOrderCanceled(ctx context.Context, orderId string) (bool, error) {
	canceled, err := s.orderRepo.IsOrderCanceled(ctx, orderId)
	if err != nil {
		return false, errors.WithMessage(err, "is order canceled")
	}
	return canceled, nil
}

func convertMapStringToInt(m map[string]int32) (map[int32]int32, error) {
	res := make(map[int32]int32)
	for k, v := range m {
		intK, err := strconv.ParseInt(k, 10, 32)
		if err != nil {
			return nil, errors.WithMessage(err, "invalid value")
		}
		res[int32(intK)] = v
	}
	return res, nil
}
