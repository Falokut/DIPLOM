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

type Dish struct {
	cli *db.Client
}

func NewDish(cli *db.Client) Dish {
	return Dish{
		cli: cli,
	}
}

func (r Dish) List(ctx context.Context, limit, offset int32) ([]entity.Dish, error) {
	query := `
	SELECT
		d.id,
		d.name,
		d.description,
	    d.price, 
		COALESCE(d.image_id,'') AS image_id,
		array_to_string(ARRAY_AGG(COALESCE(c.name,'')),',') AS categories,
		r.name AS restaurant_name 
	FROM dish AS d
	JOIN restaurants AS r ON d.restaurant_id = r.id
	LEFT JOIN dish_categories AS f_c ON d.id=f_c.dish_id
	LEFT JOIN categories AS c ON f_c.category_id=c.id
	GROUP BY d.id, d.name, d.description, d.price, d.image_id, r.name
	ORDER BY d.id
	LIMIT $1 OFFSET $2`
	var res []entity.Dish
	err := r.cli.SelectContext(ctx, &res, query, limit, offset)
	if err != nil {
		return nil, errors.WithMessage(err, "get dish list")
	}
	return res, nil
}

func (r Dish) AddDish(ctx context.Context, req *entity.AddDishRequest) error {
	tx, err := r.cli.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
	})
	if err != nil {
		return errors.WithMessage(err, "begin tx")
	}
	defer tx.Rollback() //nolint:errcheck

	query := `INSERT INTO dish(name, description, price, image_id, restaurant_id) VALUES($1,$2,$3,$4,$5) RETURNING id;`
	var id int32
	err = tx.GetContext(ctx, &id, query, req.Name, req.Description, req.Price, req.ImageId, req.RestaurantId)
	if err != nil {
		return errors.WithMessage(err, "insert dish")
	}
	query, args := getInsertDishCategoriesQuery(id, req.Categories)
	if len(args) > 0 {
		_, err := tx.ExecContext(ctx, query, args...)
		if err != nil {
			return errors.WithMessage(err, "insert dish_categories")
		}
	}
	err = tx.Commit()
	if err != nil {
		return errors.WithMessage(err, "commit tx")
	}
	return nil
}

func (r Dish) GetDishesByIds(ctx context.Context, ids []int32) ([]entity.Dish, error) {
	query := `
	SELECT 
		d.id,
		d.name,
		d.description,
		d.price,
		COALESCE(d.image_id,'') AS image_id,
		array_to_string(ARRAY_AGG(COALESCE(c.name,'')),',') AS categories,
		r.name AS restaurant_name 
	FROM dish AS d
	JOIN restaurants AS r ON d.restaurant_id = r.id
	LEFT JOIN dish_categories AS f_c ON d.id=f_c.dish_id
	LEFT JOIN categories AS c ON f_c.category_id=c.id
	WHERE d.id=ANY($1)
	GROUP BY d.id, d.name, d.description, d.price, d.image_id, r.name
	ORDER BY d.id;`

	var res []entity.Dish
	err := r.cli.SelectContext(ctx, &res, query, ids)
	if err != nil {
		return nil, errors.WithMessage(err, "get dish list")
	}
	return res, nil
}

func (r Dish) GetDishesByCategories(ctx context.Context, limit int32, offset int32, ids []int32) ([]entity.Dish, error) {
	query := `
	SELECT 
		d.id,
		d.name,
		d.description,
		d.price,
		COALESCE(d.image_id,'') AS image_id,
		array_to_string(ARRAY_AGG(COALESCE(c.name,'')),',') AS categories,
		r.name AS restaurant_name 
	FROM dish AS d
	JOIN restaurants AS r ON d.restaurant_id = r.id
	LEFT JOIN dish_categories AS f_c ON d.id=f_c.dish_id
	LEFT JOIN categories AS c ON f_c.category_id=c.id
	GROUP BY d.id, d.name, d.description, d.price, d.image_id, r.name
	HAVING array_agg(c.id) @> $1
	ORDER BY d.id
	LIMIT $2 OFFSET $3;`

	var res []entity.Dish
	err := r.cli.SelectContext(ctx, &res, query, ids, limit, offset)
	if err != nil {
		return nil, errors.WithMessage(err, "get dish list")
	}
	return res, nil
}

func (r Dish) EditDish(ctx context.Context, req *entity.EditDishRequest) error {
	tx, err := r.cli.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
	})
	if err != nil {
		return errors.WithMessage(err, "begin tx")
	}
	defer tx.Rollback() //nolint:errcheck

	query := `UPDATE dish SET name=$1, description=$2, price=$3, image_id=$4, restaurant_id=$5 WHERE id=$6`
	_, err = tx.ExecContext(ctx, query, req.Name, req.Description, req.Price, req.ImageId, req.RestaurantId, req.Id)
	if err != nil {
		return errors.WithMessage(err, "update dish")
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM dish_categories WHERE dish_id=$1", req.Id)
	if err != nil {
		return errors.WithMessage(err, "delete dish categories")
	}
	query, args := getInsertDishCategoriesQuery(req.Id, req.Categories)
	if len(args) > 0 {
		_, err := tx.ExecContext(ctx, query, args...)
		if err != nil {
			return errors.WithMessage(err, "insert dish_categories")
		}
	}
	err = tx.Commit()
	if err != nil {
		return errors.WithMessage(err, "commit tx")
	}
	return nil
}

func (r Dish) DeleteDish(ctx context.Context, id int32) error {
	_, err := r.cli.ExecContext(ctx, "DELETE FROM dish WHERE id=$1", id)
	if err != nil {
		return errors.WithMessage(err, "delete dishes")
	}
	return nil
}

func getInsertDishCategoriesQuery(id int32, categories []int32) (string, []any) {
	if len(categories) == 0 {
		return "", nil
	}
	var valuesPlaceholders = make([]string, len(categories))
	var args = make([]any, 0, len(categories)+1)
	args = append(args, id)
	for i, catId := range categories {
		valuesPlaceholders[i] = fmt.Sprintf("($1,$%d)", len(args)+1)
		args = append(args, catId)
	}
	return fmt.Sprintf(`INSERT INTO dish_categories(dish_id,category_id) VALUES %s ON CONFLICT DO NOTHING`,
		strings.Join(valuesPlaceholders, ",")), args
}
