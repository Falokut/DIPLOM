package bot

import (
	"context"
	"encoding/json"
	"fmt"

	tgbotapi "dish_as_a_service/bot/api"
	"dish_as_a_service/entity"
	"github.com/pkg/errors"
)

type PaymentBot struct {
	bot          *tgbotapi.BotAPI
	invoiceToken string
}

func NewPaymentBot(token string, bot *tgbotapi.BotAPI) PaymentBot {
	return PaymentBot{
		invoiceToken: token,
		bot:          bot,
	}
}

const rubCurrency = "RUB"

func (b PaymentBot) ProcessPayment(_ context.Context, order *entity.Order, chatId int64) error {
	payload := entity.PaymentPayload{
		ChatId:  chatId,
		OrderId: order.Id,
	}

	args, err := json.Marshal(payload)
	if err != nil {
		return errors.WithMessage(err, "marhal payload")
	}

	prices := make([]tgbotapi.LabeledPrice, len(order.Items))
	for i := range order.Items {
		label := fmt.Sprintf("%s x %d", order.Items[i].Name, order.Items[i].Count)
		prices[i] = tgbotapi.LabeledPrice{
			Label:  label,
			Amount: order.Items[i].Price,
		}
	}

	invoice := tgbotapi.NewInvoice(
		chatId,
		fmt.Sprintf("Заказ № %s", order.Id),
		"оплата заказа",
		string(args),
		b.invoiceToken,
		"invoice",
		rubCurrency,
		prices,
	)
	err = b.bot.Send(invoice)
	if err != nil {
		return errors.WithMessage(err, "send invoice")
	}

	return nil
}
