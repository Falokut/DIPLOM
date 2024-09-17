package entity

import (
	"github.com/Falokut/go-kit/json"
	"github.com/pkg/errors"
	"time"
)

const (
	OrderItemStatusProcess  = "PROCESS"
	OrderItemStatusCanceled = "CANCELED"
	OrderItemStatusPaid     = "PAID"
	OrderItemStatusSuccess  = "SUCCESS"
)

type ProcessOrderRequest struct {
	Items  []OrderItem
	UserId string
	Total  int32
}

type PaymentPayload struct {
	ChatId  int64
	OrderId string
}

type OrderItem struct {
	DishId int32
	Count  int32
	Price  int32
	Name   string
}

type Order struct {
	Id            string
	PaymentMethod string
	Items         OrderItems
	UserId        string
	Total         int32
	CreatedAt     time.Time
	Status        string
	Wishes        string
}

type OrderItems []OrderItem

func (o *OrderItems) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.Errorf("failed to scan OrderItems: %v", value)
	}
	return json.Unmarshal(bytes, o)
}
