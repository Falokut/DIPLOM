package routes

import (
	"context"
	"dish_as_a_service/domain"

	"github.com/Falokut/go-kit/telegram_bot"
	"github.com/Falokut/go-kit/telegram_bot/apierrors"
	"github.com/Falokut/go-kit/telegram_bot/router"
	"github.com/pkg/errors"
)

type UserRepo interface {
	GetUserIdByTelegramId(ctx context.Context, telegramId int64) (string, error)
	IsAdmin(ctx context.Context, userId string) (bool, error)
}

type AdminAuth struct {
	userRepo UserRepo
}

func NewAdminAuth(userRepo UserRepo) AdminAuth {
	return AdminAuth{
		userRepo: userRepo,
	}
}

func (m AdminAuth) AdminAuth(next router.HandlerFunc) router.HandlerFunc {
	return func(ctx context.Context, update telegram_bot.Update) (telegram_bot.Chattable, error) {
		from := update.SentFrom()
		if from == nil {
			return nil, apierrors.NewBusinessError(domain.ErrCodeInvalidArgument, "invalid update type",
				errors.New("invalid update type"))
		}
		userId, err := m.userRepo.GetUserIdByTelegramId(ctx, from.Id)
		if err != nil {
			return nil, err
		}
		isAdmin, err := m.userRepo.IsAdmin(ctx, userId)
		if err != nil {
			return nil, err
		}
		if !isAdmin {
			return nil, apierrors.NewBusinessError(
				domain.ErrCodeUserNotAdmin,
				domain.ErrUserOperationForbidden.Error(),
				domain.ErrUserOperationForbidden,
			)
		}
		return next(ctx, update)
	}
}
