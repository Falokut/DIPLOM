package controller

import (
	"context"
	"errors"
	"fmt"
	"strings"

	tgbotapi "dish_as_a_service/bot/api"
	"dish_as_a_service/domain"
)

type UserService interface {
	Register(ctx context.Context, user domain.RegisterUser) error
	IsAdmin(ctx context.Context, userId string) (bool, error)
	GetUserIdByTelegramId(ctx context.Context, telegramId int64) (string, error)
	List(ctx context.Context) ([]domain.User, error)
	AddAdmin(ctx context.Context, username string) error
	AddAdminSecret(ctx context.Context, req domain.AddAdminSecretRequest) error
	GetAdminSecret(ctx context.Context) (string, error)
}

type User struct {
	service UserService
	debug   bool
}

func NewUser(service UserService, debug bool) User {
	return User{
		service: service,
		debug:   debug,
	}
}

func (c User) Register(ctx context.Context, msg *tgbotapi.Message) tgbotapi.Chattable {
	user := domain.RegisterUser{
		Username: msg.From.UserName,
		Name:     msg.From.FirstName + " " + msg.From.LastName,
		Telegram: &domain.Telegram{
			ChatId: msg.Chat.ID,
			UserId: msg.From.ID,
		},
	}
	err := c.service.Register(ctx, user)
	switch {
	case errors.Is(err, domain.ErrUserAlreadyExists):
		return tgbotapi.NewMessage(msg.Chat.ID, "вы уже зарегистрированы")
	case err != nil:
		return HandleError(msg, err, c.debug)
	}
	return tgbotapi.NewMessage(msg.Chat.ID, "вы зарегистрированы")
}

func (c User) IsAdmin(ctx context.Context, msg *tgbotapi.Message) (bool, tgbotapi.Chattable) {
	userId, err := c.service.GetUserIdByTelegramId(ctx, msg.From.ID)
	if err != nil {
		return false, HandleError(msg, err, c.debug)
	}
	isAdmin, err := c.service.IsAdmin(ctx, userId)
	if err != nil {
		return false, HandleError(msg, err, c.debug)
	}
	return isAdmin, tgbotapi.NewMessage(msg.Chat.ID, "эта команда предназначена только для администраторов")
}

const userUnderline = "___________________________________________"

//nolint:mnd
func (c User) List(ctx context.Context, msg *tgbotapi.Message) tgbotapi.Chattable {
	users, err := c.service.List(ctx)
	if err != nil {
		return HandleError(msg, err, c.debug)
	}
	text := make([]string, 0, len(users)*2+3)
	text = append(text, userUnderline, "|  #  |  [NAME]  |  [USERNAME]  |  [ADMIN]  |")
	for i, user := range users {
		text = append(text, userUnderline, fmt.Sprintf("|  %d  |%s|%s|%t|",
			i+1, user.Name, user.Username, user.Admin))
	}
	return tgbotapi.NewMessage(msg.Chat.ID, strings.Join(text, "\n"))
}

func (c User) AddAdmin(ctx context.Context, msg *tgbotapi.Message) tgbotapi.Chattable {
	err := c.service.AddAdmin(ctx, msg.CommandArguments())
	switch {
	case err != nil:
		return HandleError(msg, err, c.debug)
	case errors.Is(err, domain.ErrUserNotExist):
		return tgbotapi.NewMessage(msg.Chat.ID, "пользователь не найден")
	}
	return tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("администратор с username=%s добавлен", msg.CommandArguments()))
}

func (c User) AddAdminSecret(ctx context.Context, msg *tgbotapi.Message) tgbotapi.Chattable {
	req := domain.AddAdminSecretRequest{
		ChatId: msg.Chat.ID,
		Secret: msg.CommandArguments(),
	}
	err := c.service.AddAdminSecret(ctx, req)
	switch {
	case err != nil:
		return HandleError(msg, err, c.debug)
	case errors.Is(err, domain.ErrWrongSecret):
		return tgbotapi.NewMessage(msg.Chat.ID, "вы ввели неправильный пароль")
	case errors.Is(err, domain.ErrUserNotExist):
		return tgbotapi.NewMessage(msg.Chat.ID, "пользователь не найден")
	}
	return tgbotapi.NewMessage(msg.Chat.ID, "вы стали администратором")
}

func (c User) GetAdminSecret(ctx context.Context, msg *tgbotapi.Message) tgbotapi.Chattable {
	secret, err := c.service.GetAdminSecret(ctx)
	if err != nil {
		return HandleError(msg, err, c.debug)
	}
	return tgbotapi.NewMessage(msg.Chat.ID, "пароль для администратора: "+secret)
}
