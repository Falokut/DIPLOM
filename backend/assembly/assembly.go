package assembly

import (
	"context"
	"fmt"

	"github.com/Falokut/go-kit/http"

	"dish_as_a_service/conf"

	"github.com/Falokut/go-kit/app"
	"github.com/Falokut/go-kit/client/db"
	"github.com/Falokut/go-kit/config"
	"github.com/Falokut/go-kit/healthcheck"
	"github.com/Falokut/go-kit/log"
	"github.com/Falokut/go-kit/telegram_bot"
	"github.com/pkg/errors"
	"github.com/txix-open/bgjob"
)

type Assembly struct {
	logger             log.Logger
	db                 *db.Client
	tgBot              *telegram_bot.BotAPI
	workers            []*bgjob.Worker
	bgjobCli           *bgjob.Client
	server             *http.Server
	healthcheckManager *healthcheck.Manager
	localCfg           conf.LocalConfig
}

func New(ctx context.Context, logger log.Logger) (*Assembly, error) {
	localCfg := conf.LocalConfig{}
	err := config.Read(&localCfg)
	if err != nil {
		return nil, errors.WithMessage(err, "read local config")
	}
	dbCli, err := db.NewDB(ctx, localCfg.DB, db.WithMigrationRunner("./migrations", logger))
	if err != nil {
		return nil, errors.WithMessage(err, "init db")
	}
	bgjobDb := bgjob.NewPgStore(dbCli.DB.DB)
	bgjobCli := bgjob.NewClient(bgjobDb)
	tgBot, err := telegram_bot.NewBotAPI(ctx, localCfg.Bot.Token, logger)
	if err != nil {
		return nil, errors.WithMessage(err, "init bot")
	}
	server := http.NewServer(logger)

	locatorCfg, err := Locator(ctx, logger, dbCli, tgBot, bgjobCli, localCfg)
	if err != nil {
		return nil, errors.WithMessage(err, "locator config")
	}
	server.Upgrade(locatorCfg.HttpRouter)
	tgBot.Upgrade(locatorCfg.BotRouter)

	healthcheckManager := healthcheck.NewHealthManager(logger, fmt.Sprint(localCfg.HealthcheckPort))
	healthcheckManager.Register("db", dbCli.PingContext)

	return &Assembly{
		logger:             logger,
		localCfg:           localCfg,
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
		func(ctx context.Context) error {
			return a.tgBot.Serve(ctx, telegram_bot.UpdatesConfig{
				Timeout: a.localCfg.Bot.Timeout,
				Limit:   a.localCfg.Bot.Limit,
			})
		},
		func(_ context.Context) error {
			return a.server.ListenAndServe(a.localCfg.Listen.GetAddress())
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
		func(_ context.Context) error {
			a.tgBot.StopReceivingUpdates()
			return nil
		},
		func(ctx context.Context) error {
			return a.server.Shutdown(ctx)
		},
		func(_ context.Context) error {
			for _, worker := range a.workers {
				worker.Shutdown()
			}
			return nil
		},
	}
}
