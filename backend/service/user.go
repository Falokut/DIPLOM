package service

import (
	"context"

	"dish_as_a_service/domain"
	"dish_as_a_service/entity"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type UserRepo interface {
	Register(ctx context.Context, user entity.RegisterUser) error
	IsAdmin(ctx context.Context, id string) (bool, error)
	GetUsers(ctx context.Context) ([]entity.User, error)
	GetUserChatId(ctx context.Context, userId string) (int64, error)
	GetUserIdByTelegramId(ctx context.Context, chatId int64) (string, error)
	AddAdmin(ctx context.Context, username string) error
	AddAdminChatId(ctx context.Context, chatId int64) error
}

type Secret interface {
	GetSecret() string
}

type User struct {
	repo   UserRepo
	secret Secret
}

func NewUser(repo UserRepo, secret Secret) User {
	return User{
		repo:   repo,
		secret: secret,
	}
}

func (s User) Register(ctx context.Context, user domain.RegisterUser) error {
	var telegram *entity.Telegram
	if user.Telegram != nil {
		telegram = &entity.Telegram{
			UserId: user.Telegram.UserId,
			ChatId: user.Telegram.ChatId,
		}
	}
	req := entity.RegisterUser{
		Id:       uuid.NewString(),
		Name:     user.Name,
		Username: user.Username,
		Telegram: telegram,
	}
	err := s.repo.Register(ctx, req)
	if err != nil {
		return errors.WithMessage(err, "register user")
	}
	return nil
}

func (s User) IsAdmin(ctx context.Context, userId string) (bool, error) {
	isAdmin, err := s.repo.IsAdmin(ctx, userId)
	if err != nil {
		return false, errors.WithMessage(err, "check is user admin")
	}
	return isAdmin, nil
}

func (s User) List(ctx context.Context) ([]domain.User, error) {
	users, err := s.repo.GetUsers(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "get users")
	}

	converted := make([]domain.User, len(users))
	for i, u := range users {
		converted[i] = domain.User{
			Username: u.Username,
			Name:     u.Name,
			Admin:    u.Admin,
		}
	}
	return converted, nil
}
func (s User) AddAdmin(ctx context.Context, username string) error {
	err := s.repo.AddAdmin(ctx, username)
	if err != nil {
		return errors.WithMessage(err, "add admin")
	}
	return nil
}

func (s User) GetChatId(ctx context.Context, userId string) (int64, error) {
	chatId, err := s.repo.GetUserChatId(ctx, userId)
	if err != nil {
		return -1, errors.WithMessage(err, "add admin")
	}
	return chatId, nil
}

func (s User) GetUserIdByTelegramId(ctx context.Context, telegramId int64) (string, error) {
	userId, err := s.repo.GetUserIdByTelegramId(ctx, telegramId)
	if err != nil {
		return "", errors.WithMessage(err, "get user id by telegram id")
	}
	return userId, nil
}

func (s User) AddAdminSecret(ctx context.Context, req domain.AddAdminSecretRequest) error {
	if s.secret.GetSecret() != req.Secret {
		return domain.ErrWrongSecret
	}

	err := s.repo.AddAdminChatId(ctx, req.ChatId)
	if err != nil {
		return errors.WithMessage(err, "add admin")
	}
	return nil
}

func (s User) GetAdminSecret(_ context.Context) (string, error) {
	return s.secret.GetSecret(), nil
}
