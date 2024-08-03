package payment

import (
	"dish_as_a_service/entity"
)

type PaymentPayload struct {
	Order  entity.Order
	ChatId int64
}
