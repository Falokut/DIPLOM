package routes

import (
	"dish_as_a_service/controller"
	"net/http"

	"github.com/Falokut/go-kit/http/endpoint"
	"github.com/Falokut/go-kit/http/router"
)

type Router struct {
	Dish             controller.Dish
	DishesCategories controller.DishesCategories
	Order            controller.Order
	User             controller.User
}

func (r Router) InitRoutes(authMiddleware UserAuth, wrapper endpoint.Wrapper) *router.Router {
	mux := router.New()
	for _, desc := range endpointDescriptors(r) {
		if desc.IsAdmin {
			withAuthWrapper := wrapper.WithMiddlewares(authMiddleware.UserAdminAuth)
			mux.Handler(desc.Method, desc.Path, withAuthWrapper.Endpoint(desc.Handler))
		} else {
			mux.Handler(desc.Method, desc.Path, wrapper.Endpoint(desc.Handler))
		}
	}

	return mux
}

type EndpointDescriptor struct {
	Method  string
	Path    string
	IsAdmin bool
	Handler any
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
			Handler: r.Order.ProcessOrder,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/get_by_telegram_id/:telegram_id",
			Handler: r.User.GetUserIdByTelegramId,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/is_admin",
			Handler: r.User.IsAdmin,
		},
	}
}
