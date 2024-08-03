package repository

import (
	"context"
	"database/sql"

	"dish_as_a_service/domain"
	"dish_as_a_service/entity"
	"github.com/Falokut/go-kit/client/db"
	"github.com/pkg/errors"
)

type User struct {
	cli *db.Client
}

func NewUser(cli *db.Client) User {
	return User{
		cli: cli,
	}
}

func (r User) Register(ctx context.Context, user entity.RegisterUser) error {
	tx, err := r.cli.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return errors.WithMessage(err, "begin tx")
	}
	defer tx.Rollback() //nolint:errcheck

	query := `INSERT INTO users (id, username, name) VALUES($1,$2,$3)`
	_, err = tx.ExecContext(ctx, query, user.Id, user.Username, user.Name)
	if err != nil {
		return errors.WithMessage(err, "insert users")
	}

	if user.Telegram != nil {
		query = `INSERT INTO users_telegrams (id,chat_id,telegram_id) VALUES($1,$2,$3) ON CONFLICT DO NOTHING`
		_, err = tx.ExecContext(ctx, query, user.Id, user.Telegram.ChatId, user.Telegram.UserId)
		if err != nil {
			return errors.WithMessage(err, "insert telegram users")
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.WithMessage(err, "commit tx")
	}
	return nil
}

func (r User) IsAdmin(ctx context.Context, id string) (bool, error) {
	query := `SELECT admin FROM users WHERE id=$1`

	var isAdmin bool
	err := r.cli.GetContext(ctx, &isAdmin, query, id)
	if err != nil {
		return false, errors.WithMessage(err, "get data from db")
	}

	return isAdmin, nil
}

func (r User) GetUsers(ctx context.Context) ([]entity.User, error) {
	query := "SELECT username, name, admin FROM users"
	var res []entity.User
	err := r.cli.SelectContext(ctx, &res, query)
	if err != nil {
		return nil, errors.WithMessage(err, "select users")
	}

	return res, nil
}

func (r User) GetUserChatId(ctx context.Context, userId string) (int64, error) {
	query := "SELECT chat_id FROM users_telegrams WHERE id=$1"
	var chatId int64
	err := r.cli.GetContext(ctx, &chatId, query, userId)
	if err != nil {
		return int64(0), errors.WithMessage(err, "select users")
	}

	return chatId, nil
}

func (r User) GetUserIdByTelegramId(ctx context.Context, telegramId int64) (string, error) {
	query := "SELECT id FROM users_telegrams WHERE telegram_id=$1"
	var userId string
	err := r.cli.GetContext(ctx, &userId, query, telegramId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", domain.ErrUserNotExist
		}
		return "", errors.WithMessage(err, "select users_telegrams")
	}
	return userId, nil
}

func (r User) AddAdmin(ctx context.Context, username string) error {
	query := "UPDATE users SET admin='true' WHERE username=$1"
	res, err := r.cli.ExecContext(ctx, query, username)
	if err != nil {
		return errors.WithMessage(err, "select users")
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return domain.ErrUserNotExist
	}
	return nil
}

func (r User) AddAdminChatId(ctx context.Context, chatId int64) error {
	query := "UPDATE users SET admin='true' WHERE chat_id=$1"
	res, err := r.cli.ExecContext(ctx, query, chatId)
	if err != nil {
		return errors.WithMessage(err, "select users")
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return domain.ErrUserNotExist
	}
	return nil
}
