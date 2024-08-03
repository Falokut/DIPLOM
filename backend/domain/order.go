package domain

type ProcessOrderRequest struct {
	Items         map[string]int32 `validate:"required"`
	PaymentMethod string           `validate:"required,min=1"`
	UserId        string           `validate:"required,uuid"`
	Wishes        string           `json:",omitempty"`
}

type ProcessOrderResponse struct {
	// for some payment methods may be empty
	PaymentUrl string
}
