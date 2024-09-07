package controller

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

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
//	@Success	200	{object}	domain.GetUserIdByTelegramIdResponse
//	@Failure	500	{string}	string
//	@Router		/users/get_by_telegram_id/:telegram_id [GET]
func (c User) GetUserIdByTelegramId(ctx echo.Context) error {
	req, err := bindRequest[domain.GetUserIdByTelegramIdRequest](ctx)
	if err != nil {
		return err
	}

	userId, err := c.service.GetUserIdByTelegramId(ctx.Request().Context(), req.TelegramId)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotExist) {
			return ctx.String(http.StatusNotFound, domain.ErrUserNotExist.Error())
		}
		return err
	}

	return ctx.JSON(http.StatusOK, domain.GetUserIdByTelegramIdResponse{UserId: userId})
}

// Is Admin
//
//	@Tags		users
//	@Summary	Проверить, является ли пользователь админом
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	domain.IsUserAdminResponse
//	@Failure	500	{string}	string
//	@Router		/users/:user_id/is_admin [GET]
func (c User) IsAdmin(ctx echo.Context) error {
	req, err := bindRequest[domain.IsUserAdminRequest](ctx)
	if err != nil {
		return err
	}

	isAdmin, err := c.service.IsAdmin(ctx.Request().Context(), req.UserId)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK,
		domain.IsUserAdminResponse{
			IsAdmin: isAdmin,
		})
}
