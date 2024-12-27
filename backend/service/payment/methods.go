package payment

import (
	"dish_as_a_service/repository"

	"github.com/txix-open/bgjob"
	telegram_payment "dish_as_a_service/service/payment/telegram"
)

func NewPaymentMethods(userRepo repository.User, bgJobCli *bgjob.Client) map[string]PaymentService {
	return map[string]PaymentService{
		telegram_payment.PaymentMethod: telegram_payment.NewPayment(userRepo, bgJobCli),
	}
}
