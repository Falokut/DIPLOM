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

	for _, desc := range endpointDescriptors(r) {
		if desc.IsAdmin {
			e.Add(desc.Method, desc.Path, authMiddleware.UserAdminAuth(desc.Handler))
		} else {
			e.Add(desc.Method, desc.Path, desc.Handler)
		}
	}

	return e.Server.Handler
}

type EndpointDescriptor struct {
	Method  string
	Path    string
	IsAdmin bool
	Handler echo.HandlerFunc
}

func endpointDescriptors(r Router) []EndpointDescriptor {
	return []EndpointDescriptor{
		{
			Method:  http.MethodGet,
			Path:    "/dishes",
			Handler: r.Dish.List,
		},
		{
			Method:  http.MethodPost,
			Path:    "/dishes",
			IsAdmin: true,
			Handler: r.Dish.AddDish,
		},
		{
			Method:  http.MethodGet,
			Path:    "/dishes/categories",
			Handler: r.DishesCategories.GetCategories,
		},
		{
			Method:  http.MethodGet,
			Path:    "/dishes/categories/:id",
			Handler: r.DishesCategories.GetCategory,
		},
		{
			Method:  http.MethodPost,
			Path:    "/dishes/categories",
			IsAdmin: true,
			Handler: r.DishesCategories.AddCategory,
		},
		{
			Method:  http.MethodPost,
			Path:    "/dishes/categories/:id",
			IsAdmin: true,
			Handler: r.DishesCategories.RenameCategory,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/dishes/categories/:id",
			IsAdmin: true,
			Handler: r.DishesCategories.RenameCategory,
		},
		{
			Method:  http.MethodPost,
			Path:    "/orders",
			IsAdmin: true,
			Handler: r.Order.ProcessOrder,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/get_by_telegram_id/:telegram_id",
			Handler: r.User.GetUserIdByTelegramId,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/:user_id/is_admin",
			Handler: r.User.IsAdmin,
		},
	}
}
