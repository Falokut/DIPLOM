package assembly

import (
	"context"
	"net/http"
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
	"dish_as_a_service/service/payment"
	"dish_as_a_service/service/payment/expiration"

	telegram_payment "dish_as_a_service/service/payment/telegram"

	"github.com/Falokut/go-kit/client/db"
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
	tgBot *telegram_bot.BotAPI,
	bgJobCli *bgjob.Client,
	cfg conf.LocalConfig,
) (Config, error) {
	userRepo := repository.NewUser(dbCli)
	secret := repository.NewSecret(cfg.App.AdminSecret)
	userService := service.NewUser(userRepo, secret)
	userService.SetRefreshAdminCommands(func(ctx context.Context) error {
		return broutes.RegisterRoutes(ctx, tgBot, userService)
	})
	userContr := controller.NewUser(userService)
	userBotContr := bcontroller.NewUser(userService)

	imagesRepo := repository.NewImage(http.DefaultClient, cfg.Images.BaseServiceUrl, cfg.Images.BaseImagePath)
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
	orderRepo := repository.NewOrder(dbCli)
	paymentBot := bot.NewPaymentBot(cfg.Bot.PaymentToken, tgBot)
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

	orderUserService := bot_service.NewOrderUserService(tgBot, userRepo)
	orderBotContrl := bcontroller.NewOrder(orderService, orderUserService)
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
	err := broutes.RegisterRoutes(ctx, tgBot, userService)
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
