package routes

import (
	"context"
	"dish_as_a_service/domain"
	"fmt"
	"net/http"

	http2 "github.com/Falokut/go-kit/http"
	"github.com/Falokut/go-kit/http/apierrors"
)

const (
	userIdHeader = "X-User-Id"
)

type UserAuthRepo interface {
	IsAdmin(ctx context.Context, id string) (bool, error)
}

type UserAuth struct {
	repo UserAuthRepo
}

func NewAuthMiddleware(repo UserAuthRepo) UserAuth {
	return UserAuth{repo: repo}
}

func (m UserAuth) UserAdminAuth(next http2.HandlerFunc) http2.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		userId := r.Header.Get(userIdHeader)
		if userId == "" {
			return apierrors.New(http.StatusUnauthorized, domain.ErrCodeEmptyUserIdHeader,
				fmt.Sprintf("заголовок %s не предоставлен", userIdHeader), domain.ErrUnauthorized)
		}
		isAdmin, err := m.repo.IsAdmin(r.Context(), userId)
		if err != nil {
			return err
		}
		if !isAdmin {
			return apierrors.New(http.StatusForbidden, domain.ErrCodeUserNotAdmin,
				domain.ErrUserOperationForbidden.Error(), domain.ErrUserOperationForbidden)
		}
		return next(ctx, w, r)
	}
}
