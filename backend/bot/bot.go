package bot

import (
	"context"
	tgbotapi "dish_as_a_service/bot/api"
	"fmt"
	"github.com/Falokut/go-kit/log"
	"github.com/pkg/errors"
)

type BotMux interface {
	HandleMessage(ctx context.Context, msg *tgbotapi.Message) tgbotapi.Chattable
	HandlePayment(ctx context.Context, msg *tgbotapi.Message) tgbotapi.Chattable
	HandlePreCheckout(ctx context.Context, msg *tgbotapi.PreCheckoutQuery) tgbotapi.Chattable
}

type TgBot struct {
	*tgbotapi.BotAPI
	timeout int
	mux     BotMux
	logger  log.Logger
}

type Config struct {
	Disable      bool   `yaml:"disable" env:"TG_BOT_DISABLE"`
	Token        string `yaml:"token" env:"TG_BOT_TOKEN"`
	PaymentToken string `yaml:"payment_token" env:"TG_BOT_PAYMENT_TOKEN"`
	Debug        bool   `yaml:"debug" env:"TG_BOT_DEBUG"`
	// timeout in seconds
	Timeout int `yaml:"timeout" env:"TG_BOT_TIMEOUT"`
}

func NewTgBot(ctx context.Context, token string, debug bool, logger log.Logger) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		return nil, errors.WithMessage(err, "failed to create bot")
	}
	bot.Debug = debug
	logger.Info(ctx, fmt.Sprint("Authorized on account ", bot.Self.UserName))
	return bot, nil
}

func NewBot(bot *tgbotapi.BotAPI, mux BotMux, timeout int, logger log.Logger) TgBot {
	return TgBot{
		BotAPI:  bot,
		timeout: timeout,
		mux:     mux,
		logger:  logger,
	}
}

func (b TgBot) Run(ctx context.Context) error {
	go func() {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = b.timeout
		updates := b.GetUpdatesChan(u)
		for msg := range updates {
			select {
			case <-ctx.Done():
				b.StopReceivingUpdates()
				return
			default:
			}

			var resp tgbotapi.Chattable
			switch {
			case msg.Message != nil:
				switch {
				case msg.Message.SuccessfulPayment != nil:
					resp = b.mux.HandlePayment(ctx, msg.Message)
				default:
					resp = b.mux.HandleMessage(ctx, msg.Message)
				}
			case msg.PreCheckoutQuery != nil:
				resp = b.mux.HandlePreCheckout(ctx, msg.PreCheckoutQuery)
			default:
				continue
			}

			if resp == nil {
				b.logger.Debug(ctx, "handle return nil resp")
				continue
			}

			err := b.Send(resp)
			if err != nil {
				b.logger.Error(ctx, fmt.Sprint("failed to send response: ", err.Error()))
			}
		}
	}()
	return nil
}

func (b TgBot) Close(_ context.Context) error {
	b.StopReceivingUpdates()
	return nil
}
