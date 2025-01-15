package controller

import (
	"context"
	"dish_as_a_service/domain"
	"dish_as_a_service/helpers"
	"net/http"

	"github.com/Falokut/go-kit/http/apierrors"
	"github.com/pkg/errors"
)

type AuthService interface {
	LoginByTelegram(ctx context.Context, req domain.LoginByTelegramRequest) (*domain.LoginResponse, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (*domain.TokenResponse, error)
	HasAdminPrivileges(ctx context.Context, accessToken string) (*domain.HasAdminPrivilegesResponse, error)
}

type Auth struct {
	service AuthService
}

func NewAuth(service AuthService) Auth {
	return Auth{
		service: service,
	}
}

// LoginByTelegram
//
//	@Tags		auth
//	@Summary	Войти в аккаунт
//	@Accept		json
//	@Produce	json
//
//	@Param		body	body		domain.LoginByTelegramRequest	true	"тело запроса"
//
//	@Success	200		{object}	domain.LoginResponse
//	@Failure	404		{object}	apierrors.Error
//	@Failure	500		{object}	apierrors.Error
//	@Router		/auth/login_by_telegram [POST]
func (c Auth) LoginByTelegram(ctx context.Context, req domain.LoginByTelegramRequest) (*domain.LoginResponse, error) {
	tokens, err := c.service.LoginByTelegram(ctx, req)
	switch {
	case errors.Is(err, domain.ErrUserNotFound), errors.Is(err, domain.ErrInvalidTelegramCredentials):
		return nil, apierrors.New(
			http.StatusNotFound,
			domain.ErrCodeForbidden,
			domain.ErrForbidden.Error(),
			err,
		)
	default:
		return tokens, err
	}
}

// RefreshAccessToken
//
//	@Tags		auth
//	@Summary	Обновить токен доступа
//	@Accept		json
//	@Produce	json
//
//
//	@Success	200	{object}	domain.TokenResponse
//	@Failure	404	{object}	apierrors.Error
//	@Failure	401	{object}	apierrors.Error
//	@Failure	500	{object}	apierrors.Error
//
//	@Security	Bearer
//
//	@Router		/auth/access_token [GET]
func (c Auth) RefreshAccessToken(ctx context.Context, r *http.Request) (*domain.TokenResponse, error) {
	token, err := helpers.GetBearerToken(r)
	if err != nil {
		return nil, err
	}
	return c.service.RefreshAccessToken(ctx, token)
}

// RefreshAccessToken
//
//	@Tags		auth
//	@Summary	Обновить токен доступа
//	@Accept		json
//	@Produce	json
//
//
//	@Success	200	{object}	domain.HasAdminPrivilegesResponse
//	@Failure	404	{object}	apierrors.Error
//	@Failure	401	{object}	apierrors.Error
//	@Failure	500	{object}	apierrors.Error
//
//	@Security	Bearer
//
//	@Router		/has_admin_privileges [GET]
func (c Auth) HasAdminPrivileges(ctx context.Context, r *http.Request) (*domain.HasAdminPrivilegesResponse, error) {
	token, err := helpers.GetBearerToken(r)
	if err != nil {
		return nil, err
	}
	return c.service.HasAdminPrivileges(ctx, token)
}
