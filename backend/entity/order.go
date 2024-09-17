package entity

import (
	"fmt"
	"strings"
	"time"

	"github.com/Falokut/go-kit/json"
	"github.com/pkg/errors"
)

const (
	OrderItemStatusProcess  = "PROCESS"
	OrderItemStatusCanceled = "CANCELED"
	OrderItemStatusPaid     = "PAID"
	OrderItemStatusSuccess  = "SUCCESS"
)

const (
	NotifyArrivalCommand = "notify_arrival"
	CancelOrderCommand   = "cancel_order"
	SuccessOrderCommand  = "success_order"
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
	return json.Unmarshal(bytes, o) //nolint:wrapcheck
}

type QueryCallbackPayload struct {
	Command string
	OrderId string
}

func (q QueryCallbackPayload) String() string {
	return fmt.Sprintf("%s;%s", q.Command, q.OrderId)
}

func (q *QueryCallbackPayload) FromString(str string) error {
	parts := strings.Split(str, ";")
	// nolint:mnd
	if len(parts) != 2 {
		return errors.New("invalid query payload")
	}
	q.Command = parts[0]
	q.OrderId = parts[1]
	return nil
}
