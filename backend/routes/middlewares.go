package routes

import (
	"context"
	"dish_as_a_service/domain"
	"fmt"
	"net/http"
	"time"

	"github.com/Falokut/go-kit/log"
	"github.com/labstack/echo/v4"
)

const (
	userIdHeader = "X-USER-ID"
)

type Logger struct {
	logger log.Logger
}

func NewLogger(logger log.Logger) Logger {
	return Logger{
		logger: logger,
	}
}

func (m Logger) LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		start := time.Now()
		err := next(c)
		stop := time.Now()
		fields := []log.Field{
			log.Any("start", start.Format(time.DateTime)),
			log.Any("stop", stop.Format(time.DateTime)),
			log.Any("duration_ns", stop.Sub(start).Nanoseconds()),
			log.Any("endpoint", req.URL.String()),
			log.Any("method", req.Method),
			log.Any("status", c.Response().Status),
			log.Any("status_text", http.StatusText(c.Response().Status)),
		}
		if err != nil {
			fields = append(fields, log.Any("error", err))
		}
		m.logger.Info(req.Context(), "handle http request", fields...)
		return err
	}
}

type UserAuthRepo interface {
	IsAdmin(ctx context.Context, id string) (bool, error)
}

type UserAuth struct {
	repo UserAuthRepo
}

func NewAuthMiddleware(repo UserAuthRepo) UserAuth {
	return UserAuth{repo: repo}
}

func (m UserAuth) UserAdminAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Request().Header.Get(userIdHeader)
		if userId == "" {
			return c.String(http.StatusUnauthorized, fmt.Sprintf("заголовок %s не предоставлен", userIdHeader))
		}
		isAdmin, err := m.repo.IsAdmin(c.Request().Context(), userId)
		if err != nil {
			return err
		}
		if !isAdmin {
			return c.String(http.StatusForbidden, domain.ErrUserOperationForbidden.Error())
		}
		return nil
	}
}

func HandleError(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			if c.Response() == nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
			return err
		}
		return nil
	}
}
