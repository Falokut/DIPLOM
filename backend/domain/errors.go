package domain

import (
	"github.com/pkg/errors"
)

var (
	ErrInvalidPaymentMethod   = errors.New("невалидный способ оплаты")
	ErrUserAlreadyExists      = errors.New("пользователь уже существует")
	ErrUserNotFound           = errors.New("пользователь не найден")
	ErrUserOperationForbidden = errors.New("данная операция запрещена для пользователя")
	ErrWrongSecret            = errors.New("неверный пароль")
	ErrDishNotFound           = errors.New("не все блюда были найдены")
	ErrInvalidDishCount       = errors.New("невалидное значение количества блюд")
	ErrDishCategoryNotFound   = errors.New("категория не найдена")
	ErrDishCategoryConflict   = errors.New("категория с таким именем уже существует")
	ErrUnauthorized           = errors.New("заголовок для авторизации не передан")
	ErrOrderingForbidden      = errors.New("оформление заказов приостановлено")
	ErrOrderNotFound          = errors.New("заказ не найден")
)

const (
	ErrCodeInvalidArgument = 400

	ErrCodeInvalidDishCount     = 600
	ErrCodeDishNotFound         = 601
	ErrCodeDishCategoryNotFound = 602
	ErrCodeDishCategoryConflict = 603
	ErrCodeUserNotFound         = 604
	ErrCodeUserAlreadyExists    = 605
	ErrCodeWrongSecret          = 606
	ErrCodeOrderingForbidden    = 607

	ErrCodeEmptyUserIdHeader = 700
	ErrCodeUserNotAdmin      = 701
)
