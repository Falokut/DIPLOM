package routes

import (
	"context"
	"dish_as_a_service/domain"
	"dish_as_a_service/helpers"
	"dish_as_a_service/jwt"
	"fmt"
	"net/http"
	"slices"

	http2 "github.com/Falokut/go-kit/http"
	"github.com/Falokut/go-kit/http/apierrors"
)

type AuthMiddleware struct {
	accessTokenSecret string
}

func NewAuthMiddleware(accessTokenSecret string) AuthMiddleware {
	return AuthMiddleware{
		accessTokenSecret: accessTokenSecret,
	}
}

func (m AuthMiddleware) AdminAuthToken() http2.Middleware {
	return AuthToken(m.accessTokenSecret, domain.AdminType)
}

func (m AuthMiddleware) UserAuthToken() http2.Middleware {
	return AuthToken(m.accessTokenSecret, domain.UserType, domain.AdminType)
}

func AuthToken(tokenSecret string, roles ...string) http2.Middleware {
	return func(next http2.HandlerFunc) http2.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			token, err := helpers.GetBearerToken(r)
			if err != nil {
				return err
			}
			userInfo, err := jwt.ParseToken(token, tokenSecret)
			if err != nil {
				return err
			}
			r.Header.Add(domain.UserIdHeader, fmt.Sprint(userInfo.UserId))
			if len(roles) == 0 {
				return next(ctx, w, r)
			}
			if !slices.Contains(roles, userInfo.RoleName) {
				return apierrors.New(http.StatusForbidden,
					domain.ErrCodeForbidden,
					"доступ запрещён",
					domain.ErrForbidden, // nolint:err113
				)
			}
			return next(ctx, w, r)
		}
	}
}
