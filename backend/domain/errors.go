package domain

import (
	"github.com/pkg/errors"
)

var (
	ErrUserAlreadyExists = errors.New("пользователь уже существует")
	ErrWrongSecret       = errors.New("неверный пароль")
	ErrUserNotExist      = errors.New("пользователя не существует")
)
