package domain

type User struct {
	Username string
	Name     string
	Admin    bool
}

type RegisterUser struct {
	Username string    `validate:"required,min=3"`
	Name     string    `validate:"required,min=3"`
	Telegram *Telegram `json:",omitempty"`
}

type GetUserIdByTelegramIdRequest struct {
	TelegramId int64 `path:"telegram_id" validate:"required"`
}

type IsUserAdminRequest struct {
	UserId string `path:"user_id" validate:"required"`
}

type Telegram struct {
	ChatId int64
	UserId int64
}

type AddAdminSecretRequest struct {
	ChatId int64
	Secret string
}

type GetUserIdByTelegramIdResponse struct {
	UserId string
}

type IsUserAdminResponse struct {
	IsAdmin bool
}
