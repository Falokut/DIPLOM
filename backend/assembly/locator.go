package assembly

import (
	"context"
	"time"

	"dish_as_a_service/bot"
	bcontroller "dish_as_a_service/bot/controller"
	broutes "dish_as_a_service/bot/routes"
	bot_service "dish_as_a_service/bot/service"
	"dish_as_a_service/conf"
	"dish_as_a_service/controller"
	"dish_as_a_service/repository"
	"dish_as_a_service/routes"
	"dish_as_a_service/service"
	"dish_as_a_service/service/events"
	"dish_as_a_service/service/payment"
	"dish_as_a_service/service/payment/expiration"
	telegram_payment "dish_as_a_service/service/payment/telegram"

	"github.com/Falokut/go-kit/client/db"
	"github.com/Falokut/go-kit/http/client"
	"github.com/Falokut/go-kit/http/endpoint"
	"github.com/Falokut/go-kit/http/router"
	"github.com/Falokut/go-kit/log"
	"github.com/Falokut/go-kit/telegram_bot"
	brouter "github.com/Falokut/go-kit/telegram_bot/router"
	"github.com/pkg/errors"
	"github.com/txix-open/bgjob"
)

type Config struct {
	BotRouter  *brouter.Router
	HttpRouter *router.Router
	Workers    []*bgjob.Worker
}

// nolint:funlen
func Locator(
	ctx context.Context,
	logger log.Logger,
	dbCli *db.Client,
	imagesCli *client.Client,
	tgBot *telegram_bot.BotAPI,
	bgJobCli *bgjob.Client,
	cfg conf.LocalConfig,
) (Config, error) {
	userRepo := repository.NewUser(dbCli)
	secret := repository.NewSecret(cfg.App.AdminSecret)
	adminEvents := events.NewAdminEvents(tgBot)
	userService := service.NewUser(userRepo, secret, adminEvents)
	userBotContr := bcontroller.NewUser(userService)

	authService := service.NewAuth(cfg.Auth, cfg.Bot.Token, userRepo)
	authCtrl := controller.NewAuth(authService)

	imagesRepo := repository.NewImage(imagesCli, cfg.Images.BaseImagePath)
	dishRepo := repository.NewDish(dbCli)
	dishService := service.NewDish(dishRepo, imagesRepo, logger)
	dishCtrl := controller.NewDish(dishService)

	dishesCategoriesRepo := repository.NewDishesCategories(dbCli)
	dishesCategoriesService := service.NewDishesCategories(dishesCategoriesRepo)
	dishesCategoriesCtrl := controller.NewDishesCategories(dishesCategoriesService)

	authMiddleware := routes.NewAuthMiddleware(cfg.Auth.Access.Secret)
	orderRepo := repository.NewOrder(dbCli)
	paymentBot := bot.NewPaymentBot(cfg.Bot.PaymentToken, tgBot, orderRepo)
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

	paymentMethods := payment.NewPaymentMethods(userRepo, bgJobCli)
	paymentService := payment.NewPayment(logger, paymentMethods, expirationService)

	orderService := service.NewOrder(paymentService, orderRepo, dishRepo)
	orderCtrl := controller.NewOrder(orderService)

	restaurantRepo := repository.NewRestaurant(dbCli)
	restaurantService := service.NewRestaurant(restaurantRepo)
	restaurantCtrl := controller.NewRestaurant(restaurantService)

	hrouter := routes.Router{
		Auth:             authCtrl,
		Dish:             dishCtrl,
		DishesCategories: dishesCategoriesCtrl,
		Order:            orderCtrl,
		Restaurant:       restaurantCtrl,
	}

	orderUserService := bot_service.NewOrderUserService(tgBot, userRepo, orderRepo)
	orderCsvExporter := service.NewCvsOrderExporter(orderRepo)
	orderBotContrl := bcontroller.NewOrder(orderService, orderUserService, orderCsvExporter)
	botControllers := broutes.Controllers{
		User:  userBotContr,
		Order: orderBotContrl,
	}
	botAdminAuth := broutes.NewAdminAuth(userRepo)
	brouter := broutes.InitRoutes(
		botControllers,
		brouter.DefaultMiddlewares(logger),
		botAdminAuth.AdminAuth,
	)

	expirationWorker := bgjob.NewWorker(bgJobCli,
		expiration.WorkerQueue,
		expirationController,
		bgjob.WithPollInterval(5*time.Second), // nolint:mnd
		bgjob.WithObserver(observer),
	)
	err := broutes.RegisterRoutes(ctx, tgBot, userRepo)
	if err != nil {
		return Config{}, errors.WithMessage(err, "register bot routes")
	}
	return Config{
		BotRouter:  brouter,
		HttpRouter: hrouter.InitRoutes(authMiddleware, endpoint.DefaultWrapper(logger)),
		Workers: []*bgjob.Worker{
			telegramWorker,
			expirationWorker,
		},
	}, nil
}
