package routes

import (
	"context"
	"fmt"

	tgbotapi "dish_as_a_service/bot/api"
	"dish_as_a_service/bot/controller"
	"github.com/Falokut/go-kit/log"
)

type Router struct {
	logger log.Logger
	user   controller.User
	order  controller.Order
}

type handlerFunction func(ctx context.Context, msg *tgbotapi.Message) tgbotapi.Chattable

type Endpoint struct {
	Handler handlerFunction
	Admin   bool
}

func NewRouter(logger log.Logger, user controller.User, order controller.Order) Router {
	return Router{
		logger: logger,
		user:   user,
		order:  order,
	}
}

func (r Router) HandleMessage(ctx context.Context, msg *tgbotapi.Message) tgbotapi.Chattable {
	endpoints := r.MessageEndpoints()

	endpoint, ok := endpoints[msg.Command()]
	if !ok {
		return tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("команда %s не найдена", msg.Command()))
	}

	if endpoint.Admin {
		isAdmin, resp := r.user.IsAdmin(ctx, msg)
		if !isAdmin {
			return resp
		}
	}

	return endpoint.Handler(ctx, msg)
}

func (r Router) HandlePayment(ctx context.Context, msg *tgbotapi.Message) tgbotapi.Chattable {
	return r.order.HandlePayment(ctx, msg)
}

func (r Router) HandlePreCheckout(ctx context.Context, msg *tgbotapi.PreCheckoutQuery) tgbotapi.Chattable {
	return r.order.HandlePreCheckout(ctx, msg)
}

func (r Router) MessageEndpoints() map[string]Endpoint {
	return map[string]Endpoint{
		"start": {
			Handler: r.user.Register,
		},
		"user_list": {
			Handler: r.user.List,
			Admin:   true,
		},
		"add_admin": {
			Handler: r.user.AddAdmin,
			Admin:   true,
		},
		"add_admin_secret": {
			Handler: r.user.AddAdminSecret,
		},
		"secret": {
			Handler: r.user.GetAdminSecret,
			Admin:   true,
		},
	}
}
