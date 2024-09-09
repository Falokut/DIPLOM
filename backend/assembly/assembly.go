package assembly

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"dish_as_a_service/bot"
	tgbotapi "dish_as_a_service/bot/api"
	"dish_as_a_service/conf"

	"github.com/Falokut/go-kit/app"
	"github.com/Falokut/go-kit/client/db"
	"github.com/Falokut/go-kit/config"
	"github.com/Falokut/go-kit/healthcheck"
	"github.com/Falokut/go-kit/log"
	"github.com/pkg/errors"
	"github.com/txix-open/bgjob"
)

type Assembly struct {
	logger             log.Logger
	db                 *db.Client
	tgBot              *bot.TgBot
	workers            []*bgjob.Worker
	bgjobCli           *bgjob.Client
	server             *http.Server
	healthcheckManager *healthcheck.Manager
}

func New(ctx context.Context, logger log.Logger) (*Assembly, error) {
	cfg := conf.LocalConfig{}
	err := config.Read(&cfg)
	if err != nil {
		return nil, errors.WithMessage(err, "read local config")
	}
	dbCli, err := db.NewDB(ctx, cfg.DB, db.WithMigrationRunner("./migrations", logger))
	if err != nil {
		return nil, errors.WithMessage(err, "init db")
	}
	bgjobDb := bgjob.NewPgStore(dbCli.DB.DB)
	bgjobCli := bgjob.NewClient(bgjobDb)
	var tgbotApi *tgbotapi.BotAPI
	if !cfg.Bot.Disable {
		tgbotApi, err = bot.NewTgBot(ctx, cfg.Bot.Token, cfg.Bot.Debug, logger)
		if err != nil {
			return nil, errors.WithMessage(err, "init bot")
		}
	}

	locatorCfg := Locator(ctx, logger, dbCli, tgbotApi, bgjobCli, cfg)
	var tgBot *bot.TgBot
	if locatorCfg.BotRouter != nil {
		bot := bot.NewBot(tgbotApi, locatorCfg.BotRouter, cfg.Bot.Timeout, logger)
		tgBot = &bot
	}

	healthcheckManager := healthcheck.NewHealthManager(logger, fmt.Sprint(cfg.HealthcheckPort))
	healthcheckManager.Register("db", dbCli.PingContext)

	server := &http.Server{
		Addr:              cfg.Listen.GetAddress(),
		ReadHeaderTimeout: time.Second * 15, //nolint:mnd
		Handler:           locatorCfg.HttpHandler,
	}
	return &Assembly{
		logger:             logger,
		db:                 dbCli,
		tgBot:              tgBot,
		server:             server,
		workers:            locatorCfg.Workers,
		bgjobCli:           bgjobCli,
		healthcheckManager: &healthcheckManager,
	}, nil
}

func (a *Assembly) Runners() []app.RunnerFunc {
	return []app.RunnerFunc{
		a.tgBot.Run,
		func(_ context.Context) error {
			return a.server.ListenAndServe()
		},
		func(_ context.Context) error {
			return a.healthcheckManager.RunHealthcheckEndpoint()
		},
		func(ctx context.Context) error {
			for _, worker := range a.workers {
				worker.Run(ctx)
			}
			return nil
		},
	}
}

func (a *Assembly) Closers() []app.CloserFunc {
	return []app.CloserFunc{
		func(_ context.Context) error {
			return a.db.Close()
		},
		func(ctx context.Context) error {
			if a.tgBot != nil {
				return a.tgBot.Close(ctx)
			}
			return nil
		},
		func(context.Context) error {
			return a.server.Close()
		},
		func(_ context.Context) error {
			for _, worker := range a.workers {
				worker.Shutdown()
			}
			return nil
		},
	}
}
