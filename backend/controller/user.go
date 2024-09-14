package controller

import (
	"context"
	"errors"
	"net/http"

	"github.com/Falokut/go-kit/http/apierrors"

	"dish_as_a_service/domain"
)

type UserService interface {
	GetUserIdByTelegramId(ctx context.Context, telegramId int64) (string, error)
	IsAdmin(ctx context.Context, userId string) (bool, error)
}

type User struct {
	service UserService
}

func NewUser(service UserService) User {
	return User{
		service: service,
	}
}

// Get user by telegram id
//
//	@Tags		users
//	@Summary	Получить id пользователя по telegram id
//	@Accept		json
//	@Produce	json
//
//	@Param		telegramId	path		int	true	"Идентификатор категории"
//
//	@Success	200			{object}	domain.GetUserIdByTelegramIdResponse
//	@Failure	404			{object}	apierrors.Error
//	@Failure	500			{object}	apierrors.Error
//	@Router		/users/get_by_telegram_id/:telegramId [GET]
func (c User) GetUserIdByTelegramId(
	ctx context.Context,
	req domain.GetUserIdByTelegramIdRequest,
) (*domain.GetUserIdByTelegramIdResponse, error) {
	userId, err := c.service.GetUserIdByTelegramId(ctx, req.TelegramId)
	switch {
	case errors.Is(err, domain.ErrUserNotExist):
		return nil, apierrors.New(
			http.StatusNotFound,
			domain.ErrCodeUserNotFound,
			domain.ErrUserNotExist.Error(),
			err,
		)
	case err != nil:
		return nil, err
	default:
		return &domain.GetUserIdByTelegramIdResponse{UserId: userId}, nil
	}
}

// Is Admin
//
//	@Tags		users
//	@Summary	Проверить, является ли пользователь админом
//	@Accept		json
//	@Produce	json
//	@Param		body	body		domain.IsUserAdminRequest	true	"request body"
//	@Success	200		{object}	domain.IsUserAdminResponse
//	@Failure	404		{object}	apierrors.Error
//	@Failure	500		{string}	string
//	@Router		/users/is_admin [GET]
func (c User) IsAdmin(ctx context.Context, req domain.IsUserAdminRequest) (*domain.IsUserAdminResponse, error) {
	isAdmin, err := c.service.IsAdmin(ctx, req.UserId)
	switch {
	case errors.Is(err, domain.ErrUserNotExist):
		return nil, apierrors.New(
			http.StatusNotFound,
			domain.ErrCodeUserNotFound,
			domain.ErrUserNotExist.Error(),
			err,
		)
	case err != nil:
		return nil, err
	default:
		return &domain.IsUserAdminResponse{IsAdmin: isAdmin}, nil
	}
}
