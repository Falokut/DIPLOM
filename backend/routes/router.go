package routes

import (
	"dish_as_a_service/controller"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Router struct {
	Dish  controller.Dish
	Order controller.Order
	User  controller.User
}

//nolint:gochecknoglobals
var DefaultValidator = validator.New()

func (r Router) InitRoutes(authMiddleware UserAuth, customMiddlewares ...echo.MiddlewareFunc) http.Handler {
	e := echo.New()
	e.Use(customMiddlewares...)
	e.JSONSerializer = JsonSerializer{}
	e.Validator = Validate{v: validator.New()}

	e.GET("/dishes", r.Dish.List)
	e.POST("/dishes", authMiddleware.UserAdminAuth(r.Dish.AddDish))

	e.POST("/orders", r.Order.ProcessOrder)

	e.GET("/users/get_by_telegram_id/:telegram_id", r.User.GetUserIdByTelegramId)
	e.GET("/users/:user_id/is_admin", r.User.IsAdmin)

	return e.Server.Handler
}
