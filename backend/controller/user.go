package controller

import (
	"context"
	"errors"
	"net/http"
	"strconv"

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
	telegramId, err := strconv.ParseInt(ctx.Param("telegram_id"), 10, 64)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	userId, err := c.service.GetUserIdByTelegramId(ctx.Request().Context(), telegramId)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotExist) {
			return ctx.String(http.StatusNotFound, err.Error())
		}
		return err
	}

	return ctx.JSON(http.StatusOK, domain.GetUserIdByTelegramIdResponse{UserId: userId})
}

// Get user by chat id
//
//	@Tags		users
//	@Summary	Получить id пользователя по id чата
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	domain.IsUserAdminResponse
//	@Failure	500	{string}	string
//	@Router		/users/:user_id/is_admin [GET]
func (c User) IsAdmin(ctx echo.Context) error {
	userId := ctx.Param("user_id")
	if userId == "" {
		return ctx.NoContent(http.StatusBadRequest)
	}

	isAdmin, err := c.service.IsAdmin(ctx.Request().Context(), userId)
	if err != nil {
		return err
	}
	resp := domain.IsUserAdminResponse{
		IsAdmin: isAdmin,
	}
	return ctx.JSON(http.StatusOK, resp)
}
