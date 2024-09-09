package assembly

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"dish_as_a_service/bot"
	tgbotapi "dish_as_a_service/bot/api"
	bcontroller "dish_as_a_service/bot/controller"
	broutes "dish_as_a_service/bot/routes"
	bot_service "dish_as_a_service/bot/service"
	"dish_as_a_service/conf"
	"dish_as_a_service/controller"
	"dish_as_a_service/repository"
	"dish_as_a_service/routes"
	"dish_as_a_service/service"
	"dish_as_a_service/service/payment"

	"github.com/Falokut/go-kit/app"
	"github.com/Falokut/go-kit/client/db"
	"github.com/Falokut/go-kit/config"
	"github.com/Falokut/go-kit/healthcheck"
	"github.com/Falokut/go-kit/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/txix-open/bgjob"
)

type Assembly struct {
	logger             log.Logger
	db                 *db.Client
	tgBot              bot.TgBot
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
	tgbot, err := bot.NewTgBot(ctx, cfg.Bot.Token, cfg.Bot.Debug, logger)
	if err != nil {
		return nil, errors.WithMessage(err, "init bot")
	}
	locatorCfg := Locator(ctx, logger, dbCli, tgbot, bgjobCli, cfg)
	tgBot := bot.NewBot(tgbot, locatorCfg.BotRouter, cfg.Bot.Timeout, logger)
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

type Config struct {
	BotRouter   broutes.Router
	HttpHandler http.Handler
	Workers     []*bgjob.Worker
}

func Locator(_ context.Context,
	logger log.Logger,
	dbCli *db.Client,
	tgbot *tgbotapi.BotAPI,
	bgJobCli *bgjob.Client,
	cfg conf.LocalConfig) Config {
	userRepo := repository.NewUser(dbCli)
	secret := repository.NewSecret(cfg.App.AdminSecret)
	userService := service.NewUser(userRepo, secret)
	userContr := controller.NewUser(userService)
	userBotContr := bcontroller.NewUser(userService, cfg.App.Debug)

	imagesRepo := repository.NewImage(http.DefaultClient, cfg.Images.Addr, cfg.Images.BaseImagePath)
	dishRepo := repository.NewDish(dbCli)
	dishService := service.NewDish(dishRepo, imagesRepo, logger)
	dishContrl := controller.NewDish(dishService)

	dishesCategoriesRepo := repository.NewDishesCategories(dbCli)
	dishesCategoriesService := service.NewDishesCategories(dishesCategoriesRepo)
	dishesCategoriesContrl := controller.NewDishesCategories(dishesCategoriesService)

	orderRepo := repository.NewOrder(dbCli)
	paymentBot := bot.NewPaymentBot(cfg.Bot.PaymentToken, tgbot)
	paymentService, workers := payment.NewPayment(bgJobCli, paymentBot, logger, userRepo,
		orderRepo, cfg.Payment.ExpirationDelay)

	orderService := service.NewOrder(paymentService, orderRepo, dishRepo)
	orderContr := controller.NewOrder(orderService)

	orderUserService := bot_service.NewOrderUserService(tgbot, userRepo)
	orderBotContrl := bcontroller.NewOrder(orderService, orderUserService)

	hrouter := routes.Router{
		Dish:             dishContrl,
		DishesCategories: dishesCategoriesContrl,
		Order:            orderContr,
		User:             userContr,
	}
	middlewares := []echo.MiddlewareFunc{
		middleware.Recover(),
		routes.HandleError,
		routes.NewLogger(logger).LoggerMiddleware,
	}

	authMiddleware := routes.NewAuthMiddleware(userRepo)
	return Config{
		BotRouter:   broutes.NewRouter(logger, userBotContr, orderBotContrl),
		HttpHandler: hrouter.InitRoutes(authMiddleware, middlewares...),
		Workers:     workers,
	}
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
		a.tgBot.Close,
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
