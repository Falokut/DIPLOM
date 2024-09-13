package domain

import (
	"github.com/pkg/errors"
)

var (
	ErrInvalidPaymentMethod   = errors.New("невалидный способ оплаты")
	ErrUserAlreadyExists      = errors.New("пользователь уже существует")
	ErrUserOperationForbidden = errors.New("данная операция запрещена для пользователя")
	ErrWrongSecret            = errors.New("неверный пароль")
	ErrUserNotExist           = errors.New("пользователя не существует")
	ErrDishNotFound           = errors.New("не все блюда были найдены")
	ErrInvalidDishCount       = errors.New("невалидное значение количества блюд")
	ErrDishCategoryNotFound   = errors.New("категория не найдена")
	ErrDishCategoryConflict   = errors.New("категория с таким именем уже существует")
	ErrUnauthorized           = errors.New("заголовок для авторизации не передан")
)

const (
	ErrCodeInvalidDishCount     = 600
	ErrCodeDishNotFound         = 601
	ErrCodeDishCategoryNotFound = 602
	ErrCodeDishCategoryConflict = 603
	ErrCodeUserNotFound         = 604

	ErrCodeEmptyUserIdHeader = 700
	ErrCodeUserNotAdmin      = 701
)
