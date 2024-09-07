package routes

import (
	"dish_as_a_service/controller"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Router struct {
	Dish             controller.Dish
	DishesCategories controller.DishesCategories
	Order            controller.Order
	User             controller.User
}

//nolint:gochecknoglobals
var DefaultValidator = validator.New()

func (r Router) InitRoutes(authMiddleware UserAuth, customMiddlewares ...echo.MiddlewareFunc) http.Handler {
	e := echo.New()
	e.Use(customMiddlewares...)
	e.JSONSerializer = JsonSerializer{}
	e.Validator = Validate{v: validator.New()}

	dishes := e.Group("/dishes")
	{
		dishes.GET("/", r.Dish.List)
		dishes.POST("/", authMiddleware.UserAdminAuth(r.Dish.AddDish))

		categories := dishes.Group("/categories")
		{
			categories.GET("/", r.DishesCategories.GetCategories)
			categories.GET("/:id", r.DishesCategories.GetCategory)

			categories.POST("/", authMiddleware.UserAdminAuth(r.DishesCategories.AddCategory))
			categories.POST("/:id", authMiddleware.UserAdminAuth(r.DishesCategories.RenameCategory))

			categories.DELETE("/:id", authMiddleware.UserAdminAuth(r.DishesCategories.DeleteCategory))
		}
	}

	e.POST("/orders", r.Order.ProcessOrder)

	users := e.Group("/users")
	{
		users.GET("/get_by_telegram_id/:telegram_id", r.User.GetUserIdByTelegramId)
		users.GET("/:user_id/is_admin", r.User.IsAdmin)
	}

	return e.Server.Handler
}
