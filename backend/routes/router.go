// nolint:funlen
package routes

import (
	"dish_as_a_service/controller"
	"net/http"

	"github.com/Falokut/go-kit/http/endpoint"
	"github.com/Falokut/go-kit/http/router"
)

type Router struct {
	Auth             controller.Auth
	Dish             controller.Dish
	DishesCategories controller.DishesCategories
	Order            controller.Order
	Restaurant       controller.Restaurant
}

func (r Router) InitRoutes(authMiddleware AuthMiddleware, wrapper endpoint.Wrapper) *router.Router {
	mux := router.New()
	for _, desc := range endpointDescriptors(r) {
		endpointWrapper := wrapper
		switch {
		case desc.IsAdmin:
			endpointWrapper = wrapper.WithMiddlewares(authMiddleware.AdminAuthToken())
		case desc.NeedUserAuth:
			endpointWrapper = wrapper.WithMiddlewares(authMiddleware.UserAuthToken())
		}
		mux.Handler(desc.Method, desc.Path, endpointWrapper.Endpoint(desc.Handler))
	}

	return mux
}

type EndpointDescriptor struct {
	Method       string
	Path         string
	IsAdmin      bool
	NeedUserAuth bool
	Handler      any
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
			Method:  http.MethodPost,
			Path:    "/dishes/edit/:id",
			IsAdmin: true,
			Handler: r.Dish.EditDish,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/dishes/delete/:id",
			IsAdmin: true,
			Handler: r.Dish.DeleteDish,
		},
		{
			Method:  http.MethodGet,
			Path:    "/dishes/all_categories",
			Handler: r.DishesCategories.GetAllCategories,
		},
		{
			Method:  http.MethodGet,
			Path:    "/dishes/categories",
			Handler: r.DishesCategories.GetDishesCategories,
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
			Handler: r.DishesCategories.DeleteCategory,
		},
		{
			Method:  http.MethodGet,
			Path:    "/restaurants",
			Handler: r.Restaurant.GetAllRestaurants,
		},
		{
			Method:  http.MethodGet,
			Path:    "/restaurants/:id",
			Handler: r.Restaurant.GetRestaurant,
		},
		{
			Method:  http.MethodPost,
			Path:    "/restaurants",
			IsAdmin: true,
			Handler: r.Restaurant.AddRestaurant,
		},
		{
			Method:  http.MethodPost,
			Path:    "/restaurants/:id",
			IsAdmin: true,
			Handler: r.Restaurant.RenameDishesRestaurant,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/restaurants/:id",
			IsAdmin: true,
			Handler: r.Restaurant.DeleteRestaurant,
		},

		{
			Method:       http.MethodPost,
			Path:         "/orders",
			Handler:      r.Order.ProcessOrder,
			NeedUserAuth: true,
		},
		{
			Method:       http.MethodGet,
			Path:         "/orders/my",
			Handler:      r.Order.GetUserOrders,
			NeedUserAuth: true,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/login_by_telegram",
			Handler: r.Auth.LoginByTelegram,
		},
		{
			Method:  http.MethodGet,
			Path:    "/auth/refresh_access_token",
			Handler: r.Auth.RefreshAccessToken,
		},
		{
			Method:  http.MethodGet,
			Path:    "/has_admin_privileges",
			Handler: r.Auth.HasAdminPrivileges,
		},
	}
}
