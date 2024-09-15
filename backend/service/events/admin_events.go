package events

import (
	"context"
	"dish_as_a_service/bot/routes"
	"github.com/Falokut/go-kit/telegram_bot"
	"github.com/pkg/errors"
)

type AdminEvents struct {
	tgBot *telegram_bot.BotAPI
}

func NewAdminEvents(tgBot *telegram_bot.BotAPI) AdminEvents {
	return AdminEvents{
		tgBot: tgBot,
	}
}

func (e AdminEvents) AdminAdded(ctx context.Context, chatId int64) error {
	scope := telegram_bot.NewBotCommandScopeChat(chatId)
	endpoints := routes.EndpointsDescriptors(routes.Controllers{})
	commands := make([]telegram_bot.BotCommand, 0, len(endpoints))
	for _, endpoint := range endpoints {
		if endpoint.Hide || endpoint.UpdateType != telegram_bot.MessageUpdateType {
			continue
		}
		commands = append(commands, telegram_bot.BotCommand{
			Command:     endpoint.Command,
			Description: endpoint.Description,
		})
	}
	err := e.tgBot.Send(telegram_bot.NewSetMyCommandsWithScope(scope, commands...))
	if err != nil {
		return errors.WithMessage(err, "send admin commands")
	}
	return nil
}

func (e AdminEvents) AdminRemoved(ctx context.Context, chatId int64) error {
	scope := telegram_bot.NewBotCommandScopeChat(chatId)
	err := e.tgBot.Send(telegram_bot.NewDeleteMyCommandsWithScope(scope))
	if err != nil {
		return errors.WithMessage(err, "send delete chat scope commands")
	}
	return nil
}
