package domain

import (
	"github.com/pkg/errors"
)

var (
	ErrInvalidPaymentMethod           = errors.New("невалидный способ оплаты")
	ErrUserAlreadyExists              = errors.New("пользователь уже существует")
	ErrUserOperationForbidden         = errors.New("данная операция запрещена для пользователя")
	ErrWrongSecret                    = errors.New("неверный пароль")
	ErrUserNotExist                   = errors.New("пользователя не существует")
	ErrDishNotFound                   = errors.New("не все блюда были найдены")
	ErrInvalidDishCount               = errors.New("невалидное значение количества блюд")
	ErrDishCategoryNotFound           = errors.New("категория не найдена")
	ErrDishCategoryNotFoundOrConflict = errors.New("категория не найдена или категория с таким именем уже существует")
)
