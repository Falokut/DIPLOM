package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"dish_as_a_service/entity"
	"github.com/Falokut/go-kit/client/db"
	"github.com/pkg/errors"
)

type Order struct {
	cli *db.Client
}

func NewOrder(cli *db.Client) Order {
	return Order{
		cli: cli,
	}
}

//nolint:mnd
func (r Order) ProcessOrder(ctx context.Context, order *entity.Order) error {
	tx, err := r.cli.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return errors.WithMessage(err, "begin tx")
	}
	defer tx.Rollback() //nolint:errcheck

	query := "INSERT INTO orders(id, user_id, total, created_at, wishes) VALUES($1, $2, $3, $4, $5)"
	_, err = tx.ExecContext(ctx, query, order.Id, order.UserId, order.Total, order.CreatedAt, order.Wishes)
	if err != nil {
		return errors.WithMessage(err, "insert orders")
	}

	args := make([]any, 0, len(order.Items)*4+1)
	args = append(args, order.Id)
	placeholders := make([]string, len(order.Items))
	for i, item := range order.Items {
		placeholders[i] = fmt.Sprintf("($1,$%d,$%d,$%d,$%d)",
			len(args)+1,
			len(args)+2,
			len(args)+3,
			len(args)+4,
		)
		args = append(args, item.DishId, item.Count, item.Price, item.Status)
	}

	query = fmt.Sprintf(`INSERT INTO order_items(order_id,dish_id,count,price,status) VALUES %s`,
		strings.Join(placeholders, ","))
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.WithMessage(err, "insert order items")
	}
	err = tx.Commit()
	if err != nil {
		return errors.WithMessage(err, "commit tx")
	}
	return nil
}

func (r Order) UpdateOrderStatus(ctx context.Context, orderId, newStatus string) error {
	query := "UPDATE order_items SET status=$1 WHERE order_id=$2"
	_, err := r.cli.ExecContext(ctx, query, newStatus, orderId)
	if err != nil {
		return errors.WithMessage(err, "exec update query")
	}
	return nil
}

func (r Order) IsOrderCanceled(ctx context.Context, orderId string) (bool, error) {
	query := "SELECT EXISTS(SELECT order_id FROM order_items WHERE order_id=$1 AND status=$2)"

	var canceled bool
	err := r.cli.GetContext(ctx, &canceled, query, orderId, entity.OrderItemStatusCanceled)
	if err != nil {
		return false, errors.WithMessage(err, "exec update query")
	}
	return canceled, nil
}

func (r Order) SetOrderStatus(ctx context.Context, orderId, oldStatus, newStatus string) error {
	query := "UPDATE order_items SET status=$1 WHERE order_id=$2 AND status=$3"
	_, err := r.cli.ExecContext(ctx, query, newStatus, orderId, oldStatus)
	if err != nil {
		return errors.WithMessage(err, "exec update query")
	}
	return nil
}
