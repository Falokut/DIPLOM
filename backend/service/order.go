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
	SetOrderStatus(ctx context.Context, orderId string, oldStatus string, newStatus string) error
	GetOrderStatus(ctx context.Context, orderId string) (string, error)
	SetOrderingAllowed(ctx context.Context, isAllowed bool) error
	IsOrderingAllowed(ctx context.Context) (bool, error)
	GetUserOrders(ctx context.Context, userId string, limit int32, offset int32) ([]entity.Order, error)
}

const (
	defaultUserOrdersLimit = 30
)

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

func (s Order) SetOrderStatus(ctx context.Context, orderId string, oldStatus string, newStatus string) error {
	err := s.orderRepo.SetOrderStatus(ctx, orderId, oldStatus, newStatus)
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

func (s Order) ProcessOrder(ctx context.Context, userId string, req domain.ProcessOrderRequest) (string, error) {
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
		})
		total += dishesMap[id].Price * count
	}
	order := &entity.Order{
		Id:            uuid.NewString(),
		PaymentMethod: req.PaymentMethod,
		Items:         orderItems,
		UserId:        userId,
		Total:         total,
		Wishes:        req.Wishes,
		Status:        entity.OrderItemStatusProcess,
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

func (s Order) GetOrderStatus(ctx context.Context, orderId string) (string, error) {
	orderStatus, err := s.orderRepo.GetOrderStatus(ctx, orderId)
	if err != nil {
		return "", errors.WithMessage(err, "get order status")
	}
	return orderStatus, nil
}

func (s Order) GetUserOrders(ctx context.Context, userId string, req domain.GetMyOrdersRequest) ([]domain.UserOrder, error) {
	limit := req.Limit
	if limit == 0 {
		limit = defaultUserOrdersLimit
	}
	orders, err := s.orderRepo.GetUserOrders(ctx, userId, limit, req.Offset)
	if err != nil {
		return nil, errors.WithMessage(err, "get user orders")
	}
	var userOrders = make([]domain.UserOrder, len(orders))
	for i, order := range orders {
		items := make([]domain.OrderItem, len(order.Items))
		for j, item := range order.Items {
			items[j] = domain.OrderItem{
				DishId:     item.DishId,
				Name:       item.Name,
				Price:      item.Price,
				Count:      item.Count,
				TotalPrice: item.Count * item.Price,
			}
		}
		userOrders[i] = domain.UserOrder{
			Id:            order.Id,
			Items:         items,
			PaymentMethod: order.PaymentMethod,
			Total:         order.Total,
			Wishes:        order.Wishes,
			CreatedAt:     order.CreatedAt,
			Status:        order.Status,
		}
	}
	return userOrders, nil
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
