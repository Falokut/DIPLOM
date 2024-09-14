//nolint:noctx
package tests_test

import (
	"bytes"
	"context"
	"dish_as_a_service/assembly"
	"dish_as_a_service/conf"
	"dish_as_a_service/domain"
	"dish_as_a_service/entity"
	"dish_as_a_service/repository"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Falokut/go-kit/client/db"
	"github.com/Falokut/go-kit/json"
	"github.com/Falokut/go-kit/test"
	"github.com/Falokut/go-kit/test/dbt"
	"github.com/stretchr/testify/suite"
	"github.com/txix-open/bgjob"
)

type DishSuite struct {
	suite.Suite
	test *test.Test

	db         *dbt.TestDb
	dishRepo   repository.Dish
	cli        *http.Client
	serverAddr string
}

func TestDish(t *testing.T) {
	t.Parallel()
	suite.Run(t, &DishSuite{})
}

func (t *DishSuite) SetupTest() {
	test, _ := test.New(t.T())
	t.test = test
	t.db = dbt.New(test, db.WithMigrationRunner("../migrations", test.Logger()))
	t.dishRepo = repository.NewDish(t.db.Client)

	bgjobDb := bgjob.NewPgStore(t.db.Client.DB.DB)
	bgjobCli := bgjob.NewClient(bgjobDb)

	locatorCfg := assembly.Locator(test.Logger(), t.db.Client, nil, bgjobCli, getConfig())
	server := httptest.NewServer(locatorCfg.HttpRouter)
	t.serverAddr = server.Listener.Addr().String()
	t.cli = server.Client()
	t.T().Cleanup(func() {
		server.Close()
	})
}

func (t *DishSuite) getServerUrl(endpoint string) string {
	return fmt.Sprintf("http://%s/%s", t.serverAddr, endpoint)
}

func (t *DishSuite) Test_List_ByLimitOffset_HappyPath() {
	var addDish = entity.AddDishRequest{
		Name:        "dish_name",
		Description: "dish_desc",
		Price:       350,
		ImageId:     "some_id",
		Categories:  []int32{1, 2},
	}
	err := t.dishRepo.AddDish(context.Background(), &addDish)
	t.Require().NoError(err)

	resp, err := t.cli.Get(t.getServerUrl("dishes"))
	t.Require().NoError(err)
	defer resp.Body.Close()

	dishes := []domain.Dish{}
	err = json.NewDecoder(resp.Body).Decode(&dishes)
	t.Require().NoError(err)
	t.Require().Len(dishes, 1)

	dish := dishes[0]
	t.Require().Equal(addDish.Name, dish.Name)
	t.Require().Equal(addDish.Description, dish.Description)
	t.Require().Equal("my_image_path/dish/"+addDish.ImageId, dish.Url)
	t.Require().ElementsMatch([]string{"Горячее", "Холодное"}, dish.Categories)
}

func (t *DishSuite) Test_List_ByIds_HappyPath() {
	var addDish = entity.AddDishRequest{
		Name:        "dish_name",
		Description: "dish_desc",
		Price:       350,
		ImageId:     "some_id",
		Categories:  []int32{1, 2},
	}
	err := t.dishRepo.AddDish(context.Background(), &addDish)
	t.Require().NoError(err)

	addDish = entity.AddDishRequest{
		Name:        "dish_name_2",
		Description: "dish_desc_2",
		Price:       350,
		ImageId:     "some_id_2",
		Categories:  []int32{3, 2},
	}
	err = t.dishRepo.AddDish(context.Background(), &addDish)
	t.Require().NoError(err)

	resp, err := t.cli.Get(t.getServerUrl("dishes?ids=2"))
	t.Require().NoError(err)
	defer resp.Body.Close()

	dishes := []domain.Dish{}
	err = json.NewDecoder(resp.Body).Decode(&dishes)
	t.Require().NoError(err)
	t.Require().Len(dishes, 1)

	dish := dishes[0]
	t.Require().EqualValues(2, dish.Id)
	t.Require().Equal(addDish.Name, dish.Name)
	t.Require().Equal(addDish.Description, dish.Description)
	t.Require().Equal("my_image_path/dish/"+addDish.ImageId, dish.Url)
	t.Require().ElementsMatch([]string{"Напиток", "Холодное"}, dish.Categories)
}

func (t *DishSuite) Test_AddDish_HappyPath() {
	var userId string
	err := t.db.Get(&userId,
		`INSERT INTO users(username,name,admin)
		VALUES($1,$2,$3)
		RETURNING id;`,
		"@test",
		"test",
		true,
	)
	t.Require().NoError(err)

	addDishReq := domain.AddDishRequest{
		Name:        "dish",
		Description: "test_desc",
		Price:       344,
		Categories:  []int32{5, 6},
	}
	body, err := json.Marshal(addDishReq)
	t.Require().NoError(err)

	reqBody := bytes.NewReader(body)
	req, err := http.NewRequest(http.MethodPost, t.getServerUrl("dishes"), reqBody)
	t.Require().NoError(err)

	req.Header.Add("X-User-Id", userId)
	req.Header.Add("Content-Type", "application/json")
	resp, err := t.cli.Do(req)
	t.Require().NoError(err)
	defer resp.Body.Close()

	var ids []int32
	err = t.db.Select(&ids, "SELECT id FROM dish WHERE name=$1 AND description=$2 AND price=$3 AND image_id=$4",
		addDishReq.Name, addDishReq.Description, addDishReq.Price, "")
	t.Require().NoError(err)
	t.Require().Len(ids, 1)

	var categoriesIds []int32
	err = t.db.Select(&categoriesIds, "SELECT category_id FROM dish_categories WHERE dish_id=$1", ids[0])
	t.Require().NoError(err)
	t.Require().ElementsMatch(addDishReq.Categories, categoriesIds)
}

func (t *DishSuite) Test_AddDish_Forbidden() {
	var userId string
	err := t.db.Get(&userId,
		`INSERT INTO users(username,name,admin)
		VALUES($1,$2,$3)
		RETURNING id;`,
		"@test",
		"test",
		false,
	)
	t.Require().NoError(err)

	addDishReq := domain.AddDishRequest{
		Name:        "dish",
		Description: "test_desc",
		Price:       344,
		Categories:  []int32{5, 6},
	}
	body, err := json.Marshal(addDishReq)
	t.Require().NoError(err)

	reqBody := bytes.NewReader(body)
	req, err := http.NewRequest(http.MethodPost, t.getServerUrl("dishes"), reqBody)
	t.Require().NoError(err)

	req.Header.Add("X-User-Id", userId)
	req.Header.Add("Content-Type", "application/json")
	resp, err := t.cli.Do(req)
	t.Require().NoError(err)
	defer resp.Body.Close()
	t.Require().Equal(http.StatusForbidden, resp.StatusCode)

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
	}
}
