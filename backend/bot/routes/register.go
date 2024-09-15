package routes

import (
	"context"
	"dish_as_a_service/entity"

	"github.com/Falokut/go-kit/telegram_bot"
	"github.com/pkg/errors"
)

type AdminUsersProvides interface {
	GetTelegramUsersInfo(ctx context.Context) ([]entity.TelegramUser, error)
}

func RegisterRoutes(ctx context.Context, tgBot *telegram_bot.BotAPI, usersProvider AdminUsersProvides) error {
	users, err := usersProvider.GetTelegramUsersInfo(ctx)
	if err != nil {
		return errors.WithMessage(err, "get telegram users info")
	}
	err = clearCommands(tgBot, users)
	if err != nil {
		return errors.WithMessage(err, "clear commands")
	}
	err = registerDefaultCommands(tgBot)
	if err != nil {
		return errors.WithMessage(err, "register default commands")
	}
	err = registerAdminCommands(tgBot, users)
	if err != nil {
		return errors.WithMessage(err, "register admin commands")
	}
	return nil
}

func clearCommandForChat(tgBot *telegram_bot.BotAPI, chatId int64) error {
	scope := telegram_bot.NewBotCommandScopeChat(chatId)
	err := tgBot.Send(telegram_bot.NewDeleteMyCommandsWithScope(scope))
	if err != nil {
		return errors.WithMessage(err, "send delete chat scope commands")
	}
	return nil
}

func clearCommands(tgBot *telegram_bot.BotAPI, users []entity.TelegramUser) error {
	err := tgBot.Send(telegram_bot.NewDeleteMyCommands())
	if err != nil {
		return errors.WithMessage(err, "send delete default scope commands")
	}
	for _, user := range users {
		err = clearCommandForChat(tgBot, user.ChatId)
		if err != nil {
			return errors.WithMessagef(err, "clear commands for chatId=%d", user.ChatId)
		}
	}
	return nil
}

func registerDefaultCommands(tgBot *telegram_bot.BotAPI) error {
	endpoints := EndpointsDescriptors(Controllers{})
	commands := make([]telegram_bot.BotCommand, 0, len(endpoints))
	for _, endpoint := range endpoints {
		if endpoint.Hide || endpoint.Admin || endpoint.UpdateType != telegram_bot.MessageUpdateType {
			continue
		}
		commands = append(commands, telegram_bot.BotCommand{
			Command:     endpoint.Command,
			Description: endpoint.Description,
		})
	}
	err := tgBot.Send(telegram_bot.NewSetMyCommands(commands...))
	if err != nil {
		return errors.WithMessage(err, "send bot commands")
	}
	return nil
}

func registerAdminCommands(tgBot *telegram_bot.BotAPI, users []entity.TelegramUser) error {
	endpoints := EndpointsDescriptors(Controllers{})
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
	scopes := make([]telegram_bot.BotCommandScope, 0)
	for _, user := range users {
		if !user.Admin {
			continue
		}
		scopes = append(scopes, telegram_bot.NewBotCommandScopeChat(user.ChatId))
	}
	for _, scope := range scopes {
		err := tgBot.Send(telegram_bot.NewSetMyCommandsWithScope(scope, commands...))
		if err != nil {
			return errors.WithMessagef(err, "send bot commands for chatId=%d", scope.ChatID)
		}
	}
	return nil
}
