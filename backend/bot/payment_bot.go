package bot

import (
	"context"
	"encoding/json"
	"fmt"

	"dish_as_a_service/entity"

	"github.com/Falokut/go-kit/telegram_bot"
	"github.com/pkg/errors"
)

type PaymentBot struct {
	bot          *telegram_bot.BotAPI
	invoiceToken string
}

func NewPaymentBot(token string, bot *telegram_bot.BotAPI) PaymentBot {
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

	prices := make([]telegram_bot.LabeledPrice, len(order.Items))
	for i := range order.Items {
		label := fmt.Sprintf("%s x %d", order.Items[i].Name, order.Items[i].Count)
		prices[i] = telegram_bot.LabeledPrice{
			Label:  label,
			Amount: order.Items[i].Price,
		}
	}

	invoice := telegram_bot.NewInvoice(
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
