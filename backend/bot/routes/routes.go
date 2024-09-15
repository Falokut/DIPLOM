package routes

import (
	"context"
	"dish_as_a_service/bot/controller"
	"dish_as_a_service/entity"

	"github.com/Falokut/go-kit/telegram_bot"
	"github.com/Falokut/go-kit/telegram_bot/router"
	"github.com/pkg/errors"
)

type Controllers struct {
	User  controller.User
	Order controller.Order
}

type RegisterDestination interface {
	UpsertCommands(commands []telegram_bot.SetMyCommandsConfig) error
}
type AdminUsersProvides interface {
	GetAdminsTelegramInfo(ctx context.Context) ([]entity.Telegram, error)
}

func RegisterRoutes(ctx context.Context, dest RegisterDestination, adminUsersProvider AdminUsersProvides) error {
	adminTelegrams, err := adminUsersProvider.GetAdminsTelegramInfo(ctx)
	if err != nil {
		return errors.WithMessage(err, "get admins telegram info")
	}
	adminScopes := make([]*telegram_bot.BotCommandScope, 0, len(adminTelegrams))
	for _, telegramInfo := range adminTelegrams {
		adminScopes = append(adminScopes, &telegram_bot.BotCommandScope{
			Type:   "chat",
			ChatID: telegramInfo.ChatId,
			UserID: telegramInfo.UserId,
		})
	}
	endpoints := endpointsDescriptors(Controllers{})
	publicCommands := make([]telegram_bot.BotCommand, 0, len(endpoints))
	for _, endpoint := range endpoints {
		if endpoint.Hide || endpoint.Admin || endpoint.UpdateType != telegram_bot.MessageUpdateType {
			continue
		}
		publicCommands = append(publicCommands, telegram_bot.BotCommand{
			Command:     endpoint.Command,
			Description: endpoint.Description,
		})
	}
	commandsConfig := make([]telegram_bot.SetMyCommandsConfig, 0, len(adminScopes)+1)
	if len(publicCommands) > 0 {
		commandsConfig = append(commandsConfig, telegram_bot.SetMyCommandsConfig{Commands: publicCommands})
	}

	allCommands := make([]telegram_bot.BotCommand, 0, len(endpoints))
	for _, endpoint := range endpoints {
		if endpoint.Hide || endpoint.UpdateType != telegram_bot.MessageUpdateType {
			continue
		}
		allCommands = append(allCommands, telegram_bot.BotCommand{
			Command:     endpoint.Command,
			Description: endpoint.Description,
		})
	}
	for _, scope := range adminScopes {
		commandsConfig = append(commandsConfig, telegram_bot.SetMyCommandsConfig{
			Commands: allCommands,
			Scope:    scope,
		})
	}

	err = dest.UpsertCommands(commandsConfig)
	if err != nil {
		return errors.WithMessage(err, "set commands")
	}
	return nil
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
	for _, command := range endpointsDescriptors(c) {
		if command.UpdateType == telegram_bot.MessageUpdateType && command.Admin {
			router.Handler(command.UpdateType, command.Command, adminAuthMiddleware(command.Handler))
			continue
		}
		router.Handler(command.UpdateType, command.Command, command.Handler)
	}
	return router
}

func endpointsDescriptors(c Controllers) []Endpoint {
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
