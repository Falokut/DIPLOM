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
	SELECT d.id, d.name, d.description, d.price, COALESCE(d.image_id,'') AS image_id,
	array_to_string(ARRAY_AGG(COALESCE(c.name,'')),',') AS categories
	FROM dish AS d
	LEFT JOIN dish_categories AS f_c ON d.id=f_c.dish_id
	LEFT JOIN categories AS c ON f_c.category_id=c.id
	GROUP BY d.id, d.name, d.description, d.price, d.image_id
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

	query := `INSERT INTO dish(name, description, price, image_id) VALUES($1,$2,$3,$4) RETURNING id;`
	var id int64
	err = tx.GetContext(ctx, &id, query, req.Name, req.Description, req.Price, req.ImageId)
	if err != nil {
		return errors.WithMessage(err, "insert dish")
	}

	if len(req.Categories) > 0 {
		var valuesPlaceholders = make([]string, len(req.Categories))
		var args = make([]any, 0, len(req.Categories)+1)
		args = append(args, id)
		for i, catId := range req.Categories {
			valuesPlaceholders[i] = fmt.Sprintf("($1,$%d)", len(args)+i+1)
			args = append(args, catId)
		}
		query = fmt.Sprintf(`INSERT INTO dish_categories(dish_id,category_id) VALUES %s ON CONFLICT DO NOTHING`,
			strings.Join(valuesPlaceholders, ","))
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
	SELECT d.id, d.name, d.description, d.price, COALESCE(d.image_id,'') AS image_id,
	array_to_string(ARRAY_AGG(COALESCE(c.name,'')),',') AS categories
	FROM dish AS d
	LEFT JOIN dish_categories AS f_c ON d.id=f_c.dish_id
	LEFT JOIN categories AS c ON f_c.category_id=c.id
	WHERE d.id=ANY($1)
	GROUP BY d.id,d.name,d.description,d.price,d.image_id
	ORDER BY d.id`

	var res []entity.Dish
	err := r.cli.SelectContext(ctx, &res, query, ids)
	if err != nil {
		return nil, errors.WithMessage(err, "get dish list")
	}
	return res, nil
}
