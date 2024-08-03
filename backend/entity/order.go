package entity

import "time"

const (
	OrderItemStatusProcess  = "PROCESS"
	OrderItemStatusCanceled = "CANCELED"
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
	Status string
	Name   string
}

type Order struct {
	Id        string
	Items     []OrderItem
	UserId    string
	Total     int32
	CreatedAt time.Time
	Wishes    string
}
