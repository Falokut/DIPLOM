package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"dish_as_a_service/domain"
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

	query := `INSERT INTO 
	orders(id, user_id, total, created_at, wishes, payment_method, status)
	VALUES($1,$2,$3,$4,$5,$6,$7)`
	_, err = tx.ExecContext(ctx, query,
		order.Id,
		order.UserId,
		order.Total,
		order.CreatedAt,
		order.Wishes,
		order.PaymentMethod,
		order.Status,
	)
	if err != nil {
		return errors.WithMessage(err, "insert orders")
	}

	args := make([]any, 0, len(order.Items)*3+1)
	args = append(args, order.Id)
	placeholders := make([]string, len(order.Items))
	for i, item := range order.Items {
		placeholders[i] = fmt.Sprintf("($1,$%d,$%d,$%d)",
			len(args)+1,
			len(args)+2,
			len(args)+3,
		)
		args = append(args, item.DishId, item.Count, item.Price)
	}

	query = fmt.Sprintf(`INSERT INTO order_items(order_id,dish_id,count,price) VALUES %s`,
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
	query := "UPDATE orders SET status=$1 WHERE id=$2"
	_, err := r.cli.ExecContext(ctx, query, newStatus, orderId)
	if err != nil {
		return errors.WithMessage(err, "exec update query")
	}
	return nil
}

func (r Order) GetOrder(ctx context.Context, orderId string) (*entity.Order, error) {
	query := `
	SELECT
		o.id,
		o.payment_method,
		o.user_id,
		o.total, 
		o.created_at,
		o.wishes,
		o.status,
		json_agg(
			json_build_object(
			'dishId', oi.dish_id,
			'count', oi.count,
			'price', oi.price,
			'name', d.name
			)
		) AS items
    FROM orders o
    JOIN order_items oi ON o.id = oi.order_id
	JOIN dish d ON oi.dish_id = d.id
    WHERE o.id = $1
	GROUP BY o.id`
	var order entity.Order
	err := r.cli.GetContext(ctx, &order, query, orderId)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, domain.ErrOrderNotFound
	case err != nil:
		return nil, errors.WithMessage(err, "query context")
	default:
		return &order, nil
	}
}

func (r Order) GetOrderStatus(ctx context.Context, orderId string) (string, error) {
	query := "SELECT status FROM orders WHERE id=$1"
	var status string
	err := r.cli.GetContext(ctx, &status, query, orderId)
	if err != nil {
		return "", errors.WithMessage(err, "get order")
	}
	return status, nil
}

func (r Order) SetOrderStatus(ctx context.Context, orderId, oldStatus, newStatus string) error {
	query := "UPDATE orders SET status=$1 WHERE id=$2 AND status=$3"
	_, err := r.cli.ExecContext(ctx, query, newStatus, orderId, oldStatus)
	if err != nil {
		return errors.WithMessage(err, "exec update query")
	}
	return nil
}

func (r Order) SetOrderingAllowed(ctx context.Context, isAllowed bool) error {
	tx, err := r.cli.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		return errors.WithMessage(err, "begin tx")
	}
	defer tx.Rollback() //nolint:errcheck
	query := "SELECT EXISTS(SELECT id FROM allow_ordering_audit WHERE end_period IS NULL)"
	current := false
	err = tx.GetContext(ctx, &current, query)
	if err != nil {
		return errors.WithMessage(err, "execute query")
	}
	if isAllowed == current {
		return nil
	}

	if isAllowed {
		_, err := tx.ExecContext(ctx, "INSERT INTO allow_ordering_audit DEFAULT VALUES")
		if err != nil {
			return errors.WithMessage(err, "execute query")
		}
	} else {
		_, err := tx.ExecContext(ctx,
			"UPDATE allow_ordering_audit SET end_period=$1 WHERE end_period IS NULL",
			time.Now(),
		)
		if err != nil {
			return errors.WithMessage(err, "execute query")
		}
	}
	err = tx.Commit()
	if err != nil {
		return errors.WithMessage(err, "commit tx")
	}

	return nil
}
func (r Order) IsOrderingAllowed(ctx context.Context) (bool, error) {
	query := "SELECT EXISTS(SELECT id FROM allow_ordering_audit WHERE end_period IS NULL)"
	isAllowed := false
	err := r.cli.GetContext(ctx, &isAllowed, query)
	if err != nil {
		return false, errors.WithMessage(err, "execute query")
	}
	return isAllowed, nil
}

func (r Order) GetUserOrders(ctx context.Context, userId string, limit int32, offset int32) ([]entity.Order, error) {
	query := `
	SELECT
		o.id,
		o.payment_method,
		o.user_id,
		o.total, 
		o.created_at,
		o.wishes,
		o.status,
		json_agg(
			json_build_object(
			'dishId', oi.dish_id,
			'count', oi.count,
			'price', oi.price,
			'name', d.name
			)
		) AS items
    FROM orders o
    JOIN order_items oi ON o.id = oi.order_id
	JOIN dish d ON oi.dish_id = d.id
    WHERE o.user_id = $1
	GROUP BY o.id
    ORDER BY o.created_at DESC
	LIMIT $2
	OFFSET $3`
	var orders []entity.Order
	err := r.cli.SelectContext(ctx, &orders, query, userId, limit, offset)
	if err != nil {
		return nil, errors.WithMessage(err, "query context")
	}
	return orders, nil
}
