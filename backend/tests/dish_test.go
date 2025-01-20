// //nolint:noctx
package tests_test

import (
	"context"
	"database/sql"
	"dish_as_a_service/assembly"
	"dish_as_a_service/conf"
	"dish_as_a_service/domain"
	"dish_as_a_service/entity"
	"dish_as_a_service/jwt"
	"dish_as_a_service/repository"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Falokut/go-kit/client/db"
	"github.com/Falokut/go-kit/http/client"
	"github.com/Falokut/go-kit/test"
	"github.com/Falokut/go-kit/test/dbt"
	"github.com/Falokut/go-kit/test/fake"
	"github.com/Falokut/go-kit/test/telegramt"
	"github.com/stretchr/testify/suite"
	"github.com/txix-open/bgjob"
)

type DishSuite struct {
	suite.Suite
	test             *test.Test
	adminAccessToken string
	userAccessToken  string

	db       *dbt.TestDb
	dishRepo repository.Dish
	cli      *client.Client
}

func TestDish(t *testing.T) {
	t.Parallel()
	suite.Run(t, &DishSuite{})
}

// nolint:dupl
func (t *DishSuite) SetupTest() {
	test, _ := test.New(t.T())
	t.test = test
	t.db = dbt.New(test, db.WithMigrationRunner("../migrations", test.Logger()))
	t.dishRepo = repository.NewDish(t.db.Client)

	bgjobDb := bgjob.NewPgStore(t.db.Client.DB.DB)
	bgjobCli := bgjob.NewClient(bgjobDb)
	tgBot, _ := telegramt.TestBot(test)
	cfg := getConfig()
	locatorCfg, err := assembly.Locator(
		context.Background(),
		test.Logger(),
		t.db.Client,
		nil,
		tgBot,
		bgjobCli,
		cfg,
	)
	t.Require().NoError(err)

	server := httptest.NewServer(locatorCfg.HttpRouter)
	t.cli = client.NewWithClient(server.Client())
	t.cli.GlobalRequestConfig().BaseUrl = fmt.Sprintf("http://%s", server.Listener.Addr())

	var userId string
	err = t.db.Get(&userId,
		`INSERT INTO users(username,name,admin)
		VALUES($1,$2,$3)
		RETURNING id;`,
		"@admin",
		"test",
		true,
	)
	t.Require().NoError(err)

	jwtGen, err := jwt.GenerateToken(cfg.Auth.Access.Secret, cfg.Auth.Access.TTL, entity.TokenUserInfo{
		UserId:   userId,
		RoleName: domain.AdminType,
	})
	t.Require().NoError(err)

	t.adminAccessToken = domain.BearerToken + " " + jwtGen.Token

	err = t.db.Get(&userId,
		`INSERT INTO users(username,name,admin)
		VALUES($1,$2,$3)
		RETURNING id;`,
		"@user",
		"test",
		false,
	)
	t.Require().NoError(err)

	jwtGen, err = jwt.GenerateToken(cfg.Auth.Access.Secret, cfg.Auth.Access.TTL, entity.TokenUserInfo{
		UserId:   userId,
		RoleName: domain.UserType,
	})
	t.Require().NoError(err)

	t.userAccessToken = domain.BearerToken + " " + jwtGen.Token

	t.T().Cleanup(func() {
		server.Close()
	})
}

func (t *DishSuite) Test_List_ByLimitOffset_HappyPath() {
	var addDish = entity.AddDishRequest{
		Name:        fake.It[string](),
		Description: fake.It[string](),
		Price:       350,
		ImageId:     fake.It[string](),
		Categories:  []int32{1, 2},
	}
	err := t.dishRepo.AddDish(context.Background(), &addDish)
	t.Require().NoError(err)

	var dishes []domain.Dish
	_, err = t.cli.Get("/dishes").
		QueryParams(map[string]any{
			"limit":  1,
			"offset": 0,
		}).
		JsonResponseBody(&dishes).
		Do(context.Background())
	t.Require().NoError(err)
	t.Require().Len(dishes, 1)

	dish := dishes[0]
	t.Require().Equal(addDish.Name, dish.Name)
	t.Require().Equal(addDish.Description, dish.Description)
	t.Require().Equal("my_image_path/image-dish/"+addDish.ImageId, dish.Url)
	t.Require().ElementsMatch([]string{"Горячее", "Холодное"}, dish.Categories)
}

// nolint:dupl
func (t *DishSuite) Test_List_ByIds_HappyPath() {
	var addDish = entity.AddDishRequest{
		Name:        fake.It[string](),
		Description: fake.It[string](),
		Price:       350,
		ImageId:     fake.It[string](),
		Categories:  []int32{1, 2},
	}
	err := t.dishRepo.AddDish(context.Background(), &addDish)
	t.Require().NoError(err)

	addDish = entity.AddDishRequest{
		Name:        fake.It[string](),
		Description: fake.It[string](),
		Price:       544,
		ImageId:     fake.It[string](),
		Categories:  []int32{3, 2},
	}
	err = t.dishRepo.AddDish(context.Background(), &addDish)
	t.Require().NoError(err)

	var dishes []domain.Dish
	_, err = t.cli.Get("/dishes").
		JsonResponseBody(&dishes).
		QueryParams(map[string]any{
			"ids": "2",
		}).
		StatusCodeToError().
		Do(context.Background())
	t.Require().NoError(err)
	t.Require().Len(dishes, 1)

	dish := dishes[0]
	t.Require().EqualValues(2, dish.Id)
	t.Require().Equal(addDish.Name, dish.Name)
	t.Require().Equal(addDish.Description, dish.Description)
	t.Require().Equal("my_image_path/image-dish/"+addDish.ImageId, dish.Url)
	t.Require().ElementsMatch([]string{"Напиток", "Холодное"}, dish.Categories)
}

// nolint:dupl
func (t *DishSuite) Test_List_ByCategories_HappyPath() {
	var addDish = entity.AddDishRequest{
		Name:        fake.It[string](),
		Description: fake.It[string](),
		Price:       350,
		ImageId:     fake.It[string](),
		Categories:  []int32{4, 5},
	}
	err := t.dishRepo.AddDish(context.Background(), &addDish)
	t.Require().NoError(err)

	addDish = entity.AddDishRequest{
		Name:        fake.It[string](),
		Description: fake.It[string](),
		Price:       544,
		ImageId:     fake.It[string](),
		Categories:  []int32{4, 2},
	}
	err = t.dishRepo.AddDish(context.Background(), &addDish)
	t.Require().NoError(err)

	var dishes []domain.Dish
	_, err = t.cli.Get("/dishes").
		JsonResponseBody(&dishes).
		QueryParams(map[string]any{
			"categoriesIds": "2,4",
			"limit":         2,
		}).
		StatusCodeToError().
		Do(context.Background())
	t.Require().NoError(err)
	t.Require().Len(dishes, 1)

	dish := dishes[0]
	t.Require().EqualValues(2, dish.Id)
	t.Require().Equal(addDish.Name, dish.Name)
	t.Require().Equal(addDish.Description, dish.Description)
	t.Require().Equal("my_image_path/image-dish/"+addDish.ImageId, dish.Url)
	t.Require().ElementsMatch([]string{"Острое", "Холодное"}, dish.Categories)

	_, err = t.cli.Get("/dishes").
		JsonResponseBody(&dishes).
		QueryParams(map[string]any{
			"categoriesIds": "6,2",
			"limit":         2,
		}).
		Do(context.Background())
	t.Require().NoError(err)
	t.Require().Empty(dishes)
}

func (t *DishSuite) Test_AddDish_HappyPath() {
	req := domain.AddDishRequest{
		Name:        fake.It[string](),
		Description: fake.It[string](),
		Price:       900,
		Categories:  []int32{5, 6},
	}
	_, err := t.cli.Post("/dishes").
		Header(domain.AuthHeaderName, t.adminAccessToken).
		JsonRequestBody(req).
		Do(context.Background())
	t.Require().NoError(err)

	var ids []int32
	err = t.db.Select(&ids, "SELECT id FROM dish WHERE name=$1 AND description=$2 AND price=$3 AND image_id=$4",
		req.Name, req.Description, req.Price, "")
	t.Require().NoError(err)
	t.Require().Len(ids, 1)

	var categoriesIds []int32
	err = t.db.Select(&categoriesIds, "SELECT category_id FROM dish_categories WHERE dish_id=$1", ids[0])
	t.Require().NoError(err)
	t.Require().ElementsMatch(req.Categories, categoriesIds)
}

// nolint:funlen
func (t *DishSuite) Test_EditDish_HappyPath() {
	req := domain.AddDishRequest{
		Name:        fake.It[string](),
		Description: fake.It[string](),
		Price:       900,
		Categories:  []int32{5, 6},
	}
	_, err := t.cli.Post("/dishes").
		Header(domain.AuthHeaderName, t.adminAccessToken).
		JsonRequestBody(req).
		StatusCodeToError().
		Do(context.Background())
	t.Require().NoError(err)

	var ids []int32
	err = t.db.Select(&ids, "SELECT id FROM dish WHERE name=$1 AND description=$2 AND price=$3 AND image_id=$4",
		req.Name, req.Description, req.Price, "")
	t.Require().NoError(err)
	t.Require().Len(ids, 1)

	var categoriesIds []int32
	err = t.db.Select(&categoriesIds, "SELECT category_id FROM dish_categories WHERE dish_id=$1", ids[0])
	t.Require().NoError(err)
	t.Require().ElementsMatch(req.Categories, categoriesIds)

	editDishReq := domain.EditDishRequest{
		Name:        fake.It[string](),
		Description: fake.It[string](),
		Price:       800,
		Categories:  []int32{1, 2, 5},
	}
	t.Require().NoError(err)

	_, err = t.cli.Post(fmt.Sprintf("dishes/edit/%d", ids[0])).
		Header(domain.AuthHeaderName, t.adminAccessToken).
		JsonRequestBody(editDishReq).
		StatusCodeToError().
		Do(context.Background())
	t.Require().NoError(err)

	var dish entity.Dish
	err = t.db.Get(&dish, "SELECT id, name, description, price FROM dish WHERE id=$1", ids[0])
	t.Require().NoError(err)
	expectedDish := entity.Dish{
		Id:          ids[0],
		Name:        editDishReq.Name,
		Description: editDishReq.Description,
		Price:       editDishReq.Price,
	}
	t.Require().Equal(expectedDish, dish)

	err = t.db.Select(&categoriesIds, "SELECT category_id FROM dish_categories WHERE dish_id=$1", expectedDish.Id)
	t.Require().NoError(err)
	t.Require().ElementsMatch(editDishReq.Categories, categoriesIds)
}

func (t *DishSuite) Test_DeleteDish_HappyPath() {
	req := domain.AddDishRequest{
		Name:        fake.It[string](),
		Description: fake.It[string](),
		Price:       900,
		Categories:  []int32{5, 6},
	}
	_, err := t.cli.Post("/dishes").
		Header(domain.AuthHeaderName, t.adminAccessToken).
		JsonRequestBody(req).
		Do(context.Background())
	t.Require().NoError(err)

	var ids []int32
	err = t.db.Select(&ids, "SELECT id FROM dish WHERE name=$1 AND description=$2 AND price=$3 AND image_id=$4",
		req.Name, req.Description, req.Price, "")
	t.Require().NoError(err)
	t.Require().Len(ids, 1)

	var categoriesIds []int32
	err = t.db.Select(&categoriesIds, "SELECT category_id FROM dish_categories WHERE dish_id=$1", ids[0])
	t.Require().NoError(err)
	t.Require().ElementsMatch(req.Categories, categoriesIds)

	err = t.cli.Delete(fmt.Sprintf("dishes/delete/%d", ids[0])).
		Header(domain.AuthHeaderName, t.adminAccessToken).
		DoWithoutResponse(context.Background())
	t.Require().NoError(err)

	var dish entity.Dish
	err = t.db.Get(&dish, "SELECT id, name, description, price FROM dish WHERE id=$1", ids[0])
	t.Require().ErrorIs(err, sql.ErrNoRows)

	err = t.db.Select(&categoriesIds, "SELECT category_id FROM dish_categories WHERE dish_id=$1", ids[0])
	t.Require().NoError(err)
	t.Require().ElementsMatch([]int32{}, categoriesIds)
}

func (t *DishSuite) Test_AddDish_Forbidden() {
	addDishReq := domain.AddDishRequest{
		Name:        fake.It[string](),
		Description: fake.It[string](),
		Price:       800,
		Categories:  []int32{5, 6},
	}

	resp, err := t.cli.Post("/dishes").
		JsonRequestBody(addDishReq).
		Header(domain.AuthHeaderName, t.userAccessToken).
		Do(context.Background())
	t.Require().NoError(err)
	t.Require().Equal(http.StatusForbidden, resp.StatusCode())

	var ids []int32
	err = t.db.Select(&ids, "SELECT id FROM dish WHERE name=$1 AND description=$2 AND price=$3 AND image_id=$4",
		addDishReq.Name, addDishReq.Description, addDishReq.Price, "")
	t.Require().NoError(err)
	t.Require().Empty(ids)
}

func getConfig() conf.LocalConfig {
	return conf.LocalConfig{
		Images: conf.Images{
			BaseImagePath: "my_image_path",
		},
		App: conf.App{
			AdminSecret: "secret",
		},
		Auth: conf.Auth{
			Access: conf.JwtToken{
				TTL:    time.Hour,
				Secret: fake.It[string](),
			},
		},
	}
}
