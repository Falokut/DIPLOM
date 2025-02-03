package repository

import (
	"context"
	"database/sql"
	"dish_as_a_service/domain"
	"dish_as_a_service/entity"

	"github.com/Falokut/go-kit/client/db"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
)

type Restaurant struct {
	cli *db.Client
}

func NewRestaurant(cli *db.Client) Restaurant {
	return Restaurant{cli: cli}
}
func (r Restaurant) GetAllRestaurants(ctx context.Context) ([]entity.Restaurant, error) {
	var restaurants []entity.Restaurant
	err := r.cli.SelectContext(ctx, &restaurants, "SELECT id, name FROM restaurants ORDER BY id")
	if err != nil {
		return nil, errors.WithMessage(err, "execute query")
	}
	return restaurants, nil
}

func (r Restaurant) GetRestaurant(ctx context.Context, id int32) (entity.Restaurant, error) {
	var restaurantName entity.Restaurant
	err := r.cli.GetContext(ctx, &restaurantName, "SELECT id, name FROM restaurants WHERE id=$1", id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return entity.Restaurant{}, domain.ErrRestaurantNotFound
	case err != nil:
		return entity.Restaurant{}, errors.WithMessage(err, "execute query")
	default:
		return restaurantName, nil
	}
}

func (r Restaurant) AddRestaurant(ctx context.Context, restaurantName string) (int32, error) {
	query := `WITH e AS(
    INSERT INTO restaurants (name) 
           VALUES ($1) 
    ON CONFLICT DO NOTHING
    RETURNING id
	)
	SELECT * FROM e UNION SELECT id FROM restaurants WHERE name=$1;`

	var id int32
	err := r.cli.GetContext(ctx, &id, query, restaurantName)
	if err != nil {
		return 0, errors.WithMessage(err, "execute query")
	}
	return id, nil
}

func (r Restaurant) RenameRestaurant(ctx context.Context, id int32, newName string) error {
	_, err := r.cli.ExecContext(ctx, "UPDATE restaurants SET name = $1 WHERE id = $2", newName, id)
	var pgErr *pgconn.PgError
	switch {
	case errors.As(err, &pgErr) && pgErr.SQLState() == pgerrcode.UniqueViolation:
		return domain.ErrRestaurantConflict
	case err != nil:
		return errors.WithMessage(err, "execute query")
	default:
		return nil
	}
}

func (r Restaurant) DeleteRestaurant(ctx context.Context, id int32) error {
	_, err := r.cli.ExecContext(ctx, "DELETE FROM restaurants WHERE id=$1", id)
	if err != nil {
		return errors.WithMessage(err, "execute query")
	}
	return nil
}
