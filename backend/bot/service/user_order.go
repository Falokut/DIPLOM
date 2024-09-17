package service

import (
	"context"
	"dish_as_a_service/entity"
	"fmt"
	"strings"
	"time"

	"github.com/Falokut/go-kit/telegram_bot"
	"github.com/pkg/errors"
)

type UserRepo interface {
	GetUserInfo(ctx context.Context, userId string) (entity.User, error)
	GetAdminsChatsIds(ctx context.Context) ([]int64, error)
	GetUserChatId(ctx context.Context, userId string) (int64, error)
}

type OrderRepo interface {
	GetOrderedChatId(ctx context.Context, orderId string) (int64, error)
	SetOrderStatus(ctx context.Context, orderId, oldStatus, newStatus string) error
}

type UserOrder struct {
	bot       *telegram_bot.BotAPI
	userRepo  UserRepo
	orderRepo OrderRepo
}

func NewOrderUserService(bot *telegram_bot.BotAPI, userRepo UserRepo, orderRepo OrderRepo) UserOrder {
	return UserOrder{
		bot:       bot,
		userRepo:  userRepo,
		orderRepo: orderRepo,
	}
}

func (s UserOrder) NotifySuccessPayment(ctx context.Context, order *entity.Order) error {
	adminIds, err := s.userRepo.GetAdminsChatsIds(ctx)
	if err != nil {
		return errors.WithMessage(err, "get admins chats ids")
	}
	user, err := s.userRepo.GetUserInfo(ctx, order.UserId)
	if err != nil {
		return errors.WithMessage(err, "get user info")
	}

	orderInfoString := s.getOrderInfoString(order, &user)
	markup := s.getMarkupForOrder(order)
	for _, chatId := range adminIds {
		message := telegram_bot.NewMessage(chatId, orderInfoString)
		message.ReplyMarkup = markup
		err = s.bot.Send(message)
		if err != nil {
			return errors.WithMessagef(err, "send notification to chat: %d", chatId)
		}
	}
	return nil
}

func (s UserOrder) getMarkupForOrder(order *entity.Order) telegram_bot.InlineKeyboardMarkup {
	arrivalPayload := entity.QueryCallbackPayload{
		Command: entity.NotifyArrivalCommand,
		OrderId: order.Id,
	}
	notifyArrivalButton := telegram_bot.NewInlineKeyboardButtonData(
		"оповестить о прибытии заказа",
		arrivalPayload.String(),
	)
	cancelPayload := entity.QueryCallbackPayload{
		Command: entity.CancelOrderCommand,
		OrderId: order.Id,
	}
	cancelButton := telegram_bot.NewInlineKeyboardButtonData(
		"отменить заказ",
		cancelPayload.String(),
	)

	markup := telegram_bot.NewInlineKeyboardMarkup(
		[]telegram_bot.InlineKeyboardButton{notifyArrivalButton},
		[]telegram_bot.InlineKeyboardButton{cancelButton},
	)
	return markup
}

func (s UserOrder) NotifyOrderArrival(ctx context.Context, req entity.QueryCallbackPayload) error {
	chatId, err := s.orderRepo.GetOrderedChatId(ctx, req.OrderId)
	if err != nil {
		return errors.WithMessage(err, "get user chat id")
	}

	button := telegram_bot.NewInlineKeyboardButtonData("подтвердить получение",
		entity.QueryCallbackPayload{Command: entity.SuccessOrderCommand, OrderId: req.OrderId}.String(),
	)
	message := telegram_bot.NewMessage(chatId, fmt.Sprintf("Заказ №%s прибыл", req.OrderId))
	message.ReplyMarkup = telegram_bot.NewInlineKeyboardMarkup([]telegram_bot.InlineKeyboardButton{button})
	err = s.bot.Send(message)
	if err != nil {
		return errors.WithMessagef(err, "send notification to chat: %d", chatId)
	}
	return nil
}

func (s UserOrder) CancelOrder(ctx context.Context, req entity.QueryCallbackPayload) error {
	err := s.orderRepo.SetOrderStatus(ctx, req.OrderId, entity.OrderItemStatusPaid, entity.CancelOrderCommand)
	if err != nil {
		return errors.WithMessage(err, "update order status")
	}
	chatId, err := s.orderRepo.GetOrderedChatId(ctx, req.OrderId)
	if err != nil {
		return errors.WithMessage(err, "get user chat id")
	}
	err = s.bot.Send(telegram_bot.NewMessage(chatId, fmt.Sprintf("Заказ №%s отменён", req.OrderId)))
	if err != nil {
		return errors.WithMessagef(err, "send notification to chat: %d", chatId)
	}
	return nil
}

func (s UserOrder) getOrderInfoString(order *entity.Order, user *entity.User) string {
	template := `Заказ №%s
	Состав заказа: 
	%v
	Имя заказавшего: %s
	Telegram ник заказавшего: @%s
	Стоимость заказа: %d.%d руб
	Пожелания: '%s'
	Дата заказа: %s
	`
	items := make([]string, len(order.Items))
	for i, item := range order.Items {
		items[i] = fmt.Sprintf("[%d] dish_id=%d %s x %d", i+1, item.DishId, item.Name, item.Count)
	}
	return fmt.Sprintf(template,
		order.Id, strings.Join(items, "\n"),
		user.Name, user.Username, order.Total/100, order.Total%100, // nolint:mnd
		order.Wishes, order.CreatedAt.Local().Format(time.DateTime), // nolint:gosmopolitan
	)
}
