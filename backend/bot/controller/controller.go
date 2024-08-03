package controller

import (
	tgbotapi "dish_as_a_service/bot/api"
)

func HandleError(msg *tgbotapi.Message, err error, debug bool) tgbotapi.Chattable {
	if debug {
		return tgbotapi.NewMessage(msg.Chat.ID, "произошла ошибка текст ошибки: "+err.Error())
	}
	return nil
}
