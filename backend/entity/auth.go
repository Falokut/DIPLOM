package entity

type TokenUserInfo struct {
	UserId   string `json:"user_id"`
	RoleName string `json:"role"`
}

type UserAuthInfo struct {
	UserId   string
	RoleName string
}
