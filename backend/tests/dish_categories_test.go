// nolint:noctx,funlen
package tests_test

import (
	"bytes"
	"context"
	"dish_as_a_service/assembly"
	"dish_as_a_service/domain"
	"dish_as_a_service/entity"
	"dish_as_a_service/repository"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Falokut/go-kit/client/db"
	"github.com/Falokut/go-kit/http/apierrors"
	"github.com/Falokut/go-kit/json"
	"github.com/Falokut/go-kit/test"
	"github.com/Falokut/go-kit/test/dbt"
	"github.com/Falokut/go-kit/test/fake"
	"github.com/Falokut/go-kit/test/telegramt"
	"github.com/stretchr/testify/suite"
	"github.com/txix-open/bgjob"
)

type DishCategoriesSuite struct {
	suite.Suite
	test *test.Test

	db                 *dbt.TestDb
	dishCategoriesRepo repository.DishesCategories
	dishRepo           repository.Dish
	cli                *http.Client
	serverAddr         string
}

func TestDishCategories(t *testing.T) {
	t.Parallel()
	suite.Run(t, &DishCategoriesSuite{})
}

// nolint:dupl
func (t *DishCategoriesSuite) SetupTest() {
	test, _ := test.New(t.T())
	t.test = test
	t.db = dbt.New(test, db.WithMigrationRunner("../migrations", test.Logger()))
	t.dishCategoriesRepo = repository.NewDishesCategories(t.db.Client)
	t.dishRepo = repository.NewDish(t.db.Client)

	bgjobDb := bgjob.NewPgStore(t.db.Client.DB.DB)
	bgjobCli := bgjob.NewClient(bgjobDb)
	tgBot, _ := telegramt.TestBot(test)
	locatorCfg, err := assembly.Locator(
		context.Background(),
		test.Logger(),
		t.db.Client,
		tgBot,
		bgjobCli,
		getConfig(),
	)
	t.Require().NoError(err)
	server := httptest.NewServer(locatorCfg.HttpRouter)
	t.serverAddr = server.Listener.Addr().String()
	t.cli = server.Client()
	t.T().Cleanup(func() {
		server.Close()
	})
}

func (t *DishCategoriesSuite) getServerUrl(endpoint string) string {
	return fmt.Sprintf("http://%s/%s", t.serverAddr, endpoint)
}

func (t *DishCategoriesSuite) Test_GetAllCategories_HappyPath() {
	resp, err := t.cli.Get(t.getServerUrl("/dishes/all_categories"))
	t.Require().NoError(err)
	defer resp.Body.Close()

	var categories []domain.DishCategory
	err = json.NewDecoder(resp.Body).Decode(&categories)
	t.Require().NoError(err)
	t.Require().ElementsMatch([]domain.DishCategory{
		{
			Id:   1,
			Name: "Горячее",
		},
		{
			Id:   2,
			Name: "Холодное",
		},
		{
			Id:   3,
			Name: "Напиток",
		},
		{
			Id:   4,
			Name: "Острое",
		},
		{
			Id:   5,
			Name: "Рыба",
		},
		{
			Id:   6,
			Name: "Вегетарианское",
		},
		{
			Id:   7,
			Name: "Мясное",
		},
	}, categories)
}

func (t *DishCategoriesSuite) Test_GetDishesCategories_HappyPath() {
	resp, err := t.cli.Get(t.getServerUrl("/dishes/categories"))
	t.Require().NoError(err)
	defer resp.Body.Close()

	var categories []domain.DishCategory
	err = json.NewDecoder(resp.Body).Decode(&categories)
	t.Require().NoError(err)
	t.Require().Empty(categories)

	err = t.dishRepo.AddDish(context.Background(), &entity.AddDishRequest{
		Name:        "dish",
		Description: "desc",
		ImageId:     "image_id",
		Price:       1000,
		Categories:  []int32{1, 2},
	})
	t.Require().NoError(err)

	resp, err = t.cli.Get(t.getServerUrl("/dishes/categories"))
	t.Require().NoError(err)
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&categories)
	t.Require().NoError(err)

	t.Require().ElementsMatch([]domain.DishCategory{
		{
			Id:   1,
			Name: "Горячее",
		},
		{
			Id:   2,
			Name: "Холодное",
		},
	}, categories)
}

func (t *DishCategoriesSuite) Test_GetCategory_HappyPath() {
	const categoryId = 3
	resp, err := t.cli.Get(
		t.getServerUrl(fmt.Sprintf("/dishes/categories/%d", categoryId)),
	)
	t.Require().NoError(err)
	defer resp.Body.Close()

	var category domain.DishCategory
	err = json.NewDecoder(resp.Body).Decode(&category)
	t.Require().NoError(err)
	t.Require().Equal(domain.DishCategory{Id: categoryId, Name: "Напиток"}, category)
}

func (t *DishCategoriesSuite) Test_AddCategory_HappyPath() {
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

	categoryName := fake.It[string]()
	body, err := json.Marshal(domain.AddCategoryRequest{Name: categoryName})
	t.Require().NoError(err)

	reqBody := bytes.NewReader(body)
	req, err := http.NewRequest(http.MethodPost,
		t.getServerUrl("dishes/categories"), reqBody)
	t.Require().NoError(err)

	req.Header.Add("X-User-Id", userId)
	req.Header.Add("Content-Type", "application/json")
	resp, err := t.cli.Do(req)
	t.Require().NoError(err)
	defer resp.Body.Close()

	var categoryId domain.AddCategoryResponse
	err = json.NewDecoder(resp.Body).Decode(&categoryId)
	t.Require().NoError(err)

	category, err := t.dishCategoriesRepo.GetCategory(context.Background(), categoryId.Id)
	t.Require().NoError(err)
	t.Require().Equal(categoryName, category.Name)
}

func (t *DishCategoriesSuite) Test_DeleteCategory_HappyPath() {
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

	categoryName := fake.It[string]()
	body, err := json.Marshal(domain.AddCategoryRequest{Name: categoryName})
	t.Require().NoError(err)

	reqBody := bytes.NewReader(body)
	req, err := http.NewRequest(http.MethodPost,
		t.getServerUrl("dishes/categories"), reqBody)
	t.Require().NoError(err)

	req.Header.Add("X-User-Id", userId)
	req.Header.Add("Content-Type", "application/json")
	resp, err := t.cli.Do(req)
	t.Require().NoError(err)
	defer resp.Body.Close()

	var categoryId domain.AddCategoryResponse
	err = json.NewDecoder(resp.Body).Decode(&categoryId)
	t.Require().NoError(err)

	category, err := t.dishCategoriesRepo.GetCategory(context.Background(), categoryId.Id)
	t.Require().NoError(err)
	t.Require().Equal(categoryName, category.Name)

	req, err = http.NewRequest(
		http.MethodDelete,
		t.getServerUrl(fmt.Sprintf("dishes/categories/%d", categoryId.Id)),
		http.NoBody,
	)
	t.Require().NoError(err)

	req.Header.Add("X-User-Id", userId)
	resp, err = t.cli.Do(req)
	t.Require().NoError(err)
	defer resp.Body.Close()

	resp, err = t.cli.Get(t.getServerUrl(fmt.Sprintf("/dishes/categories/%d", categoryId.Id)))
	t.Require().NoError(err)
	defer resp.Body.Close()
	t.Require().NoError(err)
	t.Require().EqualValues(http.StatusNotFound, resp.StatusCode)

	var errorResp apierrors.Error
	err = json.NewDecoder(resp.Body).Decode(&errorResp)
	t.Require().NoError(err)
	t.Require().EqualValues(domain.ErrCodeDishCategoryNotFound, errorResp.ErrorCode)
}

func (t *DishCategoriesSuite) Test_RenameCategory_HappyPath() {
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

	categoryName := fake.It[string]()
	body, err := json.Marshal(domain.AddCategoryRequest{Name: categoryName})
	t.Require().NoError(err)

	reqBody := bytes.NewReader(body)
	req, err := http.NewRequest(
		http.MethodPost,
		t.getServerUrl("dishes/categories"),
		reqBody,
	)
	t.Require().NoError(err)

	req.Header.Add("X-User-Id", userId)
	req.Header.Add("Content-Type", "application/json")
	resp, err := t.cli.Do(req)
	t.Require().NoError(err)
	defer resp.Body.Close()

	var categoryId domain.AddCategoryResponse
	err = json.NewDecoder(resp.Body).Decode(&categoryId)
	t.Require().NoError(err)

	category, err := t.dishCategoriesRepo.GetCategory(context.Background(), categoryId.Id)
	t.Require().NoError(err)
	t.Require().Equal(categoryName, category.Name)

	newCategoryName := fake.It[string]()
	body, err = json.Marshal(domain.RenameCategoryRequest{Name: newCategoryName})
	t.Require().NoError(err)

	reqBody = bytes.NewReader(body)
	req, err = http.NewRequest(
		http.MethodPost,
		t.getServerUrl(fmt.Sprintf("dishes/categories/%d", categoryId.Id)),
		reqBody,
	)
	t.Require().NoError(err)

	req.Header.Add("X-User-Id", userId)
	req.Header.Add("Content-Type", "application/json")
	resp, err = t.cli.Do(req)
	t.Require().NoError(err)
	defer resp.Body.Close()
	t.Require().EqualValues(http.StatusOK, resp.StatusCode)

	category, err = t.dishCategoriesRepo.GetCategory(context.Background(), categoryId.Id)
	t.Require().NoError(err)
	t.Require().Equal(newCategoryName, category.Name)
}

func (t *DishCategoriesSuite) Test_RenameCategory_Conflict() {
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

	categoryName := fake.It[string]()
	body, err := json.Marshal(domain.AddCategoryRequest{Name: categoryName})
	t.Require().NoError(err)

	reqBody := bytes.NewReader(body)
	req, err := http.NewRequest(
		http.MethodPost,
		t.getServerUrl("dishes/categories"),
		reqBody,
	)
	t.Require().NoError(err)

	req.Header.Add("X-User-Id", userId)
	req.Header.Add("Content-Type", "application/json")
	resp, err := t.cli.Do(req)
	t.Require().NoError(err)
	defer resp.Body.Close()

	var categoryId domain.AddCategoryResponse
	err = json.NewDecoder(resp.Body).Decode(&categoryId)
	t.Require().NoError(err)

	category, err := t.dishCategoriesRepo.GetCategory(context.Background(), categoryId.Id)
	t.Require().NoError(err)
	t.Require().Equal(categoryName, category.Name)

	categoryName = fake.It[string]()
	body, err = json.Marshal(domain.AddCategoryRequest{Name: categoryName})
	t.Require().NoError(err)

	reqBody = bytes.NewReader(body)
	req, err = http.NewRequest(
		http.MethodPost,
		t.getServerUrl("dishes/categories"),
		reqBody,
	)
	t.Require().NoError(err)

	req.Header.Add("X-User-Id", userId)
	req.Header.Add("Content-Type", "application/json")
	resp, err = t.cli.Do(req)
	t.Require().NoError(err)
	defer resp.Body.Close()

	body, err = json.Marshal(domain.RenameCategoryRequest{Name: categoryName})
	t.Require().NoError(err)

	req, err = http.NewRequest(
		http.MethodPost,
		t.getServerUrl(fmt.Sprintf("dishes/categories/%d", categoryId.Id)),
		bytes.NewReader(body),
	)
	t.Require().NoError(err)

	req.Header.Add("X-User-Id", userId)
	req.Header.Add("Content-Type", "application/json")
	resp, err = t.cli.Do(req)
	t.Require().NoError(err)
	defer resp.Body.Close()
	t.Require().EqualValues(http.StatusConflict, resp.StatusCode)
}
