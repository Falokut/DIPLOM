package controller

import "github.com/Falokut/go-kit/client/telegram_bot"

func HandleError(msg *telegram_bot.Message, err error, debug bool) telegram_bot.Chattable {
	if debug {
		return telegram_bot.NewMessage(msg.Chat.ID, "произошла ошибка текст ошибки: "+err.Error())
	}
	return nil
}
