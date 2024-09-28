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

type DishesCategories struct {
	cli *db.Client
}

func NewDishesCategories(cli *db.Client) DishesCategories {
	return DishesCategories{cli: cli}
}
func (r DishesCategories) GetAllCategories(ctx context.Context) ([]entity.DishCategory, error) {
	var categories []entity.DishCategory
	err := r.cli.SelectContext(ctx, &categories, "SELECT id, name FROM categories ORDER BY id")
	if err != nil {
		return nil, errors.WithMessage(err, "execute query")
	}
	return categories, nil
}

func (r DishesCategories) GetDishesCategories(ctx context.Context) ([]entity.DishCategory, error) {
	var categories []entity.DishCategory
	query := `
	SELECT DISTINCT c.id, c.name
	FROM dish_categories dc
	JOIN categories c ON dc.category_id = c.id
	ORDER BY c.id;`
	err := r.cli.SelectContext(ctx, &categories, query)
	if err != nil {
		return nil, errors.WithMessage(err, "execute query")
	}
	return categories, nil
}

func (r DishesCategories) GetCategory(ctx context.Context, id int32) (entity.DishCategory, error) {
	var category entity.DishCategory
	err := r.cli.GetContext(ctx, &category, "SELECT id, name FROM categories WHERE id=$1", id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return entity.DishCategory{}, domain.ErrDishCategoryNotFound
	case err != nil:
		return entity.DishCategory{}, errors.WithMessage(err, "execute query")
	default:
		return category, nil
	}
}

func (r DishesCategories) AddCategory(ctx context.Context, category string) (int32, error) {
	query := `WITH e AS(
    INSERT INTO categories (name) 
           VALUES ($1) 
    ON CONFLICT DO NOTHING
    RETURNING id
	)
	SELECT * FROM e UNION SELECT id FROM categories WHERE name=$1;`

	var id int32
	err := r.cli.GetContext(ctx, &id, query, category)
	if err != nil {
		return 0, errors.WithMessage(err, "execute query")
	}
	return id, nil
}

func (r DishesCategories) RenameCategory(ctx context.Context, id int32, newName string) error {
	_, err := r.cli.ExecContext(ctx, "UPDATE categories SET name = $1 WHERE id = $2", newName, id)
	var pgErr *pgconn.PgError
	switch {
	case errors.As(err, &pgErr) && pgErr.SQLState() == pgerrcode.UniqueViolation:
		return domain.ErrDishCategoryConflict
	case err != nil:
		return errors.WithMessage(err, "execute query")
	default:
		return nil
	}
}

func (r DishesCategories) DeleteCategory(ctx context.Context, id int32) error {
	_, err := r.cli.ExecContext(ctx, "DELETE FROM categories WHERE id=$1", id)
	if err != nil {
		return errors.WithMessage(err, "execute query")
	}
	return nil
}
