package bot

import (
	"context"
	"github.com/Falokut/go-kit/client/telegram_bot"
	"github.com/Falokut/go-kit/log"
	"github.com/pkg/errors"
)

type BotMux interface {
	HandleMessage(ctx context.Context, msg *telegram_bot.Message) telegram_bot.Chattable
	HandlePayment(ctx context.Context, msg *telegram_bot.Message) telegram_bot.Chattable
	HandlePreCheckout(ctx context.Context, msg *telegram_bot.PreCheckoutQuery) telegram_bot.Chattable
}

type TgBot struct {
	*telegram_bot.BotAPI
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

func NewTgBot(ctx context.Context, token string, logger log.Logger) (*telegram_bot.BotAPI, error) {
	bot, err := telegram_bot.NewBotAPI(ctx, token, logger)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create bot")
	}
	logger.Debug(ctx, "bot authorized on account", log.Any("accountName", bot.Self.UserName))
	return bot, nil
}

func NewBot(bot *telegram_bot.BotAPI, mux BotMux, timeout int, logger log.Logger) TgBot {
	return TgBot{
		BotAPI:  bot,
		timeout: timeout,
		mux:     mux,
		logger:  logger,
	}
}

func (b TgBot) Run(ctx context.Context) error {
	go func() {
		u := telegram_bot.NewUpdate(0)
		u.Timeout = b.timeout
		updates := b.GetUpdatesChan(u)
		for msg := range updates {
			var resp telegram_bot.Chattable
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
			// nothing to send
			if resp == nil {
				continue
			}

			err := b.Send(resp)
			if err != nil {
				b.logger.Error(ctx, "failed to send response", log.Any("error", err))
			}
		}
	}()
	return nil
}
