package assembly

import (
	"context"
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
	"dish_as_a_service/service/payment/expiration"

	telegram_payment "dish_as_a_service/service/payment/telegram"

	"github.com/Falokut/go-kit/client/db"
	"github.com/Falokut/go-kit/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/txix-open/bgjob"
)

type Config struct {
	BotRouter   *broutes.Router
	HttpHandler http.Handler
	Workers     []*bgjob.Worker
}

func Locator(_ context.Context,
	logger log.Logger,
	dbCli *db.Client,
	tgbot *tgbotapi.BotAPI,
	bgJobCli *bgjob.Client,
	cfg conf.LocalConfig,
) Config {
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

	hrouter := routes.Router{
		Dish:             dishContrl,
		DishesCategories: dishesCategoriesContrl,
		User:             userContr,
	}
	authMiddleware := routes.NewAuthMiddleware(userRepo)
	middlewares := []echo.MiddlewareFunc{
		middleware.Recover(),
		routes.HandleError,
		routes.NewLogger(logger).LoggerMiddleware,
	}

	if cfg.Bot.Disable || tgbot == nil {
		return Config{
			HttpHandler: hrouter.InitRoutes(authMiddleware, middlewares...),
		}
	}
	orderRepo := repository.NewOrder(dbCli)
	paymentBot := bot.NewPaymentBot(cfg.Bot.PaymentToken, tgbot)
	telegramWorkerService := telegram_payment.NewWorker(paymentBot)
	telegramController := telegram_payment.NewWorkerController(telegramWorkerService)

	observer := payment.NewObserver(logger)
	telegramWorker := bgjob.NewWorker(bgJobCli,
		telegram_payment.WorkerQueue,
		telegramController,
		bgjob.WithPollInterval(5*time.Second), // nolint:mnd
		bgjob.WithObserver(observer),
	)
	expirationService := expiration.NewExpiration(bgJobCli, cfg.Payment.ExpirationDelay)
	expirationWorkerService := expiration.NewWorker(orderRepo)
	expirationController := expiration.NewWorkerController(expirationWorkerService)

	var paymentMethods = map[string]payment.PaymentService{
		telegram_payment.PaymentMethod: telegram_payment.NewPayment(userRepo, bgJobCli),
	}
	paymentService := payment.NewPayment(logger, paymentMethods, expirationService)

	orderService := service.NewOrder(paymentService, orderRepo, dishRepo)
	hrouter.Order = controller.NewOrder(orderService)

	orderUserService := bot_service.NewOrderUserService(tgbot, userRepo)
	orderBotContrl := bcontroller.NewOrder(orderService, orderUserService)
	brouter := broutes.NewRouter(logger, userBotContr, orderBotContrl)

	expirationWorker := bgjob.NewWorker(bgJobCli,
		expiration.WorkerQueue,
		expirationController,
		bgjob.WithPollInterval(5*time.Second), // nolint:mnd
		bgjob.WithObserver(observer),
	)

	return Config{
		BotRouter:   &brouter,
		HttpHandler: hrouter.InitRoutes(authMiddleware, middlewares...),
		Workers: []*bgjob.Worker{
			telegramWorker,
			expirationWorker,
		},
	}
}
