package domain

type User struct {
	Username string
	Name     string
	Admin    bool
}

type RegisterUser struct {
	Username string
	Name     string
	Telegram *Telegram `json:",omitempty"`
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
