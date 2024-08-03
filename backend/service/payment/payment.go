package payment

import (
	"context"
	"time"

	"dish_as_a_service/domain"
	"dish_as_a_service/entity"

	"github.com/pkg/errors"

	"dish_as_a_service/repository"
	"dish_as_a_service/service/payment/expiration"
	telegram_payment "dish_as_a_service/service/payment/telegram"

	"github.com/Falokut/go-kit/log"
	"github.com/txix-open/bgjob"
)

type PaymentService interface {
	Process(ctx context.Context, order *entity.Order) (string, error)
}

type ExpirationService interface {
	AddOrder(ctx context.Context, orderId string) error
}

type Payment struct {
	paymentMethods map[string]PaymentService
	expiration     ExpirationService
}

// nolint:mnd,nonamedreturns
func NewPayment(
	bgJobCli *bgjob.Client,
	paymentBot telegram_payment.PaymentBot,
	logger log.Logger,
	userRepo repository.User,
	orderRepo repository.Order,
	expirationDelay time.Duration,
) (payment Payment, workers []*bgjob.Worker) {
	telegramWorkerService := telegram_payment.NewWorker(paymentBot)
	telegramController := telegram_payment.NewWorkerController(telegramWorkerService)

	observer := NewObserver(logger)
	telegramWorker := bgjob.NewWorker(bgJobCli,
		telegram_payment.WorkerQueue,
		telegramController,
		bgjob.WithPollInterval(5*time.Second),
		bgjob.WithObserver(observer),
	)

	expirationService := expiration.NewExpiration(bgJobCli, expirationDelay)
	expirationWorkerService := expiration.NewWorker(orderRepo)
	expirationController := expiration.NewWorkerController(expirationWorkerService)
	expirationWorker := bgjob.NewWorker(bgJobCli,
		expiration.WorkerQueue,
		expirationController,
		bgjob.WithPollInterval(5*time.Second),
		bgjob.WithObserver(observer),
	)

	var paymentMethods = map[string]PaymentService{
		telegram_payment.PaymentMethod: telegram_payment.NewPayment(userRepo, bgJobCli),
	}

	return Payment{
			paymentMethods: paymentMethods,
			expiration:     expirationService,
		},
		[]*bgjob.Worker{telegramWorker, expirationWorker}
}

func (s Payment) Process(ctx context.Context, order *entity.Order, method string) (string, error) {
	paymentService, ok := s.paymentMethods[method]
	if !ok {
		return "", domain.ErrInvalidPaymentMethod
	}
	err := s.expiration.AddOrder(ctx, order.Id)
	if err != nil {
		return "", errors.WithMessage(err, "add order to expiration")
	}

	url, err := paymentService.Process(ctx, order)
	if err != nil {
		return "", errors.WithMessage(err, "process payment")
	}

	return url, nil
}

func (s Payment) IsPaymentMethodValid(method string) bool {
	_, ok := s.paymentMethods[method]
	return ok
}
