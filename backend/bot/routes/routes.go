package routes

import (
	"dish_as_a_service/bot/controller"

	"github.com/Falokut/go-kit/telegram_bot"
	"github.com/Falokut/go-kit/telegram_bot/router"
)

type Controllers struct {
	User  controller.User
	Order controller.Order
}

type Endpoint struct {
	Handler     router.HandlerFunc
	Command     string
	Description string
	UpdateType  string
	Admin       bool
	Hide        bool
}

func InitRoutes(c Controllers, middlewares []router.Middleware, adminAuthMiddleware router.Middleware) *router.Router {
	router := router.NewRouter(middlewares...)
	for _, command := range EndpointsDescriptors(c) {
		if command.UpdateType == telegram_bot.MessageUpdateType && command.Admin {
			router.Handler(command.UpdateType, command.Command, adminAuthMiddleware(command.Handler))
			continue
		}
		router.Handler(command.UpdateType, command.Command, command.Handler)
	}
	return router
}

func EndpointsDescriptors(c Controllers) []Endpoint {
	return []Endpoint{
		{
			Handler:     c.User.Register,
			Command:     "start",
			Description: "зарегистрироваться",
			UpdateType:  telegram_bot.MessageUpdateType,
		},
		{
			Handler:     c.User.List,
			Command:     "user_list",
			Description: "получить список пользователей",
			UpdateType:  telegram_bot.MessageUpdateType,
			Admin:       true,
		},
		{
			Command:     "add_admin",
			Handler:     c.User.AddAdmin,
			UpdateType:  telegram_bot.MessageUpdateType,
			Description: "добавить админа по username",
			Admin:       true,
		},
		{
			Command:     "remove_admin",
			Handler:     c.User.RemoveAdminByUsername,
			UpdateType:  telegram_bot.MessageUpdateType,
			Description: "удалить админа по username",
			Admin:       true,
		},
		{
			Command:    "pass_by_secret",
			Handler:    c.User.AddAdminSecret,
			UpdateType: telegram_bot.MessageUpdateType,
			Hide:       true,
		},
		{
			Command:     "secret",
			Description: "Получить значение secret для становления админом",
			Handler:     c.User.GetAdminSecret,
			UpdateType:  telegram_bot.MessageUpdateType,
			Admin:       true,
		},
		{
			Command:     "allow_ordering",
			Description: "Разрешить заказывать",
			Handler:     c.Order.AllowOrdering,
			UpdateType:  telegram_bot.MessageUpdateType,
			Admin:       true,
		},
		{
			Command:     "forbid_ordering",
			Description: "Запретить заказывать",
			Handler:     c.Order.ForbidOrdering,
			UpdateType:  telegram_bot.MessageUpdateType,
			Admin:       true,
		},
		{
			Handler:    c.Order.HandleCallbackQuery,
			UpdateType: telegram_bot.CallbackQueryUpdateType,
		},
		{
			Handler:    c.Order.HandlePreCheckout,
			UpdateType: telegram_bot.PreCheckoutQueryUpdateType,
		},
		{
			Handler:    c.Order.HandlePayment,
			UpdateType: telegram_bot.SuccessfulPaymentMessageUpdateType,
			Hide:       true,
		},
	}
}
