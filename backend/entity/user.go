package entity

type User struct {
	Username string
	Name     string
	Admin    bool
}

type RegisterUser struct {
	Id       string
	Username string
	Name     string
	Telegram *Telegram
}

type Telegram struct {
	ChatId int64
	UserId int64
}

type TelegramUser struct {
	ChatId int64
	Admin  bool
}
