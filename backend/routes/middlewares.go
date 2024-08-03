package routes

import (
	"net/http"
	"time"

	"github.com/Falokut/go-kit/log"
	"github.com/labstack/echo/v4"
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
