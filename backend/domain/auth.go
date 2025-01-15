package domain

import "time"

const (
	AuthHeaderName = "Authorization"
	UserIdHeader   = "X-User-Id"
	BearerToken    = "Bearer"
)

const (
	AdminType = "admin"
	UserType  = "user"
)

type LoginByTelegramRequest struct {
	InitTelegramData string
}

type LoginResponse struct {
	AccessToken  TokenResponse
	RefreshToken TokenResponse
}

type TokenResponse struct {
	Token     string
	ExpiresAt time.Time
}

type HasAdminPrivilegesResponse struct {
	HasAdminPrivileges bool
}
