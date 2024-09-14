package service

import (
	"context"
	"dish_as_a_service/entity"
	"fmt"
	"strings"
	"time"

	"github.com/Falokut/go-kit/client/telegram_bot"
	"github.com/pkg/errors"
)

type UserRepo interface {
	GetUserInfo(ctx context.Context, userId string) (entity.User, error)
	GetAdminsChatsIds(ctx context.Context) ([]int64, error)
}

type UserOrder struct {
	bot  *telegram_bot.BotAPI
	repo UserRepo
}

func NewOrderUserService(bot *telegram_bot.BotAPI, repo UserRepo) UserOrder {
	return UserOrder{
		bot:  bot,
		repo: repo,
	}
}

func (s UserOrder) NotifySuccessPayment(ctx context.Context, order *entity.Order) error {
	adminIds, err := s.repo.GetAdminsChatsIds(ctx)
	if err != nil {
		return errors.WithMessage(err, "get admins chats ids")
	}
	user, err := s.repo.GetUserInfo(ctx, order.UserId)
	if err != nil {
		return errors.WithMessage(err, "get user info")
	}

	orderInfoString := s.getOrderInfoString(order, &user)
	for _, chatId := range adminIds {
		err = s.bot.Send(telegram_bot.NewMessage(chatId, orderInfoString))
		if err != nil {
			return errors.WithMessagef(err, "send notification to chat: %d", chatId)
		}
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
